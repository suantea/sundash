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

func TestLoginDefaultAdmin(t *testing.T) {
	svc := setup(t)
	res, err := svc.Login("admin", "admin")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if res.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if res.User.Username != "admin" {
		t.Fatalf("expected admin, got %s", res.User.Username)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	svc := setup(t)
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
	if err := svc.ChangePassword("admin", "admin", "newpass123"); err != nil {
		t.Fatalf("change password failed: %v", err)
	}
	if _, err := svc.Login("admin", "admin"); err == nil {
		t.Fatal("old password should no longer work")
	}
	if _, err := svc.Login("admin", "newpass123"); err != nil {
		t.Fatalf("new password should work: %v", err)
	}
}
