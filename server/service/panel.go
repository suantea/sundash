package service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"sundash/models"
	"sundash/repository"

	"github.com/google/uuid"
)

type PanelService struct {
	panels *repository.PanelRepo
	users  *repository.UserRepo
	mu     sync.RWMutex
	cache  map[string]*cachedPanel
}

type cachedPanel struct {
	data    *models.PanelData
	expires time.Time
}

func NewPanelService(panels *repository.PanelRepo, users *repository.UserRepo) *PanelService {
	return &PanelService{
		panels: panels,
		users:  users,
		cache:  make(map[string]*cachedPanel),
	}
}

// GetPanel loads groups and their cards. Uses a simple cache to avoid repeated DB queries.
// Cache key is userID, cache expires after 5 minutes.
func (s *PanelService) GetPanel(userID string) (*models.PanelData, error) {
	s.mu.RLock()
	if cached, found := s.cache[userID]; found && time.Now().Before(cached.expires) {
		s.mu.RUnlock()
		return cached.data, nil
	}
	s.mu.RUnlock()

	groups, err := s.panels.GroupsByUser(userID)
	if err != nil {
		return nil, err
	}
	cards, err := s.panels.CardsByUser(userID)
	if err != nil {
		return nil, err
	}

	byGroup := make(map[string][]models.Card, len(groups))
	for _, c := range cards {
		byGroup[c.GroupID] = append(byGroup[c.GroupID], c)
	}
	for i := range groups {
		groups[i].Cards = byGroup[groups[i].ID]
	}
	panelData := &models.PanelData{Groups: groups}

	s.mu.Lock()
	s.cache[userID] = &cachedPanel{
		data:    panelData,
		expires: time.Now().Add(5 * time.Minute),
	}
	s.mu.Unlock()
	return panelData, nil
}

// invalidate drops the cached panel for a user. Every mutation calls it —
// otherwise a change stays invisible to the home page (and MCP list tools)
// until the 5-minute read cache expires.
func (s *PanelService) invalidate(userID string) {
	s.mu.Lock()
	delete(s.cache, userID)
	s.mu.Unlock()
}

func (s *PanelService) CreateGroup(userID, name string) (*models.PanelGroup, error) {
	maxOrder, err := s.panels.MaxGroupOrder(userID)
	if err != nil {
		return nil, err
	}
	group := models.PanelGroup{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		SortOrder: maxOrder + 1,
	}
	if err := s.panels.CreateGroup(group.ID, userID, name, group.SortOrder); err != nil {
		return nil, err
	}
	s.invalidate(userID)
	return &group, nil
}

func (s *PanelService) UpdateGroup(userID, groupID string, name *string, sortOrder *int) error {
	owner, err := s.panels.GroupOwner(groupID)
	if err != nil {
		return err
	}
	if owner == "" || owner != userID {
		return NewError(ErrNotFound, "Group not found")
	}
	if err := s.panels.UpdateGroup(groupID, name, sortOrder); err != nil {
		return err
	}
	s.invalidate(userID)
	return nil
}

func (s *PanelService) DeleteGroup(userID, groupID string) error {
	owner, err := s.panels.GroupOwner(groupID)
	if err != nil {
		return err
	}
	if owner == "" || owner != userID {
		return NewError(ErrNotFound, "Group not found")
	}
	if err := s.panels.DeleteGroup(groupID); err != nil {
		return err
	}
	s.invalidate(userID)
	return nil
}

func (s *PanelService) CreateCard(userID string, req models.CreateCardRequest) (*models.Card, error) {
	owner, err := s.panels.GroupOwner(req.GroupID)
	if err != nil {
		return nil, err
	}
	if owner == "" || owner != userID {
		return nil, NewError(ErrBadRequest, "Invalid group")
	}

	openType := req.OpenType
	if openType == "" {
		openType = "new_tab"
	}
	maxOrder, err := s.panels.MaxCardOrder(req.GroupID)
	if err != nil {
		return nil, err
	}

	card := models.Card{
		ID:          uuid.New().String(),
		GroupID:     req.GroupID,
		UserID:      userID,
		Title:       req.Title,
		URL:         req.URL,
		URLInternal: req.URLInternal,
		Icon:        req.Icon,
		IconColor:   req.IconColor,
		BgColor:     req.BgColor,
		Description: req.Description,
		OpenType:    openType,
		SortOrder:   maxOrder + 1,
	}
	if err := s.panels.CreateCard(card); err != nil {
		return nil, err
	}
	s.invalidate(userID)
	return &card, nil
}

func (s *PanelService) UpdateCard(userID, cardID string, req models.UpdateCardRequest) error {
	owner, err := s.panels.CardOwner(cardID)
	if err != nil {
		return err
	}
	if owner == "" || owner != userID {
		return NewError(ErrNotFound, "Card not found")
	}
	if err := s.panels.UpdateCard(cardID, req); err != nil {
		return err
	}
	s.invalidate(userID)
	return nil
}

func (s *PanelService) DeleteCard(userID, cardID string) error {
	owner, err := s.panels.CardOwner(cardID)
	if err != nil {
		return err
	}
	if owner == "" || owner != userID {
		return NewError(ErrNotFound, "Card not found")
	}
	if err := s.panels.DeleteCard(cardID); err != nil {
		return err
	}
	s.invalidate(userID)
	return nil
}

// Reorder applies ordering; entries not owned by the user are silently ignored
// so a caller can never reorder another user's data.
func (s *PanelService) Reorder(userID string, req models.ReorderRequest) error {
	groups, err := s.panels.GroupsByUser(userID)
	if err != nil {
		return err
	}
	ownedGroups := make(map[string]bool, len(groups))
	for _, g := range groups {
		ownedGroups[g.ID] = true
	}
	ownedCards := make(map[string]bool)
	cards, err := s.panels.CardsByUser(userID)
	if err != nil {
		return err
	}
	for _, c := range cards {
		ownedCards[c.ID] = true
	}

	groupOrders := make([]models.GroupOrder, 0, len(req.GroupOrders))
	for _, g := range req.GroupOrders {
		if ownedGroups[g.ID] {
			groupOrders = append(groupOrders, g)
		}
	}
	cardOrders := make([]models.CardOrder, 0, len(req.CardOrders))
	for _, c := range req.CardOrders {
		if ownedCards[c.ID] && ownedGroups[c.GroupID] {
			cardOrders = append(cardOrders, c)
		}
	}
	if err := s.panels.Reorder(groupOrders, cardOrders); err != nil {
		return err
	}
	s.invalidate(userID)
	return nil
}

// BatchReorganize applies multiple create/rename/delete group and move/create card
// operations atomically. Used by MCP AI tools for natural-language reorganizing.
type BatchOp struct {
	// Action: "create_group" | "rename_group" | "delete_group" | "move_card" | "create_card"
	Action string `json:"action"`
	// For create_group/rename_group/delete_group
	GroupName string `json:"group_name,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
	// For move_card
	CardID        string `json:"card_id,omitempty"`
	TargetGroupID string `json:"target_group_id,omitempty"`
	// For create_card
	CardTitle string `json:"card_title,omitempty"`
	CardURL   string `json:"card_url,omitempty"`
}

type BatchReorganizeRequest struct {
	Operations []BatchOp `json:"operations"`
}

func (s *PanelService) BatchReorganize(userID string, req BatchReorganizeRequest) error {
	for _, op := range req.Operations {
		switch op.Action {
		case "create_group":
			if _, err := s.CreateGroup(userID, op.GroupName); err != nil {
				return fmt.Errorf("create_group %q: %w", op.GroupName, err)
			}
		case "rename_group":
			if err := s.UpdateGroup(userID, op.GroupID, &op.GroupName, nil); err != nil {
				return fmt.Errorf("rename_group %s: %w", op.GroupID, err)
			}
		case "delete_group":
			if err := s.DeleteGroup(userID, op.GroupID); err != nil {
				return fmt.Errorf("delete_group %s: %w", op.GroupID, err)
			}
		case "move_card":
			if err := s.Reorder(userID, models.ReorderRequest{
				CardOrders: []models.CardOrder{{ID: op.CardID, GroupID: op.TargetGroupID, SortOrder: 0}},
			}); err != nil {
				return fmt.Errorf("move_card %s: %w", op.CardID, err)
			}
		case "create_card":
			if _, err := s.CreateCard(userID, models.CreateCardRequest{
				GroupID: op.TargetGroupID,
				Title:   op.CardTitle,
				URL:     op.CardURL,
			}); err != nil {
				return fmt.Errorf("create_card %q: %w", op.CardTitle, err)
			}
		default:
			return fmt.Errorf("unknown batch operation: %s", op.Action)
		}
	}
	return nil
}

// BatchSetIcons applies icon updates to multiple cards in one call.
// Used by MCP AI tools for batch icon management.
type BatchSetIconsRequest struct {
	// CardID -> Icon mapping
	Icons map[string]string `json:"icons"`
}

func (s *PanelService) BatchSetIcons(userID string, req BatchSetIconsRequest) error {
	for cardID, icon := range req.Icons {
		if err := s.UpdateCard(userID, cardID, models.UpdateCardRequest{Icon: &icon}); err != nil {
			return fmt.Errorf("set icon on card %s: %w", cardID, err)
		}
	}
	return nil
}

// GetUserPanel returns another user's panel data (admin only).
func (s *PanelService) GetUserPanel(targetID string) (map[string]any, error) {
	user, err := s.users.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewError(ErrNotFound, "User not found")
	}

	data, err := s.GetPanel(targetID)
	if err != nil {
		return nil, fmt.Errorf("load panel data: %w", err)
	}
	return map[string]any{
		"user_id":  targetID,
		"username": user.Username,
		"groups":   data.Groups,
	}, nil
}

// SuggestCategorize analyzes existing group patterns and suggests group assignments
// for cards in low-value groups ("暂停使用", "其他"). Returns suggested operations
// that can be applied via BatchReorganize.
type CategorizeSuggestion struct {
	CardID      string `json:"card_id"`
	CardTitle   string `json:"card_title"`
	SourceGroup string `json:"source_group"`
	TargetGroup string `json:"target_group"`
	TargetID    string `json:"target_id"`
	Confidence  string `json:"confidence"` // "high" | "medium"
}

type SuggestCategorizeResult struct {
	Suggestions []CategorizeSuggestion `json:"suggestions"`
	Summary     map[string]int         `json:"summary"` // groupName -> count of suggested moves
}

// keywordRules defines category classification rules based on domain/title keywords.
// Each rule maps keywords to a target group name.
var keywordRules = map[string][]string{
	"APP":             {"nas", "qnap", "fnos", "transmission", "emby", "jellyfin", "qbittorrent", "docker", "alist", "home assistant", "ha.", "pve", "proxmox", "openwrt", "virtualization", "1panel"},
	"商标检索":         {"商标", "标局", "专利", "版权", "检索", "查询", "企查", "权大师", "马德里", "分类", "知产", "ip.", "lanternfish", "标库"},
	"公司门户":         {"chery", "奇瑞", "合同", "采购", "门户", "办公", "oa", "crm", "hr.", "erp"},
	"工具":            {"翻译", "转换", "pdf", "图片", "压缩", "字体", "icon", "logo", "图表", "chart", "放大", "remove", "tinify", "photopea", "vectormagic", "工具箱"},
	"新闻咨询类":       {"新闻", "资讯", "快讯", "36氪", "汽车", "媒体", "数据", "新闻", "新浪", "so.car", "艾媒", "autothinker"},
	"开发":            {"git", "github", "gitee", "code", "api", "docs", "nocobase", "wordpress", "bi ", "dashboard", "rsshub", "freshrss", "searxng", "kodbox", "mrdoc", "签到"},
	"影音":            {"music", "音乐", "navidrome", "video", "视频", "tv", "电影", "刮削", "audiobook", "播客", "podcast"},
	"网络":            {"vpn", "加速", "proxy", "tunnel", "ssh", "网盘", "同步", "sync", "百度网盘", "github加速", "docker加速"},
}

func (s *PanelService) SuggestCategorize(userID string) (*SuggestCategorizeResult, error) {
	data, err := s.GetPanel(userID)
	if err != nil {
		return nil, err
	}

	// Build group name->id map
	groupMap := make(map[string]string)
	for _, g := range data.Groups {
		groupMap[g.Name] = g.ID
	}

	// Find low-value groups (source candidates)
	lowValueGroups := map[string]bool{"暂停使用": true, "其他": true}

	var suggestions []CategorizeSuggestion

	for _, group := range data.Groups {
		if !lowValueGroups[group.Name] {
			continue
		}
		for _, card := range group.Cards {
			bestGroup := ""
			bestConfidence := ""
			titleLower := strings.ToLower(card.Title)
			urlLower := strings.ToLower(card.URL)

			for targetGroup, keywords := range keywordRules {
				for _, kw := range keywords {
					if strings.Contains(titleLower, kw) || strings.Contains(urlLower, kw) {
						if bestGroup == "" {
							bestGroup = targetGroup
							bestConfidence = "medium"
						}
						// If matched by title (not just URL), higher confidence
						if strings.Contains(titleLower, kw) {
							bestConfidence = "high"
							bestGroup = targetGroup
							break
						}
					}
				}
				if bestConfidence == "high" {
					break
				}
			}

			if bestGroup != "" {
				if targetID, ok := groupMap[bestGroup]; ok {
					suggestions = append(suggestions, CategorizeSuggestion{
						CardID:      card.ID,
						CardTitle:   card.Title,
						SourceGroup: group.Name,
						TargetGroup: bestGroup,
						TargetID:    targetID,
						Confidence:  bestConfidence,
					})
				}
			}
		}
	}

	// Build summary
	summary := make(map[string]int)
	for _, s := range suggestions {
		summary[s.TargetGroup]++
	}

	return &SuggestCategorizeResult{
		Suggestions: suggestions,
		Summary:     summary,
	}, nil
}
