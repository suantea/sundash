package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"sundash/database"
	"sundash/middleware"
	"sundash/repository"
	"sundash/service"

	"github.com/gin-gonic/gin"
)

// Regression: the auth middleware stores the authenticated user under the
// "user_id" context key. Memo/search/rss handlers once read "userID" instead,
// so every one of those protected routes returned 401 with a perfectly valid
// token — invisible for weeks because the build was broken at the same time.
// This pins the end-to-end contract: login → Authorization header → 200.
func TestMemoRoutesAuthorizeWithJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	users := repository.NewUserRepo(db)
	userSvc := service.NewUserService(users, repository.NewSettingsRepo(db), []byte("test-secret"))
	login, err := userSvc.Login("admin", "admin")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	memoSvc := service.NewMemoService(repository.NewMemoRepo(db))
	memoH := NewMemoHandler(memoSvc)

	r := gin.New()
	r.Use(middleware.AuthMiddleware([]byte("test-secret")))
	r.GET("/api/memo", memoH.ListMemos)
	r.POST("/api/memo", memoH.CreateMemo)

	// Create through the protected route.
	body := strings.NewReader(`{"content":"hello from test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/memo", body)
	req.Header.Set("Authorization", "Bearer "+login.Token)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("create memo: got %d: %s", res.Code, res.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &created); err != nil || created.ID == "" {
		t.Fatalf("create memo response: %v (%s)", err, res.Body.String())
	}

	// List through the protected route.
	req = httptest.NewRequest(http.MethodGet, "/api/memo", nil)
	req.Header.Set("Authorization", "Bearer "+login.Token)
	res = httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("list memos: got %d: %s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "hello from test") {
		t.Fatalf("created memo missing from list: %s", res.Body.String())
	}
}
