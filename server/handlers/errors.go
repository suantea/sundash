package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"sundash/service"

	"github.com/gin-gonic/gin"
)

// respondError maps a service error to an HTTP response. Internal errors are
// logged and returned as a generic message, never leaking internals.
func respondError(c *gin.Context, err error) {
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		body := gin.H{"error": appErr.Msg}
		if appErr.Status != "" {
			body["status"] = appErr.Status
		}
		c.JSON(appErr.HTTPStatus(), body)
		return
	}
	slog.Error("internal error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

// bind decodes the JSON body and emits a 400 on failure.
func bind(c *gin.Context, v any) bool {
	if err := c.ShouldBindJSON(v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
