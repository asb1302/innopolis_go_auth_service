FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app
RUN go build -o main ./cmd

# Начинаем новый этап с нуля (multi-stage build)
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY ./docs ./docs

EXPOSE 8000

CMD ["./main"]
