# Go Package Release Guide - gooDrive

Complete guide for releasing gooDrive as a Go package with proper versioning and distribution.

**Current Status**: 🟡 Beta Release (v0.1.0-beta.1)

---

## 1. Pre-Release Checklist

### Code Quality
- [ ] Run `go fmt ./...` - Code formatting
- [ ] Run `go vet ./...` - Code analysis
- [ ] Run `go test ./...` - Unit tests pass
- [ ] Run `golangci-lint run` - Linter checks
- [ ] Update `CHANGELOG.md` with changes
- [ ] Update version string
- [ ] Review all dependencies: `go mod tidy`

### Documentation
- [ ] Update `README.md`
- [ ] Add/update API documentation (godoc)
- [ ] Create/update `INSTALL.md`
- [ ] Write `SECURITY.md` if applicable
- [ ] Add `CONTRIBUTING.md` guidelines
- [ ] Update examples in `docs/`

### Testing
- [ ] Unit tests written (80%+ coverage target)
- [ ] Integration tests pass
- [ ] Manual testing on all platforms
- [ ] Test with real Google Drive credentials

### CI/CD
- [ ] GitHub Actions workflow configured
- [ ] GoReleaser configuration updated
- [ ] All CI checks passing

---

## 2. Version Numbering Strategy

### Semantic Versioning (MAJOR.MINOR.PATCH)

**Current: v0.1.0-beta.1**

```
v0.1.0-beta.1 → v0.1.0 → v0.2.0 → v1.0.0
 └─ Beta      └─ Stable  └─ Minor   └─ Major
```

### Version Progression Guide

```
Beta/Pre-release:
  v0.1.0-alpha.1   - Early development
  v0.1.0-beta.1    - Feature complete, testing phase
  v0.1.0-rc.1      - Release candidate
  v0.1.0           - First stable release

Minor Updates (0.MINOR.0):
  v0.2.0           - New features (backward compatible)
  v0.3.0           - More features

Major Updates (MAJOR.0.0):
  v1.0.0           - Breaking changes / stable API

Patch Updates (0.0.PATCH):
  v0.1.1           - Bug fixes only
  v0.1.2           - More bug fixes
```

### Current Release Status

**v0.1.0-beta.1** means:
- Beta phase (not production-ready)
- First minor version (0)
- First patch level (1)
- Pre-release suffix (beta.1)

---

## 3. Release Steps

### Step 1: Prepare Release

```bash
# 1. Update version in code
vim cmd/version.go  # Update Version variable

# 2. Update changelog
vim CHANGELOG.md    # Add new version section

# 3. Run final checks
go fmt ./...
go vet ./...
golangci-lint run

# 4. Commit changes
git add CHANGELOG.md cmd/version.go
git commit -m "chore: Prepare v0.2.0 release"
git push origin main
```

### Step 2: Create Release Tag

```bash
# Create annotated tag
git tag -a v0.2.0 -m "Release v0.2.0

- Feature 1
- Feature 2
- Bug fix

This release includes...
"

# Push tag to GitHub
git push origin v0.2.0
```

**Tag Format:**
- `v0.1.0` - Stable release
- `v0.1.0-beta.1` - Pre-release
- `v0.1.0-rc.1` - Release candidate

### Step 3: Verify GitHub Actions

```bash
# Watch the workflow
# https://github.com/mayura-andrew/gooDrive/actions

# Wait for GoReleaser to:
# 1. Build binaries (6 platforms)
# 2. Create checksums
# 3. Generate release notes
# 4. Create GitHub release
```

### Step 4: Publish Release

1. **GitHub Releases** (Automated by GoReleaser)
   - Release page: https://github.com/mayura-andrew/gooDrive/releases
   - Binary downloads available
   - Checksums included

2. **Manual GitHub Release Notes**
   ```bash
   # Edit release manually if needed
   # Add release highlights
   # Update download instructions
   ```

---

## 4. How Users Install Your Package

### Method 1: Go Install (Recommended)

```bash
# Users run:
go install github.com/mayura-andrew/goodrive@v0.1.0-beta.1

# Or latest:
go install github.com/mayura-andrew/goodrive@latest

# Or local development:
go install github.com/mayura-andrew/goodrive@main
```

**How it works:**
- Go downloads source from GitHub
- Builds binary using your code
- Installs to `$GOPATH/bin/gooDrive`

### Method 2: Binary Download

```bash
# Users download pre-built binary from:
https://github.com/mayura-andrew/gooDrive/releases/tag/v0.1.0-beta.1

# Verify checksum
sha256sum -c checksums.txt

# Make executable
chmod +x gooDrive
sudo mv gooDrive /usr/local/bin/
```

### Method 3: Package Managers (Future)

```bash
# Homebrew (macOS)
brew install mayura-andrew/gooDrive/gooDrive

# AUR (Arch Linux)
yay -S goodrive

# Snap
sudo snap install goodrive

# APT (Ubuntu/Debian)
sudo apt install goodrive
```

### Method 4: Using as Library

```go
// Users can import your packages:
import "github.com/mayura-andrew/goodrive/internal/drive"
import "github.com/mayura-andrew/goodrive/internal/auth"

// Then use functions
client, err := drive.InitDriveClient()
files, err := client.ListFiles()
```

---

## 5. Package Structure for Users

When users do `go install`, they get:

```
$GOPATH/bin/
└── gooDrive              # Binary (executable)

~/.config/gooDrive/       # Config directory
├── config.json           # Configuration
├── .tokens.json          # Stored OAuth tokens

.drive-cli-meta.json      # Sync metadata (in project dir)
oauth.json                # OAuth credentials (in project dir)
```

---

## 6. Release Automation with GoReleaser

Your `.goreleaser.yml` configuration handles:

### What Gets Built
```yaml
builds:
  - main: ./main.go
    goos:
      - linux      # Linux 64-bit
      - windows    # Windows 64-bit
      - darwin     # macOS 64-bit
    goarch:
      - amd64      # Intel/AMD 64-bit
      - arm64      # ARM 64-bit (Apple Silicon, Raspberry Pi)
```

**Result: 6 Binaries**
- `gooDrive_Linux_x86_64`
- `gooDrive_Linux_arm64`
- `gooDrive_Windows_x86_64`
- `gooDrive_Windows_arm64`
- `gooDrive_Darwin_x86_64`
- `gooDrive_Darwin_arm64`

### What Gets Created
- Compressed archives (tar.gz, zip)
- Checksums file
- GitHub release page
- Release notes

---

## 7. Testing Before Release

### Local Testing

```bash
# Build locally
go build -o gooDrive .

# Test each command
./gooDrive --help
./gooDrive version
./gooDrive auth
./gooDrive list
./gooDrive search "test"
```

### Cross-Platform Testing

```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gooDrive.exe .

# Build for macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o gooDrive_arm64 .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gooDrive_linux .
```

### Dry Run Release

```bash
# Test GoReleaser without publishing
goreleaser release --clean --snapshot

# Outputs to ./dist/ without creating GitHub release
```

---

## 8. Post-Release Tasks

### Update Documentation

1. **Update README.md**
```markdown
## Installation

Latest version: **v0.1.0-beta.1**

### Quick Install
\`\`\`bash
go install github.com/mayura-andrew/goodrive@latest
\`\`\`
```

2. **Update website**
- Update docs/index.html with latest version
- Add release notes to website
- Update quick start guide

3. **Announce Release**
```bash
# Social media posts
# Email newsletter
# GitHub discussions
# Reddit /r/golang
# Hacker News (if significant)
```

### Monitor Usage

```bash
# Track downloads
# Monitor GitHub issues
# Respond to feedback
# Plan next release
```

---

## 9. Publishing to Go Package Registry

### pkg.go.dev (Automatic)

When you tag and push:
```bash
git tag v0.1.0
git push origin v0.1.0
```

Within 5-15 minutes:
- Documentation appears at: https://pkg.go.dev/github.com/mayura-andrew/goodrive@v0.1.0
- Automatically parsed from your godoc comments
- Full API documentation available

### Go Proxy

Users can install immediately:
```bash
go install github.com/mayura-andrew/goodrive@v0.1.0
```

The Go proxy (proxy.golang.org) caches:
- Source code
- go.mod file
- go.sum file

---

## 10. Version Matrix

### Current Status

| Component | Status | Version |
|-----------|--------|---------|
| CLI Binary | Beta | v0.1.0-beta.1 |
| Go Package | Beta | v0.1.0-beta.1 |
| Documentation | Beta | Partial |
| Tests | Beta | Minimal |
| API Stability | Beta | May change |

### Roadmap to v1.0.0

```
v0.1.0-beta.1  ──→  v0.1.0  ──→  v0.2.0  ──→  v0.3.0  ──→  v1.0.0
  Current        Stable       Features    More        Stable
                              Added       Stable      API
```

---

## 11. Release Checklist for v0.1.0 Stable

- [ ] All beta issues resolved
- [ ] Tests written for core functionality
- [ ] Documentation complete
- [ ] No breaking changes from beta
- [ ] Security audit completed
- [ ] Dependencies reviewed
- [ ] Performance tested
- [ ] Example code works
- [ ] README updated
- [ ] CHANGELOG updated
- [ ] GitHub Actions passing
- [ ] Tag created: `git tag v0.1.0`
- [ ] Release announced

---

## 12. Emergency Hotfix Release

If critical bug found in v0.1.0:

```bash
# Create hotfix branch
git checkout -b hotfix/0.1.1 v0.1.0

# Fix the bug
git commit -am "fix: Critical bug in X"

# Merge back
git checkout main
git merge hotfix/0.1.1

# Tag hotfix
git tag v0.1.1
git push origin v0.1.1
```

**Hotfix Versions:**
- v0.1.1 - Security or critical fixes
- v0.1.2 - More critical fixes
- After v0.2.0, use standard semver

---

## 13. User Documentation for Each Release

### For v0.1.0-beta.1 Users

```markdown
# Installation

This is a **beta release**. Use at your own risk.

\`\`\`bash
go install github.com/mayura-andrew/goodrive@v0.1.0-beta.1
\`\`\`

## Known Issues
- Limited error handling
- Single-threaded operations only
- Token refresh may fail in edge cases

## Reporting Issues
- GitHub Issues: https://github.com/mayura-andrew/gooDrive/issues
- Email: mayura-andrew@example.com
```

### For v0.1.0 Stable Users

```markdown
# Installation

This is the **first stable release**.

\`\`\`bash
go install github.com/mayura-andrew/goodrive@v0.1.0
\`\`\`

## What's New
- OAuth authentication
- File list, download, upload
- Search functionality
- Auto token refresh

## Upgrade from beta
Just run: `go install github.com/mayura-andrew/goodrive@v0.1.0`
```

---

## 14. Security & Signing Releases

### GPG Signing (Optional but Recommended)

```bash
# Generate GPG key
gpg --gen-key

# Sign release tag
git tag -s v0.1.0 -m "Release v0.1.0"

# Users verify
git verify-tag v0.1.0
```

### Checksums Verification

GoReleaser creates `checksums.txt`:
```
abc123def456...  gooDrive_Linux_x86_64.tar.gz
def789ghi012...  gooDrive_Windows_x86_64.zip
...
```

Users verify:
```bash
sha256sum -c checksums.txt
```

---

## 15. Complete Release Command Sequence

For your next release (v0.1.0):

```bash
# 1. Update files
echo 'Version = "v0.1.0"' > cmd/version.go
# Update CHANGELOG.md

# 2. Commit
git add -A
git commit -m "chore: Release v0.1.0"
git push origin main

# 3. Create tag
git tag -a v0.1.0 -m "Release v0.1.0: First stable release

Features:
- OAuth authentication
- File operations
- Auto token refresh
"

# 4. Push tag (triggers GitHub Actions)
git push origin v0.1.0

# 5. Monitor at:
# https://github.com/mayura-andrew/gooDrive/actions

# 6. When ready, update docs
# https://github.com/mayura-andrew/gooDrive/releases/tag/v0.1.0
```

---

## Next Steps

1. **Immediate**: Add unit tests and documentation
2. **Short-term**: Fix critical issues from best practices analysis
3. **Medium-term**: Release v0.1.0 stable
4. **Long-term**: Plan v1.0.0 with multi-cloud support

---

## References

- [Semantic Versioning](https://semver.org/)
- [Go Release Checklist](https://golang.org/wiki/Release)
- [GoReleaser Documentation](https://goreleaser.com/)
- [pkg.go.dev Documentation](https://pkg.go.dev/)
- [Go Module Publishing](https://go.dev/doc/modules/publishing)

