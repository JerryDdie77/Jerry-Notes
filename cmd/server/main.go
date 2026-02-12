package main

import (
	"database/sql"
	"jerry-notes/configs"
	"jerry-notes/internal/jwt"
	"jerry-notes/internal/service"
	"jerry-notes/internal/storage"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Getting cfg variable with secret data
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("loadConfig: %v", err)
	}

	// Connecting to DB
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("sql Open: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("db Ping: %v", err)
	}

	// Creating Storages
	postgresNoteStorage := storage.NewPostgresNoteStorage(db)
	postgresUserStorage := storage.NewPostgresUserStorage(db)

	// Creating JWT manager
	jwtManager := jwt.NewManager(cfg.JWTSecret, 20*time.Minute)

	// Creating services
	emailService := service.NewEmailService(cfg.GmailToken)
	authService := service.NewAuthService(postgresUserStorage, *emailService, *jwtManager, 5*time.Minute)
	userService := service.NewUserService(postgresUserStorage)
	noteService := service.NewNoteService(postgresNoteStorage)

}
