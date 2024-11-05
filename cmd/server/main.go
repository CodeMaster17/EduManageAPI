package main

import (
	"log"
	"net/http"
	"student-api/config"
	"student-api/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    routes.RegisterRoutes(r)

    log.Printf("Starting server on %s...", config.ServerPort)
    http.ListenAndServe(config.ServerPort, r)
}
