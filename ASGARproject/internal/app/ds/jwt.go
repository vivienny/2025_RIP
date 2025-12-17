package ds

import (
	"time"

	"github.com/golang-jwt/jwt/v5" // ТОЛЬКО v5!
)

type JWTClaims struct {
	jwt.RegisteredClaims        // Вместо StandardClaims!
	UserID               uint   `json:"user_id"`
	Login                string `json:"login"`
	Role                 string `json:"role"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	UserID      uint      `json:"user_id"`
	Login       string    `json:"login"`
	Role        string    `json:"role"`
}

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type AuthUser struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}
