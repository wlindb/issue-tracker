BINARY := issue-tracker
CMD     := ./cmd/server

.PHONY: build run dev test lint fmt sqlc migrate-up migrate-down generate tidy docker-up docker-down gen-jwt-key

build:
	go build -o bin/$(BINARY) $(CMD)

run:
	go run $(CMD)

dev: docker-up
	set -a && . ./.env && air

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

sqlc:
	sqlc generate

migrate-up: docker-up
	set -a && . ./.env && goose -dir db/migrations postgres "$$DATABASE_URL" up

migrate-down:
	set -a && . ./.env && goose -dir db/migrations postgres "$$DATABASE_URL" down

generate:
	mkdir -p internal/api/generated
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		-generate types,strict-server,server,spec \
		-package generated \
		-o internal/api/generated/api.gen.go \
		api/openapi.yaml

tidy:
	go mod tidy

docker-up:
	docker compose -f ../infrastructure/docker-compose.yml up -d

docker-down:
	docker compose -f ../infrastructure/docker-compose.yml down

gen-jwt-key:
	openssl genrsa 2048 | awk 'NF {sub(/\r/, ""); printf "%s\\n", $$0}' > jwt_private_key.txt
	@echo "Paste the contents of jwt_private_key.txt into .env as JWT_PRIVATE_KEY=\"...\""
