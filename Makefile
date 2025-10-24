.PHONY: build install clean test run help deps auth-check

BINARY_NAME=gooDrive
INSTALL_PATH=/usr/local/bin
GO_FILES=$(shell find . -name '*.go' -type f)

help:
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║           gooDrive - CLI for Google Drive                 ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available commands:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Build and install to $(INSTALL_PATH)"
	@echo "  make clean       - Remove built files and cache"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run without building"
	@echo "  make deps        - Download dependencies"
	@echo "  make auth-check  - Check if OAuth credentials exist"
	@echo "  make help        - Show this help message"
	@echo ""
	@echo "Usage:"
	@echo "  ./$(BINARY_NAME) download <file-id>    - Download a file"
	@echo "  ./$(BINARY_NAME) upload <file-path>    - Upload a file"
	@echo "  ./$(BINARY_NAME) list                  - List files"
	@echo "  ./$(BINARY_NAME) search <query>        - Search files"
	@echo "  ./$(BINARY_NAME) share <file-id>       - Share a file"
	@echo ""

deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies ready"

auth-check:
	@if [ ! -f oauth.json ]; then \
		echo "⚠️  Warning: oauth.json not found"; \
		echo "Please download OAuth credentials from Google Cloud Console"; \
		echo "and save as oauth.json in the project root"; \
		exit 1; \
	else \
		echo "✓ OAuth credentials found"; \
	fi

build: deps
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "✓ Build complete: ./$(BINARY_NAME)"
	@echo ""
	@echo "Run './$(BINARY_NAME) --help' to get started"

install: build
	@echo "📦 Installing to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Installed successfully"
	@echo ""
	@echo "You can now run '$(BINARY_NAME)' from anywhere!"

clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean -cache -modcache -testcache
	@echo "✓ Clean complete"

test:
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1
	@echo "✓ Tests complete"

run: auth-check
	@echo "🚀 Running $(BINARY_NAME)..."
	@go run main.go

# Development commands
dev-download:
	@go run main.go download $(ARGS)

dev-upload:
	@go run main.go upload $(ARGS)

dev-list:
	@go run main.go list

dev-search:
	@go run main.go search $(ARGS)

.DEFAULT_GOAL := help
