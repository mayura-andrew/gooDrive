# gooDrive 

A powerful, minimalist CLI tool for Google Drive operations. Manage your Drive files directly from the terminal with ease.

[![GitHub stars](https://img.shields.io/github/stars/mayura-andrew/gooDrive?style=social)](https://github.com/mayura-andrew/gooDrive/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/mayura-andrew/gooDrive?style=social)](https://github.com/mayura-andrew/gooDrive/network/members)

[![Website](https://img.shields.io/badge/Website-gooDrive-blue)](https://mayura-andrew.github.io/gooDrive/)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mayura-andrew/gooDrive)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/mayura-andrew/gooDrive?include_prereleases)](https://github.com/mayura-andrew/gooDrive/releases)
[![Downloads](https://img.shields.io/github/downloads/mayura-andrew/gooDrive/total)](https://github.com/mayura-andrew/gooDrive/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/mayura-andrew/goodrive)](https://goreportcard.com/report/github.com/mayura-andrew/goodrive)
[![License](https://img.shields.io/github/license/mayura-andrew/gooDrive)](LICENSE)
[![Issues](https://img.shields.io/github/issues/mayura-andrew/gooDrive)](https://github.com/mayura-andrew/gooDrive/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/mayura-andrew/gooDrive)](https://github.com/mayura-andrew/gooDrive/pulls)

**[📖 Documentation](https://mayura-andrew.github.io/gooDrive/)** | **[🚀 Getting Started](#installation)** | **[🤝 Contributing](CONTRIBUTING.md)**

---

## ✨ Features

- 📥 **Download** - Clone files and folders from Google Drive  
- 📤 **Upload** - Upload files to your Drive with tracking  
- 📋 **List** - Browse and filter your Drive contents  
- 🔍 **Search** - Find files quickly with advanced queries  
- 🔗 **Share** - Generate shareable links and manage permissions  
- 🔄 **Sync** - Synchronize local folders with Drive  

---

## 📦 Installation

### Quick Install

```bash
# Clone and build
git clone https://github.com/mayura-andrew/gooDrive.git
cd gooDrive
make build

# Run
./gooDrive --help
```

### Install Globally

```bash
make install
gooDrive --help
```

---

## 🔧 Setup OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Enable **Google Drive API**
3. Create **OAuth 2.0 Client ID** (Desktop app)
4. Download credentials as `oauth.json`
5. Place in project root

---

## 🚀 Usage

```bash
# Download
gooDrive download <file-id>

# Upload
gooDrive upload document.pdf

# List
gooDrive list --name "report"

# Search
gooDrive search "presentation"

# Share
gooDrive share <file-id> --email user@example.com
```

See full documentation at [gooDrive Docs](https://mayura-andrew.github.io/gooDrive/)

---

## 📄 License

MIT License - see [LICENSE](LICENSE)

## 👤 Author

**Mayura Andrew** - [@mayura-andrew](https://github.com/mayura-andrew)
