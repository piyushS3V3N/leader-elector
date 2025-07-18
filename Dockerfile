# Use a minimal Go base image
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies and build binary
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o leader-elector ./cmd/main.go

# Create final image
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/leader-elector /app/leader-elector

ENTRYPOINT ["/app/leader-elector"]

