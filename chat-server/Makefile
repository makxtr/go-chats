include local.env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN=$(PG_DSN)

.PHONY: help install-deps get-deps generate-chat-server-api run build docker-build docker-run test test-coverage test-api generate-mocks ci-test ci-lint verify-mocks verify-proto

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-chat-server-api:
	protoc --proto_path api/chat_server_v1 --proto_path vendor.protogen \
	--go_out=pkg/chat_server_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_server_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/chat_server_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--validate_out=lang=go:pkg/chat_server_v1 --validate_opt=paths=source_relative \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
    --plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/chat_server_v1/chat_server.proto

vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
    	mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
    	git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
    	mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
    	rm -rf vendor.protogen/openapiv2 ;\
    fi

run:
	go run cmd/server/main.go

build:
	go build -o bin/chat-server cmd/server/main.go

docker-build:
	docker build -t chat-server-service .

docker-run:
	docker run -p 50052:50052 chat-server-service

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

verify-mocks: generate-mocks
	@if ! git diff --exit-code internal/repository/mocks/; then \
		echo "Mocks are outdated! Run 'make generate-mocks'"; \
		exit 1; \
	else \
		echo "Mocks are up to date"; \
	fi

verify-proto: generate-chat-server-api
	@if ! git diff --exit-code pkg/chat_server_v1/; then \
		echo "Proto files are outdated! Run 'make generate-chat-server-api' and commit changes"; \
		exit 1; \
	else \
		echo "Proto files are up to date"; \
	fi

generate:
	mkdir -p pkg/chat_server_v1
	mkdir -p pkg/swagger
	make generate-chat-server-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'