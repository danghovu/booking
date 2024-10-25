package emailsender

import "context"

type noopEmailService struct{}

func NewNoopEmailService() EmailService {
	return &noopEmailService{}
}

func (s *noopEmailService) SendEmail(ctx context.Context, email *Email) error {
	// do nothing
	return nil
}
