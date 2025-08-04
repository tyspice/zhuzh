# Makefile for zhuzh - A ChatGPT TUI client

# Basic configuration
BINARY_NAME = zhuzh
BUILD_DIR = build
CONFIG_DIR = $(HOME)/.config/zhuzh

# Overwrite these with environment variables if need be
VERSION ?= 0.1.0
INSTALL_DIR ?= /usr/local/bin
GO ?= /usr/local/go/bin/go

# Go build flags
LDFLAGS = -ldflags "-X main.version=$(VERSION)"
GOFLAGS = -trimpath

.PHONY: all build clean install uninstall config help show-go-path

all: build

# Build the binary
build:
	@if [ -z "$(GO)" ]; then \
		echo "Error: Go is not installed or not in PATH."; \
		echo "Please install Go from https://golang.org/dl/ and make sure it's in your PATH."; \
		exit 1; \
	fi
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Install the binary and config
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "Error: Binary not found at $(BUILD_DIR)/$(BINARY_NAME)"; \
		echo "Please make sure the build was successful before installing."; \
		exit 1; \
	fi
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Setting up configuration directory..."
	@mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_DIR)/config.yml ]; then \
		cp config.yml.example $(CONFIG_DIR)/config.yml; \
		echo "Copied example config to $(CONFIG_DIR)/config.yml"; \
		echo "IMPORTANT: Please edit $(CONFIG_DIR)/config.yml and add your API key!"; \
	else \
		echo "Config file already exists at $(CONFIG_DIR)/config.yml"; \
	fi
	@echo "Installation complete! Run '$(BINARY_NAME)' to start."

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) has been uninstalled."
	@echo "Note: Configuration directory $(CONFIG_DIR) has been preserved."
	@echo "To remove it completely, run: rm -rf $(CONFIG_DIR)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Setup config only
config:
	@echo "Setting up configuration directory..."
	@mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_DIR)/config.yml ]; then \
		cp config.yml.example $(CONFIG_DIR)/config.yml; \
		echo "Copied example config to $(CONFIG_DIR)/config.yml"; \
		echo "IMPORTANT: Please edit $(CONFIG_DIR)/config.yml and add your API key!"; \
	else \
		echo "Config file already exists at $(CONFIG_DIR)/config.yml"; \
	fi

# Help information
help:
	@echo "zhuzh - ChatGPT Terminal UI Client"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  build          Build the application"
	@echo "  install        Install the application to $(INSTALL_DIR)"
	@echo "  uninstall      Remove the application from $(INSTALL_DIR)"
	@echo "  clean          Remove build artifacts"
	@echo "  config         Set up configuration files only"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  VERSION        Set application version (default: $(VERSION))"
	@echo "  INSTALL_DIR    Set installation directory (default: $(INSTALL_DIR))"
	@echo "  GO             Set Go binary path (default: $(GO))"
	@echo ""
	@echo "Example:"
	@echo "  GO=/usr/bin/go INSTALL_DIR=~/bin make install"
	@echo ""
	@echo "Use 'make install' to build and install the application."
