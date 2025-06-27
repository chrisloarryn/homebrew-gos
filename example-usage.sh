#!/bin/bash
# Example usage script for GOS CLI

echo "🚀 GOS CLI - Example Usage"
echo "=========================="
echo ""

# Build the CLI if it doesn't exist
if [ ! -f "./gos" ]; then
    echo "📦 Building GOS CLI..."
    go build -o gos main.go
    echo "✅ Build complete!"
    echo ""
fi

echo "1. 📊 Checking system status:"
./gos status
echo ""

echo "2. 📋 Listing installed versions:"
./gos list
echo ""

echo "3. 🆘 Showing help for install command:"
./gos help install
echo ""

echo "4. 🆘 Showing general help:"
./gos help
echo ""

echo "✅ Example usage complete!"
echo ""
echo "💡 To setup 'g' version manager (if not installed):"
echo "   ./gos setup"
echo ""
echo "💡 To install a specific Go version:"
echo "   ./gos install 1.21.5"
echo ""
echo "💡 To switch versions:"
echo "   ./gos use 1.21.5"
echo ""
echo "💡 To clean all installations:"
echo "   ./gos clean"
