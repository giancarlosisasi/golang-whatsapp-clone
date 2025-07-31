package main

import (
	"fmt"
	"golang-whatsapp-clone/auth"
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/graph"
	"golang-whatsapp-clone/graphql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang-whatsapp-clone/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	appConfig := config.SetupAppConfig()

	dbpool := database.SetupDatabase(appConfig.DatabaseURL)

	dbQueries := db.New(dbpool)

	jwtService := auth.NewJWTService(appConfig.JWTSecret)
	oauthService := auth.NewOAuthService(appConfig, jwtService)
	authHandlers := auth.NewAuthHandlers(appConfig, oauthService, jwtService, dbQueries)
	// authMW := auth.NewAuthMiddleware(jwtService, appConfig)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	if appConfig.AppEnv == "development" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:3000,http://localhost:4000",
			AllowCredentials: true,
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		}))
	}

	// auth routes
	auth := app.Group("/api/v1/auth")
	auth.Get("/google", authHandlers.GoogleLogin)
	auth.Get("/google/callback", authHandlers.GoogleCallback)
	auth.Get("/logout", authHandlers.Logout)

	if appConfig.AppEnv == "development" {
		h := playground.ApolloSandboxHandler("GraphQL playground", "/graphql")
		app.Get("/playground", adaptor.HTTPHandlerFunc(h.ServeHTTP))
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DBQueries: dbQueries,
		AppConfig: appConfig,
	}}))
	graphqlHandler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract user info from headers (set by auth middleware)
		userID := r.Header.Get("X-User-ID")
		userEmail := r.Header.Get("X-User-Email")

		// create graphql context with user info
		ctx := r.Context()
		if userID != "" && userEmail != "" {
			ctx = graphql.WithUserContext(ctx, userID, userEmail)
			r = r.WithContext(ctx)
		}

		srv.ServeHTTP(w, r)
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	// app.All("/graphql", authMW.RequireAuth(), graphqlHandler)
	app.All("/graphql", graphqlHandler)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	go func() {

		if err := app.Listen(":" + appConfig.Port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// cleanup tasks
	dbpool.Close()
	fmt.Println("Fiber was successful shutdown.")
}
