package handlers

import (
	"net/http"
	"time"

	"sundash/models"
	"sundash/service"

	"github.com/gin-gonic/gin"
)

type BmsyncHandler struct {
	bmsync *service.BmsyncService
}

func NewBmsyncHandler(bmsync *service.BmsyncService) *BmsyncHandler {
	return &BmsyncHandler{bmsync: bmsync}
}

// GET /api/bmsync/status  — 配置状态与是否已同步过
func (h *BmsyncHandler) Status(c *gin.Context) {
	cfg, err := h.bmsync.GetConfig(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	rev, _ := h.bmsync.Repo().GetRev(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"configured": cfg.ServerURL != "" && cfg.Token != "",
		"serverUrl":  cfg.ServerURL,
		"hasSynced":  rev > 0,
		"rev":        rev,
	})
}

// GET /api/bmsync/pull  — 拉取全量规范状态（仅读，不推本地变更）
func (h *BmsyncHandler) Pull(c *gin.Context) {
	if !h.bmsync.IsConfigured(c.Request.Context()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrBmsyncNotConfigured.Error()})
		return
	}
	out, err := h.bmsync.Pull(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"rev":   out.Rev,
		"nodes": out.Nodes,
	})
}

// POST /api/bmsync/push  — 推送变更并收到新的全量规范状态
func (h *BmsyncHandler) Push(c *gin.Context) {
	if !h.bmsync.IsConfigured(c.Request.Context()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrBmsyncNotConfigured.Error()})
		return
	}
	var req struct {
		Changes []*models.ChangeIntent `json:"changes" binding:"required"`
	}
	if !bind(c, &req) {
		return
	}
	// Stamp timestamps on each change according to LWW+tombstone protocol.
	now := time.Now()
	changes := make([]*models.SyncNode, 0, len(req.Changes))
	for _, intent := range req.Changes {
		changes = append(changes, service.ApplyChange(now, *intent))
	}
	out, err := h.bmsync.Push(c.Request.Context(), changes)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"rev":   out.Rev,
		"nodes": out.Nodes,
	})
}

// GET /api/bmsync/tree  — 返回本地镜像的树（仅供前端展示/搜索，不触发网络）
func (h *BmsyncHandler) Tree(c *gin.Context) {
	nodes, err := h.bmsync.Repo().ListAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	rev, _ := h.bmsync.Repo().GetRev(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"rev":   rev,
		"nodes": nodes,
	})
}
