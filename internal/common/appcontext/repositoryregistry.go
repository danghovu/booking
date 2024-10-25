package appcontext

import (
	"booking-event/config"
	authRepo "booking-event/internal/modules/auth/repository/store"
	emailRepo "booking-event/internal/modules/booking/repository/client/email"
	bookingRepo "booking-event/internal/modules/booking/repository/store"
)

type RepositoryRegistry interface {
	EventRepository() *bookingRepo.EventRepository
	UserRepository() *authRepo.UserRepository
	BookingRepository() *bookingRepo.BookingRepository
	BookingEventTokenRepository() *bookingRepo.TokenRepository
	BookingItemRepository() *bookingRepo.BookingItemRepository
	BookingEmailRepository() *emailRepo.EmailClient
}

type repositoryRegistry struct {
	eventRepository             *bookingRepo.EventRepository
	userRepository              *authRepo.UserRepository
	bookingRepository           *bookingRepo.BookingRepository
	bookingEventTokenRepository *bookingRepo.TokenRepository
	bookingItemRepository       *bookingRepo.BookingItemRepository
	bookingEmailRepository      *emailRepo.EmailClient
}

func NewRepositoryRegistry(
	config config.Config,
	infraRegistry InfraRegistry,
) RepositoryRegistry {
	bookingTokenRepo := bookingRepo.NewTokenRepository(infraRegistry.DB(), bookingRepo.TokenConfig{LockedDuration: config.Token.LockedDuration})
	return &repositoryRegistry{
		eventRepository: bookingRepo.NewEventRepository(
			infraRegistry.DB(),
			bookingTokenRepo,
		),
		userRepository: authRepo.NewUserRepository(infraRegistry.DB()),
		bookingRepository: bookingRepo.NewBookingRepository(
			infraRegistry.DB(),
			infraRegistry.PaymentService(),
			bookingRepo.NewBookingItemRepository(infraRegistry.DB()),
			bookingTokenRepo,
			infraRegistry.AsyncTaskEnqueueClient(),
		),
		bookingEventTokenRepository: bookingTokenRepo,
		bookingItemRepository:       bookingRepo.NewBookingItemRepository(infraRegistry.DB()),
		bookingEmailRepository:      emailRepo.NewEmailClient(infraRegistry.EmailService()),
	}
}

func (r *repositoryRegistry) EventRepository() *bookingRepo.EventRepository {
	return r.eventRepository
}
func (r *repositoryRegistry) UserRepository() *authRepo.UserRepository {
	return r.userRepository
}

func (r *repositoryRegistry) BookingRepository() *bookingRepo.BookingRepository {
	return r.bookingRepository
}

func (r *repositoryRegistry) BookingEventTokenRepository() *bookingRepo.TokenRepository {
	return r.bookingEventTokenRepository
}

func (r *repositoryRegistry) BookingItemRepository() *bookingRepo.BookingItemRepository {
	return r.bookingItemRepository
}

func (r *repositoryRegistry) BookingEmailRepository() *emailRepo.EmailClient {
	return r.bookingEmailRepository
}
