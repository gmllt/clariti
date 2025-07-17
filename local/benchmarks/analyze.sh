#!/bin/bash
# analyze.sh - Advanced benchmark analysis without hardcoded values

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
BASELINE_FILE="$RESULTS_DIR/baseline.txt"

# Performance thresholds (configurable)
PERFORMANCE_REGRESSION_THRESHOLD=10  # 10% slowdown
MEMORY_REGRESSION_THRESHOLD=20       # 20% more memory
ALLOCATION_REGRESSION_THRESHOLD=25   # 25% more allocations

# Parse options
GENERATE_BASELINE="false"
COMPARE_MODE="false"
VERBOSE="false"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  -b, --baseline            Generate new baseline"
    echo "  -c, --compare             Compare with baseline"
    echo "  -v, --verbose             Enable verbose output"
    echo "  -t, --threshold PERCENT   Set regression threshold (default: $PERFORMANCE_REGRESSION_THRESHOLD%)"
    echo "  -h, --help                Show this help"
    exit 1
}

while [[ $# -gt 0 ]]; do
    case $1 in
        -b|--baseline)
            GENERATE_BASELINE="true"
            shift
            ;;
        -c|--compare)
            COMPARE_MODE="true"
            shift
            ;;
        -v|--verbose)
            VERBOSE="true"
            shift
            ;;
        -t|--threshold)
            PERFORMANCE_REGRESSION_THRESHOLD="$2"
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

# Create results directory
mkdir -p "$RESULTS_DIR"

log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo "[$(date '+%H:%M:%S')] $*" >&2
    fi
}

# Extract benchmark metrics dynamically
extract_metrics() {
    local input="$1"
    local pattern="$2"
    
    echo "$input" | grep "$pattern" | while read -r line; do
        local name=$(echo "$line" | awk '{print $1}')
        local ns_per_op=$(echo "$line" | awk '{print $3}' | sed 's/ns\/op//')
        local bytes_per_op=$(echo "$line" | awk '{print $4}' | sed 's/B\/op//')
        local allocs_per_op=$(echo "$line" | awk '{print $5}' | sed 's/allocs\/op//')
        
        printf "%-50s %10.1f ns/op %8s B/op %6s allocs/op\n" \
               "$name" "$ns_per_op" "$bytes_per_op" "$allocs_per_op"
    done
}

# Generate performance baseline
generate_baseline() {
    echo "[BASELINE] Generating performance baseline..."
    cd "$PROJECT_ROOT"
    
    go test -bench=. -benchmem ./... 2>/dev/null | \
        grep "^Benchmark" > "$BASELINE_FILE"
    
    echo "[SUCCESS] Baseline saved to: $BASELINE_FILE"
    echo "[INFO] $count benchmarks recorded"
}

# Compare current performance with baseline
compare_with_baseline() {
    if [[ ! -f "$BASELINE_FILE" ]]; then
        echo "❌ No baseline found. Run with --baseline first."
        exit 1
    fi
    
    echo "[COMPARE] Comparing current performance with baseline..."
    
    local current_file="$RESULTS_DIR/current_$(date +%Y%m%d_%H%M%S).txt"
    cd "$PROJECT_ROOT"
    go test -bench=. -benchmem ./... 2>/dev/null | \
        grep "^Benchmark" > "$current_file"
    
    echo "🔍 Performance Analysis:"
    echo "========================"
    
    local regressions=0
    local improvements=0
    local stable=0
    
    while read -r baseline_line; do
        local bench_name=$(echo "$baseline_line" | awk '{print $1}')
        local baseline_ns=$(echo "$baseline_line" | awk '{print $3}' | sed 's/ns\/op//')
        
        local current_line=$(grep "^$bench_name" "$current_file" || echo "")
        if [[ -n "$current_line" ]]; then
            local current_ns=$(echo "$current_line" | awk '{print $3}' | sed 's/ns\/op//')
            
            if [[ -n "$baseline_ns" && -n "$current_ns" ]]; then
                local change=$(echo "$current_ns $baseline_ns" | awk '{
                    if ($2 == 0) print "N/A"
                    else printf "%.1f", (($1-$2)/$2)*100
                }')
                
                local status="📈"
                local category="stable"
                
                if [[ "$change" != "N/A" ]]; then
                    if (( $(echo "$change < -5" | bc -l 2>/dev/null || echo "0") )); then
                        status="🚀"
                        category="improvement"
                        ((improvements++))
                    elif (( $(echo "$change > $PERFORMANCE_REGRESSION_THRESHOLD" | bc -l 2>/dev/null || echo "0") )); then
                        status="⚠️"
                        category="regression"
                        ((regressions++))
                    else
                        ((stable++))
                    fi
                    
                    printf "%-45s %s %+6s%% (%s)\n" "$bench_name" "$status" "$change" "$category"
                fi
            fi
        fi
    done < "$BASELINE_FILE"
    
    echo ""
    echo "📈 Summary:"
    echo "  🚀 Improvements: $improvements"
    echo "  📈 Stable: $stable"
    echo "  ⚠️  Regressions: $regressions"
    
    if [[ $regressions -gt 0 ]]; then
        echo ""
        echo "⚠️  Found $regressions performance regressions (>$PERFORMANCE_REGRESSION_THRESHOLD%)"
        echo "   Consider investigating these benchmarks."
        return 1
    fi
}

# Analyze current benchmarks
analyze_current() {
    echo "🔬 Analyzing current benchmark performance..."
    echo "==========================================="
    
    cd "$PROJECT_ROOT"
    local output=$(go test -bench=. -benchmem ./... 2>/dev/null)
    
    echo "⚡ Fast operations (< 100 ns/op):"
    extract_metrics "$output" "Benchmark" | awk '$2 < 100'
    echo ""
    
    echo "📊 Medium operations (100-1000 ns/op):"
    extract_metrics "$output" "Benchmark" | awk '$2 >= 100 && $2 < 1000'
    echo ""
    
    echo "🐌 Slow operations (>= 1000 ns/op):"
    extract_metrics "$output" "Benchmark" | awk '$2 >= 1000'
    echo ""
    
    echo "💾 High memory operations (>= 1000 B/op):"
    extract_metrics "$output" "Benchmark" | awk '$4 >= 1000'
    echo ""
}

# Main execution
main() {
    if [[ "$GENERATE_BASELINE" == "true" ]]; then
        generate_baseline
    elif [[ "$COMPARE_MODE" == "true" ]]; then
        compare_with_baseline
    else
        analyze_current
    fi
}

main "$@"
