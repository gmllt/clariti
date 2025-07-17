# Clariti CLI

**Command Line Interface for Clariti** - Manage incidents and maintenance from your terminal.

## Installation

### Using Go
```bash
go install github.com/gmllt/clariti/cli@latest
```

### Using GoReleaser builds
Download the appropriate binary for your platform from the releases page.

## Configuration

The CLI can be configured using environment variables or command line flags.

### Environment Variables

All environment variables start with `CLARITI_`:

- `CLARITI_SERVER_URL` - Clariti server URL (default: http://localhost:8080)
- `CLARITI_USERNAME` - Basic auth username 
- `CLARITI_PASSWORD` - Basic auth password
- `CLARITI_OUTPUT_FORMAT` - Output format: json, yaml, table (default: table)
- `CLARITI_TRACE_ENABLED` - Enable trace output (default: false)
- `CLARITI_TRACE_FILE` - Trace output file (default: stdout)

### Example Configuration

```bash
export CLARITI_SERVER_URL="https://clariti.yourcompany.com"
export CLARITI_USERNAME="admin"
export CLARITI_PASSWORD="your-secure-password"
export CLARITI_OUTPUT_FORMAT="json"
export CLARITI_TRACE_ENABLED="true"
export CLARITI_TRACE_FILE="/tmp/clariti-cli.log"
```

## Commands

### System Status

```bash
# Check server health
clariti-cli health

# Get service weather overview
clariti-cli weather
```

### Components

```bash
# List all components
clariti-cli components list

# Show component hierarchy tree
clariti-cli components tree

# List platforms only
clariti-cli platforms

# List instances only  
clariti-cli instances
```

### Incidents

```bash
# List all incidents
clariti-cli incident list

# Get specific incident
clariti-cli incident get [incident-id]

# Create new incident
clariti-cli incident create --title "API Down" --description "API server not responding" --severity major --component api-01

# Update incident
clariti-cli incident update [incident-id] --status resolved

# Delete incident
clariti-cli incident delete [incident-id]
```

### Planned Maintenance

```bash
# List all maintenances
clariti-cli maintenance list

# Get specific maintenance
clariti-cli maintenance get [maintenance-id]

# Create new maintenance
clariti-cli maintenance create --title "Database Upgrade" --description "Upgrading to new version" --component db-01 --start-time "2024-01-15T02:00:00Z" --end-time "2024-01-15T04:00:00Z"

# Update maintenance
clariti-cli maintenance update [maintenance-id] --status in-progress

# Delete maintenance
clariti-cli maintenance delete [maintenance-id]
```

## Output Formats

### Table (default)
```bash
clariti-cli --output table health
```

### JSON
```bash
clariti-cli --output json incident list
```

### YAML
```bash
clariti-cli --output yaml components tree
```

## Authentication

### Using Environment Variables
```bash
export CLARITI_USERNAME="admin"
export CLARITI_PASSWORD="your-password"
clariti-cli incident create --title "Test Incident"
```

### Using Command Line Flags
```bash
clariti-cli --username admin --password your-password incident list
```

## Tracing

Enable tracing to debug CLI operations:

```bash
# Trace to stdout
export CLARITI_TRACE_ENABLED="true"
clariti-cli health

# Trace to file
export CLARITI_TRACE_ENABLED="true"
export CLARITI_TRACE_FILE="/tmp/clariti-debug.log"
clariti-cli incident list
```

## Examples

### Create an incident for a critical outage
```bash
clariti-cli incident create \
  --title "Database Server Down" \
  --description "Primary database server is not responding" \
  --severity critical \
  --component db-primary
```

### Schedule a maintenance window
```bash
clariti-cli maintenance create \
  --title "Security Patches" \
  --description "Installing security updates on web servers" \
  --component web-cluster \
  --start-time "2024-01-20T03:00:00Z" \
  --end-time "2024-01-20T05:00:00Z"
```

### Monitor service status
```bash
# Quick health check
clariti-cli health

# Detailed service overview
clariti-cli weather --output json | jq '.overall_status'

# List active incidents
clariti-cli incident list --output table
```

## Cross-Platform Support

The CLI supports all platforms covered by the main Clariti build:

- **Linux** (AMD64, ARM64, 386, ARM v6/v7)
- **Windows** (AMD64, 386)
- **macOS** (Intel, Apple Silicon)
- **FreeBSD** (AMD64, ARM64, 386, ARM v6/v7)

---

**Ready to manage!** Control your Clariti server from anywhere with the CLI.
