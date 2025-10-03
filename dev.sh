#!/bin/bash

# 🚀 Expenses API - Development Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}🚀 Expenses API - Development${NC}"
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

# Set environment variables
export ENV=development
export PORT=8080
export DB_HOST=localhost:3306
export DB_USER=root
export DB_PASSWORD=${DB_PASSWORD:-"your_password_here"}
export DB_NAME=expenses_db
export JWT_SECRET=dev-secret
export CORS_ALLOWED_ORIGIN=http://localhost:4200

# Main script
print_header

case "${1:-run}" in
    "build")
        echo "🔨 Building application..."
        go build -o main .
        print_success "Build completed!"
        ;;
    "run")
        echo "🔪 Killing any running instances..."
        ./scripts/kill-instances.sh
        echo ""
        echo "🚀 Starting Expenses API..."
        echo "📍 Server will run on: http://localhost:$PORT"
        echo "🗄️  Database: $DB_HOST/$DB_NAME"
        echo ""
        print_warning "Make sure your database is running and configured!"
        echo ""
        go run .
        ;;
    "kill")
        echo "🔪 Killing all running instances..."
        ./scripts/kill-instances.sh
        ;;
    "debug")
        echo "🔪 Killing any running instances..."
        ./scripts/kill-instances.sh
        echo ""
        echo "🐛 Starting in debug mode..."
        export DEBUG=true
        export LOG_LEVEL=debug
        go run .
        ;;
    "test")
        echo "🧪 Running tests..."
        go test ./... -v
        ;;
    "clean")
        echo "🧹 Cleaning build artifacts..."
        rm -f main
        go clean
        print_success "Clean completed!"
        ;;
    "setup")
        echo "⚙️  Setting up development environment..."
        
        # Check if Go is installed
        if ! command -v go &> /dev/null; then
            print_error "Go is not installed!"
            exit 1
        fi
        
        # Tidy modules
        echo "📦 Tidying Go modules..."
        go mod tidy
        
        # Build to check for errors
        echo "🔨 Building to check for errors..."
        go build -o main .
        
        print_success "Setup completed!"
        print_warning "Don't forget to:"
        echo "  1. Configure your database connection"
        echo "  2. Run database setup: mysql -u root -p < sql/database/setup_database.sql"
        echo "  3. Update DB_PASSWORD in your environment"
        ;;
    "db-setup")
        echo "🗄️  Setting up database..."
        print_warning "This will create/reset the database. Continue? (y/N)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            mysql -u root -p < sql/database/setup_database.sql
            print_success "Database setup completed!"
        else
            echo "Database setup cancelled."
        fi
        ;;
    "help"|"-h"|"--help")
        echo "Available commands:"
        echo "  build     - Build the application"
        echo "  run       - Run the application (default)"
        echo "  debug     - Run with debug logging"
        echo "  kill      - Kill all running instances"
        echo "  test      - Run tests"
        echo "  clean     - Clean build artifacts"
        echo "  setup     - Setup development environment"
        echo "  db-setup  - Setup database"
        echo "  help      - Show this help"
        ;;
    *)
        print_error "Unknown command: $1"
        echo "Use './dev.sh help' for available commands"
        exit 1
        ;;
esac
