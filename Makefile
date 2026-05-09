-include .env
export

run:
	go run ./cmd/api
test:
	go test ./...
tidy:
	go mod tidy
db-up:
	docker compose up -d
db-down:
	docker compose down
migrate-up:
	goose up
migrate-down:
	goose down