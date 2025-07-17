#!/bin/bash
# tools.sh - Display all available development tools

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOCAL_DIR="$(dirname "$SCRIPT_DIR")"

# Source shared logging library
source "$SCRIPT_DIR/log.sh"

log_header "Clariti Development Tools"

log_section "Available Scripts"
log_list "Development Tools" \
    "TEST: ./local/dev/test.sh - Run all tests with coverage (if exists)" \
    "LINT: ./local/dev/lint.sh - Run golangci-lint code quality checks" \
    "VENDOR: ./local/dev/vendor.sh - Manage Go module vendoring" \
    "BENCHMARK: ./local/benchmarks/benchmark.sh - Run configurable performance benchmarks" \
    "ANALYZE: ./local/benchmarks/analyze.sh - Advanced benchmark analysis & comparison" \
    "HELP: ./local/dev/tools.sh - Show this help menu"

log_section "HTTPS & Security"
log_list "Security Tools" \
    "CERTS: ./local/certs/generate-certs.sh - Generate development certificates" \
    "TEST: ./local/certs/test-https.sh - Test HTTPS functionality"

log_section "Configuration Examples"
log_list "Config Templates" \
    "BASIC: ./local/config/config.example.yaml - Basic configuration template" \
    "HTTPS: ./local/config/config.https.yaml - HTTPS configuration example"

log_section "Quick Commands"
log_info "Run quality checks:"
echo "  ./local/dev/lint.sh"

log_info "Manage dependencies:"
echo "  ./local/dev/vendor.sh status"
echo "  ./local/dev/vendor.sh update"
echo "  ./local/dev/vendor.sh build"

log_info "Run benchmarks with options:"
echo "  ./local/benchmarks/benchmark.sh --help"
echo "  ./local/benchmarks/analyze.sh --baseline"
echo "  ./local/benchmarks/analyze.sh --compare"

log_section "Development Workflow"
log_list "Step by Step" \
    "1. Code changes - Edit your Go files" \
    "2. ./local/dev/vendor.sh verify - Check dependencies" \
    "3. ./local/dev/lint.sh - Check code quality" \
    "4. go test ./... - Verify functionality" \
    "5. ./local/benchmarks/benchmark.sh - Check performance" \
    "6. git commit - Commit your changes"

log_section "Vendor Commands"
log_info "Check vendor status:"
echo "  ./local/dev/vendor.sh status"

log_info "Update dependencies:"
echo "  ./local/dev/vendor.sh update"

log_info "Build with vendor:"
echo "  ./local/dev/vendor.sh build"

log_info "Verify vendor integrity:"
echo "  ./local/dev/vendor.sh verify"

log_section "Benchmark Commands"
log_info "Generate performance baseline:"
echo "  ./local/benchmarks/analyze.sh --baseline"

log_info "Run benchmarks with custom timeout:"
echo "  ./local/benchmarks/benchmark.sh --timeout 2m"

log_info "Compare current performance:"
echo "  ./local/benchmarks/analyze.sh --compare"

log_section "Directory Structure"
echo "local/"
echo "├── dev/           # Development tools (lint, test scripts)"
echo "├── benchmarks/    # Performance testing tools and results"
echo "├── config/        # Configuration examples"
echo "├── certs/         # HTTPS certificates and tools"
echo "└── README.md      # Documentation"

log_success "All tools are ready to use!"
log_info "Most scripts support --help for detailed options"
