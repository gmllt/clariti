# Dockerfile for Clariti server
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 clariti && \
    adduser -D -s /bin/sh -u 1001 -G clariti clariti

# Set working directory
WORKDIR /app

# Copy binary from build context
COPY clariti-server /app/clariti-server

# Copy default configuration
COPY local/config/config.s3.yaml /app/config.yaml

# Change ownership
RUN chown -R clariti:clariti /app

# Switch to non-root user
USER clariti

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
ENTRYPOINT ["/app/clariti-server"]
CMD ["--config", "/app/config.yaml"]
