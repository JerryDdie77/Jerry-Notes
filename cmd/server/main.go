package main

import (
	"database/sql"
	"jerry-notes/internal/api"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Getting environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Can not Load .env variables", err)
	}

	connStr := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("sql Open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("db Ping: %v", err)
	}

	// Creating the route and register all handlers in them
	r := mux.NewRouter()
	api.RegisterRoutes(r)

	// Stating a server in a port :8080
	http.ListenAndServe(":8080", r)
}
