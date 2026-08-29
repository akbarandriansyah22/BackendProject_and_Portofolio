FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o app \
    ./server/cmd

FROM alpine:3.20.3

WORKDIR /app

RUN apk --no-cache add ca-certificates curl && \
    addgroup -S appgroup && \
    adduser -S -G appgroup -H -D appuser

COPY --from=builder /app/app .

RUN chmod 550 /app/app && \
    chown appuser:appgroup /app/app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD curl -f http://localhost:8080/live || exit 1

CMD ["/app/app"]