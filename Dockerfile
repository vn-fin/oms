# ---- Build stage ----
FROM golang:1.24.1 AS builder

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Install swag CLI (only in builder image)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source
COPY . .

# Generate Swagger docs
RUN $(go env GOPATH)/bin/swag init -g cmd/main.go

# Build static binary
RUN go build -o app ./cmd/main.go

# ---- Run stage ----
FROM alpine:3.20

# Add non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary
COPY --from=builder /app/app .
# Copy generated Swagger docs if needed (optional)
COPY --from=builder /app/docs ./docs

USER appuser

EXPOSE 3000

ENTRYPOINT ["./app"]