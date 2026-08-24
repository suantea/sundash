package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"sundash/models"
)

// WeatherData represents the current weather data from Open-Meteo.
type WeatherData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	GenerationtimeMs float64 `json:"generationtime_ms"`
	UtcOffsetSeconds int `json:"utc_offset_seconds"`
	Timezone string `json:"timezone"`
	TimezoneAbbreviation string `json:"timezone_abbreviation"`
	Elevation float64 `json:"elevation"`
	CurrentWeather CurrentWeather `json:"current_weather"`
}

// CurrentWeather represents the current weather conditions.
type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	Windspeed float64 `json:"windspeed"`
	Winddirection float64 `json:"winddirection"`
	Weathercode int `json:"weathercode"`
	Time string `json:"time"`
}

// WeatherService provides weather data fetching and caching.
type WeatherService struct {
	mu        sync.RWMutex
	cache     *weatherCache
	httpClient *http.Client
}

// weatherCacheItem represents a single cached weather entry.
type weatherCacheItem struct {
	data      WeatherData
	expiry    time.Time
}

// weatherCache is a simple in-memory cache with expiration.
type weatherCache struct {
	items map[string]*weatherCacheItem
	mu    sync.RWMutex
}

// newWeatherCache creates a new weather cache.
func newWeatherCache() *weatherCache {
	return &weatherCache{
		items: make(map[string]*weatherCacheItem),
	}
}

// get retrieves an item from the cache if it exists and is not expired.
func (c *weatherCache) get(key string) (WeatherData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if item, ok := c.items[key]; ok {
		if time.Now().Before(item.expiry) {
			return item.data, true
		}
		// Expired, delete it
		delete(c.items, key)
	}
	return WeatherData{}, false
}

// set adds an item to the cache with a given duration.
func (c *weatherCache) set(key string, data WeatherData, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &weatherCacheItem{
		data:    data,
		expiry:  time.Now().Add(duration),
	}
}

// NewWeatherService creates a new WeatherService with a default cache duration of 10 minutes.
func NewWeatherService() *WeatherService {
	return &WeatherService{
		cache:     newWeatherCache(),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchWeather gets the weather data for the given latitude and longitude.
// It uses caching to avoid excessive API calls.
func (s *WeatherService) FetchWeather(latitude, longitude float64) (WeatherData, error) {
	// Create a cache key based on latitude and longitude (rounded to 2 decimal places for ~1km accuracy)
	key := fmt.Sprintf("%.2f-%.2f", latitude, longitude)

	// Check cache
	if data, ok := s.cache.get(key); ok {
		return data, nil
	}

	// If not in cache, fetch from Open-Meteo API
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true", latitude, longitude)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return WeatherData{}, fmt.Errorf("failed to decode weather response: %w", err)
	}

	// Store in cache for 10 minutes
	s.cache.set(key, weatherData, 10*time.Minute)

	return weatherData, nil
}

// WeatherToModel converts WeatherData to a simplified model for frontend display.
func (wd *WeatherData) WeatherToModel() models.Weather {
	return models.Weather{
		Temperature:     wd.CurrentWeather.Temperature,
		Windspeed:       wd.CurrentWeather.Windspeed,
		Winddirection:   wd.CurrentWeather.Winddirection,
		Weathercode:     wd.CurrentWeather.Weathercode,
		Time:            wd.CurrentWeather.Time,
		Timezone:        wd.Timezone,
		Latitude:        wd.Latitude,
		Longitude:       wd.Longitude,
	}
}