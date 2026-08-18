package handlers

import (
	"net/http"

	"sundash/models"
	"sundash/service"

	"github.com/gin-gonic/gin"
)

type PanelHandler struct {
	panels   *service.PanelService
	settings *service.SettingsService
}

func NewPanelHandler(panels *service.PanelService, settings *service.SettingsService) *PanelHandler {
	return &PanelHandler{panels: panels, settings: settings}
}

func (h *PanelHandler) GetPanel(c *gin.Context) {
	data, err := h.panels.GetPanel(c.GetString("user_id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *PanelHandler) CreateGroup(c *gin.Context) {
	var req models.CreateGroupRequest
	if !bind(c, &req) {
		return
	}
	group, err := h.panels.CreateGroup(c.GetString("user_id"), req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, group)
}

func (h *PanelHandler) UpdateGroup(c *gin.Context) {
	var req models.UpdateGroupRequest
	if !bind(c, &req) {
		return
	}
	if err := h.panels.UpdateGroup(c.GetString("user_id"), c.Param("id"), req.Name, req.SortOrder); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group updated"})
}

func (h *PanelHandler) DeleteGroup(c *gin.Context) {
	if err := h.panels.DeleteGroup(c.GetString("user_id"), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

func (h *PanelHandler) CreateCard(c *gin.Context) {
	var req models.CreateCardRequest
	if !bind(c, &req) {
		return
	}
	card, err := h.panels.CreateCard(c.GetString("user_id"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, card)
}

func (h *PanelHandler) UpdateCard(c *gin.Context) {
	var req models.UpdateCardRequest
	if !bind(c, &req) {
		return
	}
	if err := h.panels.UpdateCard(c.GetString("user_id"), c.Param("id"), req); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card updated"})
}

func (h *PanelHandler) DeleteCard(c *gin.Context) {
	if err := h.panels.DeleteCard(c.GetString("user_id"), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card deleted"})
}

func (h *PanelHandler) Reorder(c *gin.Context) {
	var req models.ReorderRequest
	if !bind(c, &req) {
		return
	}
	if err := h.panels.Reorder(c.GetString("user_id"), req); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated"})
}

func (h *PanelHandler) GetSettings(c *gin.Context) {
	settings, err := h.settings.GetForUser(c.GetString("user_id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *PanelHandler) UpdateSetting(c *gin.Context) {
	var req models.UpdateSettingRequest
	if !bind(c, &req) {
		return
	}
	if err := h.settings.UpdateForUser(c.GetString("user_id"), req.Key, req.Value); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Setting updated"})
}

func (h *PanelHandler) BatchUpdateSettings(c *gin.Context) {
	var req struct {
		Settings map[string]string `json:"settings" binding:"required"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.settings.BatchUpdateForUser(c.GetString("user_id"), req.Settings); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}
