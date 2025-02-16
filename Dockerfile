FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/merchshop ./cmd/main.go ./cmd/app.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/merchshop .
CMD ["/app/merchshop"]
LABEL authors="danya"
