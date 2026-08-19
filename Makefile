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
	cd backend/services/api && air

run:
	cd backend/services/api && go run ./cmd/ava

build:
	cd backend/services/api && go build -ldflags="-s -w" -o ./tmp/main ./cmd/ava

test:
	cd backend/services/api && go test ./...
	cd backend/services/hub && go test ./...

lint:
	cd backend/services/api && golangci-lint run
	cd backend/services/hub && golangci-lint run

fmt:
	cd backend/services/api && gofumpt -w .
	cd backend/services/hub && gofumpt -w .

tidy:
	cd backend/services/api && go mod tidy
	cd backend/services/hub && go mod tidy

hub-run:
	cd backend/services/hub && go run ./cmd/avahub

hub-build:
	cd backend/services/hub && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/avahub-linux-arm64 ./cmd/avahub

hub-build-dev:
	cd backend/services/hub && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/avahub-darwin-arm64 ./cmd/avahub

check-boundary:
	cd backend/services/hub && GOWORK=off go build ./...
	cd backend/services/hub && ! go list -deps ./... | grep -q '^ava/api' || (echo "hub must not import api" && exit 1)
	cd backend/services/api && ! go list -deps ./... | grep -q '^ava/hub' || (echo "api must not import hub" && exit 1)
