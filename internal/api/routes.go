package api

import (
	"jerry-notes/internal/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}
