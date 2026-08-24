package handlers

import (
	"fmt"
	"net/http"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	weather *service.WeatherService
}

// NewWeatherHandler creates a new WeatherHandler.
func NewWeatherHandler(weather *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{weather: weather}
}

// GetWeather returns the current weather for the default location (can be configured via settings).
// For simplicity, we use a fixed location (e.g., user's IP-based location or a default).
// In a real scenario, we might get location from user settings or IP geolocation.
// Here we use a default (e.g., 0,0) but we can change to a configurable one.
// For now, we'll use a hardcoded example: Beijing (39.9042, 116.4074).
func (h *WeatherHandler) GetWeather(c *gin.Context) {
	// Default location: Beijing, China
	latitude := 39.9042
	longitude := 116.4074

	// Allow override via query parameters for testing
	if lat := c.Query("lat"); lat != "" {
		if parsed, err := parseFloat(lat); err == nil {
			latitude = parsed
		}
	}
	if lon := c.Query("lon"); lon != "" {
		if parsed, err := parseFloat(lon); err == nil {
			longitude = parsed
		}
	}

	weatherData, err := h.weather.FetchWeather(latitude, longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to model for frontend
	c.JSON(http.StatusOK, weatherData.WeatherToModel())
}

// Helper to parse float64
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}