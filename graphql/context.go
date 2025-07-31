package graphql

import (
	"context"
	"golang-whatsapp-clone/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) *auth.UserContext {
	user, ok := ctx.Value(UserContextKey).(*auth.UserContext)
	if !ok {
		return nil
	}

	return user
}

// adds user context to GraphQL context
func WithUserContext(ctx context.Context, userID string, email string) context.Context {
	userCtx := &auth.UserContext{
		UserID: userID,
		Email:  email,
	}

	return context.WithValue(ctx, UserContextKey, userCtx)
}
