package entity

import (
	"booking-event/internal/common/util"
	"booking-event/internal/modules/booking/model"
	"database/sql"
)

func ConvertBookingToModel(booking *Booking) *model.Booking {
	return &model.Booking{
		ID:              booking.ID,
		UserID:          booking.UserID,
		EventID:         booking.EventID,
		Status:          model.BookingStatus(booking.Status),
		Quantity:        booking.Quantity,
		InitialQuantity: booking.InitialQuantity,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
	}
}

func ConvertBookingToEntity(booking *model.Booking) *Booking {
	return &Booking{
		ID:              booking.ID,
		UserID:          booking.UserID,
		EventID:         booking.EventID,
		Status:          string(booking.Status),
		Quantity:        booking.Quantity,
		InitialQuantity: booking.InitialQuantity,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
	}
}

func ConvertEventTokenToEntity(token model.EventToken) *EventToken {
	out := &EventToken{
		ID:        token.ID,
		EventID:   token.EventID,
		Token:     token.Token,
		Status:    string(token.Status),
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
	if token.HolderID != nil {
		out.HolderID = sql.NullInt64{Int64: int64(*token.HolderID), Valid: true}
	}
	if token.LockedUntil != nil {
		out.LockedUntil = sql.NullTime{Time: *token.LockedUntil, Valid: true}
	}
	return out
}

func ConvertEventTokensToEntities(tokens []model.EventToken) []EventToken {
	entityTokens := make([]EventToken, len(tokens))
	for i, token := range tokens {
		entityTokens[i] = *ConvertEventTokenToEntity(token)
	}
	return entityTokens
}

func ConvertEventTokenToModel(token EventToken) *model.EventToken {
	out := &model.EventToken{
		ID:        token.ID,
		EventID:   token.EventID,
		Token:     token.Token,
		Status:    model.TokenStatus(token.Status),
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
	if token.HolderID.Valid {
		out.HolderID = util.ToPtr(int(token.HolderID.Int64))
	}
	if token.LockedUntil.Valid {
		out.LockedUntil = &token.LockedUntil.Time
	}
	return out
}

func ConvertEventTokensToModels(tokens []EventToken) []*model.EventToken {
	models := make([]*model.EventToken, len(tokens))
	for i, token := range tokens {
		models[i] = ConvertEventTokenToModel(token)
	}
	return models
}

func ConvertBookingItemsToModels(bookingItems []BookingItem) []model.BookingItem {
	models := make([]model.BookingItem, len(bookingItems))
	for i, item := range bookingItems {
		models[i] = model.BookingItem{ID: item.ID, BookingID: item.BookingID, Token: item.Token, CreatedAt: item.CreatedAt}
	}
	return models
}

func ConvertBookingItemsToEntities(bookingItems []model.BookingItem) []BookingItem {
	entities := make([]BookingItem, len(bookingItems))
	for i, item := range bookingItems {
		entities[i] = BookingItem{ID: item.ID, BookingID: item.BookingID, Token: item.Token, CreatedAt: item.CreatedAt}
	}
	return entities
}

func ConvertEventToEntity(event model.Event) *Event {
	return &Event{
		ID:             event.ID,
		Name:           event.Name,
		AvailableSeats: event.AvailableSeats,
		StartAt:        event.StartAt,
		Location:       event.Location,
		Category:       string(event.Category),
		Price:          event.Price,
		CreatorID:      event.CreatorID,
		Currency:       event.Currency,
		Status:         string(event.Status),
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func ConvertEventsToEntities(events []*model.Event) []Event {
	entities := make([]Event, len(events))
	for i, event := range events {
		entities[i] = *ConvertEventToEntity(*event)
	}
	return entities
}

func ConvertEventToModel(event Event) *model.Event {
	return &model.Event{
		ID:             event.ID,
		Name:           event.Name,
		AvailableSeats: event.AvailableSeats,
		StartAt:        event.StartAt,
		Location:       event.Location,
		Category:       model.EventCategory(event.Category),
		Price:          event.Price,
		Currency:       event.Currency,
		Status:         model.EventStatus(event.Status),
		CreatorID:      event.CreatorID,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func ConvertEventsToModels(events []Event) []model.Event {
	models := make([]model.Event, len(events))
	for i, event := range events {
		models[i] = *ConvertEventToModel(event)
	}
	return models
}
