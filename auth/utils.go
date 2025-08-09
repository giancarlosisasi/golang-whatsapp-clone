package auth

import (
	"context"
)

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) *UserContext {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	if !ok {
		return nil
	}

	return user
}

// adds user context to GraphQL context
func WithUserContext(ctx context.Context, userID string, email string) context.Context {
	userCtx := &UserContext{
		UserID: userID,
		Email:  email,
	}

	return context.WithValue(ctx, UserContextKey, userCtx)
}
