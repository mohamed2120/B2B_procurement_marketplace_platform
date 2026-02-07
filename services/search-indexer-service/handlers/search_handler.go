package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/b2b-platform/search-indexer-service/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

type SearchRequest struct {
	Query    string            `form:"q" binding:"required"`
	Type     string            `form:"type"` // all, part, equipment, company, listing, service
	Page     int               `form:"page"`
	PageSize int               `form:"page_size"`
	Filters  map[string]string `form:"filters"`
	Sort     string            `form:"sort"` // relevance, rating, price, eta
}

type SearchResponse struct {
	Results []SearchResult          `json:"results"`
	Facets  map[string]interface{} `json:"facets"`
	Total   int64                   `json:"total"`
	Page    int                     `json:"page"`
	PageSize int                    `json:"page_size"`
}

type SearchResult struct {
	Type        string                 `json:"type"` // part, equipment, company, listing, service
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Fields      map[string]interface{} `json:"fields"`
	Score       float64                `json:"score,omitempty"`
}

// Search handles the unified search endpoint
func (h *SearchHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Type == "" {
		req.Type = "all"
	}
	if req.Sort == "" {
		req.Sort = "relevance"
	}

	// Get user context from JWT (if present)
	userIDVal, _ := c.Get("user_id")
	tenantIDVal, _ := c.Get("tenant_id")
	rolesVal, _ := c.Get("roles")
	
	userID := ""
	tenantID := ""
	var roles []string
	
	if userIDVal != nil {
		userID = userIDVal.(string)
	}
	if tenantIDVal != nil {
		tenantID = tenantIDVal.(string)
	}
	if rolesVal != nil {
		if rolesSlice, ok := rolesVal.([]string); ok {
			roles = rolesSlice
		}
	}
	
	isGuest := userID == ""

	// Apply guest restrictions
	if isGuest {
		if req.PageSize > 10 {
			req.PageSize = 10
		}
	}

	// Perform search
	results, total, facets, err := h.searchService.Search(
		req.Query,
		req.Type,
		req.Page,
		req.PageSize,
		req.Sort,
		req.Filters,
		isGuest,
		tenantID,
		roles,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed: " + err.Error()})
		return
	}

	// Filter sensitive fields for guests
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = r
		if isGuest {
			// Remove sensitive fields
			delete(searchResults[i].Fields, "stock")
			delete(searchResults[i].Fields, "internal_notes")
			delete(searchResults[i].Fields, "email")
			delete(searchResults[i].Fields, "phone")
			// Mask price if restricted
			if price, ok := searchResults[i].Fields["price"].(float64); ok {
				if restricted, _ := searchResults[i].Fields["price_restricted"].(bool); restricted {
					searchResults[i].Fields["price"] = nil
					searchResults[i].Fields["price_display"] = "Contact for pricing"
				}
			}
		}
	}

	c.JSON(http.StatusOK, SearchResponse{
		Results:  searchResults,
		Facets:   facets,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// Autocomplete provides search suggestions
func (h *SearchHandler) Autocomplete(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	suggestions, err := h.searchService.Autocomplete(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Autocomplete failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
