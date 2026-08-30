package handlers

import (
	"net/http"
	"strconv"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

type RSSHandler struct {
	rss *service.RSSService
}

// NewRSSHandler creates a new RSSHandler.
func NewRSSHandler(rss *service.RSSService) *RSSHandler {
	return &RSSHandler{rss: rss}
}

// ListFeeds handles GET /api/rss
func (h *RSSHandler) ListFeeds(c *gin.Context) {
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

	feeds, err := h.rss.GetFeeds(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feeds)
}

// AddFeed handles POST /api/rss
func (h *RSSHandler) AddFeed(c *gin.Context) {
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

	var req struct {
		URL string `json:"url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed, err := h.rss.AddFeed(c.Request.Context(), uid, req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, feed)
}

// UpdateFeed handles PUT /api/rss/:id
func (h *RSSHandler) UpdateFeed(c *gin.Context) {
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

	feedID := c.Param("id")
	if feedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feed ID is required"})
		return
	}

	var req struct {
		URL string `json:"url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feed, err := h.rss.UpdateFeed(c.Request.Context(), uid, feedID, req.URL)
	if err != nil {
		if err.Error() == "feed not found or unauthorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, feed)
}

// DeleteFeed handles DELETE /api/rss/:id
func (h *RSSHandler) DeleteFeed(c *gin.Context) {
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

	feedID := c.Param("id")
	if feedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feed ID is required"})
		return
	}

	err := h.rss.RemoveFeed(c.Request.Context(), uid, feedID)
	if err != nil {
		if err.Error() == "feed not found or unauthorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Feed deleted"})
}

// GetFeedItems handles GET /api/rss/:id/items
func (h *RSSHandler) GetFeedItems(c *gin.Context) {
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

	feedID := c.Param("id")
	if feedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feed ID is required"})
		return
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	items, err := h.rss.GetFeedItems(c.Request.Context(), uid, feedID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to frontend format
	var rssItems []map[string]interface{}
	for _, item := range items {
		rssItems = append(rssItems, map[string]interface{}{
			"id":          item.ID,
			"feed_id":     item.FeedID,
			"title":       item.Title,
			"link":        item.Link,
			"description": item.Description,
			"pub_date":    item.PubDate,
			"author":      item.Author,
			"guid":        item.Guid,
			"created_at":  item.CreatedAt,
			"updated_at":  item.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, rssItems)
}
