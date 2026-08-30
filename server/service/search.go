package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// SearchResult represents a search result item.
type SearchResult struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description,omitempty"`
	Icon        string  `json:"icon,omitempty"`
	GroupID     string  `json:"group_id"`
	GroupName   string  `json:"group_name,omitempty"`
	Score       float64 `json:"score,omitempty"` // FTS5 rank score
}

// SearchService provides search functionality.
type SearchService struct {
	db *sql.DB
}

// NewSearchService creates a new SearchService.
func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{db: db}
}

// SearchCards searches for cards matching the query.
// Returns results sorted by relevance.
func (s *SearchService) SearchCards(ctx context.Context, userID string, query string, limit int) ([]SearchResult, error) {
	// Clean up query - remove extra spaces
	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchResult{}, nil
	}

	// Prepare FTS5 query
	// Split by spaces and combine with AND for better matching
	terms := strings.Fields(query)
	ftsQuery := strings.Join(terms, " AND ")

	// Search using FTS5
	querySQL := `
		SELECT 
			c.id,
			c.title,
			c.url,
			c.description,
			c.icon,
			c.group_id,
			g.name as group_name,
			bm25(cards_fts) as score
		FROM cards_fts
		JOIN cards c ON c.id = cards_fts.card_id
		JOIN panel_groups g ON g.id = c.group_id
		WHERE cards_fts MATCH ? 
			AND c.user_id = ?
			AND g.user_id = ?
		ORDER BY score
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, querySQL, ftsQuery, userID, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("search query: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		var desc sql.NullString
		var icon sql.NullString
		var groupName sql.NullString
		if err := rows.Scan(
			&r.ID,
			&r.Title,
			&r.URL,
			&desc,
			&icon,
			&r.GroupID,
			&groupName,
			&r.Score,
		); err != nil {
			return nil, fmt.Errorf("scan result: %w", err)
		}
		if desc.Valid {
			r.Description = desc.String
		}
		if icon.Valid {
			r.Icon = icon.String
		}
		if groupName.Valid {
			r.GroupName = groupName.String
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

// GetSearchSuggestions returns search suggestions based on partial input.
func (s *SearchService) GetSearchSuggestions(ctx context.Context, userID string, prefix string, limit int) ([]string, error) {
	if prefix == "" {
		return []string{}, nil
	}

	// Simple prefix search on titles
	querySQL := `
		SELECT DISTINCT c.title
		FROM cards c
		JOIN panel_groups g ON g.id = c.group_id
		WHERE c.title LIKE ? 
			AND c.user_id = ?
			AND g.user_id = ?
		ORDER BY c.title
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, querySQL, prefix+"%", userID, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("suggestions query: %w", err)
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("scan suggestion: %w", err)
		}
		suggestions = append(suggestions, title)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return suggestions, nil
}
