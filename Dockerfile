# Runtime stage - goreleaser provides pre-built binary
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy pre-built binary (goreleaser injects this)
COPY phntm /usr/local/bin/phntm

# Create non-root user
RUN adduser -D -u 1000 phntm
USER phntm

WORKDIR /home/phntm

ENTRYPOINT ["phntm"]
CMD ["--help"]