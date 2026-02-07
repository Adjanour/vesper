#!/bin/bash
# Quick setup script for Vesper
# This script automates the initial setup process

set -e  # Exit on error

echo "🚀 Vesper Quick Setup Script"
echo "=============================="
echo ""

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24 or higher."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Found Go version: $GO_VERSION"

# Check if we're in the project directory
if [ ! -f "go.mod" ]; then
    echo "❌ This script must be run from the project root directory"
    exit 1
fi

# Install dependencies
echo ""
echo "📦 Installing dependencies..."
go mod download
go mod verify
echo "✓ Dependencies installed"

# Run migrations
echo ""
echo "🗄️  Setting up database..."
mkdir -p ./data
go run ./internal/database/migrate/migrate.go up
echo "✓ Database migrations applied"

# Build the project
echo ""
echo "🔨 Building the application..."
go build -o vesper ./cmd/server
echo "✓ Build complete"

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo ""
    if [ -f ".env.example" ]; then
        echo "📝 Creating .env file..."
        cp .env.example .env
        echo "✓ .env file created (you can customize it later)"
    else
        echo "⚠️  No .env.example found, skipping .env creation"
    fi
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Run 'make run' or './vesper' to start the server"
echo "  2. Test the API: curl http://localhost:8080/api/health"
echo "  3. Read API.md for available endpoints"
echo ""
echo "For development with hot reload:"
echo "  1. Install Air: go install github.com/cosmtrek/air@latest"
echo "  2. Run: make dev"
echo ""
echo "Happy coding! 🎯"
