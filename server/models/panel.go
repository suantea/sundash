package models

import (
	"time"
)

type PanelGroup struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	Cards     []Card    `json:"cards,omitempty"`
}

type Card struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	URLInternal string    `json:"url_internal,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	IconColor   string    `json:"icon_color,omitempty"`
	BgColor     string    `json:"bg_color,omitempty"`
	Description string    `json:"description,omitempty"`
	OpenType    string    `json:"open_type"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

type PanelData struct {
	Groups []PanelGroup `json:"groups"`
}

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGroupRequest struct {
	Name      *string `json:"name"`
	SortOrder *int    `json:"sort_order"`
}

type CreateCardRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url" binding:"required"`
	URLInternal string `json:"url_internal"`
	Icon        string `json:"icon"`
	IconColor   string `json:"icon_color"`
	BgColor     string `json:"bg_color"`
	Description string `json:"description"`
	OpenType    string `json:"open_type"`
}

type UpdateCardRequest struct {
	Title       *string `json:"title"`
	URL         *string `json:"url"`
	URLInternal *string `json:"url_internal"`
	Icon        *string `json:"icon"`
	IconColor   *string `json:"icon_color"`
	BgColor     *string `json:"bg_color"`
	Description *string `json:"description"`
	OpenType    *string `json:"open_type"`
	SortOrder   *int    `json:"sort_order"`
}

type ReorderRequest struct {
	GroupOrders []GroupOrder `json:"group_orders"`
	CardOrders  []CardOrder  `json:"card_orders"`
}

type GroupOrder struct {
	ID        string `json:"id"`
	SortOrder int    `json:"sort_order"`
}

type CardOrder struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	SortOrder int    `json:"sort_order"`
}

type Setting struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	UserID    string `json:"user_id,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

// Setting Keys - 个性化设置
const (
	// 时钟设置
	SettingClockShow         = "clock_show"
	SettingClockShowSeconds  = "clock_show_seconds"
	SettingClockFormat       = "clock_format" // 12/24

	// Logo 设置
	SettingLogoType          = "logo_type"      // image/text
	SettingLogoImageUrl      = "logo_image_url"
	SettingLogoText          = "logo_text"

	// 卡片风格
	SettingCardStyle         = "card_style" // default/round/square

	// 布局设置
	SettingContentMaxWidth   = "content_max_width"
	SettingContentPaddingX   = "content_padding_x"
	SettingContentPaddingTop = "content_padding_top"
	SettingContentPaddingBottom = "content_padding_bottom"

	// 颜色自定义
	SettingBgColor           = "bg_color"
	SettingTextColor         = "text_color"
	SettingBorderColor       = "border_color"
	SettingPrimaryColor      = "primary_color"

	// 页脚
	SettingFooterHtml        = "footer_html"

	// 站点设置
	SettingSiteTitle         = "site_title"
	SettingSiteIconUrl       = "site_icon_url"
	SettingLoginBgUrl        = "login_bg_url"
	SettingEnableCaptcha     = "enable_captcha"
	SettingEnableCache       = "enable_cache"

	// 组件设置
	SettingShowSystemStatus  = "show_system_status"
	SettingNetworkModeAuto   = "network_mode_auto"
)

// Weather represents simplified weather data for frontend display.
type Weather struct {
	Temperature     float64 `json:"temperature"` // Celsius
	Windspeed       float64 `json:"windspeed"`   // km/h
	Winddirection   float64 `json:"winddirection"` // degrees
	Weathercode     int     `json:"weathercode"` // WMO weather code
	Time            string  `json:"time"`        // ISO time string
	Timezone        string  `json:"timezone"`    // IANA timezone
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

// Memo represents a quick note/memo.
type Memo struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	IsArchived bool     `json:"is_archived"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RSSFeed represents an RSS feed subscription.
type RSSFeed struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	LastFetched time.Time `json:"last_fetched,omitempty"`
	UpdateInterval int    `json:"update_interval_minutes"` // How often to fetch (in minutes)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RSSItem represents a single item in an RSS feed.
type RSSItem struct {
	ID        string    `json:"id"`
	FeedID    string    `json:"feed_id"`
	Title     string    `json:"title"`
	Link      string    `json:"link"`
	Description string    `json:"description,omitempty"`
	PubDate   time.Time `json:"pub_date,omitempty"`
	Author    string    `json:"author,omitempty"`
	Guid      string    `json:"guid,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
