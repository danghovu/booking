package entity

import "booking-event/internal/modules/auth/model"

func ConvertUserToModel(user User) *model.User {
	return &model.User{
		ID:             user.ID,
		Email:          user.Email,
		Role:           model.UserRole(user.Role),
		HashedPassword: user.Password,
	}
}
