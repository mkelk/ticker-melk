# Ticker installation script for Windows
# Install: irm https://raw.githubusercontent.com/mkelk/ticker-melk/main-melk/scripts/install.ps1 | iex
# Upgrade: ticker upgrade (or re-run the install script)

$ErrorActionPreference = "Stop"

$Repo = "mkelk/ticker-melk"
$BinaryName = "ticker"

function Write-Info($msg) {
    Write-Host "info: " -ForegroundColor Blue -NoNewline
    Write-Host $msg
}

function Write-Success($msg) {
    Write-Host "success: " -ForegroundColor Green -NoNewline
    Write-Host $msg
}

function Write-Warning($msg) {
    Write-Host "warn: " -ForegroundColor Yellow -NoNewline
    Write-Host $msg
}

function Write-Error($msg) {
    Write-Host "error: " -ForegroundColor Red -NoNewline
    Write-Host $msg
    exit 1
}

function Get-Architecture {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    switch ($arch) {
        "X64" { return "amd64" }
        "Arm64" { return "arm64" }
        default { Write-Error "Unsupported architecture: $arch" }
    }
}

function Get-InstallDir {
    # Check for custom install dir
    if ($env:TICKER_INSTALL_DIR) {
        return $env:TICKER_INSTALL_DIR
    }

    # Default to user's local bin directory
    $localBin = Join-Path $env:LOCALAPPDATA "Programs\ticker"
    if (-not (Test-Path $localBin)) {
        New-Item -ItemType Directory -Path $localBin -Force | Out-Null
    }
    return $localBin
}

function Get-LatestVersion {
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Headers @{ "User-Agent" = "ticker-installer" }
        return $release.tag_name
    } catch {
        Write-Error "Could not determine latest version. Check your internet connection or GitHub API limits."
    }
}

function Get-InstalledVersion {
    try {
        $output = & $BinaryName --version 2>$null
        if ($output -match '(\d+\.\d+\.\d+)') {
            return $matches[1]
        }
    } catch {
        # Not installed
    }
    return $null
}

function Install-Ticker {
    Write-Info "Installing $BinaryName..."

    $arch = Get-Architecture
    Write-Info "Detected platform: windows/$arch"

    # Get versions
    $latestVersion = Get-LatestVersion
    Write-Info "Latest version: $latestVersion"

    $installedVersion = Get-InstalledVersion
    if ($installedVersion) {
        if ("v$installedVersion" -eq $latestVersion -or $installedVersion -eq $latestVersion) {
            Write-Success "$BinaryName $latestVersion is already installed and up to date"
            return
        }
        Write-Info "Upgrading from $installedVersion to $latestVersion"
    }

    # Determine install location
    $installDir = Get-InstallDir
    Write-Info "Install location: $installDir"

    # Create temp directory
    $tempDir = Join-Path $env:TEMP "ticker-install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    try {
        # Construct download URL
        $versionNum = $latestVersion -replace '^v', ''
        $archiveName = "${BinaryName}_${versionNum}_windows_${arch}.zip"
        $downloadUrl = "https://github.com/$Repo/releases/download/$latestVersion/$archiveName"

        Write-Info "Downloading $downloadUrl..."
        $archivePath = Join-Path $tempDir $archiveName
        Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -UseBasicParsing

        # Extract
        Write-Info "Extracting..."
        Expand-Archive -Path $archivePath -DestinationPath $tempDir -Force

        # Install binary
        $sourcePath = Join-Path $tempDir "$BinaryName.exe"
        $destPath = Join-Path $installDir "$BinaryName.exe"

        # If ticker is currently running, we may need to rename the old exe first
        if (Test-Path $destPath) {
            $oldPath = "$destPath.old"
            if (Test-Path $oldPath) {
                Remove-Item $oldPath -Force
            }
            try {
                Rename-Item $destPath $oldPath -Force
            } catch {
                # File might be in use, try to remove it anyway
            }
        }

        Copy-Item $sourcePath $destPath -Force

        Write-Success "$BinaryName $latestVersion installed successfully to $destPath"

        # Check if install dir is in PATH
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($userPath -notlike "*$installDir*") {
            Write-Warning "$installDir is not in your PATH"
            Write-Host ""
            Write-Host "Add it to your PATH by running:"
            Write-Host ""
            Write-Host '  $env:Path += ";' -NoNewline
            Write-Host $installDir -NoNewline
            Write-Host '"'
            Write-Host ""
            Write-Host "Or permanently add it:"
            Write-Host ""
            Write-Host '  [Environment]::SetEnvironmentVariable("Path", $env:Path + ";' -NoNewline
            Write-Host $installDir -NoNewline
            Write-Host '", "User")'
            Write-Host ""

            # Offer to add to PATH automatically
            $addToPath = Read-Host "Add to PATH now? (Y/n)"
            if ($addToPath -ne "n" -and $addToPath -ne "N") {
                [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
                $env:Path += ";$installDir"
                Write-Success "Added $installDir to PATH"
            }
        }

        Write-Success "Run '$BinaryName --help' to get started"

    } finally {
        # Cleanup temp directory
        if (Test-Path $tempDir) {
            Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# Run installation
Install-Ticker
