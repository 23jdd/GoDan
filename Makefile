.PHONY: build run test lint clean

APP_NAME = godan
BUILD_DIR = bin
CMD_DIR = cmd/server

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

run:
	go run $(CMD_DIR)/main.go

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR) logs/

tidy:
	go mod tidy

vet:
	go vet ./...
