package main

import (
	"log"
	"net/http"
	"os"

	"example.com/pz12-rest-graphql/graph"
	"example.com/pz12-rest-graphql/internal/rest"
	"example.com/pz12-rest-graphql/internal/service"
	"example.com/pz12-rest-graphql/internal/store"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	addr := env("SERVER_ADDR", ":8095")

	st := store.New()
	svc := service.New(st)

	mux := http.NewServeMux()
	rest.New(svc).Register(mux)

	resolver := &graph.Resolver{TaskSvc: svc}
	gql := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	gql.AddTransport(transport.Options{})
	gql.AddTransport(transport.GET{})
	gql.AddTransport(transport.POST{})
	gql.Use(extension.Introspection{})

	mux.Handle("/graphql", gql)
	mux.Handle("/", playground.Handler("pz12 REST vs GraphQL", "/graphql"))

	log.Printf("REST:     http://localhost%s/v1/tasks", addr)
	log.Printf("GraphQL:  http://localhost%s/graphql", addr)
	log.Printf("Playground: http://localhost%s/", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
