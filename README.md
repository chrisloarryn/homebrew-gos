# GOS - Enhanced Go Version Manager CLI

GOS is a comprehensive command-line tool for managing Go versions built with Cobra. It integrates all the functionality from the shell scripts in the `scripts/` directory into a single, powerful CLI application with enhanced environment management and automatic configuration.

## 🌟 Enhanced Features

- 🔧 **Setup**: Install and configure the 'g' Go version manager with automatic environment setup
- 📦 **Install**: Install specific Go versions
- 🔄 **Switch**: Switch between installed Go versions with verification
- 📋 **List**: View installed and available Go versions
- 🗑️ **Clean**: Deep clean all Go installations and configurations
- 📊 **Status**: Show comprehensive system status with environment validation
- 📁 **Project**: Configure Go version for specific projects
- 🚀 **Latest**: Install and use the latest Go version
- 🌍 **Environment**: Advanced environment configuration management
- 🔄 **Reload**: Reload and verify Go environment configuration
- ✅ **Auto-verification**: Automatic verification of installations and configurations

## 🎯 Key Improvements

- **Automatic GOPATH/GOROOT Configuration**: Sets up proper Go environment variables
- **Environment Validation**: Checks and fixes common configuration issues
- **Enhanced Diagnostics**: Detailed status reporting with color-coded output
- **Cross-platform Support**: Works on macOS, Linux, and Windows
- **Smart PATH Management**: Automatically configures PATH with Go binaries
- **Session Reload**: Apply environment changes without restarting terminal

## Installation

### Prerequisites

- macOS, Linux, or Windows (Git Bash/WSL)
- `curl` command available
- `git` command available (optional, for alternative installation methods)

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
sudo mv gos /usr/local/bin/
```

Or add to your PATH:
```bash
export PATH=$PATH:$(pwd)
```

## Quick Start

### 1. Setup the Go Version Manager (Enhanced!)

First, setup the 'g' version manager with automatic configuration:

```bash
gos setup
```

This enhanced setup will:
- Download and install the 'g' version manager
- Configure GOPATH and GOROOT environment variables
- Add Go binaries to your PATH
- Install and activate the latest stable Go version
- Verify the installation automatically
- Create helpful aliases and scripts

### 2. Verify Installation

```bash
# Check if everything is configured correctly
gos status

# Or check detailed environment configuration
gos env
```

### 3. Reload Environment (if needed)

```bash
# Reload your shell configuration
source ~/.zshrc

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

# Show comprehensive system status
gos status

# Fix any environment issues
gos env --fix
```

## 🚀 Enhanced Commands

### Environment Management (NEW!)

```bash
# Check environment configuration
gos env                     # Show detailed environment status
gos env --fix              # Fix common configuration issues
gos env --export           # Export environment variables for sourcing

# Reload environment
gos reload                 # Reload and verify Go environment

# Apply environment to current session
eval $(gos env --export)   # Quick environment setup
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
gos list --remote           # List available remote versions

# Remove versions
gos remove 1.20.10          # Remove specific version
```

### System Management

```bash
# Setup and configuration
gos setup                  # Complete setup with verification
gos status                 # Enhanced system status with validation
gos clean                  # Deep clean all installations

# Diagnostics and troubleshooting
gos env                    # Check environment configuration
gos reload                 # Fix PATH and environment issues
```

### Project Management

```bash
# Configure version for current project
gos project 1.21.5          # Creates .go-version file and switches to version
```

### Help

```bash
gos help                    # Show general help
gos help [command]          # Show help for specific command
```

## What Each Command Does

### `gos setup`
Replicates the functionality of `scripts/setup-go-version-manager.sh`:
- Downloads and installs the 'g' version manager
- Configures environment variables in shell configuration files
- Installs the latest stable Go version
- Creates helpful scripts and documentation

### `gos clean`
Replicates the functionality of `scripts/deep-clean-go.sh`:
- Cleans Go cache and modules
- Removes Homebrew Go installations
- Removes manual system installations
- Cleans user directories with special permission handling
- Removes other Go version managers (gvm, goenv, etc.)
- Cleans shell configuration files
- Creates backups of configuration files

### Version Management Commands
Replicate the functionality of `scripts/go-version-switcher.sh`:
- Install specific Go versions using 'g'
- Switch between installed versions
- List installed and available versions
- Remove unused versions
- Project-specific version configuration
- Comprehensive status reporting

## Environment Variables

After setup, these environment variables will be configured:

```bash
export GOPATH=$HOME/go
export GOROOT=$HOME/.g/go
export PATH=$HOME/.g/bin:$GOROOT/bin:$GOPATH/bin:$PATH
```

## Project Structure

```
.
├── main.go              # Main entry point
├── go.mod               # Go module definition
├── cmd/                 # Cobra commands
│   ├── root.go          # Root command and CLI setup
│   ├── version.go       # Version management commands (install, use, list, remove, latest, project)
│   ├── setup.go         # Setup command for installing 'g'
│   ├── clean.go         # Deep clean command
│   └── status.go        # Status command
└── scripts/             # Original shell scripts (reference)
    ├── deep-clean-go.sh
    ├── go-version-switcher.sh
    └── setup-go-version-manager.sh
```

## Examples

### Complete Workflow

```bash
# 1. Setup the environment
gos setup

# 2. Reload shell
source ~/.zshrc

# 3. Install specific versions
gos install 1.21.5
gos install 1.20.10

# 4. Switch between versions
gos use 1.21.5
go version  # Verify

# 5. Configure project
cd my-go-project
gos project 1.20.10  # Creates .go-version file

# 6. Check status
gos status

# 7. List all versions
gos list

# 8. Clean up when needed
gos clean
```

### Development Workflow

```bash
# For a new project requiring Go 1.21.5
mkdir my-project && cd my-project
gos project 1.21.5
go mod init my-project

# Switch to latest for general development  
gos latest

# Check what's installed
gos list
gos status
```

## Compatibility

- **macOS**: Full support (Intel and Apple Silicon)
- **Linux**: Full support
- **Windows**: Support via Git Bash or WSL

## Troubleshooting

### Command not found after setup
```bash
# Reload your shell configuration
source ~/.zshrc  # or ~/.bashrc

# Or open a new terminal window
```

### Permission issues during cleanup
The `gos clean` command handles permission issues automatically and will use `sudo` when necessary.

### 'g' manager not found
```bash
# Run setup again
gos setup

# Verify installation
which g
g --version
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License.
