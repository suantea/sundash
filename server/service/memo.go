package service

import (
	"context"
	"fmt"
	"time"

	"sundash/models"
	"sundash/repository"
)

// MemoService provides memo operations.
type MemoService struct {
	memoRepo *repository.MemoRepo
}

// NewMemoService creates a new MemoService.
func NewMemoService(memoRepo *repository.MemoRepo) *MemoService {
	return &MemoService{memoRepo: memoRepo}
}

// CreateMemo creates a new memo.
func (s *MemoService) CreateMemo(ctx context.Context, userID, content string) (*models.Memo, error) {
	memo := &models.Memo{
		ID:        generateID(),
		UserID:    userID,
		Content:   content,
		IsArchived: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.memoRepo.Create(ctx, memo)
	if err != nil {
		return nil, fmt.Errorf("create memo: %w", err)
	}
	return memo, nil
}

// GetMemos returns all non-archived memos for a user, ordered by updated_at desc.
func (s *MemoService) GetMemos(ctx context.Context, userID string) ([]*models.Memo, error) {
	memos, err := s.memoRepo.ListByUser(ctx, userID, false)
	if err != nil {
		return nil, fmt.Errorf("list memos: %w", err)
	}
	return memos, nil
}

// ArchiveMemo archives or unarchives a memo.
func (s *MemoService) ArchiveMemo(ctx context.Context, userID, id string, archived bool) error {
	err := s.memoRepo.UpdateArchive(ctx, userID, id, archived)
	if err != nil {
		return fmt.Errorf("update memo archive: %w", err)
	}
	return nil
}

// DeleteMemo deletes a memo.
func (s *MemoService) DeleteMemo(ctx context.Context, userID, id string) error {
	err := s.memoRepo.Delete(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("delete memo: %w", err)
	}
	return nil
}

// UpdateMemo updates memo content.
func (s *MemoService) UpdateMemo(ctx context.Context, userID, id, content string) error {
	err := s.memoRepo.UpdateContent(ctx, userID, id, content)
	if err != nil {
		return fmt.Errorf("update memo content: %w", err)
	}
	return nil
}

func generateID() string {
	// Simple ID generation; in production use UUID or similar
	return fmt.Sprintf("%d", time.Now().UnixNano())
}