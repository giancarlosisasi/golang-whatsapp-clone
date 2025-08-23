package graph

import (
	"context"
	"golang-whatsapp-clone/auth"
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/graph/model"
	"golang-whatsapp-clone/service"

	"github.com/rs/zerolog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	DBQueries           *db.Queries
	AppConfig           *config.AppConfig
	Logger              *zerolog.Logger
	ConversationService *service.ConversationService
	MessageService      *service.MessageService
}

func (r *Resolver) mustGetAuthenticatedUser(ctx context.Context) (*auth.UserContext, *model.UnauthorizedError) {
	user := auth.GetUserFromContext(ctx)
	if user == nil {
		return nil, &model.UnauthorizedError{
			ErrorMessage: "You must be authenticated to access to this data",
			Code:         "UNAUTHORIZED",
		}
	}

	return user, nil
}
