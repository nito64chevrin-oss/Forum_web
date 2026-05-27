# Build stage
FROM golang:1.21-alpine AS builder

# Install gcc and musl-dev for CGO (required by go-sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./main

# Run stage
FROM alpine:3.19

# Install libc and ca-certificates
RUN apk add --no-cache libc6-compat ca-certificates

WORKDIR /app

# Copy binary and required assets from builder
COPY --from=builder /app/forum .
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views

EXPOSE 8080

CMD ["./forum"]
