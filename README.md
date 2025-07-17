# Clariti

REST API for incident and planned maintenance management with hierarchical component structure.

## ✨ Features

- 🏗️ **Modular architecture** with hierarchical structure (Platform → Instance → Component)
- 🔐 **Basic Auth authentication** for write operations
- 🌐 **Complete REST API** with versioning (`/api/v1/`)
- 🌤️ **Weather endpoint** - service status overview
- 🔒 **Optional HTTPS/TLS support** for security
- 📚 **Auto-generated API documentation**
- ✅ **Complete test suite**

## 🚀 Quick start

### 1. Configuration

```bash
# Copy example configuration
cp config.sample.yaml config.yaml

# Customize configuration
nano config.yaml
```

### 2. Build and run

```bash
# Build server
go build -o clariti-server ./server

# Start HTTP server
./clariti-server
```

### 3. Test HTTPS (optional)

```bash
# Generate certificates and test HTTPS
cd local/
./test-https.sh
```

## 📁 Project structure

```
├── server/          # REST API server code
├── models/          # Data models (events, components)
├── utils/           # Shared utilities
├── tests/           # Integration tests
├── local/           # Local development and testing (not committed)
├── config.sample.yaml  # Example configuration
└── .gitignore       # Files ignored by Git
```

## 🔧 Configuration

See [`config.sample.yaml`](config.sample.yaml) for complete example.

### Minimal configuration

```yaml
server:
  host: "localhost"
  port: "8080"

auth:
  admin_username: "admin"
  admin_password: "your-password"

components:
  platforms:
    - name: "Production"
      code: "PROD"
      instances:
        - name: "API"
          code: "api"
          components:
            - name: "REST API"
              code: "rest"
```

### HTTPS configuration

```yaml
server:
  host: "localhost"
  port: "8443"
  cert_file: "path/to/certificate.crt"
  key_file: "path/to/private.key"
```

## 🌐 Endpoints API

- **Health**: `GET /health`
- **Documentation**: `GET /api/v1/docs`
- **Weather**: `GET /api/v1/weather`
- **Incidents**: `GET|POST /api/v1/incidents`
- **Maintenances**: `GET|POST /api/v1/planned-maintenances`

## 🧪 Testing and development

### `local/` directory

The `local/` directory contains all test and development files that are not committed:

- 📋 `config.example.yaml` - Example configuration for development
- 🔐 `generate-certs.sh` - Self-signed certificate generation
- 🌐 `test-https.sh` - Complete HTTPS test script
- 🗄️ `*.crt`, `*.key` - Generated certificates (ignored by Git)

### Automated tests

```bash
# All tests
go test ./...

# Tests with coverage
go test -cover ./...

# Specific integration tests
go test ./tests/...
```

## 📚 Detailed documentation

- [Server documentation](server/README.md) - Complete REST API server guide
- [Local documentation](local/README.md) - Development and local testing guide

## 🔒 Security

- Basic Auth authentication for write endpoints
- Optional HTTPS/TLS support with custom certificates
- CORS configured for cross-origin calls
- Input validation on all endpoints

## 🤝 Contributing

1. Files in `local/` are ignored by Git
2. Use `config.sample.yaml` as base for your configurations
3. Run tests before committing: `go test ./...`
