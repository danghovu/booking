package emailsender

import "context"

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

type EmailService interface {
	SendEmail(ctx context.Context, email *Email) error
}
