#!/bin/bash

# Script to start development server with Air
echo "🚀 Starting development server with Air hot reloading..."
echo "📁 Working directory: $(pwd)"
echo "🔧 Air will watch for changes in Go files and restart the server automatically"
echo ""

# Check if Air is installed
if ! command -v air &> /dev/null; then
    echo "❌ Air is not installed or not in PATH"
    echo "Please install it with: go install github.com/air-verse/air@latest"
    exit 1
fi

# Check if tmp directory exists
if [ ! -d "tmp" ]; then
    echo "📁 Creating tmp directory..."
    mkdir -p tmp
fi

# Start Air
echo "🌪️  Starting Air..."
air
