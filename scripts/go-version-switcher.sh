#!/usr/bin/env bash
# go-version-switcher.sh
#
# Advanced script for managing multiple Go versions
# Works with the 'g' manager and provides extra functionality

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show help
show_help() {
    echo -e "${BLUE}🐹 Go Version Switcher${NC}"
    echo ""
    echo "Usage: $0 [command] [arguments]"
    echo ""
    echo "Available commands:"
    echo "  install <version>     Install a specific Go version"
    echo "  use <version>         Switch to a specific version"
    echo "  list                  List installed versions"
    echo "  list-remote           List available remote versions"
    echo "  current               Show current version"
    echo "  remove <version>      Remove a specific version"
    echo "  latest                Install and use latest version"
    echo "  project [version]     Configure version for current project"
    echo "  cleanup               Clean unused versions"
    echo "  status                Show system status"
    echo "  help                  Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 install 1.21.5     # Install Go 1.21.5"
    echo "  $0 use 1.21.5         # Switch to Go 1.21.5"
    echo "  $0 latest             # Install latest version"
    echo "  $0 project 1.20.10    # Configure project to use Go 1.20.10"
}

# Check if 'g' is installed
check_g_installed() {
    if ! command -v g >/dev/null 2>&1; then
        echo -e "${RED}❌ Error: The 'g' manager is not installed.${NC}"
        echo -e "${YELLOW}💡 Run first: ./setup-go-version-manager.sh${NC}"
        exit 1
    fi
}

# Function to install a version
install_version() {
    local version="$1"
    echo -e "${BLUE}📦 Installing Go ${version}...${NC}"
    
    if g install "$version"; then
        echo -e "${GREEN}✅ Go ${version} installed successfully${NC}"
    else
        echo -e "${RED}❌ Error installing Go ${version}${NC}"
        return 1
    fi
}

# Function to switch version
use_version() {
    local version="$1"
    echo -e "${BLUE}🔄 Switching to Go ${version}...${NC}"
    
    if g use "$version"; then
        echo -e "${GREEN}✅ Switched to Go ${version}${NC}"
        echo -e "${BLUE}📋 Current version:${NC} $(go version)"
    else
        echo -e "${RED}❌ Error switching to Go ${version}${NC}"
        echo -e "${YELLOW}💡 Is this version installed? Use: $0 list${NC}"
        return 1
    fi
}

# Function to list installed versions
list_versions() {
    echo -e "${BLUE}📋 Installed Go versions:${NC}"
    g list 2>/dev/null || echo -e "${YELLOW}No versions installed${NC}"
}

# Function to list remote versions
list_remote_versions() {
    echo -e "${BLUE}🌐 Available versions (latest 10):${NC}"
    g list-all 2>/dev/null | head -10 || echo -e "${YELLOW}Could not get remote versions${NC}"
}

# Function to show current version
show_current() {
    echo -e "${BLUE}📍 Current Go version:${NC}"
    if command -v go >/dev/null 2>&1; then
        go version
        echo -e "${BLUE}📂 GOROOT:${NC} $(go env GOROOT)"
        echo -e "${BLUE}📂 GOPATH:${NC} $(go env GOPATH)"
    else
        echo -e "${YELLOW}⚠️  Go is not available in PATH${NC}"
    fi
}

# Function to remove a version
remove_version() {
    local version="$1"
    echo -e "${YELLOW}🗑️  Removing Go ${version}...${NC}"
    
    if g remove "$version"; then
        echo -e "${GREEN}✅ Go ${version} removed successfully${NC}"
    else
        echo -e "${RED}❌ Error removing Go ${version}${NC}"
        return 1
    fi
}

# Function to install latest version
install_latest() {
    echo -e "${BLUE}🚀 Installing latest Go version...${NC}"
    
    if g install latest; then
        echo -e "${GREEN}✅ Latest version installed${NC}"
        g use latest
        echo -e "${BLUE}📋 Current version:${NC} $(go version)"
    else
        echo -e "${RED}❌ Error installing latest version${NC}"
        return 1
    fi
}

# Function to configure project version
setup_project_version() {
    local version="$1"
    local go_version_file=".go-version"
    
    echo -e "${BLUE}📁 Configuring version ${version} for this project...${NC}"
    
    # Create .go-version file in current directory
    echo "$version" > "$go_version_file"
    
    # Switch to that version
    use_version "$version"
    
    echo -e "${GREEN}✅ Project configured to use Go ${version}${NC}"
    echo -e "${BLUE}📄 File created: ${go_version_file}${NC}"
}

# Function to clean unused versions
cleanup_versions() {
    echo -e "${YELLOW}🧹 Cleaning unused versions...${NC}"
    
    if g prune 2>/dev/null; then
        echo -e "${GREEN}✅ Cleanup completed${NC}"
    else
        echo -e "${YELLOW}ℹ️  No versions to clean or command not available${NC}"
    fi
}

# Function to show system status
show_status() {
    echo -e "${BLUE}📊 Go system status:${NC}"
    echo ""
    
    echo -e "${BLUE}🔧 'g' Manager:${NC}"
    if command -v g >/dev/null 2>&1; then
        echo -e "  ✅ Installed: $(g --version 2>/dev/null || echo 'unknown version')"
    else
        echo -e "  ❌ Not installed"
    fi
    
    echo ""
    echo -e "${BLUE}🐹 Current Go:${NC}"
    show_current
    
    echo ""
    echo -e "${BLUE}📦 Installed versions:${NC}"
    list_versions
    
    echo ""
    echo -e "${BLUE}💾 Disk space:${NC}"
    if [[ -d "$HOME/.g" ]]; then
        du -sh "$HOME/.g" 2>/dev/null || echo "  Could not calculate"
    else
        echo "  ~/.g directory not found"
    fi
}

# Main script
main() {
    # Check that g is installed
    check_g_installed
    
    case "${1:-help}" in
        "install")
            if [[ -z "${2:-}" ]]; then
                echo -e "${RED}❌ Error: Specify a version${NC}"
                echo "Example: $0 install 1.21.5"
                exit 1
            fi
            install_version "$2"
            ;;
        "use")
            if [[ -z "${2:-}" ]]; then
                echo -e "${RED}❌ Error: Specify a version${NC}"
                echo "Example: $0 use 1.21.5"
                exit 1
            fi
            use_version "$2"
            ;;
        "list")
            list_versions
            ;;
        "list-remote")
            list_remote_versions
            ;;
        "current")
            show_current
            ;;
        "remove")
            if [[ -z "${2:-}" ]]; then
                echo -e "${RED}❌ Error: Specify a version${NC}"
                echo "Example: $0 remove 1.20.10"
                exit 1
            fi
            remove_version "$2"
            ;;
        "latest")
            install_latest
            ;;
        "project")
            if [[ -z "${2:-}" ]]; then
                echo -e "${RED}❌ Error: Specify a version${NC}"
                echo "Example: $0 project 1.21.5"
                exit 1
            fi
            setup_project_version "$2"
            ;;
        "cleanup")
            cleanup_versions
            ;;
        "status")
            show_status
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            echo -e "${RED}❌ Unknown command: ${1}${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
