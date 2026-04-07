# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Download dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /app/aubergine ./main.go

# Final runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/aubergine ./aubergine

EXPOSE 8080
CMD ["./aubergine"]
