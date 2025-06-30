$destDir = "$HOME\tools\bin"
$pathUser = [Environment]::GetEnvironmentVariable("Path", "User")

if (-not ($pathUser.Split(";") -contains $destDir)) {
    [Environment]::SetEnvironmentVariable("Path", "$pathUser;$destDir", "User")
    Write-Host "✅ PATH actualizado. Reinicia PowerShell."
} else {
    Write-Host "ℹ️ Ya estaba en el PATH."
}
