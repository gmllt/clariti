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

2. **Run performance benchmarks:**
   ```bash
   ./local/benchmarks/benchmark.sh
   ```

3. **Generate performance baseline:**
   ```bash
   ./local/benchmarks/analyze.sh --baseline
   ```

4. **View all available tools:**
   ```bash
   ./local/dev/tools.sh
   ```

## Development Scripts (`dev/`)

| Script | Purpose |
|--------|---------|
| `lint.sh` | Run golangci-lint for code quality checks |
| `tools.sh` | Display all available development tools |

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
| `config.example.yaml` | Basic configuration template |
| `config.https.yaml` | HTTPS-enabled configuration example |

## HTTPS Tools (`certs/`)

| Script | Purpose |
|--------|---------|
| `generate-certs.sh` | Generate development certificates |
| `test-https.sh` | Test HTTPS functionality |

## Development Workflow

1. **Make changes** to your Go code
2. **Check quality:** `./local/dev/lint.sh`
3. **Test functionality:** `go test ./...`
4. **Verify performance:** `./local/benchmarks/benchmark.sh`
5. **Commit changes:** `git commit`

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
