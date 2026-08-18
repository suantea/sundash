package handlers

import (
	"net/http"

	"sundash/models"
	"sundash/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	users    *service.UserService
	panels   *service.PanelService
	settings *service.SettingsService
}

func NewAuthHandler(users *service.UserService, panels *service.PanelService, settings *service.SettingsService) *AuthHandler {
	return &AuthHandler{users: users, panels: panels, settings: settings}
}

// --- public auth ---

func (h *AuthHandler) GetAuthSettings(c *gin.Context) {
	settings, err := h.users.AuthSettings()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if !bind(c, &req) {
		return
	}
	res, err := h.users.Login(req.Username, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if !bind(c, &req) {
		return
	}
	res, err := h.users.Register(req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// --- profile (authenticated) ---

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.users.GetProfile(c.GetString("user_id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		DisplayName *string `json:"display_name"`
		Avatar      *string `json:"avatar"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.users.UpdateProfile(c.GetString("user_id"), req.DisplayName, req.Avatar); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.users.ChangePassword(c.GetString("user_id"), req.OldPassword, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// --- admin user management ---

func (h *AuthHandler) ListUsers(c *gin.Context) {
	users, err := h.users.ListUsers()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	targetID := c.Param("id")
	var req struct {
		DisplayName *string `json:"display_name"`
		Password    *string `json:"password"`
		Role        *string `json:"role"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.users.UpdateUser(targetID, c.GetString("user_id"), c.GetString("role"), req.DisplayName, req.Password, req.Role); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	if err := h.users.DeleteUser(c.Param("id"), c.GetString("user_id")); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (h *AuthHandler) AdminResetPassword(c *gin.Context) {
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.users.AdminResetPassword(c.Param("id"), req.NewPassword); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset for user " + c.Param("id")})
}

func (h *AuthHandler) ApproveUser(c *gin.Context) {
	if err := h.users.SetUserStatus(c.Param("id"), service.StatusApproved); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User approved"})
}

func (h *AuthHandler) RejectUser(c *gin.Context) {
	if err := h.users.SetUserStatus(c.Param("id"), service.StatusRejected); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User rejected"})
}

func (h *AuthHandler) GetUserPanel(c *gin.Context) {
	data, err := h.panels.GetUserPanel(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

// --- settings ---

func (h *AuthHandler) GetGlobalSettings(c *gin.Context) {
	settings, err := h.settings.GetGlobal()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *AuthHandler) UpdateGlobalSettings(c *gin.Context) {
	var req struct {
		Settings map[string]string `json:"settings" binding:"required"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.settings.UpdateGlobal(req.Settings); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}

func (h *AuthHandler) GetSiteConfig(c *gin.Context) {
	settings, err := h.settings.GetSiteConfig()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}
