APP_NAME := shift-manager-api
MAIN_PKG := ./cmd/server

.PHONY: run build test lint fmt tidy

run:
	go run $(MAIN_PKG)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PKG)

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -type f)

tidy:
	go mod tidy
