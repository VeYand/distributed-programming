FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go mod vendor

RUN go test -v ./...
RUN go build -o ./bin/rankcalculator ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/rankcalculator ./bin/rankcalculator

EXPOSE 8082

CMD ["./bin/rankcalculator"]