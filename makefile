.DEFAULT_GOAL := build

.PHONY: clean critic security lint test build run update

APP_NAME = finz
BUILD_DIR = $(PWD)/.build

update:
	go get -u ./...

clean:
	rm -rf ./.build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	mkdir -p $(BUILD_DIR)
	go test -v -timeout 30s -coverprofile=$(BUILD_DIR)/cover.out -cover ./...
	go tool cover -func=$(BUILD_DIR)/cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go
