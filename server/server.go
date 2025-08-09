package server

import (
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

func NewServer() (*App, *http.Server) {
	log := logger.NewLogger()
	appConfig := config.SetupAppConfig()

	dbpool := database.SetupDatabase(appConfig.DatabaseURL)

	dbQueries := db.New(dbpool)

	jwtService := auth.NewJWTService(appConfig.JWTSecret)
	oauthService := auth.NewOAuthService(appConfig, jwtService)

	handler := handler.NewHandler(log, appConfig, dbQueries, oauthService, jwtService)
	// authHandlers := auth.NewAuthHandlers(appConfig, oauthService, jwtService, dbQueries)
	// authMW := auth.NewAuthMiddleware(jwtService, appConfig)

	// app := fiber.New(fiber.Config{
	// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 		code := fiber.StatusInternalServerError
	// 		if e, ok := err.(*fiber.Error); ok {
	// 			code = e.Code
	// 		}

	// 		return c.Status(code).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	},
	// })

	// switch appConfig.AppEnv {
	// case "development":
	// 	app.Use(cors.New(cors.Config{
	// 		AllowOrigins:     "http://localhost:3000,http://localhost:4000",
	// 		AllowCredentials: true,
	// 		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	// 		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	// 	}))
	// case "production":
	// 	app.Use(cors.New(cors.Config{
	// 		AllowOrigins:     "https://studio.apollographql.com",
	// 		AllowCredentials: true,
	// 		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	// 		AllowMethods:     "GET, POST, OPTIONS",
	// 	}))
	// }

	// // authApiGroup routes
	// authApiGroup := app.Group("/api/v1/auth")
	// authApiGroup.Get("/google", authHandlers.GoogleLogin)
	// authApiGroup.Get("/google/callback", authHandlers.GoogleCallback)
	// authApiGroup.Get("/logout", authHandlers.Logout)

	// app.Get("/chats", func(c *fiber.Ctx) error {
	// 	return c.SendFile("./views/chats.html")
	// })

	// if appConfig.AppEnv == "development" {
	// 	h := playground.ApolloSandboxHandler("GraphQL playground", "/graphql")
	// 	app.Get("/playground", adaptor.HTTPHandlerFunc(h.ServeHTTP))
	// }

	// srv := gqlHandler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
	// 	DBQueries: dbQueries,
	// 	AppConfig: appConfig,
	// }}))
	// graphqlHandler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// extract user info from headers (set by auth middleware)
	// 	userID := r.Header.Get("X-User-ID")
	// 	userEmail := r.Header.Get("X-User-Email")

	// 	// create graphql context with user info
	// 	ctx := r.Context()
	// 	if userID != "" && userEmail != "" {
	// 		ctx = auth.WithUserContext(ctx, userID, userEmail)
	// 		r = r.WithContext(ctx)
	// 	}

	// 	srv.ServeHTTP(w, r)
	// })

	// srv.AddTransport(transport.Websocket{
	// 	// Keep-alives are important for WebSockets to detect dead connections. This is
	// 	// not unlike asking a partner who seems to have zoned out while you tell them
	// 	// a story crucial to understanding the dynamics of your workplace: "Are you
	// 	// listening to me?"
	// 	//
	// 	// Failing to set a keep-alive interval can result in the connection being held
	// 	// open and the server expending resources to communicate with a client that has
	// 	// long since walked to the kitchen to make a sandwich instead.
	// 	KeepAlivePingInterval: 10 * time.Second,

	// 	// The `github.com/gorilla/websocket.Upgrader` is used to handle the transition
	// 	// from an HTTP connection to a WebSocket connection. Among other options, here
	// 	// you must check the origin of the request to prevent cross-site request forgery
	// 	// attacks.
	// 	Upgrader: websocket.Upgrader{
	// 		CheckOrigin: func(r *http.Request) bool {
	// 			// Allow exact match on host.
	// 			origin := r.Header.Get("Origin")

	// 			fiberlog.Infof("origin: %s", origin)

	// 			if origin == "" || origin == r.Header.Get("Host") {
	// 				return true
	// 			}

	// 			// Match on allow-listed origins.
	// 			return slices.Contains(
	// 				[]string{
	// 					":5000",
	// 					":4000",
	// 					"https://sandbox.embed.apollographql.com",
	// 				}, origin)
	// 		},
	// 	},
	// })
	// srv.AddTransport(transport.Options{})
	// srv.AddTransport(transport.GET{})
	// srv.AddTransport(transport.POST{})

	// srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// srv.Use(extension.Introspection{})
	// srv.Use(extension.AutomaticPersistedQuery{
	// 	Cache: lru.New[string](100),
	// })
	// app.All(
	// 	"/graphql",
	// 	authMW.AuthenticateUser(),
	// 	graphqlHandler,
	// )

	// app.Get("/health", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"status": "ok",
	// 	})
	// })
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
		log.Debug().Msgf("user is: %+v", user)
		io.WriteString(w, "Chats page")
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

	log.Info().Msgf("> The server is running in http://localhost:%s", appConfig.Port)

	return app, server

}
