FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --chown=65534:65534 --from=builder /app/app ./app
USER 65534

ENTRYPOINT ["./app"]