SERVICES := auth chat message media
REGISTRY := localhost:5000

up:
	bash ./scripts/up.sh

down:
	docker-compose -f deploy/docker-compose.yml down

migrate-up:
	bash ./scripts/migrate.sh up

migrate-down:
	bash ./scripts/migrate.sh down

init-kafka:
	bash ./scripts/init-kafka.sh

PROTO_ROOT := .
PROTO_FILES := \
	protos/common/v1/common.proto \
	protos/auth/v1/auth.proto \
	protos/chat/v1/chat.proto \
	protos/message/v1/message.proto \
	protos/presence/v1/presence.proto

SDK_DIR := sdk

generate-protos:
	protoc \
		--proto_path=$(PROTO_ROOT) \
		--go_out=$(SDK_DIR) --go_opt=module=sdk \
		--go-grpc_out=$(SDK_DIR) --go-grpc_opt=module=sdk \
		$(PROTO_FILES)
	cd $(SDK_DIR) && go mod tidy