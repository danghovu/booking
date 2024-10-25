package entity

import (
	"time"
)

type Event struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	AvailableSeats int       `db:"available_seats"`
	StartAt        time.Time `db:"start_at"`
	Location       string    `db:"location"`
	Category       string    `db:"category"`
	Price          int64     `db:"price"`
	Currency       string    `db:"currency"`
	Status         string    `db:"status"`
	CreatorID      int       `db:"creator_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
