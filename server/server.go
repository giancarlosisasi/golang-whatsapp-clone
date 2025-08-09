package server

import (
	"errors"
	"fmt"
	"golang-whatsapp-clone/auth"
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/handler"
	"golang-whatsapp-clone/logger"
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

	jwtService := auth.NewJWTService(appConfig.JWTSecret)
	oauthService := auth.NewOAuthService(appConfig, jwtService)

	handler := handler.NewHandler(log, appConfig, dbQueries, oauthService, jwtService)

	app := &App{
		AppConfig: appConfig,
		DBpool:    dbpool,
		Logger:    log,
		Handler:   handler,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/api/v1/health", handler.HealthCheckHandler)
	mux.HandleFunc("/api/v1/auth/google", handler.GoogleLoginHandler)
	mux.HandleFunc("/api/v1/auth/google/callback", handler.GoogleCallbackHandler)
	mux.HandleFunc("/api/v1/auth/logout", handler.LogoutHandler)
	mux.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r.Context())
		log.Info().Msgf("user is: %+v", user)
		_, err := io.WriteString(w, "Chats page")
		if err != nil {
			handler.ServerErrorResponse(w, r, errors.New("internal server error"))
		}
	})
	mux.HandleFunc("/graphql", handler.GraphqlHandler)
	mux.HandleFunc("/playground", handler.GraphqlPlaygroundHandler)

	rootHandler := handler.AuthenticateUserMiddleware(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", appConfig.Port),
		Handler:      rootHandler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return app, server, rootHandler

}
