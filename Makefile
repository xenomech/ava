.PHONY: setup dev run build test lint fmt tidy

setup:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go mod download

dev:
	air

run:
	go run ./cmd/ava

build:
	go build -ldflags="-s -w" -o ./tmp/main ./cmd/ava

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofumpt -w .

tidy:
	go mod tidy
