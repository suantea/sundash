package handlers

import (
	"net/http"
	"strings"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

type FaviconHandler struct {
	favicon *service.FaviconService
}

func NewFaviconHandler(favicon *service.FaviconService) *FaviconHandler {
	return &FaviconHandler{favicon: favicon}
}

func (h *FaviconHandler) FetchFavicon(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	resp, err := h.favicon.FetchFavicon(url)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// BatchFetchFavicons handles POST /api/favicons with JSON body { "urls": ["url1", "url2", ...] }
func (h *FaviconHandler) BatchFetchFavicons(c *gin.Context) {
	var req struct {
		URLs []string `json:"urls" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "urls array is required"})
		return
	}

	results := h.favicon.BatchFetchFavicons(req.URLs)
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// ExtractDomain 辅助函数：从 URL 提取域名（供前端导入时使用）
func extractDomain(rawURL string) string {
	domain := strings.TrimPrefix(rawURL, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	if idx := strings.IndexAny(domain, "/:?#"); idx != -1 {
		domain = domain[:idx]
	}
	return domain
}
