package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"booking-event/internal/modules/booking/model"
)

func TestEventTokenService_ReleaseToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		token         string
		mockTokenRepo func(ctrl *gomock.Controller) *MockTokenRepository
		expectedError error
	}{
		{
			name:  "Successful token release",
			token: "token1",
			mockTokenRepo: func(ctrl *gomock.Controller) *MockTokenRepository {
				mock := NewMockTokenRepository(ctrl)
				mock.EXPECT().GetByToken(gomock.Any(), "token1").Return(&model.EventToken{Token: "token1"}, nil)
				mock.EXPECT().ReleaseToken(gomock.Any(), &model.EventToken{Token: "token1"}).Return(nil)
				return mock
			},
			expectedError: nil,
		},
		{
			name:  "Token not found",
			token: "nonexistent",
			mockTokenRepo: func(ctrl *gomock.Controller) *MockTokenRepository {
				mock := NewMockTokenRepository(ctrl)
				mock.EXPECT().GetByToken(gomock.Any(), "nonexistent").Return(nil, errors.New("not found"))
				return mock
			},
			expectedError: errors.New("token not found"),
		},
		{
			name:  "Error releasing token",
			token: "token2",
			mockTokenRepo: func(ctrl *gomock.Controller) *MockTokenRepository {
				mock := NewMockTokenRepository(ctrl)
				mock.EXPECT().GetByToken(gomock.Any(), "token2").Return(&model.EventToken{Token: "token2"}, nil)
				mock.EXPECT().ReleaseToken(gomock.Any(), &model.EventToken{Token: "token2"}).Return(errors.New("database error"))
				return mock
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenRepo := tt.mockTokenRepo(ctrl)
			service := NewEventTokenService(mockTokenRepo, nil)

			err := service.ReleaseToken(context.Background(), tt.token)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEventTokenService_SelectAvailableToken(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		holderID       int32
		eventID        int
		quantity       int
		mockTokenRepo  func(ctrl *gomock.Controller) *MockTokenRepository
		expectedTokens []string
		expectedError  error
	}{
		{
			name:     "Successful token selection",
			holderID: 1,
			eventID:  100,
			quantity: 2,
			mockTokenRepo: func(ctrl *gomock.Controller) *MockTokenRepository {
				mock := NewMockTokenRepository(ctrl)
				mock.EXPECT().SelectAvailableToken(gomock.Any(), int32(1), 100, 2).Return([]string{"token1", "token2"}, nil)
				return mock
			},
			expectedTokens: []string{"token1", "token2"},
			expectedError:  nil,
		},
		{
			name:     "No tokens available",
			holderID: 2,
			eventID:  101,
			quantity: 3,
			mockTokenRepo: func(ctrl *gomock.Controller) *MockTokenRepository {
				mock := NewMockTokenRepository(ctrl)
				mock.EXPECT().SelectAvailableToken(gomock.Any(), int32(2), 101, 3).Return([]string{}, nil)
				return mock
			},
			expectedTokens: []string{},
			expectedError:  nil,
		},
		{
			name:     "Error selecting tokens",
			holderID: 3,
			eventID:  102,
			quantity: 1,
			mockTokenRepo: func(ctrl *gomock.Controller) *MockTokenRepository {
				mock := NewMockTokenRepository(ctrl)
				mock.EXPECT().SelectAvailableToken(gomock.Any(), int32(3), 102, 1).Return(nil, errors.New("database error"))
				return mock
			},
			expectedTokens: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenRepo := tt.mockTokenRepo(ctrl)
			service := NewEventTokenService(mockTokenRepo, nil)

			tokens, err := service.SelectAvailableToken(context.Background(), tt.holderID, tt.eventID, tt.quantity)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTokens, tokens)
			}
		})
	}
}
