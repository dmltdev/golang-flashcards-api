.PHONY: db-up db-down db-logs db-reset

# Database commands
db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

db-logs:
	docker-compose logs -f postgres

db-reset:
	docker-compose down -v
	docker-compose up -d postgres

# Development commands
dev:
	go run cmd/api/main.go

run:
	./bin/api

migrate-up:
	go run cmd/migrations/main.go up

migrate-down:
	go run cmd/migrations/main.go down

# Build commands
build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy
