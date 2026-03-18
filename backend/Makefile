APP_NAME := motivator-backend
BUILD_DIR := bin

.PHONY: install dev build test clean run swagger

install:
	go mod download

dev: swagger
	go run ./cmd/server

build: swagger
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

swagger:
	swag init -g cmd/server/main.go -o docs
