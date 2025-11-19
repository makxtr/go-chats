LOCAL_BIN:=$(CURDIR)/bin

.PHONY: help install-deps get-deps generate-chat-server-api run build docker-build docker-run

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-chat-server-api:
	mkdir -p pkg/chat_server_v1
	protoc --proto_path api/chat_server_v1 \
	--go_out=pkg/chat_server_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_server_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_server_v1/chat_server.proto

run:
	go run cmd/server/main.go

build:
	go build -o bin/chat-server cmd/server/main.go

docker-build:
	docker build -t chat-server-service .

docker-run:
	docker run -p 50052:50052 chat-server-service
