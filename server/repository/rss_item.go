package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sundash/models"
)

// RSSItemRepo handles database operations for RSS items.
type RSSItemRepo struct {
	db *sql.DB
}

// NewRSSItemRepo creates a new RSSItemRepo.
func NewRSSItemRepo(db *sql.DB) *RSSItemRepo {
	return &RSSItemRepo{db: db}
}

// Create inserts a new RSS item into the database.
func (r *RSSItemRepo) Create(ctx context.Context, item *models.RSSItem) error {
	query := `
		INSERT INTO rss_items (id, feed_id, title, link, description, pub_date, author, guid, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.FeedID,
		item.Title,
		item.Link,
		item.Description,
		item.PubDate,
		item.Author,
		item.Guid,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert rss item: %w", err)
	}
	return nil
}

// ListByFeed returns items for a feed, optionally limited.
func (r *RSSItemRepo) ListByFeed(ctx context.Context, feedID string, limit int) ([]*models.RSSItem, error) {
	query := `
		SELECT id, feed_id, title, link, description, pub_date, author, guid, created_at, updated_at
		FROM rss_items
		WHERE feed_id = ?
		ORDER BY pub_date DESC
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := r.db.QueryContext(ctx, query, feedID)
	if err != nil {
		return nil, fmt.Errorf("list rss items: %w", err)
	}
	defer rows.Close()

	var items []*models.RSSItem
	for rows.Next() {
		var i models.RSSItem
		var pubDate, createdAt, updatedAt []byte
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Link,
			&i.Description,
			&pubDate,
			&i.Author,
			&i.Guid,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rss item: %w", err)
		}
		if len(pubDate) > 0 {
			i.PubDate, _ = time.Parse(time.RFC3339, string(pubDate))
		}
		if len(createdAt) > 0 {
			i.CreatedAt, _ = time.Parse(time.RFC3339, string(createdAt))
		}
		if len(updatedAt) > 0 {
			i.UpdatedAt, _ = time.Parse(time.RFC3339, string(updatedAt))
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return items, nil
}

// DeleteByFeed deletes all items for a feed.
func (r *RSSItemRepo) DeleteByFeed(ctx context.Context, feedID string) error {
	query := `
		DELETE FROM rss_items
		WHERE feed_id = ?
	`
	result, err := r.db.ExecContext(ctx, query, feedID)
	if err != nil {
		return fmt.Errorf("delete rss items: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		// No items to delete is not an error
		return nil
	}
	return nil
}