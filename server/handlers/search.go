package handlers

import (
	"context"
	"fmt"
	"net/http"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	search *service.SearchService
}

// NewSearchHandler creates a new SearchHandler.
func NewSearchHandler(search *service.SearchService) *SearchHandler {
	return &SearchHandler{search: search}
}

// Search handles the search request.
// GET /api/search?q=query&limit=10
func (h *SearchHandler) Search(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uid, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := parseInt(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	ctx := context.Background()
	results, err := h.search.SearchCards(ctx, uid, query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

// Suggestions handles search suggestions.
// GET /api/search/suggestions?q=prefix&limit=5
func (h *SearchHandler) Suggestions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uid, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{"suggestions": []string{}})
		return
	}

	limit := 5
	if l := c.Query("limit"); l != "" {
		if parsed, err := parseInt(l); err == nil && parsed > 0 && parsed <= 20 {
			limit = parsed
		}
	}

	ctx := context.Background()
	suggestions, err := h.search.GetSearchSuggestions(ctx, uid, query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":       query,
		"suggestions": suggestions,
	})
}

// Helper to parse int
func parseInt(s string) (int, error) {
	var n int
	fmt.Sscanf(s, "%d", &n)
	if n == 0 && s != "0" {
		return 0, fmt.Errorf("not a valid integer")
	}
	return n, nil
}
