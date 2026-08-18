.PHONY: setup dev run build test lint fmt tidy \
	adapter-run adapter-build adapter-build-dev check-boundary

setup:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/evilmartians/lefthook@latest
	go work sync
	lefthook install

dev:
	cd api && air

run:
	cd api && go run ./cmd/ava

build:
	cd api && go build -ldflags="-s -w" -o ./tmp/main ./cmd/ava

test:
	cd api && go test ./...

lint:
	cd api && golangci-lint run

fmt:
	cd api && gofumpt -w .

tidy:
	cd api && go mod tidy
