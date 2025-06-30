# Test script for Windows setup
Write-Host "Testing gos Windows setup..." -ForegroundColor Blue
Write-Host "======================================" -ForegroundColor Cyan

# Check if gos.exe exists
if (!(Test-Path ".\gos_windows.exe")) {
    Write-Host "❌ gos_windows.exe not found. Please build it first:" -ForegroundColor Red
    Write-Host "   go build -o gos_windows.exe ." -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Found gos_windows.exe" -ForegroundColor Green

# Show current Go status (if any)
Write-Host "`n📊 Current Go status:" -ForegroundColor Blue
try {
    $goVersion = go version 2>$null
    if ($goVersion) {
        Write-Host "✅ Go is already installed: $goVersion" -ForegroundColor Green
    } else {
        Write-Host "ℹ️  Go is not currently installed or not in PATH" -ForegroundColor Yellow
    }
} catch {
    Write-Host "ℹ️  Go is not currently installed or not in PATH" -ForegroundColor Yellow
}

# Check for existing version managers
Write-Host "`n🔍 Checking for existing version managers:" -ForegroundColor Blue

try {
    $gobrewVersion = gobrew version 2>$null
    if ($gobrewVersion) {
        Write-Host "✅ gobrew is already installed: $gobrewVersion" -ForegroundColor Green
    }
} catch {
    Write-Host "ℹ️  gobrew not found" -ForegroundColor Gray
}

try {
    $gVersion = g version 2>$null
    if ($gVersion) {
        Write-Host "✅ voidint/g is already installed: $gVersion" -ForegroundColor Green
    }
} catch {
    Write-Host "ℹ️  voidint/g not found" -ForegroundColor Gray
}

# Run the setup command
Write-Host "`n🚀 Running gos setup (first time)..." -ForegroundColor Blue
Write-Host "======================================" -ForegroundColor Cyan

# Simulate user input by piping "Y" to the setup command
"Y" | .\gos_windows.exe setup

Write-Host "`n🔄 Running gos setup (second time to test duplicate detection)..." -ForegroundColor Blue
Write-Host "======================================" -ForegroundColor Cyan

# Run setup again to test if it detects existing installation
.\gos_windows.exe setup

Write-Host "`n💪 Testing force reinstallation flag..." -ForegroundColor Blue
Write-Host "======================================" -ForegroundColor Cyan

# Test the --force flag (but don't actually run it to avoid reinstalling)
Write-Host "ℹ️  Force flag test: .\gos_windows.exe setup --force" -ForegroundColor Yellow
Write-Host "   (Skipped to avoid actual reinstallation)" -ForegroundColor Gray

Write-Host "`n======================================" -ForegroundColor Cyan
Write-Host "🎉 Setup completed!" -ForegroundColor Green

# Post-setup checks
Write-Host "`n📋 Post-setup verification:" -ForegroundColor Blue

# Check if Go is now available
try {
    $goVersion = go version 2>$null
    if ($goVersion) {
        Write-Host "✅ Go is now available: $goVersion" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Go command not found - you may need to restart your terminal" -ForegroundColor Yellow
    }
} catch {
    Write-Host "⚠️  Go command not found - you may need to restart your terminal" -ForegroundColor Yellow
}

# Check version managers
try {
    $gobrewVersion = gobrew version 2>$null
    if ($gobrewVersion) {
        Write-Host "✅ gobrew is available: $gobrewVersion" -ForegroundColor Green
    }
} catch {
    # Ignore
}

try {
    $gVersion = g version 2>$null
    if ($gVersion) {
        Write-Host "✅ voidint/g is available: $gVersion" -ForegroundColor Green
    }
} catch {
    # Ignore
}

Write-Host "`n💡 Next steps:" -ForegroundColor Yellow
Write-Host "1. Restart your PowerShell/Terminal" -ForegroundColor White
Write-Host "2. Run: go version" -ForegroundColor White
Write-Host "3. Run: gos status" -ForegroundColor White

Write-Host "`n🎯 Test completed!" -ForegroundColor Green