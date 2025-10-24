# Installation Guide

Complete guide to install gooDrive on all platforms.

⚠️ **BETA RELEASE**: This is version v0.1.0-beta.1. Expect changes and potential issues. [Report bugs here](https://github.com/mayura-andrew/gooDrive/issues)

---

## Quick Start (Recommended)

### 1. Go Install (Fastest)

```bash
# Install latest version
go install github.com/mayura-andrew/goodrive@latest

# Verify installation
gooDrive version
```

**Requirements:**
- Go 1.24.0 or higher installed
- `$GOPATH/bin` in your `$PATH`

**What happens:**
- Downloads source from GitHub
- Builds binary for your OS/architecture
- Installs to `$GOPATH/bin/gooDrive`

**Verify PATH is set:**
```bash
echo $GOPATH
# Should output something like /home/username/go or /Users/username/go

echo $PATH
# Should contain $GOPATH/bin path

# If not, add to ~/.bashrc or ~/.zshrc:
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Platform-Specific Installation

### macOS

#### Option 1: Go Install
```bash
go install github.com/mayura-andrew/goodrive@latest
gooDrive --help
```

#### Option 2: Homebrew (When Available)
```bash
brew install mayura-andrew/gooDrive/gooDrive
gooDrive --help
```

#### Option 3: Binary Download
```bash
# Download latest release
curl -L https://github.com/mayura-andrew/gooDrive/releases/latest/download/gooDrive_Darwin_x86_64.tar.gz | tar xz

# Move to PATH
sudo mv gooDrive /usr/local/bin/
gooDrive --help

# For Apple Silicon (M1/M2):
curl -L https://github.com/mayura-andrew/gooDrive/releases/latest/download/gooDrive_Darwin_arm64.tar.gz | tar xz
sudo mv gooDrive /usr/local/bin/
```

---

### Linux

#### Option 1: Go Install
```bash
go install github.com/mayura-andrew/goodrive@latest
gooDrive --help
```

#### Option 2: AUR (Arch Linux)
```bash
yay -S goodrive
# or
pacman -S goodrive
```

#### Option 3: Binary Download
```bash
# Download latest release
wget https://github.com/mayura-andrew/gooDrive/releases/latest/download/gooDrive_Linux_x86_64.tar.gz
tar xzf gooDrive_Linux_x86_64.tar.gz

# Move to PATH
sudo mv gooDrive /usr/local/bin/
gooDrive --help

# For ARM64 (Raspberry Pi):
wget https://github.com/mayura-andrew/gooDrive/releases/latest/download/gooDrive_Linux_arm64.tar.gz
tar xzf gooDrive_Linux_arm64.tar.gz
sudo mv gooDrive /usr/local/bin/
```

#### Option 4: Build from Source
```bash
# Clone repository
git clone https://github.com/mayura-andrew/gooDrive.git
cd gooDrive

# Build
make build

# Install globally
sudo make install

# Verify
gooDrive --help
```

---

### Windows

#### Option 1: Go Install
```powershell
go install github.com/mayura-andrew/goodrive@latest
gooDrive version
```

#### Option 2: Scoop
```powershell
scoop bucket add mayura-andrew https://github.com/mayura-andrew/scoop-bucket
scoop install goodrive
gooDrive version
```

#### Option 3: Binary Download
```powershell
# Download latest release (PowerShell)
curl -o gooDrive.zip https://github.com/mayura-andrew/gooDrive/releases/latest/download/gooDrive_Windows_x86_64.zip

# Extract
Expand-Archive gooDrive.zip -DestinationPath .

# Move to PATH (requires admin)
Move-Item gooDrive.exe $env:ProgramFiles\gooDrive\

# Verify
gooDrive version

# For ARM64:
# Replace 'gooDrive_Windows_x86_64.zip' with 'gooDrive_Windows_arm64.zip'
```

#### Option 4: Build from Source
```powershell
# Clone repository
git clone https://github.com/mayura-andrew/gooDrive.git
cd gooDrive

# Build (requires Go 1.24.0+)
go build -o gooDrive.exe .

# Run
.\gooDrive.exe --help
```

---

## Post-Installation Setup

### 1. Create OAuth Credentials

**Step 1: Google Cloud Console**
1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project (or select existing)
3. Enable "Google Drive API"
4. Create OAuth 2.0 credentials (Desktop application)
5. Download credentials as JSON

**Step 2: Setup gooDrive**
```bash
# For macOS/Linux:
cp ~/Downloads/client_secret_*.json ~/path/to/project/oauth.json

# For Windows:
copy %UserProfile%\Downloads\client_secret_*.json %CD%\oauth.json

# Or place in config directory:
mkdir -p ~/.config/gooDrive
cp oauth.json ~/.config/gooDrive/
```

### 2. Authenticate

```bash
# First time using gooDrive
gooDrive auth

# Browser will open for authentication
# Grant permissions when prompted
# Token saved automatically

# Verify it worked
gooDrive list
```

### 3. Optional: Configure

Create `~/.config/gooDrive/config.json`:

```json
{
  "drive_path": "~/GoogleDrive",
  "sync": false,
  "log_level": "info"
}
```

---

## Verification

### Check Installation

```bash
# Version
gooDrive version

# Help
gooDrive --help

# List available commands
gooDrive

# Test authentication
gooDrive list
```

### Verify Authenticity (Checksums)

```bash
# Download checksums
wget https://github.com/mayura-andrew/gooDrive/releases/latest/download/checksums.txt

# Verify your binary
sha256sum -c checksums.txt

# Should output: gooDrive_Linux_x86_64: OK
```

---

## Upgrading

### From Version A to Version B

#### Using Go Install
```bash
# Automatic upgrade to latest
go install github.com/mayura-andrew/goodrive@latest

# Or specific version
go install github.com/mayura-andrew/goodrive@v0.2.0

# Verify
gooDrive version
```

#### From Binary
```bash
# Download new version
# See binary download section above

# Backup old version
cp $(which gooDrive) ~/gooDrive.backup

# Replace
sudo mv gooDrive_new /usr/local/bin/gooDrive

# Verify
gooDrive version
```

### Data Preservation

Upgrading preserves:
- Authentication tokens (in `~/.config/gooDrive/.tokens.json`)
- Configuration files
- Sync metadata
- Download history

---

## Troubleshooting

### "gooDrive: command not found"

**Solution 1: Add to PATH**
```bash
# For Go install, ensure $GOPATH/bin in PATH
echo $PATH | grep $(go env GOPATH)/bin

# If not present, add to ~/.bashrc or ~/.zshrc:
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc  # or ~/.zshrc
```

**Solution 2: Manual PATH**
```bash
# macOS/Linux - Find binary location
which gooDrive

# Add to PATH manually
sudo ln -s /path/to/gooDrive /usr/local/bin/gooDrive
```

**Solution 3: Windows PATH**
```powershell
# Add to system PATH (requires admin)
$env:Path += ";C:\Program Files\gooDrive"
```

### "oauth.json not found"

```bash
# Make sure oauth.json exists
ls -la oauth.json

# Or in config directory
ls -la ~/.config/gooDrive/

# Get credentials from Google Cloud Console
# See "Post-Installation Setup" section above
```

### "Token expired" Error

```bash
# Refresh authentication
gooDrive auth refresh

# Or re-authenticate
gooDrive auth
```

### Build Failures

```bash
# Check Go version
go version
# Should be 1.24.0 or higher

# Update Go if needed
# Download from golang.org

# Try updating dependencies
go mod tidy
go mod download
```

### Permission Denied (Linux/macOS)

```bash
# Make executable
chmod +x /usr/local/bin/gooDrive

# Or for current user
chmod +x ~/gooDrive
```

---

## System Requirements

| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| Go Version | 1.24.0 | 1.24.0+ |
| RAM | 50 MB | 256 MB+ |
| Disk Space | 50 MB | 500 MB+ |
| Internet | Required | Recommended |
| OS | Any POSIX | macOS, Linux, Windows |

### Supported Platforms

| OS | Arch | Status | Notes |
|----|------|--------|-------|
| Linux | x86_64 | ✅ | Ubuntu, Debian, Fedora |
| Linux | ARM64 | ✅ | Raspberry Pi 4+ |
| macOS | x86_64 | ✅ | Intel Macs |
| macOS | ARM64 | ✅ | Apple Silicon (M1/M2) |
| Windows | x86_64 | ✅ | Windows 10+ |
| Windows | ARM64 | ✅ | Windows on ARM |

---

## Uninstallation

### Go Install
```bash
# Remove binary
rm $(which gooDrive)

# Remove config directory (optional)
rm -rf ~/.config/gooDrive/
```

### Homebrew (macOS)
```bash
brew uninstall mayura-andrew/gooDrive/gooDrive

# Remove config
rm -rf ~/.config/gooDrive/
```

### AUR (Linux)
```bash
yay -R goodrive

# Remove config
rm -rf ~/.config/gooDrive/
```

### Manual Binary
```bash
# Remove binary
sudo rm /usr/local/bin/gooDrive

# Remove config
rm -rf ~/.config/gooDrive/
```

---

## Getting Help

- **GitHub Issues**: https://github.com/mayura-andrew/gooDrive/issues
- **Documentation**: https://mayura-andrew.github.io/gooDrive/
- **Discussion Forum**: https://github.com/mayura-andrew/gooDrive/discussions

