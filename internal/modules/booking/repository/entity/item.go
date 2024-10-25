package entity

import "time"

type BookingItem struct {
	ID        int       `db:"id"`
	BookingID int       `db:"booking_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}
