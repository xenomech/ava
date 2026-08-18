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
	cd adapter && go test ./...

lint:
	cd api && golangci-lint run
	cd adapter && golangci-lint run

fmt:
	cd api && gofumpt -w .
	cd adapter && gofumpt -w .

tidy:
	cd api && go mod tidy
	cd adapter && go mod tidy

adapter-run:
	cd adapter && go run ./cmd/adapter

adapter-build:
	cd adapter && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/adapter-linux-arm64 ./cmd/adapter

adapter-build-dev:
	cd adapter && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/adapter-darwin-arm64 ./cmd/adapter

check-boundary:
	cd adapter && GOWORK=off go build ./...
	cd adapter && ! go list -deps ./... | grep -q '^ava/api' || (echo "adapter must not import api" && exit 1)
	cd api && ! go list -deps ./... | grep -q '^ava/adapter' || (echo "api must not import adapter" && exit 1)
