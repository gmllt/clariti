#!/bin/bash
# log.sh - Shared logging library for all development scripts
# Usage: source this file in other scripts to get colored logging functions

# Color definitions
readonly LOG_RED='\033[0;31m'
readonly LOG_GREEN='\033[0;32m'
readonly LOG_YELLOW='\033[1;33m'
readonly LOG_BLUE='\033[0;34m'
readonly LOG_PURPLE='\033[0;35m'
readonly LOG_CYAN='\033[0;36m'
readonly LOG_WHITE='\033[1;37m'
readonly LOG_GRAY='\033[0;90m'
readonly LOG_NC='\033[0m' # No Color

# Log level definitions
readonly LOG_LEVEL_DEBUG=0
readonly LOG_LEVEL_INFO=1
readonly LOG_LEVEL_WARN=2
readonly LOG_LEVEL_ERROR=3

# Default log level (can be overridden by setting LOG_LEVEL environment variable)
LOG_LEVEL=${LOG_LEVEL:-$LOG_LEVEL_INFO}

# Timestamp format
log_timestamp() {
    date '+%H:%M:%S'
}

# Core logging function
log_message() {
    local level="$1"
    local color="$2"
    local tag="$3"
    local message="$4"
    local timestamp=$(log_timestamp)
    
    if [ "$level" -ge "$LOG_LEVEL" ]; then
        echo -e "${LOG_GRAY}[$timestamp]${LOG_NC} ${color}[$tag]${LOG_NC} $message" >&2
    fi
}

# Debug logging (gray/dimmed)
log_debug() {
    log_message $LOG_LEVEL_DEBUG "$LOG_GRAY" "DEBUG" "$1"
}

# Info logging (blue)
log_info() {
    log_message $LOG_LEVEL_INFO "$LOG_BLUE" "INFO" "$1"
}

# Success logging (green)
log_success() {
    log_message $LOG_LEVEL_INFO "$LOG_GREEN" "SUCCESS" "$1"
}

# Warning logging (yellow)
log_warn() {
    log_message $LOG_LEVEL_WARN "$LOG_YELLOW" "WARNING" "$1"
}

# Error logging (red)
log_error() {
    log_message $LOG_LEVEL_ERROR "$LOG_RED" "ERROR" "$1"
}

# Section headers (purple/cyan)
log_section() {
    log_message $LOG_LEVEL_INFO "$LOG_PURPLE" "SECTION" "$1"
}

# Command execution logging (cyan)
log_command() {
    log_message $LOG_LEVEL_INFO "$LOG_CYAN" "CMD" "$1"
}

# Status logging (white/bold)
log_status() {
    log_message $LOG_LEVEL_INFO "$LOG_WHITE" "STATUS" "$1"
}

# Specialized logging functions for common use cases

# Vendor operations
log_vendor() {
    log_message $LOG_LEVEL_INFO "$LOG_BLUE" "VENDOR" "$1"
}

# Build operations
log_build() {
    log_message $LOG_LEVEL_INFO "$LOG_CYAN" "BUILD" "$1"
}

# Test operations
log_test() {
    log_message $LOG_LEVEL_INFO "$LOG_PURPLE" "TEST" "$1"
}

# Benchmark operations
log_benchmark() {
    log_message $LOG_LEVEL_INFO "$LOG_YELLOW" "BENCHMARK" "$1"
}

# Lint operations
log_lint() {
    log_message $LOG_LEVEL_INFO "$LOG_GREEN" "LINT" "$1"
}

# Analysis operations
log_analysis() {
    log_message $LOG_LEVEL_INFO "$LOG_CYAN" "ANALYSIS" "$1"
}

# Progress indicators
log_progress() {
    local current="$1"
    local total="$2"
    local task="$3"
    log_message $LOG_LEVEL_INFO "$LOG_BLUE" "PROGRESS" "[$current/$total] $task"
}

# File operations
log_file() {
    local operation="$1"
    local file="$2"
    log_message $LOG_LEVEL_INFO "$LOG_GRAY" "FILE" "$operation: $file"
}

# Network operations
log_network() {
    log_message $LOG_LEVEL_INFO "$LOG_CYAN" "NETWORK" "$1"
}

# Configuration logging
log_config() {
    log_message $LOG_LEVEL_INFO "$LOG_BLUE" "CONFIG" "$1"
}

# Utility functions

# Print a separator line
log_separator() {
    echo -e "${LOG_GRAY}================================================${LOG_NC}" >&2
}

# Print a header with separator
log_header() {
    local title="$1"
    echo >&2
    log_separator
    log_section "$title"
    log_separator
}

# Print configuration summary
log_config_summary() {
    local title="$1"
    shift
    
    log_section "$title Configuration:"
    while [ $# -gt 0 ]; do
        local key="$1"
        local value="$2"
        echo -e "  ${LOG_GRAY}•${LOG_NC} ${LOG_BLUE}$key:${LOG_NC} $value" >&2
        shift 2
    done
}

# Print a list of items
log_list() {
    local title="$1"
    shift
    
    log_section "$title:"
    for item in "$@"; do
        echo -e "  ${LOG_GRAY}-${LOG_NC} $item" >&2
    done
}

# Conditional logging based on command success
log_result() {
    local command="$1"
    local success_msg="$2"
    local error_msg="$3"
    
    if eval "$command" >/dev/null 2>&1; then
        log_success "$success_msg"
        return 0
    else
        log_error "$error_msg"
        return 1
    fi
}

# Execute command with logging
log_exec() {
    local cmd="$1"
    local description="$2"
    
    log_command "Executing: $description"
    log_debug "Command: $cmd"
    
    if eval "$cmd"; then
        log_success "Completed: $description"
        return 0
    else
        log_error "Failed: $description"
        return 1
    fi
}

# Spinner for long operations (optional)
log_spinner() {
    local pid=$1
    local message="$2"
    local spin='-\|/'
    local i=0
    
    echo -n -e "${LOG_BLUE}[WAIT]${LOG_NC} $message " >&2
    
    while kill -0 $pid 2>/dev/null; do
        i=$(( (i+1) %4 ))
        printf "\r${LOG_BLUE}[WAIT]${LOG_NC} $message ${spin:$i:1}" >&2
        sleep 0.1
    done
    
    printf "\r${LOG_GREEN}[DONE]${LOG_NC} $message ✓\n" >&2
}

# Enable debug mode if DEBUG environment variable is set
if [ "${DEBUG:-}" = "true" ] || [ "${DEBUG:-}" = "1" ]; then
    LOG_LEVEL=$LOG_LEVEL_DEBUG
    log_debug "Debug logging enabled"
fi

# Export functions so they're available in subshells
export -f log_message log_debug log_info log_success log_warn log_error
export -f log_section log_command log_status log_vendor log_build log_test
export -f log_benchmark log_lint log_analysis log_progress log_file
export -f log_network log_config log_separator log_header log_config_summary
export -f log_list log_result log_exec log_timestamp

# Export color constants
export LOG_RED LOG_GREEN LOG_YELLOW LOG_BLUE LOG_PURPLE LOG_CYAN LOG_WHITE LOG_GRAY LOG_NC
