package handlers

import (
	"net/http"

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
