package appcontext

import (
	"github.com/google/uuid"

	"booking-event/config"
	authServices "booking-event/internal/modules/auth/services"
	bookingServices "booking-event/internal/modules/booking/services"
)

type ServiceRegistry interface {
	EventService() *bookingServices.EventService
	AuthService() *authServices.AuthService
	BookingService() *bookingServices.BookingService
	BookingEventTokenService() *bookingServices.EventTokenService
	EmailService() *bookingServices.EmailService
}

type serviceRegistry struct {
	eventService      *bookingServices.EventService
	authService       *authServices.AuthService
	bookingService    *bookingServices.BookingService
	eventTokenService *bookingServices.EventTokenService
	emailService      *bookingServices.EmailService
}

func NewServiceRegistry(
	config config.Config,
	infraRegistry InfraRegistry,
	repositoryRegistry RepositoryRegistry,
) ServiceRegistry {
	bookingEventTokenService := bookingServices.NewEventTokenService(
		repositoryRegistry.BookingEventTokenRepository(),
		func() string {
			return uuid.New().String()
		},
	)
	return &serviceRegistry{
		eventService: bookingServices.NewEventService(
			repositoryRegistry.EventRepository(),
			config.SupportingMoney.Currency,
			func() string {
				return uuid.New().String()
			},
		),
		authService: authServices.NewAuthService(
			repositoryRegistry.UserRepository(),
			authServices.AuthServiceConfig{
				SecretKey:       config.JWT.SecretKey,
				AccessTokenExp:  config.JWT.AccessTokenExp,
				RefreshTokenExp: config.JWT.RefreshTokenExp,
			},
		),
		bookingService: bookingServices.NewBookingService(
			repositoryRegistry.EventRepository(),
			repositoryRegistry.BookingEventTokenRepository(),
			repositoryRegistry.BookingRepository(),
			repositoryRegistry.BookingItemRepository(),
			infraRegistry.PaymentService(),
			bookingServices.BookingConfig{
				MaxBookingPerUser: config.Booking.MaxBookingPerUser,
			},
		),
		eventTokenService: bookingEventTokenService,
		emailService: bookingServices.NewEmailService(
			repositoryRegistry.BookingRepository(),
			repositoryRegistry.BookingEmailRepository(),
		),
	}
}

func (s *serviceRegistry) EventService() *bookingServices.EventService {
	return s.eventService
}

func (s *serviceRegistry) AuthService() *authServices.AuthService {
	return s.authService
}

func (s *serviceRegistry) BookingService() *bookingServices.BookingService {
	return s.bookingService
}

func (s *serviceRegistry) BookingEventTokenService() *bookingServices.EventTokenService {
	return s.eventTokenService
}

func (s *serviceRegistry) EmailService() *bookingServices.EmailService {
	return s.emailService
}
