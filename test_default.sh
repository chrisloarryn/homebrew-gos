#!/bin/bash
# Test script for gos default functionality

echo "🧪 Testing gos default functionality..."
echo ""

# Set PATH to include gobrew
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"

# Check if gobrew is available now
echo "Checking if gobrew is available..."
which gobrew || echo "gobrew not found in PATH"
echo ""

# List available Go versions using gos
echo "Listing available Go versions..."
./gos list
echo ""

# Try to install the latest version if none is installed
echo "Installing latest Go version if needed..."
./gos latest
echo ""

# List versions again to confirm installation
echo "Listing available Go versions after installation..."
./gos list
echo ""

# Use a specific version (1.21.5 or whatever is available)
echo "Setting a specific version as default..."
./gos default 1.21.5 || ./gos default latest
echo ""

# Check if default version is set correctly
echo "Checking default version..."
./gos default
echo ""

# Verify Go version
echo "Verifying Go version..."
go version
echo ""

echo "✅ Test completed"
