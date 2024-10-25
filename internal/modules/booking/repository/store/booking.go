package store

import (
	"context"
	"errors"

	"github.com/Rhymond/go-money"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	bookingasynq "booking-event/internal/infra/asynq"
	"booking-event/internal/infra/paymentgateway"
	postgresql "booking-event/internal/infra/posgresql"
	"booking-event/internal/modules/booking/model"
	"booking-event/internal/modules/booking/repository/entity"
)

type PaymentClient interface {
	CreatePayment(ctx context.Context, payment *paymentgateway.Payment) error
}

type BookingItemRepositoryForBooking interface {
	CreateBookingItemsTx(ctx context.Context, tx postgresql.ExecerContext, bookingItems []model.BookingItem) error
	GetBookingItemsByBookingID(ctx context.Context, bookingID int) ([]model.BookingItem, error)
}

type EventTokenRepositoryForBooking interface {
	ReleaseTokensByTx(ctx context.Context, tx postgresql.ExecerContext, tokens []string) error
	ConfirmUsedTokensByTx(ctx context.Context, tx postgresql.ExecerContext, tokens []model.ConfirmingToken) error
}

type BookingRepository struct {
	db              *sqlx.DB
	paymentClient   PaymentClient
	bookingItemRepo BookingItemRepositoryForBooking
	tokenRepo       EventTokenRepositoryForBooking
	asynqClient     bookingasynq.AsyncTaskEnqueueClient
}

func NewBookingRepository(
	db *sqlx.DB,
	paymentClient PaymentClient,
	bookingItemRepo BookingItemRepositoryForBooking,
	tokenRepo EventTokenRepositoryForBooking,
	asynqClient bookingasynq.AsyncTaskEnqueueClient,
) *BookingRepository {
	return &BookingRepository{db: db, paymentClient: paymentClient, bookingItemRepo: bookingItemRepo, tokenRepo: tokenRepo, asynqClient: asynqClient}
}

const ()

func (c *BookingRepository) CreateBooking(ctx context.Context, booking *model.Booking, bookingItems []model.BookingItem) error {
	if len(bookingItems) == 0 {
		return errors.New("booking items are required")
	}
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	entityBooking := entity.ConvertBookingToEntity(booking)
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO bookings (user_id, event_id, status, initial_quantity, quantity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, entityBooking.UserID, entityBooking.EventID, entityBooking.Status, entityBooking.InitialQuantity, entityBooking.Quantity).Scan(&booking.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for i := range bookingItems {
		bookingItems[i].BookingID = booking.ID
	}

	err = c.bookingItemRepo.CreateBookingItemsTx(ctx, tx, bookingItems)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *BookingRepository) CountBookingByUserID(ctx context.Context, eventID int, userID int) (int, error) {
	var count int
	err := c.db.GetContext(ctx, &count, "SELECT COALESCE(SUM(quantity), 0) FROM bookings WHERE event_id = $1 AND user_id = $2 AND status = ANY($3)", eventID, userID, pq.Array([]model.BookingStatus{model.BookingStatusPending, model.BookingStatusConfirmed, model.BookingStatusPaid}))
	return count, err
}

func (c *BookingRepository) GetBookingByID(ctx context.Context, id int) (*model.Booking, error) {
	entityBooking := &entity.Booking{}
	err := c.db.QueryRowxContext(ctx, "SELECT id, user_id, event_id, status, initial_quantity, quantity FROM bookings WHERE id = $1", id).StructScan(entityBooking)
	if err != nil {
		return nil, err
	}

	return entity.ConvertBookingToModel(entityBooking), nil
}

func (c *BookingRepository) ConfirmBooking(ctx context.Context, booking *model.Booking, event *model.Event, bookingItems []model.BookingItem) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", string(model.BookingStatusConfirmed), booking.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	tokens := make([]model.ConfirmingToken, len(bookingItems))
	for i, item := range bookingItems {
		tokens[i] = model.ConfirmingToken{Token: item.Token, HolderID: int32(booking.UserID), EventID: event.ID}
	}

	err = c.tokenRepo.ConfirmUsedTokensByTx(ctx, tx, tokens)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	amount := money.New(event.Price, event.Currency)

	payment := &paymentgateway.Payment{
		Amount:   float64(amount.Amount()),
		Currency: amount.Currency().Code,
	}

	err = c.paymentClient.CreatePayment(ctx, payment)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// _ = c.asynqClient.Enqueue(ctx, asynq.NewTask(string(model.TaskTypeSendConfirmationEmail), []byte("{}")))
	// _ = c.asynqClient.Enqueue(ctx, asynq.NewTask(string(model.TaskTypeSendReminderEmail), []byte("{}")), asynq.ProcessAt(event.StartAt))

	return tx.Commit()
}

func (c *BookingRepository) CancelBooking(ctx context.Context, bookingID int) error {
	bookingItems, err := c.bookingItemRepo.GetBookingItemsByBookingID(ctx, bookingID)
	if err != nil {
		return err
	}

	if len(bookingItems) == 0 {
		return nil
	}

	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.NamedExecContext(ctx, "UPDATE bookings SET status = :status, updated_at = CURRENT_TIMESTAMP WHERE id = :id", map[string]interface{}{
		"id":     bookingID,
		"status": string(model.BookingStatusCanceled),
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	bookingItemsTokens := make([]string, len(bookingItems))
	for i, item := range bookingItems {
		bookingItemsTokens[i] = item.Token
	}

	err = c.tokenRepo.ReleaseTokensByTx(ctx, tx, bookingItemsTokens)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
