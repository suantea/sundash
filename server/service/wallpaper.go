package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// BingWallpaper mirrors the Bing wallpaper API response.
type BingWallpaper struct {
	Images []struct {
		URL         string `json:"url"`
		Copyright   string `json:"copyright"`
		Date        string `json:"date"`
		Title       string `json:"title"`
		Hdpicture   string `json:"hdpicture"`
		Fullpicture string `json:"fullpicture"`
	} `json:"images"`
}

type WallpaperService struct {
	client *http.Client

	mu      sync.Mutex
	cached  *BingWallpaper
	fetched time.Time
}

func NewWallpaperService() *WallpaperService {
	return &WallpaperService{client: &http.Client{Timeout: 10 * time.Second}}
}

// GetToday returns the daily Bing wallpaper, cached for one hour (mutex-guarded).
func (s *WallpaperService) GetToday() (*BingWallpaper, error) {
	s.mu.Lock()
	if s.cached != nil && time.Since(s.fetched) < time.Hour {
		s.mu.Unlock()
		return s.cached, nil
	}
	s.mu.Unlock()

	wallpaper, err := s.fetch("")
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cached = wallpaper
	s.fetched = time.Now()
	s.mu.Unlock()
	return wallpaper, nil
}

func (s *WallpaperService) GetByDate(date string) (*BingWallpaper, error) {
	return s.fetch(date)
}

func (s *WallpaperService) fetch(date string) (*BingWallpaper, error) {
	url := "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"
	if date != "" {
		url = fmt.Sprintf("%s&date=%s", url, date)
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch Bing wallpaper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bing api returned status %d", resp.StatusCode)
	}

	var result struct {
		Images []struct {
			URL       string `json:"url"`
			Copyright string `json:"copyright"`
			Enddate   string `json:"enddate"`
			Title     string `json:"title"`
		} `json:"images"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return nil, fmt.Errorf("parse Bing response: %w", err)
	}
	if len(result.Images) == 0 {
		return nil, fmt.Errorf("no wallpaper data found")
	}

	img := result.Images[0]
	base := fmt.Sprintf("https://cn.bing.com%s", img.URL)
	return &BingWallpaper{Images: []struct {
		URL         string `json:"url"`
		Copyright   string `json:"copyright"`
		Date        string `json:"date"`
		Title       string `json:"title"`
		Hdpicture   string `json:"hdpicture"`
		Fullpicture string `json:"fullpicture"`
	}{{
		URL:         base,
		Copyright:   img.Copyright,
		Date:        img.Enddate,
		Title:       img.Title,
		Hdpicture:   base,
		Fullpicture: base,
	}}}, nil
}
