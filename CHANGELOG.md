# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0-beta.1] - 2025-10-24

### Added
- OAuth 2.0 authentication with Google Drive
- File listing with filtering by name and type
- File download functionality with folder support
- File upload with local file tracking
- Search functionality with file queries
- File sharing with email support and role-based permissions
- Directory synchronization with file watcher
- Token auto-refresh mechanism for expired tokens
- Authentication command (`auth`, `auth refresh`)
- Configuration management with settings support
- Multi-platform binary support (Linux, Windows, macOS)
- ARM64 and AMD64 architecture support
- GoReleaser workflow for automated releases
- GitHub Pages deployment for documentation
- Beautiful CLI output with formatted tables
- Progress indicators for file operations

### Features
- `gooDrive list` - List files in your Drive
- `gooDrive download <file-id>` - Download files or folders
- `gooDrive upload <file>` - Upload local files
- `gooDrive search <query>` - Search for files
- `gooDrive share <file-id>` - Share files with users
- `gooDrive sync` - Synchronize local directories
- `gooDrive auth` - Authenticate with Google
- `gooDrive version` - Show version information

### Security
- Token refresh with refresh tokens
- Secure token file storage (mode 0600)
- HTTPS-only OAuth flow
- No hardcoded credentials in code
- Environment variable support for sensitive data

### Documentation
- Comprehensive README with examples
- OAuth setup guide
- CLI usage documentation
- GitHub Pages documentation site
- Build and installation instructions

### Known Limitations
- Single-threaded sync operations (no parallelization)
- No built-in conflict resolution strategy
- Limited error recovery mechanisms
- Token refresh only on API calls (not proactive)

## [v0.2.0] - Planned

### Planned Features
- Parallel file transfers
- Advanced conflict resolution for sync
- Batch operations support
- Custom metadata support
- File versioning support
- Improved error messages
- Comprehensive test suite
- Shell completion (bash, zsh, fish)

## [v1.0.0] - Future

### Planned for Stable Release
- Production-ready stability
- Complete test coverage (80%+)
- Performance optimizations
- Multi-cloud integration roadmap
- API stability guarantee

---

For installation, see [INSTALLATION.md](INSTALLATION.md)
For usage, see [README.md](README.md)
For security, see [SECURITY.md](SECURITY.md)
