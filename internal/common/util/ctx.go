package util

import (
	"context"
)

type UserContextKey struct{}

func GetUserIDContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserContextKey{}).(int)
	if !ok {
		return 0
	}
	return userID
}

func SetUserIDContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserContextKey{}, userID)
}
