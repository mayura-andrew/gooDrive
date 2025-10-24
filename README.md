# gooDrive 🚀

A powerful, minimalist CLI tool for Google Drive operations. Manage your Drive files directly from the terminal with ease.

[![Website](https://img.shields.io/badge/Website-gooDrive-blue)](https://mayura-andrew.github.io/gooDrive/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)

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
