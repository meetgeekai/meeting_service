BINARY := main
BUILD_DIR := bin

.PHONY: all clean build run run_local deps test test_local db

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./

run: build
	$(BUILD_DIR)/$(BINARY)

deps:
	go mod tidy
	go mod download

db:
	sqlc generate

test:
	go test ./...

# Use this rule to test locally, it is necessary so all vars are visible at package init time
test_local:
	test -f .env || (echo ".env not found" && exit 1)
	bash -c 'set -a && source .env && set +a && go test ./...'

# Run the service locally with env vars loaded from .env
run_local:
	test -f .env || (echo ".env not found" && exit 1)
	bash -c 'set -a && source .env && set +a && go run .'

clean:
	rm -rf $(BUILD_DIR)
