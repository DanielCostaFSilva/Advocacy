.PHONY: up down migrate-up migrate-down test run

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down

run:
	go run ./cmd/api

test:
	go test ./... -v -count=1
