FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/bin/ava ./cmd/ava

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S ava && adduser -S ava -G ava

WORKDIR /app

COPY --from=build /app/bin/ava /app/bin/ava

RUN chown -R ava:ava /app

USER ava

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8000/api/v1/health || exit 1

CMD ["/app/bin/ava"]
