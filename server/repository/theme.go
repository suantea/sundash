package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CSSContent  string    `json:"css_content"`
	Author      string    `json:"author"`
	IsBuiltin   bool      `json:"is_builtin"`
	CreatedAt   time.Time `json:"created_at"`
}

type ThemeRepo struct {
	db *sql.DB
}

func NewThemeRepo(db *sql.DB) *ThemeRepo {
	return &ThemeRepo{db: db}
}

func (r *ThemeRepo) List() ([]Theme, error) {
	rows, err := r.db.Query("SELECT id, name, description, css_content, author, is_builtin, created_at FROM themes ORDER BY is_builtin DESC, name")
	if err != nil {
		return nil, fmt.Errorf("list themes: %w", err)
	}
	defer rows.Close()
	var themes []Theme
	for rows.Next() {
		var t Theme
		var builtin int
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CSSContent, &t.Author, &builtin, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan theme: %w", err)
		}
		t.IsBuiltin = builtin == 1
		themes = append(themes, t)
	}
	return themes, rows.Err()
}

func (r *ThemeRepo) GetByID(id string) (*Theme, error) {
	var t Theme
	var builtin int
	err := r.db.QueryRow("SELECT id, name, description, css_content, author, is_builtin, created_at FROM themes WHERE id = ?", id).
		Scan(&t.ID, &t.Name, &t.Description, &t.CSSContent, &t.Author, &builtin, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get theme: %w", err)
	}
	t.IsBuiltin = builtin == 1
	return &t, nil
}

func (r *ThemeRepo) Create(t *Theme) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	_, err := r.db.Exec("INSERT INTO themes (id, name, description, css_content, author, is_builtin) VALUES (?, ?, ?, ?, ?, ?)",
		t.ID, t.Name, t.Description, t.CSSContent, t.Author, boolToInt(t.IsBuiltin))
	return err
}

func (r *ThemeRepo) Update(t *Theme) error {
	_, err := r.db.Exec("UPDATE themes SET name=?, description=?, css_content=?, author=? WHERE id=?",
		t.Name, t.Description, t.CSSContent, t.Author, t.ID)
	return err
}

func (r *ThemeRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM themes WHERE id = ? AND is_builtin = 0", id)
	return err
}
