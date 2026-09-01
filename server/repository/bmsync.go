package repository

import (
	"context"
	"database/sql"
	"fmt"

	"sundash/models"
)

// BmsyncRepo stores the local mirror of the canonical bookmark tree plus sync meta.
type BmsyncRepo struct {
	db *sql.DB
}

func NewBmsyncRepo(db *sql.DB) *BmsyncRepo {
	return &BmsyncRepo{db: db}
}

// ReplaceAll atomically replaces the whole mirror with `nodes` and sets `rev`.
func (r *BmsyncRepo) ReplaceAll(ctx context.Context, nodes []*models.SyncNode, rev int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin replace: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM bmsync_nodes`); err != nil {
		return fmt.Errorf("clear bmsync_nodes: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO bmsync_nodes(sync_id, type, title, url, parent_sync_id, idx,
			created_at, updated_at, deleted_at, rev)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()
	for _, n := range nodes {
		if _, err := stmt.ExecContext(ctx, n.SyncID, n.Type, n.Title, nullable(n.URL),
			nullable(n.ParentSyncID), n.Index, n.CreatedAt, n.UpdatedAt,
			nullable(n.DeletedAt), n.Rev); err != nil {
			return fmt.Errorf("insert node %s: %w", n.SyncID, err)
		}
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO bmsync_meta(key, value) VALUES('rev', ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`, fmt.Sprintf("%d", rev)); err != nil {
		return fmt.Errorf("set rev: %w", err)
	}
	return tx.Commit()
}

// ListAll returns every mirror node, including tombstones.
func (r *BmsyncRepo) ListAll(ctx context.Context) ([]*models.SyncNode, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT sync_id, type, title, url, parent_sync_id, idx, created_at, updated_at, deleted_at, rev
		FROM bmsync_nodes`)
	if err != nil {
		return nil, fmt.Errorf("list bmsync nodes: %w", err)
	}
	defer rows.Close()

	var out []*models.SyncNode
	for rows.Next() {
		var n models.SyncNode
		var url, parent, deleted sql.NullString
		if err := rows.Scan(&n.SyncID, &n.Type, &n.Title, &url, &parent,
			&n.Index, &n.CreatedAt, &n.UpdatedAt, &deleted, &n.Rev); err != nil {
			return nil, fmt.Errorf("scan bmsync node: %w", err)
		}
		n.URL = url.String
		n.ParentSyncID = parent.String
		n.DeletedAt = deleted.String
		out = append(out, &n)
	}
	return out, rows.Err()
}

// GetRev returns the last synced rev (0 when never synced).
func (r *BmsyncRepo) GetRev(ctx context.Context) (int, error) {
	var v string
	err := r.db.QueryRowContext(ctx,
		`SELECT value FROM bmsync_meta WHERE key = 'rev'`).Scan(&v)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get bmsync rev: %w", err)
	}
	rev := 0
	fmt.Sscanf(v, "%d", &rev)
	return rev, nil
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}
