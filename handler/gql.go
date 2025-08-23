package handler

import (
	"golang-whatsapp-clone/graph"
	"net/http"
	"time"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/vektah/gqlparser/v2/ast"
)

func NewGraphqlHandler(logger *zerolog.Logger, resolver *graph.Resolver) *gqlHandler.Server {
	srv := gqlHandler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Websocket{
		// Keep-alives are important for WebSockets to detect dead connections. This is
		// not unlike asking a partner who seems to have zoned out while you tell them
		// a story crucial to understanding the dynamics of your workplace: "Are you
		// listening to me?"
		//
		// Failing to set a keep-alive interval can result in the connection being held
		// open and the server expending resources to communicate with a client that has
		// long since walked to the kitchen to make a sandwich instead.
		KeepAlivePingInterval: 10 * time.Second,

		// The `github.com/gorilla/websocket.Upgrader` is used to handle the transition
		// from an HTTP connection to a WebSocket connection. Among other options, here
		// you must check the origin of the request to prevent cross-site request forgery
		// attacks.
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow exact match on host.
				origin := r.Header.Get("Origin")

				logger.Info().Msgf("origin: %s", origin)

				return true
				// if origin == "" || origin == r.Header.Get("Host") {
				// 	return true
				// }

				// // Match on allow-listed origins.
				// return slices.Contains(
				// 	[]string{
				// 		":5000",
				// 		":4000",
				// 		"https://sandbox.embed.apollographql.com",
				// 	}, origin)
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	srv.Use(extension.FixedComplexityLimit(30))

	return srv
}

func NewGraphqlPlaygroundHandler() http.HandlerFunc {
	playgroundHandler := playground.ApolloSandboxHandler("GraphQL playground", "/graphql")

	// playgroundHandler.ServeHTTP(w, r)
	return playgroundHandler
}
