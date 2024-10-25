package services

import (
	"booking-event/internal/modules/booking/model"
	"context"
)

type EmailService struct {
	eventRepo   EventRepository
	bookingRepo BookingRepository
	emailClient BookingEmailRepository
}

type BookingEmailRepository interface {
	SendReminderEmail(ctx context.Context, task model.SendReminderEmailTask) error
	SendConfirmationEmail(ctx context.Context, task model.SendConfirmationEmailTask) error
}

func NewEmailService(bookingRepo BookingRepository, emailClient BookingEmailRepository) *EmailService {
	return &EmailService{bookingRepo: bookingRepo, emailClient: emailClient}
}

func (s *EmailService) SendReminderEmail(ctx context.Context, task model.SendReminderEmailTask) error {
	return s.emailClient.SendReminderEmail(ctx, task)
}

func (s *EmailService) SendConfirmationEmail(ctx context.Context, task model.SendConfirmationEmailTask) error {
	return s.emailClient.SendConfirmationEmail(ctx, task)
}
