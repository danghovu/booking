//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=services
package services

import (
	"context"
	"errors"
	"time"

	"booking-event/internal/modules/auth/model"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthServiceConfig struct {
	SecretKey       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

type AuthService struct {
	userDBRepo UserRepository
	cfg        AuthServiceConfig
}

func NewAuthService(userDBRepo UserRepository, cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		userDBRepo: userDBRepo,
		cfg:        cfg,
	}
}

func (s *AuthService) LoginByEmail(ctx context.Context, email string, password string) (*model.User, *model.TokenPair, error) {
	user, err := s.userDBRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if !s.comparePassword(password, user.HashedPassword) {
		return nil, nil, model.ErrInvalidPassword
	}

	// Generate JWT tokens
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *AuthService) comparePassword(password string, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (s *AuthService) generateTokenPair(user *model.User) (*model.TokenPair, error) {
	accessToken, err := s.generateJWTToken(user, s.cfg.AccessTokenExp)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateJWTToken(user, s.cfg.RefreshTokenExp)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessTokenExp:  time.Now().Add(s.cfg.AccessTokenExp),
		RefreshTokenExp: time.Now().Add(s.cfg.RefreshTokenExp),
	}, nil
}

func (s *AuthService) generateJWTToken(user *model.User, expiration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiration)
	claims := &model.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	// Verify the refresh token
	claims := &model.TokenClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userDBRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.generateJWTToken(user, s.cfg.AccessTokenExp)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) VerifyJWTToken(token string) (*model.TokenClaims, error) {
	claims := &model.TokenClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
