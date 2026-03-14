BINARY := issue-tracker
CMD     := ./cmd/server

.PHONY: build run dev test test-integration lint fmt generate tidy docker-up docker-down

build: lint
	go build -o bin/$(BINARY) $(CMD)

run:
	go run $(CMD)

dev: docker-up
	set -a && . ./.env && air

test:
	go test ./...

test-integration:
	go test -tags integration ./internal/api/middleware/...

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

docker-down:
	docker compose down
