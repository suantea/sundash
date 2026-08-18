package repository

import (
	"database/sql"
	"fmt"
)

type SettingsRepo struct {
	db *sql.DB
}

func NewSettingsRepo(db *sql.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

// GlobalByKeys returns the global (user_id IS NULL) values for the given keys.
func (r *SettingsRepo) GlobalByKeys(keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		out[key] = ""
		var v string
		if err := r.db.QueryRow(
			"SELECT COALESCE(value, '') FROM settings WHERE key = ? AND user_id IS NULL", key,
		).Scan(&v); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("query setting %q: %w", key, err)
		}
		out[key] = v
	}
	return out, nil
}

// AllForUser returns the merged view for a user: global settings overridden by
// the user's own values.
func (r *SettingsRepo) AllForUser(userID string) (map[string]string, error) {
	rows, err := r.db.Query(
		"SELECT key, COALESCE(value,''), COALESCE(user_id,'') FROM settings WHERE user_id = ? OR user_id IS NULL ORDER BY key",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query user settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value, uid string
		if err := rows.Scan(&key, &value, &uid); err != nil {
			return nil, err
		}
		if uid == userID || settings[key] == "" {
			settings[key] = value
		}
	}
	return settings, rows.Err()
}

// GetValue returns a single setting; userID may be nil for global.
func (r *SettingsRepo) GetValue(key string, userID *string) (string, error) {
	var v string
	var err error
	if userID == nil {
		err = r.db.QueryRow("SELECT COALESCE(value, '') FROM settings WHERE key = ? AND user_id IS NULL", key).Scan(&v)
	} else {
		err = r.db.QueryRow("SELECT COALESCE(value, '') FROM settings WHERE key = ? AND user_id = ?", key, *userID).Scan(&v)
	}
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("query setting %q: %w", key, err)
	}
	return v, nil
}

// Upsert creates or updates a single setting. userID nil → global.
func (r *SettingsRepo) Upsert(key, value string, userID *string) error {
	if userID == nil {
		if _, err := r.db.Exec(
			`INSERT INTO settings (key, value, user_id, updated_at) VALUES (?, ?, NULL, CURRENT_TIMESTAMP)
			 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`,
			key, value,
		); err != nil {
			return fmt.Errorf("upsert global setting: %w", err)
		}
		return nil
	}
	if _, err := r.db.Exec(
		`INSERT INTO settings (key, value, user_id, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(key, user_id) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`,
		key, value, *userID,
	); err != nil {
		return fmt.Errorf("upsert user setting: %w", err)
	}
	return nil
}

// UpsertBatch updates many settings in one transaction.
func (r *SettingsRepo) UpsertBatch(settings map[string]string, userID *string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin upsert batch: %w", err)
	}
	defer tx.Rollback()

	for key, value := range settings {
		if userID == nil {
			_, err = tx.Exec(
				`INSERT INTO settings (key, value, user_id, updated_at) VALUES (?, ?, NULL, CURRENT_TIMESTAMP)
				 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`,
				key, value,
			)
		} else {
			_, err = tx.Exec(
				`INSERT INTO settings (key, value, user_id, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)
				 ON CONFLICT(key, user_id) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`,
				key, value, *userID,
			)
		}
		if err != nil {
			return fmt.Errorf("upsert setting %q: %w", key, err)
		}
	}
	return tx.Commit()
}
