package service

import (
	"fmt"
	"time"

	"sundash/models"
	"sundash/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	StatusApproved = "approved"
	StatusPending  = "pending"
	StatusRejected = "rejected"
)

type UserService struct {
	users     *repository.UserRepo
	settings  *repository.SettingsRepo
	jwtSecret []byte
}

func NewUserService(users *repository.UserRepo, settings *repository.SettingsRepo, jwtSecret []byte) *UserService {
	return &UserService{users: users, settings: settings, jwtSecret: jwtSecret}
}

// AuthSettings returns the public registration policy.
func (s *UserService) AuthSettings() (map[string]bool, error) {
	v, err := s.settings.GlobalByKeys([]string{"allow_registration", "require_approval"})
	if err != nil {
		return nil, err
	}
	return map[string]bool{
		"allow_registration": v["allow_registration"] == "true",
		"require_approval":   v["require_approval"] == "true",
	}, nil
}

func (s *UserService) Login(username, password string) (*models.LoginResponse, error) {
	user, err := s.users.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewError(ErrUnauthorized, "Invalid username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, NewError(ErrUnauthorized, "Invalid username or password")
	}

	switch user.Status {
	case StatusPending:
		return nil, &AppError{Kind: ErrForbidden, Msg: "Your account is pending approval. Please wait for admin to approve.", Status: StatusPending}
	case StatusRejected:
		return nil, &AppError{Kind: ErrForbidden, Msg: "Your account has been rejected.", Status: StatusRejected}
	}

	token, err := s.generateToken(*user)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{Token: token, User: *user}, nil
}

// RegisterResult mirrors both possible registration responses.
type RegisterResult struct {
	Token   string      `json:"token,omitempty"`
	User    models.User `json:"user"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
}

func (s *UserService) Register(req models.RegisterRequest) (*RegisterResult, error) {
	policy, err := s.settings.GlobalByKeys([]string{"allow_registration", "require_approval"})
	if err != nil {
		return nil, err
	}

	userCount, err := s.users.Count()
	if err != nil {
		return nil, err
	}
	if userCount > 1 && policy["allow_registration"] != "true" {
		return nil, NewError(ErrForbidden, "Registration is disabled. Ask admin to enable it.")
	}

	exists, err := s.users.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, NewError(ErrConflict, "Username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	status := StatusApproved
	if policy["require_approval"] == "true" {
		status = StatusPending
	}
	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}

	user := models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Role:         "user",
		Status:       status,
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""

	if status == StatusPending {
		return &RegisterResult{
			User:    user,
			Status:  StatusPending,
			Message: "Registration successful. Please wait for admin approval.",
		}, nil
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	return &RegisterResult{Token: token, User: user, Status: StatusApproved}, nil
}

func (s *UserService) GetProfile(userID string) (*models.User, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewError(ErrNotFound, "User not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID string, displayName, avatar *string) error {
	return s.users.UpdateProfile(userID, displayName, avatar)
}

func (s *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	currentHash, err := s.users.PasswordHash(userID)
	if err != nil {
		return err
	}
	if currentHash == "" {
		return NewError(ErrNotFound, "User not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(oldPassword)); err != nil {
		return NewError(ErrBadRequest, "Current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}
	return s.users.UpdatePassword(userID, string(hash))
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.users.List()
}

// UpdateUser lets an admin edit any user; a normal user may only edit themselves.
func (s *UserService) UpdateUser(targetID, requesterID, requesterRole string, displayName, password, role *string) error {
	if requesterRole != "admin" && requesterID != targetID {
		return NewError(ErrForbidden, "Cannot update other users")
	}

	if displayName != nil {
		if err := s.users.UpdateProfile(targetID, displayName, nil); err != nil {
			return err
		}
	}
	if password != nil && len(*password) >= 6 {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		if err := s.users.UpdatePassword(targetID, string(hash)); err != nil {
			return err
		}
	}
	if role != nil && requesterRole == "admin" {
		if err := s.users.UpdateRole(targetID, *role); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) DeleteUser(targetID, requesterID string) error {
	if targetID == requesterID {
		return NewError(ErrBadRequest, "Cannot delete yourself")
	}
	return s.users.Delete(targetID)
}

func (s *UserService) AdminResetPassword(targetID, newPassword string) error {
	user, err := s.users.FindByID(targetID)
	if err != nil {
		return err
	}
	if user == nil {
		return NewError(ErrNotFound, "User not found")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.users.UpdatePassword(targetID, string(hash))
}

func (s *UserService) SetUserStatus(targetID, status string) error {
	user, err := s.users.FindByID(targetID)
	if err != nil {
		return err
	}
	if user == nil {
		return NewError(ErrNotFound, "User not found")
	}
	if user.Status != StatusPending {
		return NewError(ErrBadRequest, "User is not pending approval")
	}
	return s.users.SetStatus(targetID, status)
}

// NeedsSetup 报告系统是否处于"首次安装"状态（没有任何用户时返回 true）。
// 前端据此在首次访问时引导创建管理员账号，而不是使用硬编码的 admin/admin。
func (s *UserService) NeedsSetup() (bool, error) {
	n, err := s.users.Count()
	if err != nil {
		return false, err
	}
	return n == 0, nil
}

// SetupRequest 首次初始化请求。
type SetupRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=32"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name"`
	SiteTitle   string `json:"site_title"` // 可选：站点名称，空则保持默认
}

// Setup 首次初始化：仅在没有任何用户时允许，创建第一个管理员并可选设置站点名。
// 返回登录响应（token + user），前端可直接进入面板。
func (s *UserService) Setup(req SetupRequest) (*models.LoginResponse, error) {
	userCount, err := s.users.Count()
	if err != nil {
		return nil, err
	}
	if userCount > 0 {
		return nil, NewError(ErrForbidden, "System already initialized")
	}

	exists, err := s.users.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, NewError(ErrConflict, "Username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}

	user := models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Role:         "admin",
		Status:       StatusApproved,
	}
	if err := s.users.Create(&user); err != nil {
		return nil, err
	}

	// 可选：设置站点名称
	if req.SiteTitle != "" {
		if err := s.settings.Upsert(models.SettingSiteTitle, req.SiteTitle, nil); err != nil {
			return nil, err
		}
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return &models.LoginResponse{Token: token, User: user}, nil
}

func (s *UserService) generateToken(user models.User) (string, error) {
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "sundash",
			Audience:  jwt.ClaimStrings{"sundash-web"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
