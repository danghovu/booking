package email

import (
	"booking-event/internal/infra/emailsender"
	"booking-event/internal/modules/booking/model"
	"context"
)

type EmailClient struct {
	emailsender.EmailService
}

func NewEmailClient(emailService emailsender.EmailService) *EmailClient {
	return &EmailClient{EmailService: emailService}
}

func (c *EmailClient) SendReminderEmail(ctx context.Context, task model.SendReminderEmailTask) error {
	email := emailsender.Email{
		To:      task.User.Email,
		From:    "noreply@booking-event.com",
		Subject: "Reminder",
		Body:    "Reminder",
	}
	return c.EmailService.SendEmail(ctx, &email)
}

func (c *EmailClient) SendConfirmationEmail(ctx context.Context, task model.SendConfirmationEmailTask) error {
	email := emailsender.Email{
		To:      task.User.Email,
		From:    "noreply@booking-event.com",
		Subject: "Confirmation",
		Body:    "Confirmation",
	}
	return c.EmailService.SendEmail(ctx, &email)
}
