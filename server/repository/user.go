package repository

import (
	"database/sql"
	"fmt"

	"sundash/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// FindByUsername returns a user including the password hash (login path).
func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, COALESCE(display_name,''), COALESCE(avatar,''), role, COALESCE(status,'approved')
		 FROM users WHERE username = ?`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.DisplayName, &u.Avatar, &u.Role, &u.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &u, nil
}

// FindByID returns a user without the password hash.
func (r *UserRepo) FindByID(id string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		`SELECT id, username, COALESCE(display_name,''), COALESCE(avatar,''), role, COALESCE(status,'approved'), created_at
		 FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Username, &u.DisplayName, &u.Avatar, &u.Role, &u.Status, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) Count() (int, error) {
	var n int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&n); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return n, nil
}

func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var n int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&n); err != nil {
		return false, fmt.Errorf("check username: %w", err)
	}
	return n > 0, nil
}

func (r *UserRepo) Create(u *models.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, password_hash, display_name, role, status) VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.PasswordHash, u.DisplayName, u.Role, u.Status,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// PasswordHash returns the stored hash for a user (empty string if missing).
func (r *UserRepo) PasswordHash(id string) (string, error) {
	var hash string
	err := r.db.QueryRow("SELECT password_hash FROM users WHERE id = ?", id).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("query password hash: %w", err)
	}
	return hash, nil
}

func (r *UserRepo) UpdateProfile(id string, displayName, avatar *string) error {
	if displayName != nil {
		if _, err := r.db.Exec("UPDATE users SET display_name = ? WHERE id = ?", *displayName, id); err != nil {
			return fmt.Errorf("update display_name: %w", err)
		}
	}
	if avatar != nil {
		if _, err := r.db.Exec("UPDATE users SET avatar = ? WHERE id = ?", *avatar, id); err != nil {
			return fmt.Errorf("update avatar: %w", err)
		}
	}
	return nil
}

func (r *UserRepo) UpdatePassword(id, hash string) error {
	if _, err := r.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hash, id); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateRole(id, role string) error {
	if _, err := r.db.Exec("UPDATE users SET role = ? WHERE id = ?", role, id); err != nil {
		return fmt.Errorf("update role: %w", err)
	}
	return nil
}

func (r *UserRepo) SetStatus(id, status string) error {
	if _, err := r.db.Exec("UPDATE users SET status = ? WHERE id = ?", status, id); err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	return nil
}

// Delete removes a user. cards/panel_groups cascade via FK; settings has no FK
// and is deleted explicitly.
func (r *UserRepo) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin delete user: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM settings WHERE user_id = ?", id); err != nil {
		return fmt.Errorf("delete user settings: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM users WHERE id = ?", id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return tx.Commit()
}

func (r *UserRepo) List() ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, COALESCE(display_name,''), COALESCE(avatar,''), role, COALESCE(status,'approved'), created_at
		 FROM users ORDER BY created_at`,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Avatar, &u.Role, &u.Status, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
