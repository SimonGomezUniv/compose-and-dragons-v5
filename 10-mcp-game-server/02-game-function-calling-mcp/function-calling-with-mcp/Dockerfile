# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /build

# Copy from current directory (context is npc-agent-services)
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o tools-test .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /build/tools-test .

# Run the application
#CMD ["./tools-test"]
