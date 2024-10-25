package model

import (
	"booking-event/internal/common/model"
	"time"
)

type EventCategory string

const (
	EventCategoryMusic EventCategory = "music"
)

type EventStatus string

const (
	EventStatusActive   EventStatus = "active"
	EventStatusInactive EventStatus = "inactive"
)

type Event struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	AvailableSeats int           `json:"available_seats"`
	StartAt        time.Time     `json:"start_at"`
	Location       string        `json:"location"`
	Category       EventCategory `json:"category"`
	Price          int64         `json:"price"`
	Currency       string        `json:"currency"`
	Status         EventStatus   `json:"status"`
	CreatorID      int           `json:"creator_id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type EventQuery struct {
	ID         int              `json:"id"`
	Category   EventCategory    `json:"category"`
	Location   string           `json:"location"`
	StartFrom  time.Time        `json:"start_from"`
	StartTo    time.Time        `json:"start_to"`
	Name       string           `json:"name"`
	Pagination model.Pagination `json:"pagination" binding:"required"`
}

type RetrieveEventDetailRequest struct {
	EventID int `uri:"event_id"`
}

type CreateEventRequest struct {
	Name           string        `json:"name" binding:"required"`
	AvailableSeats int           `json:"available_seats" binding:"required"`
	StartAt        time.Time     `json:"start_at" binding:"required"`
	Location       string        `json:"location" binding:"required"`
	Category       EventCategory `json:"category" binding:"required"`
	Price          float64       `json:"price" binding:"required"`
	ExecutorID     int
}

type UpdateEventRequest struct {
	EventID    int
	Status     EventStatus `json:"status" binding:"required,oneof=active inactive"`
	ExecutorID int
}
