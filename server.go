package main

import (
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/alecthomas/log4go"

	"github.com/ayusvcaus/business/auth"
	"github.com/ayusvcaus/business/graph"
	"github.com/ayusvcaus/business/graph/generated"

	"github.com/ayusvcaus/business/logging"
	"github.com/ayusvcaus/business/persist"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func init() {
	logging.InitConfigs()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.Middleware())
	status := persist.InitDB()
	if !status {
		log4go.Error("Connection to Mysql failed")
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	}
	server := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log4go.Info("connect to http://localhost:%s/ for GraphQL playground", port)
	log4go.Info(http.ListenAndServe(":"+port, router))
}
