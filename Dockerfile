FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o load-tester main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/load-tester /app/load-tester
ENTRYPOINT ["/app/load-tester"]