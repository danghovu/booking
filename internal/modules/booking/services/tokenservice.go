//go:generate mockgen -source=tokenservice.go -destination=tokenservice_mock.go -package=services
package services

import (
	"context"
	"fmt"

	"booking-event/internal/modules/booking/model"
)

type TokenRepository interface {
	CreateTokens(ctx context.Context, tokens []model.EventToken) error
	SelectAvailableToken(ctx context.Context, holderID int32, eventID int, quantity int) ([]string, error)
	ReleaseToken(ctx context.Context, eventToken *model.EventToken) error
	GetByToken(ctx context.Context, token string) (*model.EventToken, error)
}

type EventTokenService struct {
	tokenRepo TokenRepository
	uuidFn    func() string
}

func NewEventTokenService(tokenRepo TokenRepository, uuidFn func() string) *EventTokenService {
	return &EventTokenService{tokenRepo: tokenRepo, uuidFn: uuidFn}
}

func (s *EventTokenService) ReleaseToken(ctx context.Context, token string) error {
	eventToken, err := s.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("token not found")
	}

	return s.tokenRepo.ReleaseToken(ctx, eventToken)
}

func (s *EventTokenService) SelectAvailableToken(ctx context.Context, holderID int32, eventID int, quantity int) ([]string, error) {
	tokens, err := s.tokenRepo.SelectAvailableToken(ctx, holderID, eventID, quantity)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
