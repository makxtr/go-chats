.PHONY: help run-auth run-chat grpcui-auth grpcui-chat grpcui-auth-prod grpcui-chat-prod infra-up infra-down infra-logs

run-auth:
	@echo "Starting auth service..."
	cd auth && go run cmd/server/main.go -config-path=local.env

run-chat:
	@echo "Starting chat-server..."
	cd chat-server && go run cmd/server/main.go -config-path=local.env

grpcui-auth:
	@cd auth && export $$(grep -v '^#' local.env | xargs) && grpcui -plaintext $$GRPC_HOST:$$GRPC_PORT

grpcui-chat:
	@cd chat-server && export $$(grep -v '^#' local.env | xargs) && grpcui -plaintext $$GRPC_HOST:$$GRPC_PORT

grpcui-auth-prod:
	grpcui -plaintext auth-service-rxpqkfxb3a-uc.a.run.app:443

grpcui-chat-prod:
	grpcui -plaintext chat-service-rxpqkfxb3a-uc.a.run.app:443

infra-up:
	docker-compose up -d

infra-down:
	docker-compose down

infra-logs:
	docker-compose logs -f
