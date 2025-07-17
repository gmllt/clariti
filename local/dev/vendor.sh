#!/bin/bash
# vendor.sh - Vendor management script for Clariti

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Source shared logging library
source "$SCRIPT_DIR/log.sh"

print_usage() {
    log_header "Vendor Management"
    log_list "Available Commands" \
        "update - Update vendor directory with latest dependencies" \
        "verify - Verify vendor directory is up to date" \
        "clean - Clean vendor directory" \
        "status - Show vendor status and dependencies" \
        "build - Build using vendor directory" \
        "test - Test using vendor directory" \
        "help - Show this help"
    
    log_section "Examples"
    log_info "Update all vendored dependencies:"
    echo "  $0 update"
    log_info "Check if vendor is up to date:"
    echo "  $0 verify"
    log_info "Build with vendored dependencies:"
    echo "  $0 build"
}

update_vendor() {
    log_vendor "Updating vendor directory..."
    
    # Ensure go.mod and go.sum are up to date
    log_info "Tidying modules..."
    go mod tidy
    
    # Download all dependencies
    log_info "Downloading dependencies..."
    go mod download
    
    # Update vendor directory
    log_info "Updating vendor directory..."
    go mod vendor
    
    log_success "Vendor directory updated successfully!"
    show_vendor_stats
}

verify_vendor() {
    log_vendor "Verifying vendor directory..."
    
    if [ ! -d "vendor" ]; then
        log_error "Vendor directory does not exist"
        return 1
    fi
    
    # Check if vendor is up to date
    if go mod verify &>/dev/null; then
        log_success "Vendor directory is up to date"
        return 0
    else
        log_error "Vendor directory is out of date"
        log_warn "Run '$0 update' to fix this"
        return 1
    fi
}

clean_vendor() {
    log_vendor "Cleaning vendor directory..."
    
    if [ -d "vendor" ]; then
        rm -rf vendor/
        log_success "Vendor directory cleaned"
    else
        log_warn "Vendor directory doesn't exist"
    fi
}

show_vendor_status() {
    log_header "Vendor Status"
    
    if [ -d "vendor" ]; then
        log_success "Vendor directory exists"
        show_vendor_stats
        
        log_section "Dependencies"
        if [ -f "vendor/modules.txt" ]; then
            grep "^# " vendor/modules.txt | sed 's/^# /  - /'
        fi
        
        log_section "Verification"
        if go mod verify &>/dev/null; then
            log_success "All dependencies verified"
        else
            log_error "Some dependencies may be corrupted"
        fi
    else
        log_error "Vendor directory does not exist"
        log_info "Run '$0 update' to create it"
    fi
}

show_vendor_stats() {
    if [ -d "vendor" ]; then
        local pkg_count=$(find vendor -name "*.go" | wc -l)
        local dir_count=$(find vendor -type d | wc -l)
        local size=$(du -sh vendor 2>/dev/null | cut -f1)
        
        log_config_summary "Vendor Statistics" \
            "Size" "$size" \
            "Directories" "$dir_count" \
            "Go files" "$pkg_count"
    fi
}

build_with_vendor() {
    log_build "Building with vendor directory..."
    
    if [ ! -d "vendor" ]; then
        log_error "Vendor directory does not exist"
        log_info "Run '$0 update' first"
        return 1
    fi
    
    # Build using vendor
    log_info "Building project..."
    go build -mod=vendor -v ./...
    
    log_success "Build completed using vendor directory"
}

test_with_vendor() {
    log_test "Testing with vendor directory..."
    
    if [ ! -d "vendor" ]; then
        log_error "Vendor directory does not exist"
        log_info "Run '$0 update' first"
        return 1
    fi
    
    # Test using vendor
    log_info "Running tests..."
    go test -mod=vendor -v ./...
    
    log_success "Tests completed using vendor directory"
}

# Main script logic
cd "$PROJECT_ROOT"

case "${1:-}" in
    "update")
        update_vendor
        ;;
    "verify")
        verify_vendor
        ;;
    "clean")
        clean_vendor
        ;;
    "status")
        show_vendor_status
        ;;
    "build")
        build_with_vendor
        ;;
    "test")
        test_with_vendor
        ;;
    "help"|"--help"|"-h")
        print_usage
        ;;
    "")
        log_warn "No command specified"
        print_usage
        exit 1
        ;;
    *)
        log_error "Unknown command: $1"
        print_usage
        exit 1
        ;;
esac
