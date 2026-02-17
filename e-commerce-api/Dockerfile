FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app ./server/cmd

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates curl

COPY --from=builder /app/app .

RUN chmod +x /app/app

EXPOSE 8080

CMD ["/app/app"]
