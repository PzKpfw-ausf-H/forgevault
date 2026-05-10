-include .env
export

export PATH := $(shell go env GOPATH)/bin:$(PATH)

run:
	go run ./cmd/api
test:
	go test ./...
tidy:
	go mod tidy
db-up:
	docker compose -f deployments/docker/docker-compose.yaml up -d
db-down:
	docker compose -f deployments/docker/docker-compose.yaml down
migrate-up:
	goose up
migrate-down:
	goose down