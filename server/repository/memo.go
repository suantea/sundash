package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sundash/models"
)

// MemoRepo handles database operations for memos.
type MemoRepo struct {
	db *sql.DB
}

// NewMemoRepo creates a new MemoRepo.
func NewMemoRepo(db *sql.DB) *MemoRepo {
	return &MemoRepo{db: db}
}

// Create inserts a new memo into the database.
func (r *MemoRepo) Create(ctx context.Context, memo *models.Memo) error {
	query := `
		INSERT INTO memos (id, user_id, content, is_archived, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		memo.ID,
		memo.UserID,
		memo.Content,
		boolToInt(memo.IsArchived),
		memo.CreatedAt,
		memo.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert memo: %w", err)
	}
	return nil
}

// ListByUser returns memos for a user, optionally filtering by archived status.
func (r *MemoRepo) ListByUser(ctx context.Context, userID string, archived bool) ([]*models.Memo, error) {
	query := `
		SELECT id, user_id, content, is_archived, created_at, updated_at
		FROM memos
		WHERE user_id = ? AND is_archived = ?
		ORDER BY updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, boolToInt(archived))
	if err != nil {
		return nil, fmt.Errorf("query memos: %w", err)
	}
	defer rows.Close()

	var memos []*models.Memo
	for rows.Next() {
		var m models.Memo
		var isArchivedInt int
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Content,
			&isArchivedInt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan memo: %w", err)
		}
		m.IsArchived = isArchivedInt == 1
		memos = append(memos, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return memos, nil
}

// UpdateArchive updates the archived status of a memo.
func (r *MemoRepo) UpdateArchive(ctx context.Context, userID, id string, archived bool) error {
	query := `
		UPDATE memos
		SET is_archived = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		boolToInt(archived),
		time.Now(),
		id,
		userID,
	)
	if err != nil {
		return fmt.Errorf("update memo archive: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("memo not found or unauthorized")
	}
	return nil
}

// Delete removes a memo from the database.
func (r *MemoRepo) Delete(ctx context.Context, userID, id string) error {
	query := `
		DELETE FROM memos
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete memo: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("memo not found or unauthorized")
	}
	return nil
}

// UpdateContent updates the content of a memo.
func (r *MemoRepo) UpdateContent(ctx context.Context, userID, id, content string) error {
	query := `
		UPDATE memos
		SET content = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.ExecContext(ctx, query, content, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("update memo content: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("memo not found or unauthorized")
	}
	return nil
}

// Helper to convert bool to int for SQLite
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}