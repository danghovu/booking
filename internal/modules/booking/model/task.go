package model

type TaskType string

const (
	TaskTypeSendReminderEmail     TaskType = "send_reminder_email"
	TaskTypeSendConfirmationEmail TaskType = "send_confirmation_email"
)

type User struct {
	ID    int
	Email string
}

type SendReminderEmailTask struct {
	User    User    `json:"user"`
	Event   Event   `json:"event"`
	Booking Booking `json:"booking"`
}

type SendConfirmationEmailTask struct {
	User    User    `json:"user"`
	Event   Event   `json:"event"`
	Booking Booking `json:"booking"`
}
