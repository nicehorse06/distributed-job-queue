APP_NAME=go-job-queue

.PHONY: run build test tidy

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy
