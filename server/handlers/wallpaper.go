package handlers

import (
	"net/http"
	"regexp"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

var datePattern = regexp.MustCompile(`^\d{8}$`)

type WallpaperHandler struct {
	wallpaper *service.WallpaperService
}

func NewWallpaperHandler(wallpaper *service.WallpaperService) *WallpaperHandler {
	return &WallpaperHandler{wallpaper: wallpaper}
}

func (h *WallpaperHandler) GetBingWallpaper(c *gin.Context) {
	wallpaper, err := h.wallpaper.GetToday()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, wallpaper)
}

func (h *WallpaperHandler) GetWallpaperByDate(c *gin.Context) {
	date := c.Param("date")
	if !datePattern.MatchString(date) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be in YYYYMMDD format"})
		return
	}
	wallpaper, err := h.wallpaper.GetByDate(date)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, wallpaper)
}
