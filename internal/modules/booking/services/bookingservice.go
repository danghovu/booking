//go:generate mockgen -source=bookingservice.go -destination=bookingservice_mock.go -package=services
package services

import (
	"context"
	"errors"
	"log"
	"time"

	"booking-event/internal/infra/paymentgateway"
	"booking-event/internal/modules/booking/model"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *model.Booking, bookingItems []model.BookingItem) error
	CountBookingByUserID(ctx context.Context, eventID int, userID int) (int, error)
	GetBookingByID(ctx context.Context, id int) (*model.Booking, error)
	ConfirmBooking(ctx context.Context, booking *model.Booking, event *model.Event, bookingItems []model.BookingItem) error
	CancelBooking(ctx context.Context, bookingID int) error
}

type BookingItemRepository interface {
	GetBookingItemsByBookingID(ctx context.Context, bookingID int) ([]model.BookingItem, error)
	CreateBookingItems(ctx context.Context, bookingItems []model.BookingItem) error
}

type EventServiceForBooking interface {
	GetEventByID(ctx context.Context, id int) (*model.Event, error)
}

type BookingEventTokenService interface {
	SelectAvailableToken(ctx context.Context, holderID int32, eventID int, quantity int) ([]string, error)
}

type BookingConfig struct {
	MaxBookingPerUser int
}

type BookingService struct {
	eventService      EventServiceForBooking
	eventTokenService BookingEventTokenService

	bookingRepository     BookingRepository
	bookingItemRepository BookingItemRepository
	cfg                   BookingConfig
}

func NewBookingService(
	eventService EventServiceForBooking,
	eventTokenService BookingEventTokenService,
	bookingRepo BookingRepository,
	bookingItemRepo BookingItemRepository,
	paymentService paymentgateway.PaymentGateway,
	cfg BookingConfig,
) *BookingService {
	return &BookingService{
		eventService:          eventService,
		eventTokenService:     eventTokenService,
		bookingRepository:     bookingRepo,
		bookingItemRepository: bookingItemRepo,
		cfg:                   cfg,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, booking model.CreateBookingRequest) (*model.Booking, error) {
	event, err := s.eventService.GetEventByID(ctx, booking.EventID)
	if err != nil {
		return nil, err
	}

	if event.Status != model.EventStatusActive {
		return nil, errors.New("event is not active")
	}

	count, err := s.bookingRepository.CountBookingByUserID(ctx, booking.EventID, booking.UserID)
	if err != nil {
		log.Println("error counting booking by user id", err)
		return nil, err
	}

	if count >= s.cfg.MaxBookingPerUser {
		return nil, errors.New("max booking per user reached")
	}

	tokens, err := s.eventTokenService.SelectAvailableToken(ctx, int32(booking.UserID), booking.EventID, booking.Quantity)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, errors.New("no available token")
	}

	bookingModel := &model.Booking{
		Status:          model.BookingStatusPending,
		UserID:          booking.UserID,
		EventID:         booking.EventID,
		InitialQuantity: booking.Quantity,
		Quantity:        len(tokens),
	}
	bookingItems := make([]model.BookingItem, len(tokens))
	for i, token := range tokens {
		bookingItems[i] = model.BookingItem{
			BookingID: bookingModel.ID,
			Token:     token,
		}
	}

	err = s.bookingRepository.CreateBooking(ctx, bookingModel, bookingItems)
	if err != nil {
		return nil, err
	}

	return bookingModel, nil
}

func (s *BookingService) ConfirmBooking(ctx context.Context, userID int, bookingID int) error {
	booking, err := s.bookingRepository.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.Status != model.BookingStatusPending {
		return errors.New("booking is not pending")
	}

	if booking.UserID != userID {
		return errors.New("unauthorized user is not allowed to confirm this booking")
	}

	event, err := s.eventService.GetEventByID(ctx, booking.EventID)
	if err != nil {
		return err
	}

	booking.Status = model.BookingStatusConfirmed

	bookingItems, err := s.bookingItemRepository.GetBookingItemsByBookingID(ctx, booking.ID)
	if err != nil {
		return err
	}

	return s.bookingRepository.ConfirmBooking(ctx, booking, event, bookingItems)
}

func (s *BookingService) GetBookingByID(ctx context.Context, id int) (*model.Booking, error) {
	return s.bookingRepository.GetBookingByID(ctx, id)
}

func (s *BookingService) CancelBooking(ctx context.Context, id int, executorID int) error {
	booking, err := s.bookingRepository.GetBookingByID(ctx, id)
	if err != nil {
		return err
	}
	if booking.UserID != executorID {
		return errors.New("unauthorized user is not allowed to cancel this booking")
	}
	if booking.Status == model.BookingStatusCanceled {
		return errors.New("booking is already canceled")
	}

	event, err := s.eventService.GetEventByID(ctx, booking.EventID)
	if err != nil {
		return err
	}

	if event.StartAt.Before(time.Now()) { // CANNOT CANCEL EVENT THAT HAS ALREADY STARTED
		return errors.New("event is already started")
	}

	return s.bookingRepository.CancelBooking(ctx, booking.ID)
}
