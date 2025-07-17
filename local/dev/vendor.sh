#!/bin/bash
# vendor.sh - Vendor management script for Clariti

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  update     - Update vendor directory with latest dependencies"
    echo "  verify     - Verify vendor directory is up to date"
    echo "  clean      - Clean vendor directory"
    echo "  status     - Show vendor status and dependencies"
    echo "  build      - Build using vendor directory"
    echo "  test       - Test using vendor directory"
    echo "  help       - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 update          # Update all vendored dependencies"
    echo "  $0 verify          # Check if vendor is up to date"
    echo "  $0 build           # Build with vendored dependencies"
}

update_vendor() {
    echo -e "${BLUE}📦 Updating vendor directory...${NC}"
    
    # Ensure go.mod and go.sum are up to date
    echo -e "${YELLOW}🔄 Tidying modules...${NC}"
    go mod tidy
    
    # Download all dependencies
    echo -e "${YELLOW}⬇️  Downloading dependencies...${NC}"
    go mod download
    
    # Update vendor directory
    echo -e "${YELLOW}📁 Updating vendor directory...${NC}"
    go mod vendor
    
    echo -e "${GREEN}✅ Vendor directory updated successfully!${NC}"
    show_vendor_stats
}

verify_vendor() {
    echo -e "${BLUE}🔍 Verifying vendor directory...${NC}"
    
    if [ ! -d "vendor" ]; then
        echo -e "${RED}❌ Vendor directory does not exist${NC}"
        return 1
    fi
    
    # Check if vendor is up to date
    if go mod verify &>/dev/null; then
        echo -e "${GREEN}✅ Vendor directory is up to date${NC}"
        return 0
    else
        echo -e "${RED}❌ Vendor directory is out of date${NC}"
        echo -e "${YELLOW}💡 Run '$0 update' to fix this${NC}"
        return 1
    fi
}

clean_vendor() {
    echo -e "${BLUE}🧹 Cleaning vendor directory...${NC}"
    
    if [ -d "vendor" ]; then
        rm -rf vendor/
        echo -e "${GREEN}✅ Vendor directory cleaned${NC}"
    else
        echo -e "${YELLOW}⚠️  Vendor directory doesn't exist${NC}"
    fi
}

show_vendor_status() {
    echo -e "${BLUE}📊 Vendor Status${NC}"
    echo "===================="
    
    if [ -d "vendor" ]; then
        echo -e "${GREEN}✅ Vendor directory exists${NC}"
        show_vendor_stats
        
        echo ""
        echo -e "${BLUE}📋 Dependencies:${NC}"
        if [ -f "vendor/modules.txt" ]; then
            grep "^# " vendor/modules.txt | sed 's/^# /  - /'
        fi
        
        echo ""
        echo -e "${BLUE}🔍 Verification:${NC}"
        if go mod verify &>/dev/null; then
            echo -e "${GREEN}✅ All dependencies verified${NC}"
        else
            echo -e "${RED}❌ Some dependencies may be corrupted${NC}"
        fi
    else
        echo -e "${RED}❌ Vendor directory does not exist${NC}"
        echo -e "${YELLOW}💡 Run '$0 update' to create it${NC}"
    fi
}

show_vendor_stats() {
    if [ -d "vendor" ]; then
        local pkg_count=$(find vendor -name "*.go" | wc -l)
        local dir_count=$(find vendor -type d | wc -l)
        local size=$(du -sh vendor 2>/dev/null | cut -f1)
        
        echo -e "${BLUE}📈 Vendor Statistics:${NC}"
        echo "  • Size: $size"
        echo "  • Directories: $dir_count"
        echo "  • Go files: $pkg_count"
    fi
}

build_with_vendor() {
    echo -e "${BLUE}🔨 Building with vendor directory...${NC}"
    
    if [ ! -d "vendor" ]; then
        echo -e "${RED}❌ Vendor directory does not exist${NC}"
        echo -e "${YELLOW}💡 Run '$0 update' first${NC}"
        return 1
    fi
    
    # Build using vendor
    echo -e "${YELLOW}🏗️  Building project...${NC}"
    go build -mod=vendor -v ./...
    
    echo -e "${GREEN}✅ Build completed using vendor directory${NC}"
}

test_with_vendor() {
    echo -e "${BLUE}🧪 Testing with vendor directory...${NC}"
    
    if [ ! -d "vendor" ]; then
        echo -e "${RED}❌ Vendor directory does not exist${NC}"
        echo -e "${YELLOW}💡 Run '$0 update' first${NC}"
        return 1
    fi
    
    # Test using vendor
    echo -e "${YELLOW}🔬 Running tests...${NC}"
    go test -mod=vendor -v ./...
    
    echo -e "${GREEN}✅ Tests completed using vendor directory${NC}"
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
        echo -e "${YELLOW}⚠️  No command specified${NC}"
        echo ""
        print_usage
        exit 1
        ;;
    *)
        echo -e "${RED}❌ Unknown command: $1${NC}"
        echo ""
        print_usage
        exit 1
        ;;
esac
