BINARY := issue-tracker
CMD     := ./cmd/server

.PHONY: build run dev test lint fmt sqlc migrate-up migrate-down generate tidy

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

generate:
	mkdir -p internal/api/generated
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		-generate types,strict-server,server,spec \
		-package generated \
		-o internal/api/generated/api.gen.go \
		api/openapi.yaml

tidy:
	go mod tidy
