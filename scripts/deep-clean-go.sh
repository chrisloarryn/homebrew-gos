#!/usr/bin/env bash
# deep-clean-go.sh
#
# Script for complete and aggressive Go cleanup
# Handles special permissions and completely cleans all traces of Go

set -euo pipefail

echo "🗑️  Complete Go system cleanup..."

# Function to recursively change permissions
fix_permissions() {
    local dir="$1"
    if [[ -d "$dir" ]]; then
        echo "  Fixing permissions in $dir..."
        find "$dir" -type f -exec chmod +w {} \; 2>/dev/null || true
        find "$dir" -type d -exec chmod +w {} \; 2>/dev/null || true
    fi
}

echo -e "\n▸ Cleaning existing Go cache and modules…"
if command -v go >/dev/null 2>&1; then
    echo "  Running go clean -modcache..."
    go clean -modcache 2>/dev/null || true
    echo "  Running go clean -cache..."
    go clean -cache 2>/dev/null || true
fi

echo -e "\n▸ Removing Homebrew installations…"
if command -v brew >/dev/null 2>&1; then
    go_formulas=$(brew list --formula 2>/dev/null | grep -E '^go(@[0-9]+(\.[0-9]+)*)?$' || true)
    if [[ -n "$go_formulas" ]]; then
        echo "$go_formulas" | while read -r f; do
            echo "  – brew uninstall $f"
            brew uninstall --ignore-dependencies --force "$f" 2>/dev/null || true
        done
    fi
fi

echo -e "\n▸ Removing manual system installations…"
sudo rm -rf /usr/local/go 2>/dev/null || true

echo -e "\n▸ Removing user directories with special permissions…"

# Clean ~/go directory
if [[ -d "$HOME/go" ]]; then
    echo "  Fixing permissions in $HOME/go..."
    fix_permissions "$HOME/go"
    rm -rf "$HOME/go" 2>/dev/null || {
        echo "  Using sudo to remove $HOME/go..."
        sudo rm -rf "$HOME/go"
    }
fi

# Clean Go cache
for cache_dir in "$HOME/.cache/go-build" "$HOME/Library/Caches/go-build"; do
    if [[ -d "$cache_dir" ]]; then
        echo "  Removing cache: $cache_dir"
        fix_permissions "$cache_dir"
        rm -rf "$cache_dir" 2>/dev/null || true
    fi
done

# Clean other directories
echo -e "\n▸ Removing other managers and directories…"
rm -rf "$HOME"/sdk/go* 2>/dev/null || true
rm -rf "$HOME/.gvm" 2>/dev/null || true
rm -rf "$HOME/.goenv" 2>/dev/null || true
rm -rf "$HOME/.g" 2>/dev/null || true  # g manager

# Clean environment variables from PATH
echo -e "\n▸ Cleaning shell configuration…"
if [[ -f "$HOME/.zshrc" ]]; then
    # Make backup
    cp "$HOME/.zshrc" "$HOME/.zshrc.backup.$(date +%Y%m%d_%H%M%S)"
    
    # Remove Go-related lines
    grep -v -E '(go/bin|GOPATH|GOROOT|\.gvm|\.goenv)' "$HOME/.zshrc" > "$HOME/.zshrc.tmp" && mv "$HOME/.zshrc.tmp" "$HOME/.zshrc"
fi

if [[ -f "$HOME/.bash_profile" ]]; then
    cp "$HOME/.bash_profile" "$HOME/.bash_profile.backup.$(date +%Y%m%d_%H%M%S)"
    grep -v -E '(go/bin|GOPATH|GOROOT|\.gvm|\.goenv)' "$HOME/.bash_profile" > "$HOME/.bash_profile.tmp" && mv "$HOME/.bash_profile.tmp" "$HOME/.bash_profile"
fi

# Clear command hash
hash -r 2>/dev/null || true

echo -e "\n✅ Complete Go cleanup finished."
echo "📋 Backups of your configuration files were created."
echo "🔄 Run 'source ~/.zshrc' or open a new terminal."
