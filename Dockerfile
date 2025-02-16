FROM golang:1.23.1-alpine AS builder
RUN apk add --no-cache tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/merchshop ./cmd/main.go ./cmd/app.go

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/merchshop .
CMD ["/app/merchshop"]
LABEL authors="danya"
