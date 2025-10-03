#!/bin/bash

# 🔪 Kill Expenses API Instances Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get the current directory to make searches more specific
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_NAME="expenses-api"

# Functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}🔪 Killing Expenses API Instances${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to safely kill processes
safe_kill() {
    local pids="$1"
    local description="$2"
    
    if [ -n "$pids" ]; then
        print_info "Found $description, killing..."
        echo "$pids" | while read -r pid; do
            if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
                kill -TERM "$pid" 2>/dev/null || kill -9 "$pid" 2>/dev/null
            fi
        done
        return 0
    fi
    return 1
}

# Main script
print_header
echo "🔍 Searching for running instances in: $PROJECT_DIR"

KILLED_COUNT=0

# Method 1: Kill by specific binary path (most precise)
MAIN_BINARY_PIDS=$(pgrep -f "$PROJECT_DIR/main$" 2>/dev/null || true)
if safe_kill "$MAIN_BINARY_PIDS" "main binary from project directory"; then
    KILLED_COUNT=$((KILLED_COUNT + 1))
fi

# Method 2: Kill Go processes running from our project directory
GO_RUN_PIDS=$(pgrep -f "go run.*$PROJECT_DIR" 2>/dev/null || true)
if safe_kill "$GO_RUN_PIDS" "go run processes from project directory"; then
    KILLED_COUNT=$((KILLED_COUNT + 1))
fi

# Method 3: Kill by port (safest method)
PORT=${PORT:-8080}
if lsof -ti:$PORT > /dev/null 2>&1; then
    PORT_PIDS=$(lsof -ti:$PORT 2>/dev/null || true)
    if safe_kill "$PORT_PIDS" "processes using port $PORT"; then
        KILLED_COUNT=$((KILLED_COUNT + 1))
    fi
fi

# Method 4: More specific patterns (only if running from our directory)
SPECIFIC_PATTERNS=(
    "$PROJECT_DIR/main"
    "go run.*$PROJECT_NAME"
)

for pattern in "${SPECIFIC_PATTERNS[@]}"; do
    PATTERN_PIDS=$(pgrep -f "$pattern" 2>/dev/null || true)
    if safe_kill "$PATTERN_PIDS" "processes matching '$pattern'"; then
        KILLED_COUNT=$((KILLED_COUNT + 1))
    fi
done

# Wait a moment for processes to die gracefully
sleep 1

# Final verification - only check our specific patterns
REMAINING_PIDS=$(pgrep -f "$PROJECT_DIR/main$|go run.*$PROJECT_DIR" 2>/dev/null || true)
if [ -n "$REMAINING_PIDS" ]; then
    REMAINING_COUNT=$(echo "$REMAINING_PIDS" | wc -l | tr -d ' ')
else
    REMAINING_COUNT=0
fi

if [ "$REMAINING_COUNT" -gt 0 ]; then
    print_warning "Some processes might still be running, checking port..."
    # Only force kill if port is still occupied
    if lsof -ti:$PORT > /dev/null 2>&1; then
        print_info "Port $PORT still occupied, force killing..."
        lsof -ti:$PORT | xargs kill -9 2>/dev/null || true
        KILLED_COUNT=$((KILLED_COUNT + 1))
    fi
fi

# Final status
if [ "$KILLED_COUNT" -gt 0 ]; then
    print_success "Killed $KILLED_COUNT process(es). All instances stopped."
else
    print_success "No running instances found. Ready to start fresh."
fi

# Clean up any leftover files in our project directory
if [ -f "$PROJECT_DIR/main" ]; then
    rm -f "$PROJECT_DIR/main"
    print_info "Cleaned up binary file"
fi

echo ""
print_success "Ready to start new instance! 🚀"


