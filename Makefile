.PHONY: setup dev run build test lint fmt tidy \
	hub-run hub-build hub-build-dev check-boundary

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
	cd hub && go test ./...

lint:
	cd api && golangci-lint run
	cd hub && golangci-lint run

fmt:
	cd api && gofumpt -w .
	cd hub && gofumpt -w .

tidy:
	cd api && go mod tidy
	cd hub && go mod tidy

hub-run:
	cd hub && go run ./cmd/avahub

hub-build:
	cd hub && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/avahub-linux-arm64 ./cmd/avahub

hub-build-dev:
	cd hub && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/avahub-darwin-arm64 ./cmd/avahub

check-boundary:
	cd hub && GOWORK=off go build ./...
	cd hub && ! go list -deps ./... | grep -q '^ava/api' || (echo "hub must not import api" && exit 1)
	cd api && ! go list -deps ./... | grep -q '^ava/hub' || (echo "api must not import hub" && exit 1)
