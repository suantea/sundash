package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	// 显式校验目录可写性：MkdirAll 对已存在的目录不会报错，若目录存在但
	// 不可写（如 NAS 卷权限异常、只读挂载），SQLite 打开数据库时会报晦涩的
	// "unable to open database file (14)"（CANTOPEN），难以定位。这里先创建
	// 并删除一个临时文件探测写权限，失败时给出清晰错误。
	probe, err := os.CreateTemp(dir, ".sundash-write-test-*")
	if err != nil {
		return nil, fmt.Errorf("data directory %q is not writable (cannot create database file): %w", dir, err)
	}
	probeName := probe.Name()
	probe.Close()
	_ = os.Remove(probeName)

	// _pragma=foreign_keys(1) enables FK cascades (ON DELETE CASCADE) which the schema relies on.
	// _pragma=synchronous(NORMAL) balances safety and performance.
	// _pragma=cache_size=-20000 sets cache to 20MB.
	dsn := "file:" + dbPath +
		"?_journal_mode=WAL" +
		"&_busy_timeout=5000" +
		"&_pragma=foreign_keys(1)" +
		"&_pragma=synchronous(NORMAL)" +
		"&_pragma=cache_size=-20000"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	// WAL mode supports concurrent readers; a small pool lets parallel API
	// requests (panel + settings + wallpaper) avoid serializing on one conn.
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, err
	}
	if err := createDefaultAdmin(db); err != nil {
		db.Close()
		return nil, err
	}

	DB = db
	log.Println("Database initialized successfully")
	return db, nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

// --- Versioned migrations ---

type migration struct {
	version int
	sql     string
	table   string // when set, the migration only runs if the column is missing
	column  string
}

var migrations = []migration{
	{1, `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		display_name TEXT,
		avatar TEXT,
		role TEXT DEFAULT 'user',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`, "", ""},
	{2, `CREATE TABLE IF NOT EXISTS panel_groups (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		name TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, "", ""},
	{3, `CREATE TABLE IF NOT EXISTS cards (
		id TEXT PRIMARY KEY,
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		url_internal TEXT,
		icon TEXT,
		icon_color TEXT,
		bg_color TEXT,
		description TEXT,
		open_type TEXT DEFAULT 'new_tab',
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES panel_groups(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, "", ""},
	// Column additions are guarded by existence checks so both fresh and legacy databases migrate cleanly.
	{4, `ALTER TABLE cards ADD COLUMN bg_color TEXT`, "cards", "bg_color"},
	{5, `ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'approved'`, "users", "status"},
	{6, `CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT,
		user_id TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`, "", ""},
	{7, `CREATE INDEX IF NOT EXISTS idx_panel_groups_user ON panel_groups(user_id)`, "", ""},
	{8, `CREATE INDEX IF NOT EXISTS idx_cards_group ON cards(group_id)`, "", ""},
	{9, `CREATE INDEX IF NOT EXISTS idx_cards_user ON cards(user_id)`, "", ""},
	{10, `CREATE INDEX IF NOT EXISTS idx_settings_user ON settings(user_id)`, "", ""},
	{11, `CREATE UNIQUE INDEX IF NOT EXISTS idx_settings_key_user ON settings(key, user_id) WHERE user_id IS NOT NULL`, "", ""},
	// FTS5 full-text search index on cards (title, url, description)
	{12, `CREATE VIRTUAL TABLE IF NOT EXISTS cards_fts USING fts5(
		card_id UNINDEXED,
		title,
		url,
		description,
		user_id UNINDEXED,
		content='',
		contentless_delete=1
	)`, "", ""},
	// Triggers to keep FTS index in sync with cards table. The FTS rowid is
	// keyed to cards.rowid so the 'delete'-by-rowid form works. Note: this is
	// a contentless_delete=1 table, where the special INSERT ... VALUES
	// ('delete', ...) command is NOT allowed — only plain DELETE statements
	// are (see migration 24, which also rebuilds the old wrong triggers).
	{13, `CREATE TRIGGER IF NOT EXISTS cards_ai AFTER INSERT ON cards BEGIN
		INSERT INTO cards_fts(rowid, card_id, title, url, description, user_id)
		VALUES (new.rowid, new.id, new.title, new.url, COALESCE(new.description,''), new.user_id);
	END`, "", ""},
	{14, `CREATE TRIGGER IF NOT EXISTS cards_ad AFTER DELETE ON cards BEGIN
		DELETE FROM cards_fts WHERE rowid = old.rowid;
	END`, "", ""},
	{15, `CREATE TRIGGER IF NOT EXISTS cards_au AFTER UPDATE ON cards BEGIN
		DELETE FROM cards_fts WHERE rowid = old.rowid;
		INSERT INTO cards_fts(rowid, card_id, title, url, description, user_id)
		VALUES (new.rowid, new.id, new.title, new.url, COALESCE(new.description,''), new.user_id);
	END`, "", ""},
	// Memo table for quick notes
	{16, `CREATE TABLE IF NOT EXISTS memos (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		is_archived INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, "", ""},
	{17, `CREATE INDEX IF NOT EXISTS idx_memos_user ON memos(user_id)`, "", ""},
	{18, `CREATE INDEX IF NOT EXISTS idx_memos_updated_at ON memos(updated_at)`, "", ""},
	// RSS feeds table
	{19, `CREATE TABLE IF NOT EXISTS rss_feeds (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT,
		url TEXT NOT NULL,
		description TEXT,
		image_url TEXT,
		last_fetched DATETIME,
		update_interval INTEGER DEFAULT 60,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`, "", ""},
	{20, `CREATE INDEX IF NOT EXISTS idx_rss_feeds_user ON rss_feeds(user_id)`, "", ""},
	// RSS items table
	{21, `CREATE TABLE IF NOT EXISTS rss_items (
		id TEXT PRIMARY KEY,
		feed_id TEXT NOT NULL,
		title TEXT,
		link TEXT,
		description TEXT,
		pub_date DATETIME,
		author TEXT,
		guid TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (feed_id) REFERENCES rss_feeds(id) ON DELETE CASCADE
	)`, "", ""},
	{22, `CREATE INDEX IF NOT EXISTS idx_rss_items_feed ON rss_items(feed_id)`, "", ""},
	{23, `CREATE INDEX IF NOT EXISTS idx_rss_items_pub_date ON rss_items(pub_date)`, "", ""},
	// Repair the FTS sync for databases that already ran migrations 13-15:
	// the old triggers used the 'delete' INSERT command, which fails on a
	// contentless_delete=1 table (every card update/reorder errored), and
	// never keyed the FTS rowid to cards.rowid. Rebuild triggers and index.
	{24, `DROP TRIGGER IF EXISTS cards_ai;
	DROP TRIGGER IF EXISTS cards_ad;
	DROP TRIGGER IF EXISTS cards_au;
	CREATE TRIGGER cards_ai AFTER INSERT ON cards BEGIN
		INSERT INTO cards_fts(rowid, card_id, title, url, description, user_id)
		VALUES (new.rowid, new.id, new.title, new.url, COALESCE(new.description,''), new.user_id);
	END;
	CREATE TRIGGER cards_ad AFTER DELETE ON cards BEGIN
		DELETE FROM cards_fts WHERE rowid = old.rowid;
	END;
	CREATE TRIGGER cards_au AFTER UPDATE ON cards BEGIN
		DELETE FROM cards_fts WHERE rowid = old.rowid;
		INSERT INTO cards_fts(rowid, card_id, title, url, description, user_id)
		VALUES (new.rowid, new.id, new.title, new.url, COALESCE(new.description,''), new.user_id);
	END;
	DELETE FROM cards_fts;
	INSERT INTO cards_fts(rowid, card_id, title, url, description, user_id)
		SELECT rowid, id, title, url, COALESCE(description,''), user_id FROM cards;`, "", ""},
	// Bookmark-sync (bmsync): sundash keeps a local mirror of the canonical
	// bookmark tree (same shape as the bookmark-sync server) plus sync meta.
	{25, `CREATE TABLE IF NOT EXISTS bmsync_nodes (
		sync_id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		url TEXT,
		parent_sync_id TEXT,
		idx INTEGER,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		deleted_at TEXT,
		rev INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS bmsync_meta (key TEXT PRIMARY KEY, value TEXT);`, "", ""},
}

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	applied := map[int]bool{}
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("read schema_migrations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}

	for _, m := range migrations {
		if applied[m.version] {
			continue
		}
		// Guarded column additions: skip silently when the column already exists.
		if m.table != "" && columnExists(db, m.table, m.column) {
			if err := markApplied(db, m.version); err != nil {
				return err
			}
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(m.sql); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d failed: %w\nSQL: %s", m.version, err, m.sql)
		}
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.version); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("Applied database migration %d", m.version)
	}
	return nil
}

func columnExists(db *sql.DB, table, column string) bool {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return false
		}
		if name == column {
			return true
		}
	}
	return false
}

func markApplied(db *sql.DB, version int) error {
	_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
	return err
}

func createDefaultAdmin(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash default admin password: %w", err)
	}
	if _, err := db.Exec(
		`INSERT INTO users (id, username, password_hash, display_name, role, status) VALUES (?, ?, ?, ?, ?, 'approved')`,
		"admin", "admin", string(hash), "Administrator", "admin",
	); err != nil {
		return fmt.Errorf("create default admin: %w", err)
	}
	log.Println("Default admin user created (username: admin, password: admin) — CHANGE THE PASSWORD IMMEDIATELY")
	return nil
}
