BINARY := issue-tracker
CMD     := ./cmd/server

.PHONY: build run dev test test-integration lint fmt generate tidy sqlc-tracker docker-up docker-down wait-keycloak

build: lint
	go build -o bin/$(BINARY) $(CMD)

run:
	go run $(CMD)

dev: docker-up wait-keycloak
	set -a && . ./.env && air

test:
	go test ./...

test-integration:
	go test -tags integration -v -json -coverprofile=coverage-integration.out ./... > test-integration-results.json 2>&1

sqlc-tracker:
	cd internal/infrastructure/tracker && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

generate:
	mkdir -p internal/api/model
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		-config api/oapi-codegen.yaml \
		api/openapi.yaml

tidy:
	go mod tidy

docker-up:
	docker compose up -d

wait-keycloak:
	@printf 'Waiting for Keycloak JWKS'
	@until curl -sf http://localhost:8180/realms/issue-tracker/protocol/openid-connect/certs > /dev/null 2>&1; do \
		printf '.'; sleep 2; \
	done
	@echo ' ready.'

docker-down:
	docker compose down
