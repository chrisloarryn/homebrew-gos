# Windows Support for gos

## Overview

The `gos` tool now includes native Windows support! Since the original 'g' version manager doesn't officially support Windows, we've implemented an intelligent approach that automatically installs and configures Windows-compatible Go version managers.

## Supported Windows Version Managers

### 1. gobrew (Recommended)
- **Best option for Windows** - Modern, fast, and reliable
- Written in Go with native Windows support
- Supports PowerShell and Command Prompt
- GitHub: https://github.com/kevincobain2000/gobrew

### 2. voidint/g (Alternative)
- Enhanced version of 'g' with Windows support
- Similar commands to the original 'g'
- Works in PowerShell and Git Bash
- GitHub: https://github.com/voidint/g

### 3. Direct Go Installation (Fallback)
- Downloads official Go binaries directly
- Manual version management through gos
- Works when other managers fail

## Installation Process

The setup command will automatically:

1. **Detect Windows** and show available options
2. **Try gobrew first** (recommended)
3. **Fall back to voidint/g** if gobrew fails
4. **Attempt direct installation** if both fail
5. **Show manual options** if automatic installation fails

```powershell
.\gos.exe setup
```

## What Gets Installed

### With gobrew:
- Go version manager at `%USERPROFILE%\.gobrew\`
- Current Go version symlinked at `%USERPROFILE%\.gobrew\current\`
- PowerShell profile configuration
- Git Bash profile configuration (if available)

### With voidint/g:
- Go version manager at `%USERPROFILE%\.g\`
- Go installations managed in `%USERPROFILE%\.g\go\`
- Environment configuration for both PowerShell and Bash

## Commands After Installation

### Using gobrew:
```powershell
# Version management
gobrew use latest          # Use latest Go version
gobrew use 1.21.6         # Use specific version
gobrew install 1.20.12    # Install specific version
gobrew ls                 # List installed versions
gobrew ls-remote          # List available versions

# Combined with gos
gos status                # Show current status
gos env                   # Show environment details
```

### Using voidint/g:
```powershell
# Version management  
g install latest          # Install latest Go version
g use 1.21.6             # Switch to specific version
g ls                     # List installed versions
g ls-remote              # List available versions

# Combined with gos
gos status               # Show current status
gos list                 # List versions (if compatible)
```

## Manual Installation Alternatives

If automatic installation fails, you have several options:

### Option 1: Chocolatey (Recommended)
```powershell
# Install Chocolatey (if not already installed)
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Go
choco install golang
```

### Option 2: Scoop
```powershell
# Install Scoop (if not already installed)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Install Go
scoop install go
```

### Option 3: Official Installer
1. Visit https://golang.org/dl/
2. Download the Windows installer (.msi)
3. Run the installer

### Option 4: WSL (Windows Subsystem for Linux)
```powershell
# Install WSL
wsl --install

# Use Linux version of gos in WSL
```

## Environment Configuration

### PowerShell Profile
Configuration is automatically added to:
```
%USERPROFILE%\Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1
```

### Git Bash Profile
Configuration is automatically added to:
```
%USERPROFILE%\.bashrc
```

### Manual Environment Setup

If you need to manually configure:

#### For gobrew:
```powershell
# PowerShell
$env:PATH = "$HOME\.gobrew\current\bin;$HOME\.gobrew\bin;" + $env:PATH
$env:GOROOT = "$HOME\.gobrew\current\go"
$env:GOPATH = "$HOME\go"
```

```bash
# Git Bash
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
export GOROOT="$HOME/.gobrew/current/go"
export GOPATH="$HOME/go"
```

#### For voidint/g:
```powershell
# PowerShell
$env:GOROOT = "$HOME\.g\go"
$env:PATH = "$HOME\.g\bin;$env:GOROOT\bin;" + $env:PATH
$env:GOPATH = "$HOME\go"
```

```bash
# Git Bash
export GOROOT="$HOME/.g/go"
export PATH="$HOME/.g/bin:$GOROOT/bin:$PATH"
export GOPATH="$HOME/go"
```

## Troubleshooting

### PowerShell Execution Policy
If you encounter execution policy errors:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Command Not Found
1. **Restart your terminal** completely
2. **Check environment variables**:
   ```powershell
   echo $env:GOROOT
   echo $env:GOPATH
   echo $env:PATH
   ```
3. **Manually reload profile**:
   ```powershell
   . $PROFILE
   ```

### Download Issues
- **Check internet connection**
- **Verify antivirus isn't blocking downloads**
- **Try running PowerShell as Administrator**
- **Check Windows firewall settings**

### Version Manager Not Working
1. **Verify installation**:
   ```powershell
   gobrew version  # or
   g version
   ```
2. **Check PATH configuration**
3. **Try manual installation** of the version manager

### gos Commands Not Working
Some gos commands may not work perfectly with Windows version managers since they expect the original 'g'. Use the native version manager commands instead:

- Instead of `gos install X.Y.Z`, use `gobrew use X.Y.Z` or `g install X.Y.Z`
- Instead of `gos use X.Y.Z`, use `gobrew use X.Y.Z` or `g use X.Y.Z`
- `gos status` and `gos env` should work fine

## Requirements

- **Windows 10 or later**
- **PowerShell 5.0+** (usually pre-installed)
- **Internet connection** for downloads
- **~500MB free space** for Go installations

## Performance Notes

- **gobrew** is generally faster and more reliable on Windows
- **voidint/g** is a good alternative if gobrew doesn't work
- **Direct installation** is the most basic fallback option

## Support

For Windows-specific issues:
1. Check PowerShell version: `$PSVersionTable.PSVersion`
2. Verify execution policy: `Get-ExecutionPolicy`
3. Test internet connectivity: `Test-NetConnection golang.org -Port 443`
4. Check available disk space

For general gos issues, refer to the main README.md

## Examples

### Complete Installation Flow
```powershell
# Download and run gos setup
.\gos.exe setup

# Follow prompts (choose Y to continue)

# After installation, restart PowerShell and test:
go version
gobrew version  # or g version

# Install and switch Go versions:
gobrew use 1.21.6
go version

# List available versions:
gobrew ls-remote
```

This approach gives you the best of both worlds: the convenience of gos with the reliability of proven Windows Go version managers!
