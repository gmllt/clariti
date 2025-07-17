#!/bin/bash
# test-minio.sh - Test MinIO connectivity and API functionality

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/log.sh"

CONFIG_FILE="local/config/config.s3.yaml"
API_PORT="8080"
API_URL="http://localhost:$API_PORT"

# Test MinIO connectivity first
test_minio_connectivity() {
    log_section "Testing MinIO Connectivity"
    
    if curl -s http://localhost:9000 >/dev/null 2>&1; then
        log_success "MinIO is accessible on localhost:9000"
    else
        log_error "MinIO is not accessible on localhost:9000"
        log_info "Please ensure MinIO is running: minio server /data"
        return 1
    fi
}

# Test configuration loading
test_config_loading() {
    log_section "Testing Configuration Loading"
    
    if [ ! -f "$CONFIG_FILE" ]; then
        log_error "Configuration file not found: $CONFIG_FILE"
        return 1
    fi
    
    log_success "Configuration file found: $CONFIG_FILE"
}

# Start the API server in background
start_api_server() {
    log_section "Starting API Server"
    
    log_info "Building Clariti server..."
    if go build -o clariti-server ./server; then
        log_success "Server built successfully"
    else
        log_error "Failed to build server"
        return 1
    fi
    
    log_info "Starting Clariti API with S3/MinIO storage..."
    CONFIG_PATH="$CONFIG_FILE" ./clariti-server &
    SERVER_PID=$!
    
    log_info "API server started with PID: $SERVER_PID"
    log_info "Waiting for server to be ready..."
    
    # Wait for server to start
    for i in {1..30}; do
        if curl -s "$API_URL/health" >/dev/null 2>&1; then
            log_success "API server is ready on $API_URL"
            return 0
        fi
        sleep 1
    done
    
    log_error "API server failed to start within 30 seconds"
    return 1
}

# Test API endpoints
test_api_endpoints() {
    log_section "Testing API Endpoints"
    
    # Test health endpoint
    log_info "Testing health endpoint..."
    if curl -s "$API_URL/health" | grep -q "healthy"; then
        log_success "Health endpoint working"
    else
        log_error "Health endpoint failed"
        return 1
    fi
    
    # Test platforms endpoint
    log_info "Testing platforms endpoint..."
    if curl -s "$API_URL/api/v1/platforms" >/dev/null 2>&1; then
        log_success "Platforms endpoint working"
    else
        log_error "Platforms endpoint failed"
        return 1
    fi
    
    # Test components endpoint
    log_info "Testing components endpoint..."
    if curl -s "$API_URL/api/v1/components" >/dev/null 2>&1; then
        log_success "Components endpoint working"
    else
        log_error "Components endpoint failed"
        return 1
    fi
}

# Test incident creation and retrieval
test_incident_crud() {
    log_section "Testing Incident CRUD Operations"
    
    # Create an incident
    log_info "Creating a test incident..."
    INCIDENT_DATA='{
        "title": "Test Incident - MinIO Storage",
        "content": "Testing incident storage with MinIO backend",
        "components": [],
        "perpetual": false,
        "criticality": 1
    }'
    
    RESPONSE=$(curl -s -X POST "$API_URL/api/v1/incidents" \
        -H "Content-Type: application/json" \
        -u "admin:password" \
        -d "$INCIDENT_DATA")
    
    if echo "$RESPONSE" | grep -q "guid"; then
        INCIDENT_GUID=$(echo "$RESPONSE" | jq -r '.guid' 2>/dev/null || echo "")
        log_success "Incident created with GUID: $INCIDENT_GUID"
    else
        log_error "Failed to create incident"
        log_info "Response: $RESPONSE"
        return 1
    fi
    
    # Retrieve the incident
    if [ -n "$INCIDENT_GUID" ]; then
        log_info "Retrieving incident $INCIDENT_GUID..."
        if curl -s "$API_URL/api/v1/incidents/$INCIDENT_GUID" | grep -q "$INCIDENT_GUID"; then
            log_success "Incident retrieved successfully"
        else
            log_error "Failed to retrieve incident"
            return 1
        fi
    fi
    
    # List all incidents
    log_info "Listing all incidents..."
    if curl -s "$API_URL/api/v1/incidents" | grep -q "$INCIDENT_GUID"; then
        log_success "Incident appears in list"
    else
        log_error "Incident not found in list"
        return 1
    fi
}

# Cleanup function
cleanup() {
    log_section "Cleanup"
    
    if [ -n "$SERVER_PID" ]; then
        log_info "Stopping API server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
        log_success "API server stopped"
    fi
    
    # Clean up the built binary
    if [ -f "clariti-server" ]; then
        log_info "Removing server binary..."
        rm -f clariti-server
        log_success "Server binary removed"
    fi
}

# Trap cleanup on exit
trap cleanup EXIT

# Main test execution
main() {
    log_header "MinIO + Clariti API Integration Test"
    
    test_minio_connectivity || exit 1
    test_config_loading || exit 1
    start_api_server || exit 1
    sleep 2  # Give server a moment to fully initialize
    test_api_endpoints || exit 1
    test_incident_crud || exit 1
    
    log_success "All tests passed! MinIO integration is working correctly."
}

# Check if jq is available for JSON parsing
if ! command -v jq >/dev/null 2>&1; then
    log_warn "jq not found - JSON parsing will be limited"
fi

main "$@"
