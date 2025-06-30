# GOS - Enhanced Go Version Manager CLI

GOS is a comprehensive command-line tool for managing Go versions built with Cobra. It integrates all the functionality from the shell scripts in the `scripts/` directory into a single, powerful CLI application with enhanced environment management and automatic configuration.

## 🌟 Enhanced Features

- 🔧 **Setup**: Install and configure Go version managers with intelligent detection and automatic environment setup
- 📦 **Install**: Install specific Go versions using multiple version manager backends
- 🔄 **Switch**: Switch between installed Go versions with verification
- 📋 **List**: View installed and available Go versions (local and remote)
- 🗑️ **Clean**: Deep clean all Go installations and configurations with automated backups
- 📊 **Status**: Show comprehensive system status with environment validation
- 📁 **Project**: Configure Go version for specific projects (.go-version files)
- 🚀 **Latest**: Install and use the latest Go version automatically
- 🌍 **Environment**: Advanced environment configuration management with diagnostic tools
- 🔄 **Reload**: Reload and verify Go environment configuration
- ✅ **Auto-verification**: Automatic verification of installations and configurations
- 🪟 **Windows Support**: Native Windows support with PowerShell and Git Bash integration

## 🎯 Key Improvements

- **Multi-Platform Support**: Full support for macOS, Linux, and Windows (including PowerShell and Git Bash)
- **Intelligent Version Manager Detection**: Automatically works with gobrew, g, voidint/g, and manual installations
- **Automatic Environment Configuration**: Sets up proper Go environment variables across all platforms
- **Smart Installation Prevention**: Prevents unnecessary reinstallations with `--force` override option
- **Enhanced Diagnostics**: Detailed status reporting with color-coded output and emoji indicators
- **Cross-Platform PATH Management**: Automatically configures PATH with Go binaries on all systems
- **Session Reload**: Apply environment changes without restarting terminal
- **Comprehensive Error Handling**: Graceful fallbacks and detailed error messages

## Installation

### Via Homebrew (macOS)

```bash
# Add the tap
brew tap chrisloarryn/homebrew-gos

# Install gos
brew install --cask gos

# Run setup
gos setup
```

### Prerequisites for Building from Source

- macOS, Linux, or Windows (PowerShell, Git Bash, or WSL)
- Go 1.21 or later
- `curl` command available
- `git` command available (optional, for cloning)

### Build from Source

1. Clone this repository:
```bash
git clone https://github.com/cristobalcontreras/homebrew-gos.git
cd homebrew-gos
```

2. Build the CLI:
```bash
go mod tidy
go build -o gos main.go
```

3. Install globally (optional):
```bash
# macOS/Linux
sudo mv gos /usr/local/bin/

# Or add to your PATH
export PATH=$PATH:$(pwd)
```

4. Verify installation:
```bash
gos --help
gos version
```

## Quick Start

### 1. Setup the Go Version Manager (Enhanced!)

First, setup Go version managers with intelligent detection and automatic configuration:

```bash
gos setup
```

This enhanced setup will:
- **Detect your platform** (Windows, macOS, Linux) automatically
- **Choose the best version manager** for your system:
  - **gobrew** (recommended for Windows and cross-platform use)
  - **g** (traditional Unix-like systems)
  - **voidint/g** (Windows-compatible alternative)
- **Install latest Go version** automatically
- **Configure environment variables** (GOPATH, GOROOT, PATH)
- **Set up shell profiles** (PowerShell, Bash, Zsh)
- **Verify installation** automatically
- **Prevent duplicate installations** (use `--force` to override)

### 2. Verify Installation

```bash
# Check comprehensive system status
gos status

# Check detailed environment configuration
gos env

# Show current Go version
go version
```

### 3. Reload Environment (if needed)

```bash
# Reload your shell configuration
source ~/.zshrc   # Linux/macOS Zsh
source ~/.bashrc  # Linux/macOS Bash

# For Windows PowerShell - restart PowerShell or:
. $PROFILE

# Or use the built-in reload command
gos reload
```

### 4. Start Using GOS

```bash
# Install a specific Go version
gos install 1.21.5

# Switch to a version
gos use 1.21.5

# List installed versions
gos list

# List available remote versions
gos list --remote

# Show comprehensive system status
gos status

# Fix any environment issues
gos env --fix

# Configure project-specific version
gos project 1.21.5   # Creates .go-version file
```

## 🚀 Enhanced Commands

### Environment Management

```bash
# Check environment configuration
gos env                     # Show detailed environment status
gos env --fix              # Fix common configuration issues
gos env --export           # Export environment variables for sourcing

# Reload environment
gos reload                 # Reload and verify Go environment

# Apply environment to current session
eval $(gos env --export)   # Quick environment setup (Unix-like)
```

### Version Management

```bash
# Install Go versions
gos install 1.21.5          # Install specific version
gos install latest          # Install latest version
gos latest                  # Install and use latest version

# Switch versions
gos use 1.21.5              # Switch to specific version

# List versions
gos list                    # List installed versions
gos list --remote           # List available remote versions (last 10)

# Remove versions
gos remove 1.20.10          # Remove specific version
```

### System Management

```bash
# Setup and configuration
gos setup                  # Intelligent setup with platform detection
gos setup --force          # Force reinstallation (bypass existing detection)
gos status                 # Enhanced system status with validation
gos clean                  # Deep clean all installations with backups
gos clean --force          # Skip confirmation prompts

# Diagnostics and troubleshooting
gos env                    # Check environment configuration
gos env --fix              # Fix common PATH and environment issues
gos reload                 # Reload environment without restarting terminal
```

### Project Management

```bash
# Configure version for current project
gos project 1.21.5          # Creates .go-version file and switches to version
```

### Help and Information

```bash
gos help                    # Show general help
gos help [command]          # Show help for specific command
gos version                 # Show gos version information
gos --version               # Show version (short form)
```

## What Each Command Does

### `gos setup`
Intelligently sets up Go version management for your platform:
- **Platform Detection**: Automatically detects Windows, macOS, or Linux
- **Version Manager Installation**:
  - **Windows**: Installs gobrew (primary) or voidint/g (fallback)
  - **Unix-like**: Installs traditional 'g' manager
- **Environment Configuration**: Sets up shell profiles (PowerShell, Bash, Zsh)
- **Latest Go Installation**: Installs and activates latest stable Go version
- **Installation Prevention**: Avoids duplicate installations (use `--force` to override)
- **Verification**: Automatically verifies successful installation

### `gos clean`
Comprehensive cleanup of Go installations and configurations:
- **Go Cache and Modules**: Cleans build cache and module cache
- **Homebrew Installations**: Removes Homebrew-installed Go versions
- **Manual System Installations**: Removes system-wide manual installations
- **User Directory Cleanup**: Cleans user-specific Go directories
- **Version Manager Cleanup**: Removes other version managers (gvm, goenv, etc.)
- **Shell Configuration**: Cleans Go-related entries from shell configs
- **Backup Creation**: Creates automatic backups before cleaning

### `gos env`
Advanced environment management and diagnostics:
- **Environment Status**: Shows detailed GOROOT, GOPATH, and PATH configuration
- **Issue Detection**: Identifies common configuration problems
- **Auto-Fix**: Repairs environment issues with `--fix` flag
- **Export Support**: Generates shell-sourceable environment variables

### `gos reload`
Environment refresh without terminal restart:
- **Environment Refresh**: Reloads Go environment variables
- **PATH Updates**: Updates PATH for current session
- **Verification**: Verifies Go is accessible after reload

### Version Management Commands
Unified interface for multiple version managers:
- **Install/Switch**: Works with gobrew, g, or voidint/g automatically
- **List Management**: Shows versions from any installed version manager
- **Remote Browsing**: Lists available versions from official sources
- **Project Configuration**: Creates .go-version files for project-specific versions
- **Comprehensive Status**: Shows detailed system information with diagnostics

## Environment Variables

After setup, these environment variables will be configured automatically:

### Unix-like Systems (macOS/Linux)
```bash
export GOPATH=$HOME/go
export GOROOT=$HOME/.g/go
export PATH=$HOME/.g/bin:$GOROOT/bin:$GOPATH/bin:$PATH
```

### Windows Systems
```powershell
# PowerShell Profile
$env:GOPATH = "$HOME\go"
$env:GOROOT = "$HOME\.g\go"  # or gobrew path if using gobrew
$env:PATH = "$HOME\.g\bin;$env:GOROOT\bin;$env:GOPATH\bin;$env:PATH"
```

## Supported Version Managers

GOS intelligently works with multiple Go version managers:

### Primary Managers
- **gobrew** - Modern, cross-platform Go version manager (recommended for Windows)
- **g** - Traditional Unix-like systems version manager
- **voidint/g** - Windows-compatible alternative to original 'g'

### Detection Logic
1. **Platform Detection**: Automatically detects your operating system
2. **Manager Priority**: Uses the best available manager for your platform
3. **Fallback Support**: Falls back to alternative managers if primary fails
4. **Manual Installation**: Direct Go installation as final fallback

## Project Structure

```
.
├── main.go                    # Main entry point with version info
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── Makefile                   # Build automation
├── README.md                  # This documentation
├── LICENSE                    # MIT License
├── Dockerfile                 # Docker image configuration
├── .goreleaser.yml           # Release automation
├── .github/workflows/        # CI/CD pipelines
│   ├── ci.yml                # Continuous integration
│   └── release.yml           # Release automation
├── cmd/                      # Cobra commands
│   ├── root.go               # Root command and CLI setup
│   ├── version.go            # Version management (install, use, list, remove, latest, project)
│   ├── setup.go              # Multi-platform setup command
│   ├── clean.go              # Deep clean with backups
│   ├── status.go             # System status and diagnostics
│   ├── env.go                # Environment management
│   ├── reload.go             # Environment reload
│   └── root_test.go          # Command tests
├── scripts/                  # Original shell scripts (reference)
│   ├── deep-clean-go.sh
│   ├── go-version-switcher.sh
│   ├── setup-go-version-manager.sh
│   ├── add_to_path.ps1       # Windows PowerShell scripts
│   └── test_windows_setup.ps1
├── Casks/                    # Homebrew cask definition
│   └── gos.rb
├── build/                    # Build artifacts
├── coverage.out              # Test coverage data
└── docs/                     # Additional documentation
    ├── WINDOWS.md            # Windows-specific guide
    ├── SETUP_IMPROVEMENTS.md
    ├── WINDOWS_IMPLEMENTATION_SUMMARY.md
    ├── RELEASE_SUMMARY.md
    └── CI_FIXES.md
```

## Examples

### Complete Workflow

```bash
# 1. Setup the environment (first time)
gos setup

# 2. Reload shell (if needed)
source ~/.zshrc              # macOS/Linux Zsh
source ~/.bashrc             # macOS/Linux Bash
# For Windows: restart PowerShell or reload profile

# 3. Install specific versions
gos install 1.21.5
gos install 1.20.10
gos install latest

# 4. Switch between versions
gos use 1.21.5
go version                   # Verify current version

# 5. Configure project-specific version
cd my-go-project
gos project 1.20.10         # Creates .go-version file

# 6. Check comprehensive status
gos status

# 7. List all versions
gos list                     # Show installed
gos list --remote            # Show available

# 8. Environment management
gos env                      # Check environment
gos env --fix                # Fix issues
gos reload                   # Refresh environment

# 9. Clean up when needed
gos clean                    # Interactive cleanup
gos clean --force            # Skip confirmations
```

### Development Workflow

```bash
# Setup for new development environment
gos setup                    # One-time setup

# Create new project with specific Go version
mkdir my-project && cd my-project
gos project 1.21.5          # Set project version
go mod init my-project       # Initialize Go module

# Switch to latest for general development  
gos latest                   # Install and use latest

# Check what's available and installed
gos list                     # Local versions
gos list --remote            # Remote versions
gos status                   # Comprehensive status
```

### Windows-Specific Examples

```powershell
# PowerShell setup
gos.exe setup               # Automatic Windows setup

# After setup, use normal commands
gos install 1.21.5
gos use 1.21.5
gos status

# Environment management
gos env                     # Check Windows environment
gos reload                  # Refresh PowerShell environment

# Using with different version managers
gobrew ls                   # If gobrew was installed
g ls                        # If voidint/g was installed
```

### CI/CD Integration

```yaml
# GitHub Actions example
steps:
  - uses: actions/checkout@v4
  - name: Setup Go with gos
    run: |
      curl -L https://github.com/chrisloarryn/homebrew-gos/releases/latest/download/gos_Linux_x86_64.tar.gz | tar xz
      ./gos setup
      ./gos install 1.21.5
      ./gos use 1.21.5
  - name: Verify Go
    run: go version
```

## Compatibility and Platform Support

### Supported Platforms
- **macOS**: Full support (Intel and Apple Silicon)
  - Native 'g' version manager
  - Homebrew integration
  - Zsh and Bash shell support
- **Linux**: Full support (AMD64 and ARM64)
  - All major distributions
  - Multiple shell environments
- **Windows**: Native support
  - PowerShell integration
  - Git Bash compatibility
  - WSL (Windows Subsystem for Linux) support
  - Automatic platform detection

### Version Manager Compatibility
- **gobrew**: Cross-platform, recommended for Windows
- **g**: Traditional Unix-like systems (stefanmaric/g)
- **voidint/g**: Windows-compatible alternative
- **Manual installations**: Direct Go binary management

### Shell Support
- **Bash**: Full support on all platforms
- **Zsh**: Full support on macOS and Linux
- **PowerShell**: Native Windows support
- **Git Bash**: Windows support for Unix-like experience

## Troubleshooting

### Command not found after setup
```bash
# Unix-like systems
source ~/.zshrc              # Reload Zsh
source ~/.bashrc             # Reload Bash

# Windows PowerShell
. $PROFILE                   # Reload PowerShell profile
# Or restart PowerShell

# Verify PATH
echo $PATH                   # Unix-like
echo $env:PATH              # PowerShell

# Alternative: open new terminal window
```

### Permission issues during cleanup
```bash
# gos clean handles permissions automatically and will use sudo when necessary
gos clean

# For manual intervention
sudo gos clean              # Unix-like systems
# Run PowerShell as Administrator on Windows
```

### Version manager not found
```bash
# Run setup again to install version manager
gos setup

# Or force reinstallation
gos setup --force

# Verify installation
which g                     # Unix-like
which gobrew               # Cross-platform
where.exe gobrew           # Windows
```

### Environment issues
```bash
# Check current environment status
gos env

# Fix common environment problems
gos env --fix

# Reload environment
gos reload

# Check comprehensive status
gos status
```

### Windows-Specific Issues
```powershell
# If setup fails, try manual alternatives:
# 1. Install via Chocolatey
choco install golang

# 2. Install via Scoop
scoop install go

# 3. Download official installer from golang.org

# Verify environment after manual installation
gos env
gos status
```

### Reset Everything
```bash
# Complete cleanup and fresh start
gos clean --force           # Remove all Go installations
gos setup --force           # Fresh setup
gos install latest          # Install latest Go
gos status                  # Verify everything works
```

## Release and Distribution

### GitHub Releases
GOS is automatically built and released using GoReleaser with support for:
- **Multiple platforms**: Windows, macOS, Linux (AMD64/ARM64)
- **Package formats**: tar.gz, zip, .deb, .rpm
- **Docker images**: Multi-platform container support
- **Homebrew integration**: Automatic formula updates

### Installation Methods
1. **Homebrew** (recommended for macOS): `brew install --cask gos`
2. **Direct download**: From GitHub releases
3. **Docker**: `docker pull ghcr.io/chrisloarryn/gos`
4. **Build from source**: Clone and build locally

## Development

### Building
```bash
# Basic build
go build -o gos main.go

# Build with version info
make build

# Build for all platforms
make build-all

# Run tests
make test
```

### Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature-name`
3. **Make your changes**: Follow Go best practices
4. **Add tests**: Ensure your changes are tested
5. **Run tests**: `make test`
6. **Submit a pull request**: With clear description of changes

### Development Guidelines
- Follow Go standard formatting (`go fmt`)
- Add tests for new functionality
- Update documentation for user-facing changes
- Ensure cross-platform compatibility
- Test on multiple platforms when possible

## Advanced Usage

### Environment Variables
You can customize gos behavior with environment variables:
```bash
export GOS_DEBUG=1           # Enable debug output
export GOS_NO_COLOR=1        # Disable colored output
export GOS_CONFIG_DIR=~/.gos # Custom config directory
```

### Configuration Files
GOS respects these configuration files:
- `.go-version`: Project-specific Go version
- Shell profiles: Automatic environment setup
- Version manager configs: Works with existing setups

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- **stefanmaric/g**: Original 'g' version manager for Unix-like systems
- **kevincobain2000/gobrew**: Modern cross-platform Go version manager
- **voidint/g**: Windows-compatible alternative to original 'g'
- **spf13/cobra**: Powerful CLI framework for Go
- **fatih/color**: Terminal color support

## Links

- **GitHub Repository**: https://github.com/chrisloarryn/homebrew-gos
- **Issue Tracker**: https://github.com/chrisloarryn/homebrew-gos/issues
- **Releases**: https://github.com/chrisloarryn/homebrew-gos/releases
- **Documentation**: See the `docs/` folder for detailed guides
