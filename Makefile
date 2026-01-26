migrate_up:
	go run cmd/migrate/main.go up

migrate_down:
	go run cmd/migrate/main.go down

run:
	go run cmd/server/main.go

db_up:
	docker-compose up -d

db_down:
	docker-compose down

db_restart:
	docker-compose down
	docker-compose up -d