package service

import (
	"path/filepath"
	"testing"

	"sundash/database"
	"sundash/models"
	"sundash/repository"
)

func setup(t *testing.T) *UserService {
	t.Helper()
	db, err := database.Init(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return NewUserService(
		repository.NewUserRepo(db),
		repository.NewSettingsRepo(db),
		[]byte("test-secret"),
	)
}

func TestNeedsSetup(t *testing.T) {
	svc := setup(t)
	needs, err := svc.NeedsSetup()
	if err != nil {
		t.Fatalf("needs setup: %v", err)
	}
	if !needs {
		t.Fatal("expected needs_setup=true with no users")
	}

	if _, err := svc.Setup(SetupRequest{Username: "admin", Password: "admin123"}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	needs, err = svc.NeedsSetup()
	if err != nil {
		t.Fatalf("needs setup after: %v", err)
	}
	if needs {
		t.Fatal("expected needs_setup=false after setup")
	}
}

func TestSetupCreatesAdmin(t *testing.T) {
	svc := setup(t)
	res, err := svc.Setup(SetupRequest{Username: "admin", Password: "admin123", DisplayName: "Administrator"})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if res.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if res.User.Username != "admin" {
		t.Fatalf("expected admin, got %s", res.User.Username)
	}
	if res.User.Role != "admin" {
		t.Fatalf("expected admin role, got %s", res.User.Role)
	}

	// 再次 Setup 应被拒绝
	if _, err := svc.Setup(SetupRequest{Username: "admin", Password: "admin123"}); err == nil {
		t.Fatal("expected error when setup called twice")
	}

	// 用设置好的账号登录
	login, err := svc.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("login after setup failed: %v", err)
	}
	if login.User.Username != "admin" {
		t.Fatalf("expected admin, got %s", login.User.Username)
	}
}

func TestLoginDefaultAdmin(t *testing.T) {
	svc := setup(t)
	// 不再有自动创建的 admin/admin：默认密码登录必须失败
	if _, err := svc.Login("admin", "admin"); err == nil {
		t.Fatal("expected error: no default admin account exists")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	svc := setup(t)
	if _, err := svc.Setup(SetupRequest{Username: "admin", Password: "admin123"}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if _, err := svc.Login("admin", "wrong"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestRegisterAndLogin(t *testing.T) {
	svc := setup(t)
	res, err := svc.Register(models.RegisterRequest{Username: "alice", Password: "secret123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if res.Status != StatusApproved {
		t.Fatalf("expected approved, got %s", res.Status)
	}
	if res.Token == "" {
		t.Fatal("expected token after register")
	}

	login, err := svc.Login("alice", "secret123")
	if err != nil {
		t.Fatalf("login after register failed: %v", err)
	}
	if login.User.Username != "alice" {
		t.Fatalf("expected alice, got %s", login.User.Username)
	}
}

func TestRegisterDuplicateUsername(t *testing.T) {
	svc := setup(t)
	// Registration is locked after the first user; enable it for this test.
	if err := svc.settings.Upsert("allow_registration", "true", nil); err != nil {
		t.Fatalf("enable registration: %v", err)
	}
	if _, err := svc.Register(models.RegisterRequest{Username: "bob", Password: "secret123"}); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if _, err := svc.Register(models.RegisterRequest{Username: "bob", Password: "other456"}); err == nil {
		t.Fatal("expected duplicate username error")
	}
}

func TestChangePassword(t *testing.T) {
	svc := setup(t)
	// 先通过首次设置创建 admin 账号（返回的 User.ID 是 uuid）
	setupRes, err := svc.Setup(SetupRequest{Username: "admin", Password: "admin123"})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	adminID := setupRes.User.ID
	if err := svc.ChangePassword(adminID, "admin123", "newpass123"); err != nil {
		t.Fatalf("change password failed: %v", err)
	}
	if _, err := svc.Login("admin", "admin123"); err == nil {
		t.Fatal("old password should no longer work")
	}
	if _, err := svc.Login("admin", "newpass123"); err != nil {
		t.Fatalf("new password should work: %v", err)
	}
}
