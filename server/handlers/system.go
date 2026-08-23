package handlers

import (
	"net/http"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	system *service.SystemService
}

func NewSystemHandler(system *service.SystemService) *SystemHandler {
	return &SystemHandler{system: system}
}

func (h *SystemHandler) GetStats(c *gin.Context) {
	stats, err := h.system.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
