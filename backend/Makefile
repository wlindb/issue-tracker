BINARY := issue-tracker
CMD     := ./cmd/server

.PHONY: build run dev test lint fmt sqlc migrate-up migrate-down swagger tidy

build:
	go build -o bin/$(BINARY) $(CMD)

run:
	go run $(CMD)

dev:
	air

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

sqlc:
	sqlc generate

migrate-up:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

swagger:
	swag init -g $(CMD)/main.go -o docs

tidy:
	go mod tidy
