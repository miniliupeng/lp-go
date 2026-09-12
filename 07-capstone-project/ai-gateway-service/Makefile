.PHONY: all build run test race bench clean

APP_NAME = ai-gateway
BIN_DIR = bin

all: test build

build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) ./cmd/gateway/main.go
	@echo "Build successful: $(BIN_DIR)/$(APP_NAME)"

run:
	go run ./cmd/gateway/main.go

test:
	go test -v ./...

race:
	go test -v -race ./...

bench:
	go test -v -bench=. -benchmem ./...

clean:
	rm -rf $(BIN_DIR)
