package store

import (
	"context"

	"github.com/jmoiron/sqlx"

	postgresql "booking-event/internal/infra/posgresql"
	"booking-event/internal/modules/booking/model"
	"booking-event/internal/modules/booking/repository/entity"
)

type BookingItemRepository struct {
	db *sqlx.DB
}

func NewBookingItemRepository(db *sqlx.DB) *BookingItemRepository {
	return &BookingItemRepository{db: db}
}

func (r *BookingItemRepository) CreateBookingItemsTx(ctx context.Context, tx postgresql.ExecerContext, bookingItems []model.BookingItem) error {
	entities := entity.ConvertBookingItemsToEntities(bookingItems)
	_, err := tx.NamedExecContext(ctx, "INSERT INTO booking_items (booking_id, token) VALUES (:booking_id, :token)", entities)
	return err
}

func (r *BookingItemRepository) CreateBookingItems(ctx context.Context, bookingItems []model.BookingItem) error {
	return r.CreateBookingItemsTx(ctx, r.db, bookingItems)
}

func (r *BookingItemRepository) GetBookingItemsByBookingID(ctx context.Context, bookingID int) ([]model.BookingItem, error) {
	var bookingItems []entity.BookingItem
	err := r.db.SelectContext(ctx, &bookingItems, "SELECT id, booking_id, token, created_at FROM booking_items WHERE booking_id = $1", bookingID)
	if err != nil {
		return nil, err
	}

	return entity.ConvertBookingItemsToModels(bookingItems), nil
}
