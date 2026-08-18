package repository

import (
	"database/sql"
	"fmt"

	"sundash/models"
)

type PanelRepo struct {
	db *sql.DB
}

func NewPanelRepo(db *sql.DB) *PanelRepo {
	return &PanelRepo{db: db}
}

func (r *PanelRepo) GroupsByUser(userID string) ([]models.PanelGroup, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, sort_order, created_at FROM panel_groups WHERE user_id = ? ORDER BY sort_order`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query groups: %w", err)
	}
	defer rows.Close()

	groups := make([]models.PanelGroup, 0)
	for rows.Next() {
		var g models.PanelGroup
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.SortOrder, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *PanelRepo) CardsByGroup(groupID string) ([]models.Card, error) {
	rows, err := r.db.Query(
		`SELECT id, group_id, user_id, title, url, COALESCE(url_internal,''), COALESCE(icon,''),
		 COALESCE(icon_color,''), COALESCE(bg_color,''), COALESCE(description,''), open_type, sort_order, created_at
		 FROM cards WHERE group_id = ? ORDER BY sort_order`,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("query cards: %w", err)
	}
	defer rows.Close()
	return scanCards(rows)
}

// CardsByUser returns all cards for a user in one query (avoids N+1 panel loading).
func (r *PanelRepo) CardsByUser(userID string) ([]models.Card, error) {
	rows, err := r.db.Query(
		`SELECT id, group_id, user_id, title, url, COALESCE(url_internal,''), COALESCE(icon,''),
		 COALESCE(icon_color,''), COALESCE(bg_color,''), COALESCE(description,''), open_type, sort_order, created_at
		 FROM cards WHERE user_id = ? ORDER BY sort_order`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query user cards: %w", err)
	}
	defer rows.Close()
	return scanCards(rows)
}

func scanCards(rows *sql.Rows) ([]models.Card, error) {
	cards := make([]models.Card, 0)
	for rows.Next() {
		var c models.Card
		if err := rows.Scan(&c.ID, &c.GroupID, &c.UserID, &c.Title, &c.URL,
			&c.URLInternal, &c.Icon, &c.IconColor, &c.BgColor, &c.Description,
			&c.OpenType, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (r *PanelRepo) MaxGroupOrder(userID string) (int, error) {
	return r.scalarInt("SELECT COALESCE(MAX(sort_order), 0) FROM panel_groups WHERE user_id = ?", userID)
}

func (r *PanelRepo) MaxCardOrder(groupID string) (int, error) {
	return r.scalarInt("SELECT COALESCE(MAX(sort_order), 0) FROM cards WHERE group_id = ?", groupID)
}

func (r *PanelRepo) scalarInt(query string, args ...any) (int, error) {
	var n int
	if err := r.db.QueryRow(query, args...).Scan(&n); err != nil {
		return 0, fmt.Errorf("scalar query: %w", err)
	}
	return n, nil
}

func (r *PanelRepo) CreateGroup(id, userID, name string, sortOrder int) error {
	_, err := r.db.Exec(
		"INSERT INTO panel_groups (id, user_id, name, sort_order) VALUES (?, ?, ?, ?)",
		id, userID, name, sortOrder,
	)
	if err != nil {
		return fmt.Errorf("create group: %w", err)
	}
	return nil
}

func (r *PanelRepo) GroupOwner(groupID string) (string, error) {
	var owner string
	err := r.db.QueryRow("SELECT user_id FROM panel_groups WHERE id = ?", groupID).Scan(&owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("query group owner: %w", err)
	}
	return owner, nil
}

func (r *PanelRepo) UpdateGroup(groupID string, name *string, sortOrder *int) error {
	if name != nil {
		if _, err := r.db.Exec("UPDATE panel_groups SET name = ? WHERE id = ?", *name, groupID); err != nil {
			return fmt.Errorf("update group name: %w", err)
		}
	}
	if sortOrder != nil {
		if _, err := r.db.Exec("UPDATE panel_groups SET sort_order = ? WHERE id = ?", *sortOrder, groupID); err != nil {
			return fmt.Errorf("update group sort_order: %w", err)
		}
	}
	return nil
}

// DeleteGroup removes a group; its cards cascade via FK.
func (r *PanelRepo) DeleteGroup(groupID string) error {
	if _, err := r.db.Exec("DELETE FROM panel_groups WHERE id = ?", groupID); err != nil {
		return fmt.Errorf("delete group: %w", err)
	}
	return nil
}

func (r *PanelRepo) CreateCard(c models.Card) error {
	_, err := r.db.Exec(
		`INSERT INTO cards (id, group_id, user_id, title, url, url_internal, icon, icon_color, bg_color, description, open_type, sort_order)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.GroupID, c.UserID, c.Title, c.URL, c.URLInternal,
		c.Icon, c.IconColor, c.BgColor, c.Description, c.OpenType, c.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("create card: %w", err)
	}
	return nil
}

func (r *PanelRepo) CardOwner(cardID string) (string, error) {
	var owner string
	err := r.db.QueryRow("SELECT user_id FROM cards WHERE id = ?", cardID).Scan(&owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("query card owner: %w", err)
	}
	return owner, nil
}

func (r *PanelRepo) UpdateCard(cardID string, req models.UpdateCardRequest) error {
	if req.Title != nil {
		if _, err := r.db.Exec("UPDATE cards SET title = ? WHERE id = ?", *req.Title, cardID); err != nil {
			return fmt.Errorf("update card title: %w", err)
		}
	}
	if req.URL != nil {
		if _, err := r.db.Exec("UPDATE cards SET url = ? WHERE id = ?", *req.URL, cardID); err != nil {
			return fmt.Errorf("update card url: %w", err)
		}
	}
	if req.URLInternal != nil {
		if _, err := r.db.Exec("UPDATE cards SET url_internal = ? WHERE id = ?", *req.URLInternal, cardID); err != nil {
			return fmt.Errorf("update card url_internal: %w", err)
		}
	}
	if req.Icon != nil {
		if _, err := r.db.Exec("UPDATE cards SET icon = ? WHERE id = ?", *req.Icon, cardID); err != nil {
			return fmt.Errorf("update card icon: %w", err)
		}
	}
	if req.IconColor != nil {
		if _, err := r.db.Exec("UPDATE cards SET icon_color = ? WHERE id = ?", *req.IconColor, cardID); err != nil {
			return fmt.Errorf("update card icon_color: %w", err)
		}
	}
	if req.BgColor != nil {
		if _, err := r.db.Exec("UPDATE cards SET bg_color = ? WHERE id = ?", *req.BgColor, cardID); err != nil {
			return fmt.Errorf("update card bg_color: %w", err)
		}
	}
	if req.Description != nil {
		if _, err := r.db.Exec("UPDATE cards SET description = ? WHERE id = ?", *req.Description, cardID); err != nil {
			return fmt.Errorf("update card description: %w", err)
		}
	}
	if req.OpenType != nil {
		if _, err := r.db.Exec("UPDATE cards SET open_type = ? WHERE id = ?", *req.OpenType, cardID); err != nil {
			return fmt.Errorf("update card open_type: %w", err)
		}
	}
	if req.SortOrder != nil {
		if _, err := r.db.Exec("UPDATE cards SET sort_order = ? WHERE id = ?", *req.SortOrder, cardID); err != nil {
			return fmt.Errorf("update card sort_order: %w", err)
		}
	}
	return nil
}

func (r *PanelRepo) DeleteCard(cardID string) error {
	if _, err := r.db.Exec("DELETE FROM cards WHERE id = ?", cardID); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	return nil
}

// Reorder applies group/card ordering in a single transaction.
func (r *PanelRepo) Reorder(groupOrders []models.GroupOrder, cardOrders []models.CardOrder) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin reorder: %w", err)
	}
	defer tx.Rollback()

	for _, g := range groupOrders {
		if _, err := tx.Exec("UPDATE panel_groups SET sort_order = ? WHERE id = ?", g.SortOrder, g.ID); err != nil {
			return fmt.Errorf("reorder group: %w", err)
		}
	}
	for _, c := range cardOrders {
		if _, err := tx.Exec("UPDATE cards SET sort_order = ?, group_id = ? WHERE id = ?", c.SortOrder, c.GroupID, c.ID); err != nil {
			return fmt.Errorf("reorder card: %w", err)
		}
	}
	return tx.Commit()
}
