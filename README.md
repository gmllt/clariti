# Clariti

**Service Status Management System** - Track incidents and maintenance with monitoring support.

## What is Clariti

Clariti helps you manage service status pages and track system health:

- **Component Management** - Organize services in Platform → Instance → Component structure  
- **Incident Tracking** - Create and manage system outages with severity levels
- **Maintenance Planning** - Schedule and track planned maintenance windows
- **Monitoring Integration** - Built-in Prometheus metrics for system health
- **REST API** - Complete API for status page integration
- **Multiple Storage** - Use RAM for testing or S3/MinIO for production
- **Authentication** - Secure admin operations with basic auth

## Quick Start

### 1. Get Configuration

```bash
# Copy example config
cp local/config/config.s3.yaml config.yaml

# Edit settings
nano config.yaml
```

### 2. Build and Run

**Using GoReleaser (recommended):**
```bash
# Build for your platform
goreleaser build --snapshot --clean --single-target

# Run server
./dist/clariti-server_linux_amd64_v1/clariti-server --config config.yaml serve
```

**Using Go directly:**
```bash
# Build manually
go build -o clariti-server ./server

# Run server
./clariti-server --config config.yaml serve
```

### 3. Test Setup

```bash
# Check health
curl http://localhost:8080/health

# View service status
curl http://localhost:8080/api/v1/weather

# See metrics
curl http://localhost:8080/metrics
```

## Project Structure

```
clariti/
├── server/               # Main REST API server
│   ├── main.go          # Application entry point
│   ├── core/            # Server setup and routing
│   ├── handlers/        # HTTP request handlers  
│   ├── middleware/      # HTTP middleware (auth, metrics)
│   ├── routes/          # API route definitions
│   ├── drivers/         # Storage backend drivers
│   ├── metrics/         # Prometheus metrics setup
│   └── config/          # Configuration loading
├── models/              # Data models and structures
│   ├── event/           # Incident and maintenance models
│   └── component/       # Component hierarchy models
├── logger/              # Application logging
├── utils/               # Helper utilities  
├── local/               # Development files (not in git)
├── .goreleaser.yaml     # Multi-platform build setup
└── config.*.yaml        # Configuration examples
```

## Configuration

Complete storage configuration guide: [Storage Drivers](STORAGE_DRIVERS.md)

### Basic Setup

```yaml
# Server settings
server:
  host: "localhost"
  port: "8080"

# Admin access
auth:
  admin_username: "admin"
  admin_password: "your-secure-password"

# Logging
logging:
  level: "info"
  format: "text"

# Storage (use "ram" for testing)
storage:
  driver: "ram"

# Your service components
components:
  platforms:
    - name: "Production Environment"
      code: "PROD"
      base_url: "https://api.yoursite.com"
      instances:
        - name: "Web Services"
          code: "web"
          components:
            - name: "Load Balancer"
              code: "lb-01"
            - name: "API Server"
              code: "api-01"
```

### Production Storage (S3/MinIO)

```yaml
storage:
  driver: "s3"
  s3:
    region: "us-east-1"
    bucket: "clariti-status"
    endpoint: "http://localhost:9000"  # For MinIO
    access_key_id: "admin"
    secret_access_key: "password"
    prefix: "clariti/"
```

## API Reference

### System Health
- `GET /health` - Check if server is running
- `GET /api/v1/weather` - Get overall service status  
- `GET /metrics` - Prometheus metrics endpoint
- `GET /api/docs` and `GET /api/v1/docs` - API documentation with version info

### Components
- `GET /api/v1/components` - List all components
- `GET /api/v1/components/hierarchy` - Get component tree
- `GET /api/v1/platforms` - List platforms only
- `GET /api/v1/instances` - List instances only

### Incidents (Authentication Required for POST/PUT/DELETE)
- `GET /api/v1/incidents` - List all incidents
- `POST /api/v1/incidents` - Create new incident
- `GET /api/v1/incidents/{id}` - Get specific incident
- `PUT /api/v1/incidents/{id}` - Update incident
- `DELETE /api/v1/incidents/{id}` - Delete incident

### Maintenance (Authentication Required for POST/PUT/DELETE)
- `GET /api/v1/planned-maintenances` - List all maintenances
- `POST /api/v1/planned-maintenances` - Schedule maintenance
- `GET /api/v1/planned-maintenances/{id}` - Get specific maintenance
- `PUT /api/v1/planned-maintenances/{id}` - Update maintenance  
- `DELETE /api/v1/planned-maintenances/{id}` - Cancel maintenance

## Monitoring

Clariti exports Prometheus metrics for monitoring:

### HTTP Metrics
- `clariti_http_requests_total` - Request count by method, path, status
- `clariti_http_request_duration_seconds` - Response time distribution

### System Metrics
- `clariti_uptime_seconds` - Server uptime since start
- `clariti_components_total` - Number of monitored components
- `clariti_application_info` - Version and build information

### Business Metrics
- `clariti_incidents_total` - Current incidents by severity and status
- `clariti_planned_maintenances_total` - Scheduled maintenances by status

## Development

See [Local Development Guide](local/README.md) for development setup.

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests only
go test ./tests/...
```

### Building

```bash
# Single platform (current OS)
goreleaser build --snapshot --clean --single-target

# All supported platforms
goreleaser build --snapshot --clean

# Test release process
goreleaser release --snapshot --clean --skip=publish
```

## Security

- **Basic Authentication** protects write operations (POST/PUT/DELETE)
- **CORS support** for web browser integration
- **Input validation** on all endpoints
- **HTTPS/TLS ready** with custom certificate support
- **Structured logging** for security audit trails

## Documentation

- **[Storage Drivers Guide](STORAGE_DRIVERS.md)** - How to configure different storage backends
- **[Local Development](local/README.md)** - Development and testing guide

## Contributing

1. Write code and documentation in simple English
2. Use the `local/` directory for development (ignored by git)
3. Run tests before committing: `go test ./...`
4. Follow existing code style and logging patterns

## Supported Platforms

**Operating Systems:**
- Linux (AMD64, ARM64, 386, ARM v6/v7)  
- Windows (AMD64, 386)
- macOS (Intel, Apple Silicon)
- FreeBSD (AMD64, ARM64, 386, ARM v6/v7)

**Release Formats:**
- TAR.GZ and ZIP archives
- SHA256 checksums included
- Sample configurations included

---

**Ready to use!** Start monitoring your services with Clariti.
