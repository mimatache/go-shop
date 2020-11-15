SHELL:=/bin/bash
TOP_DIR:=$(notdir $(CURDIR))
BUILD_DIR:=build
BIN_DIR:=$(BUILD_DIR)/_bin
PORT?=9090
DOCKER_REPO?="matache91mh"
SHOP_IMAGE?=$(DOCKER_REPO)/"go-shop"

ifeq ($(VERSION),)
	VERSION:=$(shell git describe --tags --dirty --always)
endif


all: install-go-tools lint run-test build

build: build-shop build-client
	
build-shop:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/shop ./cmd/shop

build-client:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/client ./cmd/client

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
	$(BIN_DIR)/shop -port $(PORT)

run-client: build-client
	$(BIN_DIR)/client -port $(PORT)

shop-image: build-shop
	cp $(BIN_DIR)/shop $(BUILD_DIR)/shop/ && \
	cp -r data $(BUILD_DIR)/shop/ && \
	docker build -t $(SHOP_IMAGE):$(VERSION) $(BUILD_DIR)/shop/ && \
	rm  $(BUILD_DIR)/shop/shop && \
	rm -rf $(BUILD_DIR)/shop/data

push-images: shop-image
	docker push $(SHOP_IMAGE):$(VERSION)