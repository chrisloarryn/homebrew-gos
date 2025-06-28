#!/bin/bash
# Enhanced Example usage script for GOS CLI

echo "🚀 GOS CLI - Enhanced Go Version Manager"
echo "========================================"
echo ""

# Build the CLI if it doesn't exist
if [ ! -f "./gos" ]; then
    echo "📦 Building GOS CLI..."
    go build -o gos main.go
    echo "✅ Build complete!"
    echo ""
fi

echo "📊 1. Checking system status:"
./gos status
echo ""

echo "🌍 2. Checking environment configuration:"
./gos env
echo ""

echo "� 3. Reloading environment:"
./gos reload
echo ""

echo "�📋 4. Listing installed versions:"
./gos list
echo ""

echo "🆘 5. Showing help for new commands:"
echo "   Environment commands:"
./gos help env
echo ""
echo "   Reload command:"
./gos help reload
echo ""

echo "✅ Enhanced features demo complete!"
echo ""
echo "💡 NEW: Enhanced environment management:"
echo "   ./gos env         # Check environment"
echo "   ./gos env --fix   # Fix issues"
echo "   ./gos env --export # Export for sourcing"
echo "   ./gos reload      # Reload environment"
echo ""
echo "💡 Setup and usage:"
echo "   ./gos setup       # Auto-configure everything"
echo "   ./gos install 1.21.5 # Install specific version"
echo "   ./gos use 1.21.5     # Switch versions"
echo "   ./gos clean          # Deep clean"
echo ""
echo "🎯 Quick setup workflow:"
echo "   ./gos setup && source ~/.zshrc && ./gos status"
