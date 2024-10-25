package store

import (
	"booking-event/internal/modules/auth/model"
	"booking-event/internal/modules/auth/repository/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db sqlx.ExtContext
}

func NewUserRepository(db sqlx.ExtContext) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user entity.User
	if err := r.db.QueryRowxContext(ctx, "SELECT id, email, password, role FROM users WHERE email = $1", email).StructScan(&user); err != nil {
		return nil, err
	}

	return entity.ConvertUserToModel(user), nil
}
