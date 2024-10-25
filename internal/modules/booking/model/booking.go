package model

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusCanceled  BookingStatus = "canceled"
)

type Booking struct {
	ID              int           `json:"id"`
	EventID         int           `json:"event_id"`
	UserID          int           `json:"user_id"`
	Status          BookingStatus `json:"status"`
	InitialQuantity int           `json:"initial_quantity"`
	Quantity        int           `json:"quantity"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type BookingItem struct {
	ID        int
	BookingID int
	Token     string
	CreatedAt time.Time
}

type CreateBookingRequest struct {
	EventID  int `json:"event_id" binding:"required"`
	UserID   int
	Quantity int `json:"quantity" binding:"required"`
}

type ConfirmBookingRequest struct {
	BookingID int `uri:"booking_id" binding:"required"`
}

type CancelBookingRequest struct {
	BookingID int `uri:"booking_id" binding:"required"`
}

type GetBookingByIDRequest struct {
	BookingID int `uri:"booking_id" binding:"required"`
}
