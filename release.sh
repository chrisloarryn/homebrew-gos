#!/bin/bash
# release.sh - Script to create and push a new release

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show help
show_help() {
    echo -e "${BLUE}🚀 GOS Release Script${NC}"
    echo ""
    echo "Usage: $0 [version]"
    echo ""
    echo "Examples:"
    echo "  $0 1.0.0           # Create release v1.0.0"
    echo "  $0 1.0.1           # Create release v1.0.1"
    echo "  $0 --help          # Show this help"
    echo ""
    echo "This script will:"
    echo "  1. Validate the version format"
    echo "  2. Update version in go.mod if needed"
    echo "  3. Run tests"
    echo "  4. Create and push a git tag"
    echo "  5. Trigger GitHub Actions for release"
}

# Function to validate version format
validate_version() {
    local version="$1"
    if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo -e "${RED}❌ Invalid version format: $version${NC}"
        echo -e "${YELLOW}💡 Use semantic versioning: X.Y.Z (e.g., 1.0.0)${NC}"
        exit 1
    fi
}

# Function to check if tag already exists
check_tag_exists() {
    local version="$1"
    local tag="v$version"
    
    if git tag -l | grep -q "^$tag$"; then
        echo -e "${RED}❌ Tag $tag already exists${NC}"
        echo -e "${YELLOW}💡 Use a different version number${NC}"
        exit 1
    fi
}

# Function to check git status
check_git_status() {
    if [[ -n $(git status --porcelain) ]]; then
        echo -e "${RED}❌ Working directory is not clean${NC}"
        echo -e "${YELLOW}💡 Please commit or stash your changes first${NC}"
        git status --short
        exit 1
    fi
}

# Function to run tests
run_tests() {
    echo -e "${BLUE}🧪 Running tests...${NC}"
    
    if ! go mod tidy; then
        echo -e "${RED}❌ Failed to tidy go modules${NC}"
        exit 1
    fi
    
    if ! go test -v ./...; then
        echo -e "${RED}❌ Tests failed${NC}"
        exit 1
    fi
    
    if ! go build -o gos main.go; then
        echo -e "${RED}❌ Build failed${NC}"
        exit 1
    fi
    
    # Test the binary
    if ! ./gos --help > /dev/null; then
        echo -e "${RED}❌ Binary test failed${NC}"
        exit 1
    fi
    
    # Clean up
    rm -f gos
    
    echo -e "${GREEN}✅ All tests passed${NC}"
}

# Function to create and push tag
create_and_push_tag() {
    local version="$1"
    local tag="v$version"
    
    echo -e "${BLUE}🏷️  Creating tag $tag...${NC}"
    
    # Create annotated tag
    git tag -a "$tag" -m "Release $tag

🚀 GOS CLI $tag

A comprehensive Go version manager CLI built with Cobra.

See the release notes for details about what's new in this version.
"
    
    echo -e "${BLUE}📤 Pushing tag to origin...${NC}"
    git push origin "$tag"
    
    echo -e "${GREEN}✅ Tag $tag created and pushed${NC}"
}

# Function to show next steps
show_next_steps() {
    local version="$1"
    local tag="v$version"
    
    echo ""
    echo -e "${GREEN}🎉 Release process initiated!${NC}"
    echo ""
    echo -e "${BLUE}📋 Next steps:${NC}"
    echo "  1. GitHub Actions will automatically build and create the release"
    echo "  2. Check the progress at: https://github.com/cristobalcontreras/homebrew-gos/actions"
    echo "  3. Once complete, the release will be available at:"
    echo "     https://github.com/cristobalcontreras/homebrew-gos/releases/tag/$tag"
    echo ""
    echo -e "${YELLOW}💡 The release will include:${NC}"
    echo "  - Binaries for multiple platforms (macOS, Linux, Windows)"
    echo "  - Docker images"
    echo "  - Homebrew formula updates"
    echo "  - Package manager releases (deb, rpm)"
}

# Main function
main() {
    case "${1:-}" in
        "--help"|"-h"|"help")
            show_help
            exit 0
            ;;
        "")
            echo -e "${RED}❌ Version is required${NC}"
            echo ""
            show_help
            exit 1
            ;;
        *)
            local version="$1"
            ;;
    esac
    
    echo -e "${BLUE}🚀 Starting release process for version $version${NC}"
    
    # Validation steps
    validate_version "$version"
    check_tag_exists "$version"
    check_git_status
    
    # Make sure we're on main branch
    local current_branch=$(git branch --show-current)
    if [[ "$current_branch" != "main" ]]; then
        echo -e "${YELLOW}⚠️  You're on branch '$current_branch', not 'main'${NC}"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}Release cancelled${NC}"
            exit 1
        fi
    fi
    
    # Run tests
    run_tests
    
    # Create and push tag
    create_and_push_tag "$version"
    
    # Show next steps
    show_next_steps "$version"
}

# Run main function with all arguments
main "$@"
