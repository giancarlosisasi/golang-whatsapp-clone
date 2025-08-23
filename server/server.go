package server

import (
	"errors"
	"fmt"
	"golang-whatsapp-clone/auth"
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/graph"
	"golang-whatsapp-clone/handler"
	"golang-whatsapp-clone/logger"
	"golang-whatsapp-clone/repository"
	"golang-whatsapp-clone/service"
	"io"
	"net/http"
	"time"

	"golang-whatsapp-clone/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type App struct {
	AppConfig *config.AppConfig
	DBpool    *pgxpool.Pool
	Logger    *zerolog.Logger
	Handler   *handler.Handler
}

func NewServer() (*App, *http.Server, http.Handler) {
	log := logger.NewLogger()
	appConfig := config.SetupAppConfig()

	dbpool := database.SetupDatabase(appConfig.DatabaseURL)

	dbQueries := db.New(dbpool)

	// repositories
	conversationRepository := repository.NewConversationRepository(dbQueries, log)
	participantRepository := repository.NewParticipantRepository(dbQueries, log)
	messageRepository := repository.NewMessageRepository(dbQueries)

	// services
	jwtService := auth.NewJWTService(appConfig.JWTSecret)
	oauthService := auth.NewOAuthService(appConfig, jwtService)
	conversationService := service.NewConversationService(conversationRepository, participantRepository)
	messageService := service.NewMessageService(messageRepository)

	handlers := handler.NewHandler(
		log,
		appConfig,
		dbQueries,
		oauthService,
		jwtService,
	)

	app := &App{
		AppConfig: appConfig,
		DBpool:    dbpool,
		Logger:    log,
		Handler:   handlers,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/api/v1/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/api/v1/auth/google", handlers.GoogleLoginHandler)
	mux.HandleFunc("/api/v1/auth/google/callback", handlers.GoogleCallbackHandler)
	mux.HandleFunc("/api/v1/auth/logout", handlers.LogoutHandler)
	mux.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		log.Info().Msgf("user is: %+v", user)
		_, err := io.WriteString(w, "Chats page")
		if err != nil {
			handlers.ServerErrorResponse(w, r, errors.New("internal server error"))
		}
	})

	// ------------------------- GQLGen setup handlers ---------------------- //
	gqlResolver := &graph.Resolver{
		DBQueries:           dbQueries,
		AppConfig:           appConfig,
		Logger:              log,
		ConversationService: conversationService,
		MessageService:      messageService,
	}
	graphqlHandler := handler.NewGraphqlHandler(log, gqlResolver)
	graphqlPlaygroundHandler := handler.NewGraphqlPlaygroundHandler()

	mux.HandleFunc("/graphql", graphqlHandler.ServeHTTP)
	mux.HandleFunc("/playground", graphqlPlaygroundHandler.ServeHTTP)

	// ----------------------- Middlewares ----------------------- //
	rootHandler := handlers.AuthenticateUserMiddleware(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", appConfig.Port),
		Handler:      rootHandler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return app, server, rootHandler

}
