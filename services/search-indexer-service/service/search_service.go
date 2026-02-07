package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type SearchService struct {
	opensearchURL string
	httpClient    *http.Client
}

func NewSearchService() *SearchService {
	opensearchURL := os.Getenv("OPENSEARCH_URL")
	if opensearchURL == "" {
		opensearchURL = "http://localhost:9200"
	}

	return &SearchService{
		opensearchURL: opensearchURL,
		httpClient:    &http.Client{},
	}
}

type SearchResult struct {
	Type        string
	ID          string
	Title       string
	Description string
	Fields      map[string]interface{}
	Score       float64
}

func (s *SearchService) Search(
	query string,
	searchType string,
	page int,
	pageSize int,
	sort string,
	filters map[string]string,
	isGuest bool,
	tenantID string,
	roles []string,
) ([]SearchResult, int64, map[string]interface{}, error) {
	// Build indices to search
	indices := s.getIndicesForType(searchType)

	// Build OpenSearch query
	esQuery := s.buildSearchQuery(query, searchType, page, pageSize, sort, filters, isGuest, tenantID, roles)

	// Execute multi-index search
	url := fmt.Sprintf("%s/%s/_search", s.opensearchURL, strings.Join(indices, ","))
	
	jsonData, err := json.Marshal(esQuery)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("opensearch request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, nil, fmt.Errorf("opensearch error: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse results
	hits, _ := result["hits"].(map[string]interface{})
	total, _ := hits["total"].(map[string]interface{})
	totalValue, _ := total["value"].(float64)
	
	hitsArray, _ := hits["hits"].([]interface{})
	results := make([]SearchResult, 0, len(hitsArray))

	for _, hit := range hitsArray {
		hitMap, _ := hit.(map[string]interface{})
		source, _ := hitMap["_source"].(map[string]interface{})
		score, _ := hitMap["_score"].(float64)
		index, _ := hitMap["_index"].(string)

		resultType := s.getTypeFromIndex(index)
		resultID, _ := source["id"].(string)
		
		title := s.extractTitle(source, resultType)
		description := s.extractDescription(source, resultType)

		results = append(results, SearchResult{
			Type:        resultType,
			ID:          resultID,
			Title:       title,
			Description: description,
			Fields:      source,
			Score:       score,
		})
	}

	// Extract facets
	facets := s.extractFacets(result)

	return results, int64(totalValue), facets, nil
}

func (s *SearchService) Autocomplete(query string) ([]string, error) {
	// Use edge ngram analyzer for autocomplete
	url := fmt.Sprintf("%s/parts,equipment,companies,listings/_search", s.opensearchURL)
	
	esQuery := map[string]interface{}{
		"size": 10,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name.autocomplete^2", "part_number.autocomplete^3", "model.autocomplete^2"},
				"type":   "bool_prefix",
			},
		},
		"_source": []string{"name", "part_number", "model"},
	}

	jsonData, err := json.Marshal(esQuery)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("opensearch autocomplete error: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits, _ := result["hits"].(map[string]interface{})
	hitsArray, _ := hits["hits"].([]interface{})
	suggestions := make([]string, 0, len(hitsArray))

	for _, hit := range hitsArray {
		hitMap, _ := hit.(map[string]interface{})
		source, _ := hitMap["_source"].(map[string]interface{})
		
		if name, ok := source["name"].(string); ok && name != "" {
			suggestions = append(suggestions, name)
		} else if partNumber, ok := source["part_number"].(string); ok && partNumber != "" {
			suggestions = append(suggestions, partNumber)
		} else if model, ok := source["model"].(string); ok && model != "" {
			suggestions = append(suggestions, model)
		}
	}

	return suggestions, nil
}

func (s *SearchService) getIndicesForType(searchType string) []string {
	switch searchType {
	case "part":
		return []string{"parts"}
	case "equipment":
		return []string{"equipment"}
	case "company":
		return []string{"companies"}
	case "listing":
		return []string{"listings"}
	case "service":
		return []string{"listings"} // Services are also in listings
	default:
		return []string{"parts", "equipment", "companies", "listings"}
	}
}

func (s *SearchService) buildSearchQuery(
	query string,
	searchType string,
	page int,
	pageSize int,
	sort string,
	filters map[string]string,
	isGuest bool,
	tenantID string,
	roles []string,
) map[string]interface{} {
	// Build must clauses (required)
	mustClauses := []map[string]interface{}{}

	// Guest restrictions
	if isGuest {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"visibility": "public",
			},
		})
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"status": "approved",
			},
		})
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"company_status": "approved",
			},
		})
	} else {
		// Authenticated user: apply tenant/RBAC rules
		shouldClauses := []map[string]interface{}{
			// Public and approved
			{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"term": map[string]interface{}{"visibility": "public"}},
						{"term": map[string]interface{}{"status": "approved"}},
					},
				},
			},
		}

		// Buyers can see supplier listings
		if s.hasRole(roles, "buyer") || s.hasRole(roles, "requester") || s.hasRole(roles, "procurement_manager") {
			shouldClauses = append(shouldClauses, map[string]interface{}{
				"term": map[string]interface{}{
					"type": "listing",
				},
			})
		}

		// Suppliers see own listings
		if s.hasRole(roles, "supplier") && tenantID != "" {
			shouldClauses = append(shouldClauses, map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"term": map[string]interface{}{"supplier_id": tenantID}},
					},
				},
			})
		}

		// Admins see everything
		if s.hasRole(roles, "admin") || s.hasRole(roles, "super_admin") {
			shouldClauses = []map[string]interface{}{{"match_all": map[string]interface{}{}}}
		}

		mustClauses = append(mustClauses, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": shouldClauses,
				"minimum_should_match": 1,
			},
		})
	}

	// Add filters
	for key, value := range filters {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				key: value,
			},
		})
	}

	// Build query
	esQuery := map[string]interface{}{
		"from": (page - 1) * pageSize,
		"size": pageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  query,
							"fields": s.getBoostedFields(searchType),
							"type":   "best_fields",
							"operator": "and",
						},
					},
				},
				"filter": mustClauses,
			},
		},
	}

	// Add sorting
	switch sort {
	case "rating":
		esQuery["sort"] = []map[string]interface{}{
			{"rating": map[string]interface{}{"order": "desc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "price":
		esQuery["sort"] = []map[string]interface{}{
			{"price": map[string]interface{}{"order": "asc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "eta":
		esQuery["sort"] = []map[string]interface{}{
			{"eta": map[string]interface{}{"order": "asc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	default: // relevance
		esQuery["sort"] = []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	}

	// Add aggregations for facets
	esQuery["aggs"] = map[string]interface{}{
		"types": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "type",
			},
		},
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category",
			},
		},
	}

	return esQuery
}

func (s *SearchService) getBoostedFields(searchType string) []string {
	baseFields := []string{
		"name^3",
		"description^1",
	}

	switch searchType {
	case "part":
		return append(baseFields, "part_number^5", "manufacturer_code^4", "manufacturer^2")
	case "equipment":
		return append(baseFields, "model^4", "series^3", "manufacturer^2")
	case "company":
		return []string{"name^5", "subdomain^2"}
	case "listing":
		return append(baseFields, "sku^4", "brand^2")
	default:
		return append(baseFields, "part_number^5", "model^4", "sku^4")
	}
}

func (s *SearchService) extractTitle(source map[string]interface{}, resultType string) string {
	switch resultType {
	case "part":
		if name, ok := source["name"].(string); ok {
			return name
		}
		if partNumber, ok := source["part_number"].(string); ok {
			return partNumber
		}
	case "equipment":
		if model, ok := source["model"].(string); ok {
			return model
		}
		if name, ok := source["name"].(string); ok {
			return name
		}
	case "company":
		if name, ok := source["name"].(string); ok {
			return name
		}
	case "listing":
		if title, ok := source["title"].(string); ok {
			return title
		}
		if name, ok := source["name"].(string); ok {
			return name
		}
	}
	return "Untitled"
}

func (s *SearchService) extractDescription(source map[string]interface{}, resultType string) string {
	if desc, ok := source["description"].(string); ok && desc != "" {
		return desc
	}
	return ""
}

func (s *SearchService) getTypeFromIndex(index string) string {
	if strings.HasPrefix(index, "part") {
		return "part"
	} else if strings.HasPrefix(index, "equipment") {
		return "equipment"
	} else if strings.HasPrefix(index, "compan") {
		return "company"
	} else if strings.HasPrefix(index, "listing") {
		return "listing"
	}
	return "unknown"
}

func (s *SearchService) extractFacets(result map[string]interface{}) map[string]interface{} {
	facets := make(map[string]interface{})
	
	aggs, ok := result["aggregations"].(map[string]interface{})
	if !ok {
		return facets
	}

	if types, ok := aggs["types"].(map[string]interface{}); ok {
		if buckets, ok := types["buckets"].([]interface{}); ok {
			facets["types"] = buckets
		}
	}

	if categories, ok := aggs["categories"].(map[string]interface{}); ok {
		if buckets, ok := categories["buckets"].([]interface{}); ok {
			facets["categories"] = buckets
		}
	}

	return facets
}

func (s *SearchService) hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
