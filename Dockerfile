FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY examples/ ./examples/
COPY pkg/ ./pkg/

# Build the web dashboard
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o coffee-dashboard examples/web_dashboard.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary
COPY --from=builder /app/coffee-dashboard .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the application
CMD ["./coffee-dashboard"]
