package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func Setup() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)    // logging
	r.Use(middleware.Recoverer) // recover from panics

	// Load API routes
	Routes(r)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
