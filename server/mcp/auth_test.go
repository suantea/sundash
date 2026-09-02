package mcp

import (
	"net/http"
	"path/filepath"
	"testing"

	"sundash/database"
	"sundash/repository"
)

func TestMCPAuthStaticToken(t *testing.T) {
	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	defer db.Close()
	users := repository.NewUserRepo(db)
	// database.Init no longer auto-creates an admin; static-token auth looks up
	// mcpUsername, so create the admin user first.
	ensureAdmin(t, users)

	req, _ := http.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer test-static-mcp-token")
	uid := Auth(req, users, []byte("jwt-secret"), "test-static-mcp-token", "admin")
	if uid == "" {
		t.Fatal("static token should authenticate as admin")
	}

	// 错误 token → 拒绝
	req2, _ := http.NewRequest("POST", "/mcp", nil)
	req2.Header.Set("Authorization", "Bearer wrong-token")
	if uid2 := Auth(req2, users, []byte("jwt-secret"), "test-static-mcp-token", "admin"); uid2 != "" {
		t.Fatalf("wrong token should be rejected, got %q", uid2)
	}

	// 无 header → 拒绝
	req3, _ := http.NewRequest("POST", "/mcp", nil)
	if uid3 := Auth(req3, users, []byte("jwt-secret"), "test-static-mcp-token", "admin"); uid3 != "" {
		t.Fatalf("no header should be rejected, got %q", uid3)
	}
}
