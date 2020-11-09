SHELL:=/bin/bash
TOP_DIR:=$(notdir $(CURDIR))
BUILD_DIR:=_build

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

lint:
	golangci-lint run ./...
