package entity

import (
	"database/sql"
	"time"
)

type EventToken struct {
	ID          int           `db:"id"`
	EventID     int           `db:"event_id"`
	Token       string        `db:"token"`
	Status      string        `db:"status"`
	HolderID    sql.NullInt64 `db:"holder_id"`
	LockedUntil sql.NullTime  `db:"locked_until"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}
