# Local Development Tools

This directory contains all development tools and configurations for the Clariti project.

## Directory Structure

```
local/
├── dev/           # Development tools and scripts
├── benchmarks/    # Performance testing and analysis
├── config/        # Configuration examples and templates
├── certs/         # HTTPS certificates and security tools
└── README.md      # This documentation
```

## Quick Start

1. **Check code quality:**
   ```bash
   ./local/dev/lint.sh
   ```

2. **Manage dependencies:**
   ```bash
   ./local/dev/vendor.sh status
   ./local/dev/vendor.sh update
   ```

3. **Run performance benchmarks:**
   ```bash
   ./local/benchmarks/benchmark.sh
   ```

4. **Generate performance baseline:**
   ```bash
   ./local/benchmarks/analyze.sh --baseline
   ```

5. **View all available tools:**
   ```bash
   ./local/dev/tools.sh
   ```

## Development Scripts (`dev/`)

| Script | Purpose |
|--------|---------|
| `lint.sh` | Run golangci-lint for code quality checks |
| `vendor.sh` | Manage Go module vendoring and dependencies |
| `tools.sh` | Display all available development tools |
| `log.sh` | Shared logging library for all scripts |

### Shared Logging Library

The `log.sh` file provides consistent colored logging functions for all development scripts. All scripts in this directory use this shared library for standardized output.

**Basic Functions:**
- `log_info "message"` - Informational messages (blue)
- `log_success "message"` - Success messages (green)
- `log_warn "message"` - Warning messages (yellow)
- `log_error "message"` - Error messages (red)
- `log_debug "message"` - Debug messages (gray, DEBUG mode only)

**Specialized Functions:**
- `log_vendor "message"` - Vendor operations
- `log_build "message"` - Build operations
- `log_test "message"` - Test operations
- `log_benchmark "message"` - Benchmark operations
- `log_lint "message"` - Linting operations

**Utility Functions:**
- `log_header "title"` - Section headers with separators
- `log_section "title"` - Section titles
- `log_config_summary "title" "key1" "value1" "key2" "value2"` - Configuration displays
- `log_list "title" "item1" "item2"` - Formatted lists
- `log_progress 3 10 "task"` - Progress indicators

**Usage in Scripts:**
```bash
#!/bin/bash
source "$(dirname "$0")/log.sh"
log_info "Script started"
```

**Debug Mode:**
Set `DEBUG=true` to enable debug messages:
```bash
DEBUG=true ./script.sh
```

## Benchmark Tools (`benchmarks/`)

| Script | Purpose |
|--------|---------|
| `benchmark.sh` | Configurable performance benchmarking |
| `analyze.sh` | Advanced benchmark analysis and comparison |

### Benchmark Options

The `benchmark.sh` script supports various options:

```bash
# Basic benchmark run
./local/benchmarks/benchmark.sh

# Custom timeout and iterations
./local/benchmarks/benchmark.sh --timeout 2m --count 10

# Enable memory profiling
./local/benchmarks/benchmark.sh --memory

# Generate JSON output
./local/benchmarks/benchmark.sh --json
```

### Analysis Features

The `analyze.sh` script provides advanced performance analysis:

```bash
# Generate baseline
./local/benchmarks/analyze.sh --baseline

# Compare with baseline
./local/benchmarks/analyze.sh --compare

# Analyze specific benchmark file
./local/benchmarks/analyze.sh --file results/benchmark_20240101_120000.txt
```

## Configuration Examples (`config/`)

| File | Purpose |
|------|---------|
| `config.example.yaml` | Basic development configuration with RAM storage |
| `config.https.yaml` | HTTPS-enabled configuration |
| `config.ram.yaml` | RAM storage configuration (data lost on restart) |
| `config.s3.yaml` | S3 storage configuration for persistent data |

### Storage Drivers

The application supports multiple storage drivers for event data:

**RAM Storage (`ram`):**
- In-memory storage (default)
- Fast performance
- Data lost on application restart
- Suitable for development and testing

**S3 Storage (`s3`):**
- Persistent storage using AWS S3 or S3-compatible services
- Data survives application restarts
- Supports AWS IAM roles or access keys
- Configurable bucket, region, and object prefix
- Suitable for production environments

**Storage Configuration:**
```yaml
storage:
  driver: "s3"  # Options: "ram", "s3"
  s3:
    region: "us-east-1"
    bucket: "my-clariti-bucket"
    # Optional credentials (use IAM roles when possible)
    access_key_id: "AKIAIOSFODNN7EXAMPLE"
    secret_access_key: "secret"
    # Optional S3-compatible endpoint
    endpoint: "http://localhost:9000"
    # Optional object key prefix
    prefix: "clariti/"
```

## HTTPS Tools (`certs/`)

| Script | Purpose |
|--------|---------|
| `generate-certs.sh` | Generate development certificates |
| `test-https.sh` | Test HTTPS functionality |

## Development Workflow

1. **Make changes** to your Go code
2. **Check dependencies:** `./local/dev/vendor.sh verify`
3. **Check quality:** `./local/dev/lint.sh`
4. **Test functionality:** `go test ./...`
5. **Verify performance:** `./local/benchmarks/benchmark.sh`
6. **Commit changes:** `git commit`

## Dependency Management

The project uses Go module vendoring for reproducible builds:

```bash
# Check vendor status
./local/dev/vendor.sh status

# Update all dependencies
./local/dev/vendor.sh update

# Verify vendor integrity
./local/dev/vendor.sh verify

# Build with vendored dependencies
./local/dev/vendor.sh build

# Test with vendored dependencies
./local/dev/vendor.sh test
```

## Performance Baselines

Current project performance highlights:
- **GUID generation:** ~250ns/op (memory pool optimized)
- **String normalization:** 71-473ns/op (scales with length)
- **Weather service:** ~200ns/op (efficient calculation)
- **HTTP handlers:** ~2.5μs/op (full request/response cycle)

## Tool Features

### No Hardcoded Values
All scripts are configurable and avoid hardcoded values:
- Dynamic package discovery
- Configurable timeouts and iterations
- Flexible output formats
- Environment-aware operation

### Comprehensive Coverage
- Code quality checks with golangci-lint
- Performance benchmarking across all packages
- Security testing with HTTPS
- Configuration examples for various scenarios

## Getting Help

Most scripts support a `--help` flag for detailed usage information:

```bash
./local/benchmarks/benchmark.sh --help
./local/benchmarks/analyze.sh --help
```

For a complete overview of all tools:

```bash
./local/dev/tools.sh
```

Then in your `config.yaml`:

```yaml
server:
  host: "localhost"
  port: "8443"  # Alternative HTTPS standard port
  cert_file: "local/server.crt"
  key_file: "local/server.key"
```

## File structure

- `config.example.yaml` - Example configuration for development
- `generate-certs.sh` - Script to generate test certificates
- `server.crt` / `server.key` - Generated certificates (not committed)
- `config.yaml` - Your local configuration (not committed)

## Notes

- All files in this directory are ignored by Git (`.gitignore`)
- Generated certificates are self-signed and for testing only
- Use valid certificates in production
