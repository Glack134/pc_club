.PHONY: all build-server build-client proto lint test clean

all: build-server build-client

build-server:
		go build -o bin/server ./cmd/server

build-client:
		go build -o bin/client ./cmd/client

proto:
	protoc --go_out=. --go-grpc_out=. pkg/rpc/**/*.proto

lint:
	golangci-lint run ./...

test:
	go test -v ./...

clean:
	rm -rf bin/*