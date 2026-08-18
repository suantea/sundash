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
