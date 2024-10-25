package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken     string
	RefreshToken    string
	AccessTokenExp  time.Time
	RefreshTokenExp time.Time
}
