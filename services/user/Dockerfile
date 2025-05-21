FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go mod vendor

RUN go test -v ./...
RUN go build -o ./bin/user ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/user ./bin/user

CMD ["./bin/user"]