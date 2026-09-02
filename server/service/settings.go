package service

import (
	"sync"
	"time"

	"sundash/repository"
)

// Site setting keys shared between handlers and the HTML injector.
var GlobalSettingKeys = []string{
	"allow_registration", "require_approval",
	"default_wallpaper_type", "default_wallpaper_url", "default_wallpaper_blur",
	"default_wallpaper_opacity", "default_wallpaper_copyright",
	"site_title", "site_icon_url", "site_cdn_url", "site_analytics_code",
	"site_custom_head", "site_custom_footer",
	// Bookmark-sync connection (server URL + bearer token). Device id is
	// exposed via /api/bmsync/status only and not editable from the UI.
	// bmsync_enabled 控制书签同步功能总开关（默认开启），经 site-config 下发，
	// 关闭后前端隐藏书签同步入口。
	"bmsync_server_url", "bmsync_token", "bmsync_enabled",
}

var SiteConfigKeys = []string{
	"site_title", "site_icon_url", "site_cdn_url", "site_analytics_code",
	"site_custom_head", "site_custom_footer",
	"bmsync_enabled",
}

type SettingsService struct {
	settings *repository.SettingsRepo

	// siteConfigCache caches the global site config (read on every SPA
	// page load) for a short TTL to avoid a DB query per request.
	mu             sync.Mutex
	siteConfig     map[string]string
	siteConfigExpr time.Time
}

const siteConfigTTL = 30 * time.Second

func NewSettingsService(settings *repository.SettingsRepo) *SettingsService {
	return &SettingsService{settings: settings}
}

// GetForUser returns the merged (global + user) setting map.
func (s *SettingsService) GetForUser(userID string) (map[string]string, error) {
	return s.settings.AllForUser(userID)
}

func (s *SettingsService) UpdateForUser(userID, key, value string) error {
	return s.settings.Upsert(key, value, &userID)
}

func (s *SettingsService) BatchUpdateForUser(userID string, settings map[string]string) error {
	return s.settings.UpsertBatch(settings, &userID)
}

func (s *SettingsService) GetGlobal() (map[string]string, error) {
	return s.settings.GlobalByKeys(GlobalSettingKeys)
}

// UpdateGlobal 更新全局设置并使 site-config 缓存失效。
func (s *SettingsService) UpdateGlobal(settings map[string]string) error {
	if err := s.settings.UpsertBatch(settings, nil); err != nil {
		return err
	}
	s.InvalidateSiteConfig()
	return nil
}

// InvalidateSiteConfig 清除 site-config 缓存。
// 供绕过 UpdateGlobal 直接写库（如首次 Setup 写入 site_title）后调用，
// 避免 30s TTL 内新标题不生效。
func (s *SettingsService) InvalidateSiteConfig() {
	s.mu.Lock()
	s.siteConfig = nil
	s.mu.Unlock()
}

// GetSiteConfig returns the global site config with a short TTL cache.
func (s *SettingsService) GetSiteConfig() (map[string]string, error) {
	s.mu.Lock()
	if s.siteConfig != nil && time.Now().Before(s.siteConfigExpr) {
		cached := s.siteConfig
		s.mu.Unlock()
		return cached, nil
	}
	s.mu.Unlock()

	cfg, err := s.settings.GlobalByKeys(SiteConfigKeys)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.siteConfig = cfg
	s.siteConfigExpr = time.Now().Add(siteConfigTTL)
	s.mu.Unlock()
	return cfg, nil
}
