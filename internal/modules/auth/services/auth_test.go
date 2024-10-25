package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"booking-event/internal/modules/auth/model"
)

func TestAuthService_LoginByEmail(t *testing.T) {
	password := "correctpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	t.Parallel()
	tests := []struct {
		name           string
		email          string
		password       string
		mockUserRepo   func(ctrl *gomock.Controller) *MockUserRepository
		expectedUser   *model.User
		expectedTokens *model.TokenPair
		expectedError  error
	}{
		{
			name:     "Successful login",
			email:    "test@example.com",
			password: "correctpassword",
			mockUserRepo: func(ctrl *gomock.Controller) *MockUserRepository {
				mock := NewMockUserRepository(ctrl)
				mock.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&model.User{
					ID:             1,
					Email:          "test@example.com",
					HashedPassword: string(hashedPassword),
				}, nil)
				return mock
			},
			expectedUser: &model.User{
				ID:             1,
				Email:          "test@example.com",
				HashedPassword: string(hashedPassword),
			},
			expectedTokens: &model.TokenPair{
				AccessToken:  "mocked_access_token",
				RefreshToken: "mocked_refresh_token",
			},
			expectedError: nil,
		},
		{
			name:     "User not found",
			email:    "nonexistent@example.com",
			password: "anypassword",
			mockUserRepo: func(ctrl *gomock.Controller) *MockUserRepository {
				mock := NewMockUserRepository(ctrl)
				mock.EXPECT().GetUserByEmail(gomock.Any(), "nonexistent@example.com").Return(nil, errors.New("user not found"))
				return mock
			},
			expectedUser:   nil,
			expectedTokens: nil,
			expectedError:  errors.New("user not found"),
		},
		{
			name:     "Incorrect password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockUserRepo: func(ctrl *gomock.Controller) *MockUserRepository {
				mock := NewMockUserRepository(ctrl)
				mock.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&model.User{
					ID:             1,
					Email:          "test@example.com",
					HashedPassword: string(hashedPassword),
				}, nil)
				return mock
			},
			expectedUser:   nil,
			expectedTokens: nil,
			expectedError:  errors.New("invalid password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := tt.mockUserRepo(ctrl)
			service := NewAuthService(mockUserRepo, AuthServiceConfig{
				SecretKey:       "test_secret",
				AccessTokenExp:  time.Hour,
				RefreshTokenExp: time.Hour * 24,
			})

			user, tokens, err := service.LoginByEmail(context.Background(), tt.email, tt.password)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, user)
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Parallel()
	secret := "test_secret"
	expiration := time.Hour
	expirationTime := time.Now().Add(expiration)
	claims := &model.TokenClaims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validTokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	tests := []struct {
		name          string
		refreshToken  string
		mockUserRepo  func(ctrl *gomock.Controller) *MockUserRepository
		expectedError error
	}{
		{
			name:         "Successful token refresh",
			refreshToken: validTokenString,
			mockUserRepo: func(ctrl *gomock.Controller) *MockUserRepository {
				mock := NewMockUserRepository(ctrl)
				mock.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(&model.User{
					ID:    1,
					Email: "test@example.com",
				}, nil)
				return mock
			},
			expectedError: nil,
		},
		{
			name:         "Invalid refresh token",
			refreshToken: "invalid_refresh_token",
			mockUserRepo: func(ctrl *gomock.Controller) *MockUserRepository {
				return NewMockUserRepository(ctrl)
			},
			expectedError: errors.New("invalid refresh token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := tt.mockUserRepo(ctrl)
			service := NewAuthService(mockUserRepo, AuthServiceConfig{
				SecretKey:       "test_secret",
				AccessTokenExp:  time.Hour,
				RefreshTokenExp: time.Hour * 24,
			})

			tokens, err := service.RefreshToken(context.Background(), tt.refreshToken)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.Equal(t, tt.refreshToken, tokens.RefreshToken)
			}
		})
	}
}

func TestAuthService_VerifyJWTToken(t *testing.T) {
	t.Parallel()
	secret := "test_secret"
	expiration := time.Hour
	expirationTime := time.Now().Add(expiration)
	claims := &model.TokenClaims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validTokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	tests := []struct {
		name          string
		token         string
		expectedError error
	}{
		{
			name:          "Valid token",
			token:         validTokenString,
			expectedError: nil,
		},
		{
			name:          "Invalid token",
			token:         validTokenString + "invalid",
			expectedError: errors.New("signature is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewAuthService(nil, AuthServiceConfig{
				SecretKey:       secret,
				AccessTokenExp:  expiration,
				RefreshTokenExp: time.Hour * 24,
			})

			claims, err := service.VerifyJWTToken(tt.token)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}
