-include .env
export

run:
	go run ./cmd/api
test:
	go test ./...
tidy:
	go mod tidy