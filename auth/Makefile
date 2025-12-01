include local.env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN=$(PG_DSN)

.PHONY: help install-deps get-deps generate-user-api run build docker-build docker-run


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto

run:
	go run cmd/server/main.go

build:
	go build -o bin/auth-server cmd/server/main.go

docker-build:
	docker build -t auth-service .

	docker run -p 50051:50051 auth-service

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" down -v

# Testing
test:
	go test -v -race -timeout 30s ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

test-api:
	go test -v -race ./internal/api/...

# Mock generation
generate-mocks:
	go generate ./internal/repository/...

# CI/CD commands
ci-test: generate-mocks
	@echo "Running tests for CI..."
	go test -race -timeout 30s ./...

ci-lint:
	@echo "Running linters..."
	golangci-lint run ./...

# Verify mocks are up to date
verify-mocks: generate-mocks
	@if ! git diff --exit-code internal/repository/mocks/; then \
		echo "❌ Mocks are outdated! Run 'make generate-mocks'"; \
		exit 1; \
	else \
		echo "✅ Mocks are up to date"; \
	fi
