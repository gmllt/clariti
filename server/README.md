# Clariti REST API Server

A simple REST API server for managing incidents and planned maintenances with hierarchical component structure.

## Architecture

The server is built with a clean, modular architecture:

- **`core/`** - Main server object with graceful startup/shutdown
- **`config/`** - YAML configuration management  
- **`drivers/`** - Storage drivers (currently RAM-based)
- **`handlers/`** - HTTP request handlers
- **`middleware/`** - Authentication and CORS middleware
- **`routes/`** - Route definitions and API documentation

## Features

- ✅ RESTful API for incidents and planned maintenances
- ✅ Basic authentication for write operations
- ✅ In-memory storage (RAM driver)
- ✅ YAML configuration with hierarchical component structure
- ✅ Component information from configuration (read-only)
- ✅ CORS support
- ✅ Health check endpoint
- ✅ Graceful shutdown with signal handling
- ✅ Auto-generated API documentation endpoint

## Configuration

Create a `config.yaml` file with the following hierarchical structure:

```yaml
server:
  host: localhost
  port: "8080"

auth:
  admin_username: admin
  admin_password: password123

components:
  platforms:
    - name: Production
      code: prod
      base_url: https://prod.example.com
      instances:
        - name: Web Application
          code: webapp
          components:
            - name: Frontend
              code: frontend
            - name: Authentication Service
              code: auth
        - name: API Service
          code: api
          components:
            - name: REST API
              code: rest
            - name: Database
              code: db
```

## Running the Server

1. Build the server:
```bash
go build -o clariti-server server/main.go
```

2. Run with configuration:
```bash
./clariti-server
# or specify config path
CONFIG_PATH=/path/to/config.yaml ./clariti-server
```

3. Graceful shutdown:
   - The server handles SIGINT and SIGTERM signals
   - Press Ctrl+C to stop gracefully

## API Endpoints

### Health Check
- `GET /health` - Server health status

### Documentation
- `GET /api/docs` - Auto-generated API documentation

### Components (Read-only)
- `GET /api/components` - Get all component information (flattened)
- `GET /api/components/hierarchy` - Get hierarchical component structure
- `GET /api/platforms` - Get platforms
- `GET /api/instances` - Get instances
- `GET /api/components/list` - Get components with relationships

### Incidents
- `GET /api/incidents` - List all incidents
- `GET /api/incidents/{id}` - Get specific incident
- `POST /api/incidents` - Create incident (requires auth)
- `PUT /api/incidents/{id}` - Update incident (requires auth)
- `DELETE /api/incidents/{id}` - Delete incident (requires auth)

### Planned Maintenances
- `GET /api/planned-maintenances` - List all planned maintenances
- `GET /api/planned-maintenances/{id}` - Get specific planned maintenance
- `POST /api/planned-maintenances` - Create planned maintenance (requires auth)
- `PUT /api/planned-maintenances/{id}` - Update planned maintenance (requires auth)
- `DELETE /api/planned-maintenances/{id}` - Delete planned maintenance (requires auth)

## Authentication

Write operations (POST, PUT, DELETE) require HTTP Basic Authentication using the credentials specified in the configuration file.

Example:
```bash
curl -X POST http://localhost:8080/api/incidents \
  -u admin:password123 \
  -H "Content-Type: application/json" \
  -d '{"title": "Service Outage", "content": "Database connection issues"}'
```

## Example Requests

### Create an Incident
```bash
curl -X POST http://localhost:8080/api/incidents \
  -u admin:password123 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database Issues",
    "content": "Connection timeout to primary database",
    "criticality": 2,
    "perpetual": false
  }'
```

### Get All Incidents
```bash
curl http://localhost:8080/api/incidents
```

### Create a Planned Maintenance
```bash
curl -X POST http://localhost:8080/api/planned-maintenances \
  -u admin:password123 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database Maintenance",
    "content": "Scheduled database upgrade",
    "start_planned": "2025-01-20T02:00:00Z",
    "end_planned": "2025-01-20T06:00:00Z"
  }'
```

## Response Format

All responses are in JSON format. Error responses include an `error` field:

```json
{
  "error": "Incident not found"
}
```

## Storage

Currently uses in-memory storage (data is lost when server restarts). The architecture supports pluggable storage drivers for future database implementations.

## Testing

1. Get API documentation:
```bash
curl http://localhost:8080/api/docs | jq .
```

2. Use the provided test script to test the API:
```bash
# Start the server in one terminal
./clariti-server

# Run tests in another terminal
./test_api.sh
```

3. Or test manually:
```bash
# Test hierarchical structure
curl http://localhost:8080/api/components/hierarchy | jq .

# Test creating an incident
curl -X POST http://localhost:8080/api/incidents \
  -u admin:password123 \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Incident", "content": "Test content"}'
```

## Programmatic Usage

You can also use the server as a library:

```go
package main

import (
    "log"
    "github.com/gmllt/clariti/server/core"
)

func main() {
    // Create server instance
    server, err := core.New("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // Access components if needed
    config := server.GetConfig()
    storage := server.GetStorage()

    // Run server (blocks until interrupted)
    server.Run()
}
```
