package model

import (
	"time"
)

type TokenStatus string

const (
	TokenStatusActive TokenStatus = "active"
	TokenStatusLocked TokenStatus = "locked"
	TokenStatusUsed   TokenStatus = "used"
)

type EventToken struct {
	ID          int
	EventID     int
	Token       string
	HolderID    *int
	LockedUntil *time.Time
	Status      TokenStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ConfirmingToken struct {
	Token    string
	EventID  int
	HolderID int32
}

type CreateEventTokensParams struct {
	UserID   int
	EventID  int
	Quantity int
}
