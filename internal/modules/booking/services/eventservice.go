//go:generate mockgen -source=eventservice.go -destination=eventservice_mock.go -package=services
package services

import (
	"context"
	"errors"

	"github.com/Rhymond/go-money"

	"booking-event/internal/modules/booking/model"
)

type EventRepository interface {
	GetEventByID(ctx context.Context, id int) (*model.Event, error)
	QueryEvents(ctx context.Context, query model.EventQuery) ([]model.Event, error)
	CreateEvent(ctx context.Context, event model.Event, tokens []model.EventToken) error
	UpdateEvent(ctx context.Context, event model.Event) error
}

type EventTokenServiceForEvent interface {
	CreateEventToken(ctx context.Context, eventID int, userID int) (string, error)
}

type EventService struct {
	eventRepo EventRepository
	currency  string
	uuidFn    func() string
}

func NewEventService(eventRepo EventRepository, currency string, uuidFn func() string) *EventService {
	return &EventService{eventRepo: eventRepo, currency: currency, uuidFn: uuidFn}
}

func (s *EventService) RetrieveEventDetail(ctx context.Context, eventID int) (*model.Event, error) {
	event, err := s.eventRepo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) QueryEvents(ctx context.Context, query model.EventQuery) ([]model.Event, error) {
	events, err := s.eventRepo.QueryEvents(ctx, query)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *EventService) CreateEvent(ctx context.Context, params model.CreateEventRequest) error {
	m := money.NewFromFloat(params.Price, s.currency)
	event := model.Event{
		Name:           params.Name,
		AvailableSeats: params.AvailableSeats,
		StartAt:        params.StartAt,
		Location:       params.Location,
		Category:       params.Category,
		Status:         model.EventStatusInactive,
		Currency:       s.currency,
		Price:          m.Amount(),
		CreatorID:      params.ExecutorID,
	}
	tokens := make([]model.EventToken, params.AvailableSeats)
	for i := 0; i < params.AvailableSeats; i++ {
		tokens[i] = model.EventToken{EventID: event.ID, Token: s.uuidFn(), Status: model.TokenStatusActive}
	}
	return s.eventRepo.CreateEvent(ctx, event, tokens)
}

func (s *EventService) UpdateEvent(ctx context.Context, params model.UpdateEventRequest) error {
	event, err := s.eventRepo.GetEventByID(ctx, params.EventID)
	if err != nil {
		return err
	}
	if event.CreatorID != params.ExecutorID {
		return errors.New("unauthorized to update this event")
	}
	event.Status = params.Status
	return s.eventRepo.UpdateEvent(ctx, *event)
}
