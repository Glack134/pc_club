BIN_DIR := bin
PROTO_DIR := pkg/rpc
PROTO_FILE := $(PROTO_DIR)/admin.proto
GO_SRC := $(shell find . -name '*.go')

# Добавляем пути к protoc-gen-go
PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)

.PHONY: all build clean proto test run-server run-client help

all: build

build: proto $(BIN_DIR)/server $(BIN_DIR)/client $(BIN_DIR)/gameadmin

proto: $(PROTO_DIR)/admin.pb.go $(PROTO_DIR)/admin_grpc.pb.go

$(PROTO_DIR)/%.pb.go: $(PROTO_FILE)
	protoc \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$<

$(BIN_DIR)/server: cmd/server/main.go $(GO_SRC)
	go build -o $@ ./cmd/server

$(BIN_DIR)/client: cmd/client/main.go $(GO_SRC)
	go build -o $@ ./cmd/client

$(BIN_DIR)/gameadmin: cmd/cli/main.go $(GO_SRC)
	go build -o $@ ./cmd/cli

run-server: $(BIN_DIR)/server
	./$(BIN_DIR)/server

run-client: $(BIN_DIR)/client
	./$(BIN_DIR)/client

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_DIR)/*.pb.go

help:
	@echo "Available targets:"
	@echo "  all       - Build all binaries"
	@echo "  build     - Build server, client and CLI"
	@echo "  proto     - Generate protobuf files"
	@echo "  test      - Run tests"
	@echo "  run-server- Run the server"
	@echo "  run-client- Run the client"
	@echo "  clean     - Clean build artifacts"
	@echo "  help      - Show this help"

.DEFAULT_GOAL := help