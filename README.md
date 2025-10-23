# gooDrive 🚀

A minimalist CLI tool for Google Drive operations. Sync, upload, download, and manage your Drive files from the terminal.

## Features

- 📥 Clone files/folders from Drive
- 📤 Upload and track local files
- 🔄 Sync changes bidirectionally
- 📋 List and search Drive contents
- 👁️ Preview file contents
- 🗑️ Delete files remotely

## Installation

```bash
git clone https://github.com/mayura-andrew/gooDrive.git
cd gooDrive
make install
```

Or build from source:

```bash
go build -o drive cmd/drive/main.go
sudo mv drive /usr/local/bin/
```

## Setup

1. Create a Google Cloud project and enable Drive API
2. Download OAuth credentials as `oauth.json`
3. Place `oauth.json` in the same directory
4. Run any command to authenticate

```bash
drive ls
```

## Usage

### Download from Drive
```bash
drive clone <file-id-or-link>
drive clone https://drive.google.com/file/d/abc123/view
```

### Upload to Drive
```bash
drive add-remote myfile.txt
```

### List files
```bash
drive ls                           # Current directory
drive view-files                   # All files
drive view-files --name="report"   # Filter by name
```

### Sync operations
```bash
drive status    # Check tracked files
drive pull      # Download updates
drive push      # Upload changes
```

### File operations
```bash
drive cat <file-id>    # View contents
drive rm <file-id>     # Delete file
```

## Configuration

Token and metadata files are stored in your working directory:
- `.drive-cli-token.json` - OAuth token
- `*/.drive-cli-meta.json` - File tracking metadata

## Requirements

- Go 1.19+
- Google Cloud OAuth credentials
- Internet connection

## License

MIT

## Contributing

Pull requests welcome. Please open an issue first to discuss changes.

## Author

Mayura Andrew

---

**Note:** This tool requires proper Google Drive API credentials. Never share your `oauth.json` or token files.
