#!/bin/bash
# benchmark.sh - Run comprehensive benchmarks with configurable options

set -e

# Configuration variables (no hardcoded values)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Default configuration
BENCHMARK_TIMEOUT="5m"
BENCHMARK_COUNT=""
MEMORY_PROFILING="true"
VERBOSE="false"
OUTPUT_FORMAT="standard"

# Parse command line options
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -t, --timeout DURATION    Set benchmark timeout (default: $BENCHMARK_TIMEOUT)"
    echo "  -c, --count NUMBER         Set benchmark iteration count"
    echo "  -m, --no-memory           Disable memory profiling"
    echo "  -v, --verbose             Enable verbose output"
    echo "  -f, --format FORMAT       Output format: standard|json|csv (default: $OUTPUT_FORMAT)"
    echo "  -o, --output FILE          Save results to file"
    echo "  -h, --help                Show this help"
    exit 1
}

while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--timeout)
            BENCHMARK_TIMEOUT="$2"
            shift 2
            ;;
        -c|--count)
            BENCHMARK_COUNT="-benchtime=$2x"
            shift 2
            ;;
        -m|--no-memory)
            MEMORY_PROFILING="false"
            shift
            ;;
        -v|--verbose)
            VERBOSE="true"
            shift
            ;;
        -f|--format)
            OUTPUT_FORMAT="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Create results directory if it doesn't exist
mkdir -p "$RESULTS_DIR"

# Function to log with timestamp
log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo "[$(date '+%H:%M:%S')] $*"
    fi
}

# Function to run benchmarks for a specific package
run_benchmark() {
    local package="$1"
    local name="$2"
    
    log "Running benchmarks for $name in package $package"
    
    local benchmark_args="-bench=. -timeout=$BENCHMARK_TIMEOUT $BENCHMARK_COUNT"
    
    if [[ "$MEMORY_PROFILING" == "true" ]]; then
        benchmark_args="$benchmark_args -benchmem"
    fi
    
    cd "$PROJECT_ROOT"
    go test $benchmark_args "$package" 2>/dev/null || {
        echo "❌ Failed to run benchmarks for $package" >&2
        return 1
    }
}

# Discover packages with benchmarks automatically
discover_benchmark_packages() {
    log "Discovering packages with benchmarks..."
    
    cd "$PROJECT_ROOT"
    find . -name "*_test.go" -exec grep -l "^func Benchmark" {} \; | \
        sed 's|/[^/]*$||' | \
        sed 's|^\./||' | \
        sort -u
}

# Main execution
main() {
    echo "🚀 Running Clariti Benchmarks"
    echo "=============================="
    echo "Configuration:"
    echo "  • Timeout: $BENCHMARK_TIMEOUT"
    echo "  • Memory profiling: $MEMORY_PROFILING"
    echo "  • Format: $OUTPUT_FORMAT"
    [[ -n "$OUTPUT_FILE" ]] && echo "  • Output file: $OUTPUT_FILE"
    echo ""
    
    # Discover packages with benchmarks
    local packages=($(discover_benchmark_packages))
    
    if [[ ${#packages[@]} -eq 0 ]]; then
        echo "❌ No benchmark packages found"
        exit 1
    fi
    
    log "Found ${#packages[@]} packages with benchmarks: ${packages[*]}"
    
    # Run benchmarks for each package
    local results=""
    for package in "${packages[@]}"; do
        local name=$(echo "$package" | sed 's|/| |g' | awk '{print $NF}')
        echo "📊 Running benchmarks for $name..."
        echo "----------------------------------"
        
        local output
        if output=$(run_benchmark "./$package/" "$name"); then
            echo "$output"
            results="$results$output\n"
        else
            echo "⚠️  Skipped $package due to errors"
        fi
        echo ""
    done
    
    # Save results if requested
    if [[ -n "$OUTPUT_FILE" ]]; then
        echo -e "$results" > "$OUTPUT_FILE"
        echo "📁 Results saved to: $OUTPUT_FILE"
    fi
    
    # Generate summary
    echo "✅ Benchmark run completed!"
    echo "📊 Summary: Found benchmarks in ${#packages[@]} packages"
    echo "📁 Results directory: $RESULTS_DIR"
}

# Run main function
main "$@"
