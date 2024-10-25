package entity

import "time"

type Booking struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	EventID         int       `db:"event_id"`
	Status          string    `db:"status"`
	InitialQuantity int       `db:"initial_quantity"`
	Quantity        int       `db:"quantity"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
