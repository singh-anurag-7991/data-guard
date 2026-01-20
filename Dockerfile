# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
# CGO_ENABLED=0 ensures a static binary for distroless
RUN CGO_ENABLED=0 GOOS=linux go build -o dataguard cmd/server/main.go

# Stage 2: Runner
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/dataguard .

# Expose port
EXPOSE 8080

# Run
CMD ["./dataguard"]
