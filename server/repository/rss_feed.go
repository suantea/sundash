package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sundash/models"
)

// RSSFeedRepo handles database operations for RSS feeds.
type RSSFeedRepo struct {
	db *sql.DB
}

// NewRSSFeedRepo creates a new RSSFeedRepo.
func NewRSSFeedRepo(db *sql.DB) *RSSFeedRepo {
	return &RSSFeedRepo{db: db}
}

// Create inserts a new RSS feed into the database.
func (r *RSSFeedRepo) Create(ctx context.Context, feed *models.RSSFeed) error {
	query := `
		INSERT INTO rss_feeds (id, user_id, title, url, description, image_url, last_fetched, update_interval, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		feed.ID,
		feed.UserID,
		feed.Title,
		feed.URL,
		feed.Description,
		feed.ImageURL,
		feed.LastFetched,
		feed.UpdateInterval,
		feed.CreatedAt,
		feed.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert rss feed: %w", err)
	}
	return nil
}

// GetByIDAndUser returns an RSS feed by ID and user ID.
func (r *RSSFeedRepo) GetByIDAndUser(ctx context.Context, feedID, userID string) (*models.RSSFeed, error) {
	query := `
		SELECT id, user_id, title, url, description, image_url, last_fetched, update_interval, created_at, updated_at
		FROM rss_feeds
		WHERE id = ? AND user_id = ?
	`
	row := r.db.QueryRowContext(ctx, query, feedID, userID)
	var feed models.RSSFeed
	var lastFetched, createdAt, updatedAt []byte // store as []byte to handle NULL
	if err := row.Scan(
		&feed.ID,
		&feed.UserID,
		&feed.Title,
		&feed.URL,
		&feed.Description,
		&feed.ImageURL,
		&lastFetched,
		&feed.UpdateInterval,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, fmt.Errorf("get rss feed: %w", err)
	}
	if len(lastFetched) > 0 {
		feed.LastFetched, _ = time.Parse(time.RFC3339, string(lastFetched))
	}
	if len(createdAt) > 0 {
		feed.CreatedAt, _ = time.Parse(time.RFC3339, string(createdAt))
	}
	if len(updatedAt) > 0 {
		feed.UpdatedAt, _ = time.Parse(time.RFC3339, string(updatedAt))
	}
	return &feed, nil
}

// ListByUser returns all RSS feeds for a user.
func (r *RSSFeedRepo) ListByUser(ctx context.Context, userID string) ([]*models.RSSFeed, error) {
	query := `
		SELECT id, user_id, title, url, description, image_url, last_fetched, update_interval, created_at, updated_at
		FROM rss_feeds
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list rss feeds: %w", err)
	}
	defer rows.Close()

	var feeds []*models.RSSFeed
	for rows.Next() {
		var f models.RSSFeed
		var lastFetched, createdAt, updatedAt []byte
		if err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.Title,
			&f.URL,
			&f.Description,
			&f.ImageURL,
			&lastFetched,
			&f.UpdateInterval,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rss feed: %w", err)
		}
		if len(lastFetched) > 0 {
			f.LastFetched, _ = time.Parse(time.RFC3339, string(lastFetched))
		}
		if len(createdAt) > 0 {
			f.CreatedAt, _ = time.Parse(time.RFC3339, string(createdAt))
		}
		if len(updatedAt) > 0 {
			f.UpdatedAt, _ = time.Parse(time.RFC3339, string(updatedAt))
		}
		feeds = append(feeds, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return feeds, nil
}

// ListAll returns all RSS feeds (for background fetcher).
func (r *RSSFeedRepo) ListAll(ctx context.Context) ([]*models.RSSFeed, error) {
	query := `
		SELECT id, user_id, title, url, description, image_url, last_fetched, update_interval, created_at, updated_at
		FROM rss_feeds
		ORDER BY updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list all rss feeds: %w", err)
	}
	defer rows.Close()

	var feeds []*models.RSSFeed
	for rows.Next() {
		var f models.RSSFeed
		var lastFetched, createdAt, updatedAt []byte
		if err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.Title,
			&f.URL,
			&f.Description,
			&f.ImageURL,
			&lastFetched,
			&f.UpdateInterval,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rss feed: %w", err)
		}
		if len(lastFetched) > 0 {
			f.LastFetched, _ = time.Parse(time.RFC3339, string(lastFetched))
		}
		if len(createdAt) > 0 {
			f.CreatedAt, _ = time.Parse(time.RFC3339, string(createdAt))
		}
		if len(updatedAt) > 0 {
			f.UpdatedAt, _ = time.Parse(time.RFC3339, string(updatedAt))
		}
		feeds = append(feeds, &f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return feeds, nil
}

// Update updates an existing RSS feed.
func (r *RSSFeedRepo) Update(ctx context.Context, feed *models.RSSFeed) error {
	query := `
		UPDATE rss_feeds
		SET title = ?, url = ?, description = ?, image_url = ?, last_fetched = ?, update_interval = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		feed.Title,
		feed.URL,
		feed.Description,
		feed.ImageURL,
		feed.LastFetched,
		feed.UpdateInterval,
		feed.UpdatedAt,
		feed.ID,
		feed.UserID,
	)
	if err != nil {
		return fmt.Errorf("update rss feed: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("rss feed not found or unauthorized")
	}
	return nil
}

// Delete removes an RSS feed by ID and user ID.
func (r *RSSFeedRepo) Delete(ctx context.Context, userID, feedID string) error {
	query := `
		DELETE FROM rss_feeds
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.ExecContext(ctx, query, feedID, userID)
	if err != nil {
		return fmt.Errorf("delete rss feed: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("rss feed not found or unauthorized")
	}
	return nil
}