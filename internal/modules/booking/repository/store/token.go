package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	postgresql "booking-event/internal/infra/posgresql"
	"booking-event/internal/modules/booking/model"
	"booking-event/internal/modules/booking/repository/entity"
)

type TokenConfig struct {
	LockedDuration time.Duration
}

type TokenRepository struct {
	db     *sqlx.DB
	config TokenConfig
}

const (
	tokenKey = "event:token:%d"
)

func NewTokenRepository(
	db *sqlx.DB,
	config TokenConfig,
) *TokenRepository {
	return &TokenRepository{
		db:     db,
		config: config,
	}
}

func (r *TokenRepository) CreateTokens(ctx context.Context, tokens []model.EventToken) error {
	return r.CreateTokensTX(ctx, r.db, tokens)
}

func (r *TokenRepository) CreateTokensTX(ctx context.Context, tx postgresql.ExecerContext, tokens []model.EventToken) error {
	entities := entity.ConvertEventTokensToEntities(tokens)
	_, err := tx.NamedExecContext(ctx, "INSERT INTO event_tokens (event_id, token, status) VALUES (:event_id, :token, :status)", entities)
	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) GetByToken(ctx context.Context, token string) (*model.EventToken, error) {
	var eventToken entity.EventToken
	err := r.db.GetContext(ctx, &eventToken, "SELECT event_id, token, status, holder_id, locked_until, created_at, updated_at FROM event_tokens WHERE token = ?", token)
	return entity.ConvertEventTokenToModel(eventToken), err
}

func (r *TokenRepository) ReleaseToken(ctx context.Context, eventToken *model.EventToken) error {
	entityToken := entity.ConvertEventTokenToEntity(*eventToken)
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE event_tokens SET locked_until = NULL, status = :status, holder_id = NULL, updated_at = CURRENT_TIMESTAMP 
		WHERE token = :token`, map[string]interface{}{
		"token":  entityToken.Token,
		"status": string(model.TokenStatusActive),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) ReleaseTokensByTx(ctx context.Context, tx postgresql.ExecerContext, tokens []string) error {
	_, err := tx.NamedExecContext(ctx, `
		UPDATE event_tokens SET locked_until = NULL, status = :status, holder_id = NULL, updated_at = CURRENT_TIMESTAMP 
		WHERE token = ANY(:tokens)`, map[string]interface{}{
		"tokens": pq.Array(tokens),
		"status": string(model.TokenStatusActive),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) ConfirmUsedToken(ctx context.Context, token *model.ConfirmingToken) error {
	return r.ConfirmUsedTokenByTx(ctx, r.db, token)
}

func (r *TokenRepository) ConfirmUsedTokenByTx(ctx context.Context, tx postgresql.ExecerContext, token *model.ConfirmingToken) error {
	_, err := tx.NamedExecContext(ctx, "UPDATE event_tokens SET status = :status, holder_id = :holder_id, updated_at = CURRENT_TIMESTAMP WHERE token = :token", map[string]interface{}{
		"token":     token.Token,
		"status":    string(model.TokenStatusUsed),
		"holder_id": token.HolderID,
	})
	return err
}

func (r *TokenRepository) ConfirmUsedTokensByTx(ctx context.Context, tx postgresql.ExecerContext, tokens []model.ConfirmingToken) error {
	if len(tokens) == 0 {
		return nil
	}
	_, err := tx.NamedExecContext(ctx, "UPDATE event_tokens SET status = :status, holder_id = :holder_id, updated_at = CURRENT_TIMESTAMP WHERE token = ANY(:tokens)", map[string]interface{}{
		"tokens":    pq.Array(tokens),
		"status":    string(model.TokenStatusUsed),
		"holder_id": tokens[0].HolderID,
	})
	return err
}

func (r *TokenRepository) SelectAvailableToken(ctx context.Context, holderID int32, eventID int, quantity int) ([]string, error) {
	var tokens []string
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryxContext(ctx, "SELECT token FROM event_tokens WHERE event_id = $1 AND status = $2 AND (locked_until IS NULL OR (locked_until IS NOT NULL AND locked_until < CURRENT_TIMESTAMP)) LIMIT $3 FOR UPDATE SKIP LOCKED", eventID, string(model.TokenStatusActive), quantity)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var token string
		err := rows.Scan(&token)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tokens = append(tokens, token)
	}

	_, err = tx.NamedExecContext(ctx, "UPDATE event_tokens SET locked_until = :locked_until, holder_id = :holder_id, status = :new_status, updated_at = CURRENT_TIMESTAMP WHERE token = ANY(:tokens)", map[string]interface{}{
		"tokens":       pq.Array(tokens),
		"new_status":   string(model.TokenStatusLocked),
		"locked_until": sql.NullTime{Time: time.Now().Add(r.config.LockedDuration), Valid: true},
		"holder_id":    holderID,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
