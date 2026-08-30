package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"sundash/database"

	"github.com/gin-gonic/gin"
)

// The admin backup endpoint must produce a valid, complete SQLite snapshot —
// DEPLOY.md's backup guidance depends on it, and the client-side JSON export
// only covers bookmarks (no users/settings/memos/feeds).
func TestBackupDownloadProducesSQLiteSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Init(dbPath)
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	// Some data so the snapshot has content.
	if _, err := db.Exec(`INSERT INTO memos (id, user_id, content) VALUES ('m1', 'admin', 'backup-me')`); err != nil {
		t.Fatalf("seed memo: %v", err)
	}

	h := NewBackupHandler(db)
	r := gin.New()
	r.GET("/api/admin/backup", h.Download)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/backup", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("backup: got %d: %s", res.Code, res.Body.String())
	}
	if ct := res.Header().Get("Content-Disposition"); !strings.Contains(ct, "sundash-backup-") {
		t.Fatalf("unexpected Content-Disposition: %q", ct)
	}
	body := res.Body.Bytes()
	// SQLite file magic.
	if len(body) < 16 || string(body[:15]) != "SQLite format 3" {
		t.Fatalf("snapshot is not a SQLite database (header: %q)", body[:min(len(body), 16)])
	}
	// Content survived the round trip into the snapshot.
	if !strings.Contains(string(body), "backup-me") {
		t.Fatal("seeded memo missing from snapshot")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
