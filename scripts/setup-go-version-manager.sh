#!/usr/bin/env bash
# setup-go-version-manager.sh
#
# Installs and configures 'g' - A simple and fast Go version manager
# Compatible with macOS, Linux and Windows (Git Bash/WSL)
# 
# Usage after installation:
#   g install 1.21.5    # Install Go 1.21.5
#   g use 1.21.5         # Switch to Go 1.21.5
#   g list               # List installed versions
#   g list-all           # List all available versions

set -euo pipefail

# Detect operating system
case "$(uname -s)" in
    Darwin*)    OS="macOS" ;;
    Linux*)     OS="Linux" ;;
    CYGWIN*|MINGW*|MSYS*) OS="Windows" ;;
    *)          OS="Unknown" ;;
esac

echo "🔧 Installing 'g' version manager for Go..."

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "arm64" ]]; then
    if [[ "$OS" == "macOS" ]]; then
        echo "  Detected: Apple Silicon (M1/M2/M3)"
    else
        echo "  Detected: ARM64"
    fi
elif [[ "$ARCH" == "x86_64" ]]; then
    echo "  Detected: Intel x86_64"
else
    echo "  Detected: $ARCH on $OS"
fi

echo -e "\n▸ Downloading and installing 'g'..."

# Create directory for g if it doesn't exist
mkdir -p "$HOME/.g"

# Download and install g
if curl -sSL https://git.io/g-install | bash -s -- -y; then
    echo "  ✅ 'g' installed successfully"
else
    echo "  ❌ Error installing 'g'. Trying alternative method..."
    
    # Alternative method: clone from GitHub
    if command -v git >/dev/null 2>&1; then
        cd /tmp
        git clone https://github.com/stefanmaric/g.git
        cd g
        make install PREFIX="$HOME/.g"
        cd "$HOME"
        rm -rf /tmp/g
        echo "  ✅ 'g' installed via GitHub"
    else
        echo "  ❌ Git is not available. Installing manually..."
        curl -sSL https://raw.githubusercontent.com/stefanmaric/g/main/bin/g -o "$HOME/.g/bin/g"
        chmod +x "$HOME/.g/bin/g"
        echo "  ✅ 'g' installed manually"
    fi
fi

echo -e "\n▸ Configuring PATH and environment variables..."

# Configure environment variables according to system
if [[ "$OS" == "Windows" ]]; then
  G_CONFIG="
# === Go Version Manager (g) ===
export GOPATH=\$HOME/go
export GOROOT=\$HOME/.g/go
export PATH=\$HOME/.g/bin:\$GOROOT/bin:\$GOPATH/bin:\$PATH
"
  shell_files=("$HOME/.bashrc" "$HOME/.bash_profile")
else
  G_CONFIG="
# === Go Version Manager (g) ===
export GOPATH=\$HOME/go
export GOROOT=\$HOME/.g/go
export PATH=\$HOME/.g/bin:\$GOROOT/bin:\$GOPATH/bin:\$PATH
"
  shell_files=("$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile")
fi

# Add configuration to first available shell file
config_added=false
for shell_file in "${shell_files[@]}"; do
  if [[ -f "$shell_file" ]] || [[ "$shell_file" == "${shell_files[0]}" ]]; then
    if ! grep -q "Go Version Manager (g)" "$shell_file" 2>/dev/null; then
      echo "$G_CONFIG" >> "$shell_file"
      echo "  ✅ Configuration added to $shell_file"
      config_added=true
    else
      echo "  ℹ️  Configuration already exists in $shell_file"
      config_added=true
    fi
    break
  fi
done

if [[ "$config_added" == false ]]; then
  # Create default shell file if none exists
  default_shell="${shell_files[0]}"
  echo "$G_CONFIG" > "$default_shell"
  echo "  ✅ Configuration created in $default_shell"
fi

# Export variables for current session
export GOPATH="$HOME/go"
export GOROOT="$HOME/.g/go"
export PATH="$HOME/.g/bin:$GOROOT/bin:$GOPATH/bin:$PATH"

echo -e "\n▸ Installing latest stable Go version..."

# Use g to install latest version
if "$HOME/.g/bin/g" install latest 2>/dev/null; then
    echo "  ✅ Go latest installed successfully"
else
    echo "  ℹ️  Installing known specific version..."
    "$HOME/.g/bin/g" install 1.21.5 2>/dev/null || true
fi

echo -e "\n▸ Creating help script..."

# Create help script with useful commands
cat > "$HOME/.g/go-help.sh" << 'EOF'
#!/bin/bash
# Useful commands for the 'g' version manager

echo "🐹 Go Version Manager - Useful commands:"
echo ""
echo "📦 Installation:"
echo "  g install latest        # Install latest version"
echo "  g install 1.21.5        # Install specific version"
echo "  g install 1.20.x        # Install latest 1.20.x"
echo ""
echo "🔄 Version switching:"
echo "  g use 1.21.5            # Switch to specific version"
echo "  g use latest            # Switch to latest installed"
echo ""
echo "📋 Information:"
echo "  g list                  # List installed versions"
echo "  g list-all              # List all available versions"
echo "  g current               # Show current version"
echo "  go version              # Confirm active Go version"
echo ""
echo "🗑️  Cleanup:"
echo "  g remove 1.20.10        # Remove specific version"
echo "  g prune                 # Remove unused versions"
echo ""
echo "💡 Usage examples:"
echo "  g install 1.21.5 && g use 1.21.5"
echo "  g list | head -5"
echo ""
EOF

chmod +x "$HOME/.g/go-help.sh"

echo -e "\n✅ Installation completed!"
echo ""
echo "📋 Next steps:"
if [[ "$OS" == "Windows" ]]; then
  echo "1. Run: source ~/.bashrc  (or restart Git Bash/WSL)"
else
  echo "1. Run: source ~/.zshrc  (or open a new terminal)"
fi
echo "2. Verify: g --version"
echo "3. Use: g list  (to see installed versions)"
echo ""
echo "💡 To see all available commands:"
echo "   ~/.g/go-help.sh"
echo ""
echo "🚀 Quick examples:"
echo "   g install 1.21.5     # Install Go 1.21.5"
echo "   g use 1.21.5          # Switch to Go 1.21.5"
echo "   g list                # View installed versions"
echo ""
