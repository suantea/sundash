package handlers

import (
	"net/http"

	"sundash/repository"

	"github.com/gin-gonic/gin"
)

type ThemeHandler struct {
	themes *repository.ThemeRepo
}

func NewThemeHandler(themes *repository.ThemeRepo) *ThemeHandler {
	return &ThemeHandler{themes: themes}
}

func (h *ThemeHandler) List(c *gin.Context) {
	themes, err := h.themes.List()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, themes)
}

func (h *ThemeHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		CSSContent  string `json:"css_content" binding:"required"`
		Author      string `json:"author"`
	}
	if !bind(c, &req) {
		return
	}
	theme := &repository.Theme{
		Name:        req.Name,
		Description: req.Description,
		CSSContent:  req.CSSContent,
		Author:      req.Author,
	}
	if err := h.themes.Create(theme); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, theme)
}

func (h *ThemeHandler) Update(c *gin.Context) {
	existing, err := h.themes.GetByID(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		CSSContent  *string `json:"css_content"`
		Author      *string `json:"author"`
	}
	if !bind(c, &req) {
		return
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.CSSContent != nil {
		existing.CSSContent = *req.CSSContent
	}
	if req.Author != nil {
		existing.Author = *req.Author
	}
	if err := h.themes.Update(existing); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *ThemeHandler) Delete(c *gin.Context) {
	if err := h.themes.Delete(c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Theme deleted"})
}
