SHELL:=/bin/bash
TOP_DIR:=$(notdir $(CURDIR))
BUILD_DIR:=_build
SERVER_PORT?=9090

all: install-go-tools lint run-test build

build: build-shop build-client
	
build-shop:
	go build -o $(BUILD_DIR)/shop ./cmd/shop

build-client:
	go build -o $(BUILD_DIR)/client ./cmd/client

run-tests:
	go test -v ./...

install-go-tools:
	@./scripts/install_tools.sh
	go install github.com/golang/mock/mockgen

lint:
	golangci-lint run ./...

generate:
	go generate -v ./...

run-server: build-shop
	$(BUILD_DIR)/shop -port $(SERVER_PORT)

run-client: build-client
	$(BUILD_DIR)/client -port $(SERVER_PORT)
