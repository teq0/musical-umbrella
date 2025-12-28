# Build stage
FROM golang:1.19-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o musical-umbrella .

# Final stage
FROM scratch

# Import certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy our static executable
COPY --from=builder /app/musical-umbrella /musical-umbrella

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/musical-umbrella"]