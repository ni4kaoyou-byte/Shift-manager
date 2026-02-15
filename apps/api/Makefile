APP_NAME := shift-manager-api
MAIN_PKG := ./cmd/server
GO_BIN := $(shell go env GOPATH)/bin
GOIMPORTS := $(shell command -v goimports 2>/dev/null || echo $(GO_BIN)/goimports)
GOLANGCI_LINT := $(shell command -v golangci-lint 2>/dev/null || echo $(GO_BIN)/golangci-lint)
SQLC := $(shell command -v sqlc 2>/dev/null || echo $(GO_BIN)/sqlc)
PSQL := $(shell command -v psql 2>/dev/null || echo psql)
DB_MIGRATIONS_DIR := ../../db/migrations

.PHONY: run build test test-ci lint fmt tidy tools vet db-up db-down sqlc-generate

run:
	go run $(MAIN_PKG)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PKG)

# テストはゆるゆる（速度優先）
test:
	go test ./...

# CI向け：再現性と事故検出を強める（必要なら）
test-ci:
	go test -race -count=1 ./...

vet:
	go vet ./...

# lintは golangci-lint を本丸にする（vetも内包されることが多い）
lint:
	$(GOLANGCI_LINT) run ./...

# fmtは goimports 推し（importsも直る）
fmt:
	$(GOIMPORTS) -w .

tidy:
	go mod tidy

# ツールの導入（初回だけ）
tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

db-up:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required" && exit 1)
	@for f in $$(find $(DB_MIGRATIONS_DIR) -name '*.up.sql' | sort); do \
		echo "Applying $$f"; \
		$(PSQL) "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f "$$f"; \
	done

db-down:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required" && exit 1)
	@for f in $$(find $(DB_MIGRATIONS_DIR) -name '*.down.sql' | sort -r); do \
		echo "Applying $$f"; \
		$(PSQL) "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f "$$f"; \
	done

sqlc-generate:
	$(SQLC) generate -f sqlc.yaml
