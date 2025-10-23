.PHONY: build install clean test run help

BINARY_NAME=drive
INSTALL_PATH=/usr/local/bin

help:
	@echo "gooDrive - Makefile commands:"
	@echo "  make build     - Build the binary"
	@echo "  make install   - Build and install to $(INSTALL_PATH)"
	@echo "  make clean     - Remove built files"
	@echo "  make test      - Run tests"
	@echo "  make run       - Run without installing"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) cmd/drive/main.go
	@echo "✓ Build complete: ./$(BINARY_NAME)"

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	@sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "✓ Installed successfully"
	@echo "Run 'drive --help' to get started"

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "✓ Clean complete"

test:
	@echo "Running tests..."
	@go test -v ./...

run:
	@go run cmd/drive/main.go

.DEFAULT_GOAL := help
