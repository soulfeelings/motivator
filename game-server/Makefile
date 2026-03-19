APP_NAME := motivator-game
BUILD_DIR := bin

.PHONY: install dev build test clean run

install:
	go mod download

dev:
	go run ./cmd/server

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(APP_NAME)
