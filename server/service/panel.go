package service

import (
	"fmt"

	"sundash/models"
	"sundash/repository"

	"github.com/google/uuid"
)

type PanelService struct {
	panels *repository.PanelRepo
	users  *repository.UserRepo
}

func NewPanelService(panels *repository.PanelRepo, users *repository.UserRepo) *PanelService {
	return &PanelService{panels: panels, users: users}
}

// GetPanel loads groups and their cards. Cards are fetched in a single query
// and grouped in memory, avoiding the previous N+1 query pattern.
func (s *PanelService) GetPanel(userID string) (*models.PanelData, error) {
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
	return &models.PanelData{Groups: groups}, nil
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
	return s.panels.UpdateGroup(groupID, name, sortOrder)
}

func (s *PanelService) DeleteGroup(userID, groupID string) error {
	owner, err := s.panels.GroupOwner(groupID)
	if err != nil {
		return err
	}
	if owner == "" || owner != userID {
		return NewError(ErrNotFound, "Group not found")
	}
	return s.panels.DeleteGroup(groupID)
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
	return s.panels.UpdateCard(cardID, req)
}

func (s *PanelService) DeleteCard(userID, cardID string) error {
	owner, err := s.panels.CardOwner(cardID)
	if err != nil {
		return err
	}
	if owner == "" || owner != userID {
		return NewError(ErrNotFound, "Card not found")
	}
	return s.panels.DeleteCard(cardID)
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
	return s.panels.Reorder(groupOrders, cardOrders)
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
