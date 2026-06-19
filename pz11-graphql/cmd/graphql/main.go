package main

import (
	"log"
	"net/http"
	"os"

	"example.com/pz11-graphql/graph"
	"example.com/pz11-graphql/internal/auth"
	"example.com/pz11-graphql/internal/service"
	"example.com/pz11-graphql/internal/store"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	addr := env("SERVER_ADDR", ":8094")
	if addr[0] != ':' {
		addr = ":" + addr
	}

	st := store.New()
	svc := service.New(st)
	resolver := &graph.Resolver{TaskSvc: svc}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	mux := http.NewServeMux()
	mux.Handle("/", authMiddleware(playground.Handler("pz11-graphql", "/query")))
	mux.Handle("/query", authMiddleware(srv))

	log.Printf("GraphQL Playground: http://localhost%s/", addr)
	log.Printf("GraphQL endpoint:   http://localhost%s/query", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.WithAuthorization(r.Context(), r.Header.Get("Authorization"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
