# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go.mod first for caching
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o phntm .

# Runtime stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary
COPY --from=builder /build/phntm /usr/local/bin/phntm

# Create non-root user
RUN adduser -D -u 1000 phntm
USER phntm

WORKDIR /home/phntm

ENTRYPOINT ["phntm"]
CMD ["--help"]