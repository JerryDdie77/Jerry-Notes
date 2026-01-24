package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func main() {

	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatalf("usage: go run ./cmd/migrate/main.go [up|down]")
	}

	command := os.Args[1]

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL env variable is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	goose.SetDialect("postgres")

	switch command {
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatalf("failed to run up migrations: %v", err)
		}

		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatalf("failed to get db version: %v", err)
		}
		log.Printf("migrations applied, current version: %d", version)

	case "down":
		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatalf("failed to run down migration: %v", err)
		}

		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatalf("failed to get db version: %v", err)
		}

		if version == 0 {
			log.Println("all of migrations denied")
		} else {
			log.Printf("one migration reverted, current version: %d", version)
		}

	default:
		log.Fatalf("unknown command: %s (use up or down)", command)
	}
}
