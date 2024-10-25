package model

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Role           UserRole  `json:"role"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Email           string    `json:"email"`
	UserID          int       `json:"user_id"`
	Role            string    `json:"role"`
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	ExpAccessToken  time.Time `json:"exp_access_token"`
	ExpRefreshToken time.Time `json:"exp_refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
