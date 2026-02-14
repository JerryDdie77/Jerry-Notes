package main

import (
	"database/sql"
	"fmt"
	"jerry-notes/config"
	"jerry-notes/internal/app"
	"jerry-notes/internal/httpserver"
	"jerry-notes/internal/jwt"
	"jerry-notes/internal/service"
	"jerry-notes/internal/storage"
	"log"

	_ "github.com/lib/pq"
)

func initDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("sql Open: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db Ping: %w", err)
	}

	return db, nil
}

func initApp(cfg *config.Config, db *sql.DB) *app.App {
	pgNoteStorage := storage.NewPostgresNoteStorage(db)
	pgUserStorage := storage.NewPostgresUserStorage(db)

	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTAccessTokenTTL)

	emailService := service.NewEmailService(cfg.MailToken)
	noteService := service.NewNoteService(pgNoteStorage)
	userService := service.NewUserService(pgUserStorage)
	authService := service.NewAuthService(pgUserStorage, emailService, jwtManager, cfg.CodeTTL)

	return app.New(noteService, userService, authService, jwtManager)

}

func main() {
	// Getting cfg
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("loadConfig: %v", err)
	}

	// Connecting to DB
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("initDB %v", err)
	}
	application := initApp(cfg, db)

	h := httpserver.NewHandler(application)
	r := httpserver.NewRouter(h)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}

}
