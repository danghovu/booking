package store

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"booking-event/internal/common/errors"
	postgresql "booking-event/internal/infra/posgresql"
	"booking-event/internal/modules/booking/model"
	"booking-event/internal/modules/booking/repository/entity"
)

type EventRepository struct {
	db        *sqlx.DB
	tokenRepo TokenRepositoryForEvent
}

type TokenRepositoryForEvent interface {
	CreateTokensTX(ctx context.Context, tx postgresql.ExecerContext, tokens []model.EventToken) error
}

func NewEventRepository(db *sqlx.DB, tokenRepo TokenRepositoryForEvent) *EventRepository {
	return &EventRepository{db: db, tokenRepo: tokenRepo}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event model.Event, tokens []model.EventToken) error {
	entityEvent := entity.ConvertEventToEntity(event)
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	err = tx.QueryRowxContext(ctx, "INSERT INTO events (name, available_seats, start_at, location, category, price, currency, creator_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		entityEvent.Name,
		entityEvent.AvailableSeats,
		entityEvent.StartAt,
		entityEvent.Location,
		entityEvent.Category,
		entityEvent.Price,
		entityEvent.Currency,
		entityEvent.CreatorID,
		entityEvent.Status).Scan(&entityEvent.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := range tokens {
		tokens[i].EventID = entityEvent.ID
	}
	err = r.tokenRepo.CreateTokensTX(ctx, tx, tokens)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *EventRepository) GetEventByID(ctx context.Context, id int) (*model.Event, error) {
	event := entity.Event{}
	err := r.db.GetContext(ctx, &event, "SELECT id, name, available_seats, start_at, location, category, price, currency, creator_id, status, created_at, updated_at FROM events WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return entity.ConvertEventToModel(event), nil
}

func (r *EventRepository) QueryEvents(ctx context.Context, query model.EventQuery) ([]model.Event, error) {
	queryString := `SELECT id, name, available_seats, start_at, location, category, price, currency, creator_id, status, created_at, updated_at FROM events WHERE 1=1`

	if query.ID != 0 {
		queryString += " AND id = :id"
	}
	if query.Name != "" {
		queryString += " AND name ILIKE :name"
	}
	if query.Location != "" {
		queryString += " AND location = :location"
	}
	if query.Category != "" {
		queryString += " AND category = :category"
	}
	if !query.StartFrom.IsZero() {
		queryString += " AND start_at >= :start_from"
	}
	if !query.StartTo.IsZero() {
		queryString += " AND start_at <= :start_to"
	}
	queryString += ` ORDER BY updated_at DESC`

	if query.Pagination.Limit > 0 {
		queryString += ` LIMIT :limit`
	}
	if query.Pagination.Page > 0 {
		queryString += ` OFFSET :offset`
	}

	events := []entity.Event{}
	rows, err := r.db.NamedQueryContext(ctx, queryString, map[string]interface{}{
		"id":         query.ID,
		"limit":      query.Pagination.GetLimit(),
		"offset":     query.Pagination.GetOffset(),
		"name":       "%" + query.Name + "%",
		"location":   query.Location,
		"category":   query.Category,
		"start_from": query.StartFrom,
		"start_to":   query.StartTo,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event entity.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return entity.ConvertEventsToModels(events), nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event model.Event) error {
	entityEvent := entity.ConvertEventToEntity(event)
	_, err := r.db.NamedExecContext(ctx, "UPDATE events SET status = :status, updated_at = CURRENT_TIMESTAMP WHERE id = :id", entityEvent)
	return err
}
