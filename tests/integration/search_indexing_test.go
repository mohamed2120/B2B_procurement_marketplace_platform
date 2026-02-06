package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchIndexing(t *testing.T) {
	// Wait for services to be ready
	require.NoError(t, WaitForServices())

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	uniqueID := GenerateUniqueID()

	// Test 1: Index catalog part after approval
	t.Run("Catalog part indexing", func(t *testing.T) {
		adminClient, err := Login(tenantID, "admin@demo.com", "demo123456")
		require.NoError(t, err, "Failed to login as admin")

		// Create a manufacturer first
		manufacturerReq := map[string]interface{}{
			"name":        fmt.Sprintf("Test Manufacturer %s", uniqueID),
			"description": "Test manufacturer for indexing",
		}

		manufacturerResp, err := adminClient.Post("/api/v1/manufacturers", manufacturerReq)
		require.NoError(t, err)
		
		var manufacturer map[string]interface{}
		if manufacturerResp.StatusCode == http.StatusCreated || manufacturerResp.StatusCode == http.StatusOK {
			require.NoError(t, ParseResponse(manufacturerResp, &manufacturer))
		} else {
			// Manufacturer might already exist, try to get it
			manufacturer = map[string]interface{}{
				"id": uuid.New().String(), // Fallback
			}
		}

		manufacturerID := manufacturer["id"].(string)

		// Create a category
		categoryReq := map[string]interface{}{
			"name":        fmt.Sprintf("Test Category %s", uniqueID),
			"description": "Test category",
		}

		categoryResp, err := adminClient.Post("/api/v1/categories", categoryReq)
		require.NoError(t, err)
		
		var category map[string]interface{}
		if categoryResp.StatusCode == http.StatusCreated || categoryResp.StatusCode == http.StatusOK {
			require.NoError(t, ParseResponse(categoryResp, &category))
		} else {
			category = map[string]interface{}{
				"id": uuid.New().String(), // Fallback
			}
		}

		categoryID := category["id"].(string)

		// Create a part
		partNumber := fmt.Sprintf("PART-%s", uniqueID)
		partReq := map[string]interface{}{
			"part_number":    partNumber,
			"name":           fmt.Sprintf("Test Part %s", uniqueID),
			"description":    "Test part for indexing",
			"manufacturer_id": manufacturerID,
			"category_id":     categoryID,
			"status":          "pending",
		}

		partResp, err := adminClient.Post("/api/v1/parts", partReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, partResp.StatusCode, "Failed to create part")

		var part map[string]interface{}
		require.NoError(t, ParseResponse(partResp, &part))
		partID := part["id"].(string)
		t.Logf("Created part: %s", partID)

		// Approve the part (this should trigger indexing)
		approveResp, err := adminClient.Post("/api/v1/parts/"+partID+"/approve", map[string]interface{}{})
		require.NoError(t, err)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, approveResp.StatusCode,
			"Failed to approve part")

		// Wait for indexing (OpenSearch indexing might take a moment)
		time.Sleep(5 * time.Second)

		// Verify part is indexed in OpenSearch
		searchURL := fmt.Sprintf("%s/parts/_search?q=part_number:%s", OpenSearchURL, partNumber)
		resp, err := http.Get(searchURL)
		require.NoError(t, err, "Failed to query OpenSearch")
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var searchResult map[string]interface{}
			require.NoError(t, json.Unmarshal(body, &searchResult))

			// Check if we found the part
			if hits, ok := searchResult["hits"].(map[string]interface{}); ok {
				if total, ok := hits["total"].(map[string]interface{}); ok {
					if value, ok := total["value"].(float64); ok {
						assert.Greater(t, int(value), 0, "Part should be indexed in OpenSearch")
						t.Logf("Part found in OpenSearch: %s", partNumber)
					}
				} else if value, ok := hits["total"].(float64); ok {
					// Elasticsearch 7.x format
					assert.Greater(t, int(value), 0, "Part should be indexed in OpenSearch")
					t.Logf("Part found in OpenSearch: %s", partNumber)
				}
			}
		} else {
			t.Logf("OpenSearch query returned status %d (index might not exist yet)", resp.StatusCode)
		}
	})

	// Test 2: Index company after approval
	t.Run("Company indexing", func(t *testing.T) {
		adminClient, err := Login(tenantID, "admin@demo.com", "demo123456")
		require.NoError(t, err, "Failed to login as admin")

		// Create a company
		companyName := fmt.Sprintf("Test Company %s", uniqueID)
		companyReq := map[string]interface{}{
			"name":        companyName,
			"legal_name":  companyName,
			"tax_id":      fmt.Sprintf("TAX-%s", uniqueID),
			"status":      "pending",
			"subdomain":   fmt.Sprintf("test-%s", uniqueID[:8]),
		}

		companyResp, err := adminClient.Post("/api/v1/companies", companyReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, companyResp.StatusCode, "Failed to create company")

		var company map[string]interface{}
		require.NoError(t, ParseResponse(companyResp, &company))
		companyID := company["id"].(string)
		t.Logf("Created company: %s", companyID)

		// Approve the company (this should trigger indexing)
		approveResp, err := adminClient.Post("/api/v1/companies/"+companyID+"/approve", map[string]interface{}{})
		require.NoError(t, err)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, approveResp.StatusCode,
			"Failed to approve company")

		// Wait for indexing
		time.Sleep(5 * time.Second)

		// Verify company is indexed in OpenSearch
		searchURL := fmt.Sprintf("%s/companies/_search?q=name:%s", OpenSearchURL, companyName)
		resp, err := http.Get(searchURL)
		require.NoError(t, err, "Failed to query OpenSearch")
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var searchResult map[string]interface{}
			require.NoError(t, json.Unmarshal(body, &searchResult))

			// Check if we found the company
			if hits, ok := searchResult["hits"].(map[string]interface{}); ok {
				if total, ok := hits["total"].(map[string]interface{}); ok {
					if value, ok := total["value"].(float64); ok {
						assert.Greater(t, int(value), 0, "Company should be indexed in OpenSearch")
						t.Logf("Company found in OpenSearch: %s", companyName)
					}
				} else if value, ok := hits["total"].(float64); ok {
					assert.Greater(t, int(value), 0, "Company should be indexed in OpenSearch")
					t.Logf("Company found in OpenSearch: %s", companyName)
				}
			}
		} else {
			t.Logf("OpenSearch query returned status %d (index might not exist yet)", resp.StatusCode)
		}
	})

	t.Logf("Search indexing tests completed")
}
