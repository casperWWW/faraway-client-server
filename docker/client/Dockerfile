# Stage 1: Build the client binary.
FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o client ./cmd/client/main.go

# Stage 2: Create a minimal image.
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/client /app/client

ARG SERVER_ADDR=server:8080
ENV SERVER_ADDR=${SERVER_ADDR}

CMD ["/bin/sh", "-c", "./client -server ${SERVER_ADDR}"]
