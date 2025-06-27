# Makefile for GOS - Go Version Manager CLI

BINARY_NAME=gos
MAIN_FILE=main.go
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean install uninstall test deps help release

all: clean deps build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Built successfully: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: clean deps
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	
	# Linux (Intel)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	
	# Linux (ARM)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	
	@echo "✅ Built for all platforms in $(BUILD_DIR)/"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Install the binary globally
install: build
	@echo "📥 Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Installed successfully. Run '$(BINARY_NAME) --help' to get started."

# Uninstall the binary
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Uninstalled successfully."

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Run the application locally
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development: build and run with arguments
dev: build
	@echo "🚀 Running $(BINARY_NAME) $(ARGS)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Format the code
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...

# Lint the code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Show help
help:
	@echo "GOS - Go Version Manager CLI"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build the binary for current platform"
	@echo "  build-all  - Build for all supported platforms"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  install    - Install binary globally (requires sudo)"
	@echo "  uninstall  - Remove installed binary (requires sudo)"
	@echo "  test       - Run tests"
	@echo "  run        - Build and run the application"
	@echo "  dev        - Build and run with arguments (use ARGS='...')"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code (requires golangci-lint)"
	@echo "  release    - Create a new release (use VERSION='x.y.z')"
	@echo "  help       - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make install"
	@echo "  make dev ARGS='status'"
	@echo "  make dev ARGS='install 1.21.5'"
	@echo "  make release VERSION='1.0.0'"

# Quick development commands
.PHONY: quick-test quick-setup quick-status

quick-test: build
	@echo "🧪 Quick test - showing help..."
	./$(BUILD_DIR)/$(BINARY_NAME) help

quick-setup: build
	@echo "🔧 Quick test - setup (dry run)..."
	./$(BUILD_DIR)/$(BINARY_NAME) setup --help

quick-status: build
	@echo "📊 Quick test - status..."
	./$(BUILD_DIR)/$(BINARY_NAME) status || true

# Release management
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Use: make release VERSION='1.0.0'"; \
		exit 1; \
	fi
	@echo "🚀 Creating release $(VERSION)..."
	./release.sh $(VERSION)
