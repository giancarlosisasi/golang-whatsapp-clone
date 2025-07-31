package main

import (
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"golang-whatsapp-clone/graph"
	"log"
	"net/http"

	"golang-whatsapp-clone/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	appConfig := config.SetupAppConfig()

	dbpool := database.SetupDatabase(appConfig.DatabaseURL)
	defer dbpool.Close()

	dbQueries := db.New(dbpool)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DBQueries: dbQueries,
		AppConfig: appConfig,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", appConfig.Port)
	log.Fatal(http.ListenAndServe(":"+appConfig.Port, nil))
}
