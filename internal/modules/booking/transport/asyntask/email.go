package asyntask

import (
	"booking-event/internal/modules/booking/model"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type EmailService interface {
	SendReminderEmail(ctx context.Context, task model.SendReminderEmailTask) error
	SendConfirmationEmail(ctx context.Context, task model.SendConfirmationEmailTask) error
}

type EmailTaskHandler struct {
	emailService EmailService
}

func NewEmailTaskHandler(emailService EmailService) *EmailTaskHandler {
	return &EmailTaskHandler{emailService: emailService}
}

func (h *EmailTaskHandler) HandleReminderEmail(ctx context.Context, t *asynq.Task) error {
	var task model.SendReminderEmailTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return err
	}
	return h.emailService.SendReminderEmail(ctx, task)
}

func (h *EmailTaskHandler) HandleConfirmationEmail(ctx context.Context, t *asynq.Task) error {
	var task model.SendConfirmationEmailTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return err
	}
	return h.emailService.SendConfirmationEmail(ctx, task)
}

func (h *EmailTaskHandler) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(string(model.TaskTypeSendReminderEmail), h.HandleReminderEmail)
	mux.HandleFunc(string(model.TaskTypeSendConfirmationEmail), h.HandleConfirmationEmail)
}
