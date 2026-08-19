package mcp

import (
	"context"
	"net/http"
	"strings"

	"sundash/models"
	"sundash/repository"

	"github.com/golang-jwt/jwt/v5"
)

// Auth verifies the request's Authorization header against either the static
// MCP token (SUNDASH_MCP_TOKEN) or a regular sundash JWT, and returns the
// authenticated user id.
//
//   - Bearer <MCPToken>        → binds to MCPUsername's user (default admin)
//   - Bearer <JWT>             → validated with the sundash JWT secret
//
// Returns the resolved user id, or "" when unauthenticated.
func Auth(r *http.Request, users *repository.UserRepo, jwtSecret []byte, mcpToken, mcpUsername string) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	token := parts[1]

	// 1) Static MCP token (highest priority when configured).
	if mcpToken != "" && token == mcpToken {
		u, err := users.FindByUsername(mcpUsername)
		if err != nil || u == nil {
			return ""
		}
		return u.ID
	}

	// 2) Regular JWT.
	claims := &models.Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("sundash"), jwt.WithAudience("sundash-web"))
	if err != nil || !parsed.Valid {
		return ""
	}
	return claims.UserID
}

// WithUserID returns a copy of ctx carrying the authenticated user id.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
