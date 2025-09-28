# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o caching-proxy .

# Final image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/caching-proxy .

EXPOSE 3000
ENTRYPOINT ["./caching-proxy"]
