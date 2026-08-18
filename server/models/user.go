package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	Avatar       string    `json:"avatar"`
	Role         string    `json:"role"`
	Status       string    `json:"status"` // approved, pending, rejected
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=32"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name"`
}
