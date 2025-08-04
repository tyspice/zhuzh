# Makefile for zhuzh - A ChatGPT TUI client

BINARY_NAME = zhuzh
BUILD_DIR = build
CONFIG_DIR = $(HOME)/.config/zhuzh

# Get version from git tag if available else current commit else "unknown"
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo "unknown")

ZHUZH_VERSION ?= $(GIT_TAG)
ZHUZH_INSTALL_DIR ?= /usr/local/bin
GO ?= /usr/local/go/bin/go

LDFLAGS = -ldflags "-X main.version=$(ZHUZH_VERSION)"
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
	@echo "Installing $(BINARY_NAME) to $(ZHUZH_INSTALL_DIR)..."
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "Error: Binary not found at $(BUILD_DIR)/$(BINARY_NAME)"; \
		echo "Please make sure the build was successful before installing."; \
		exit 1; \
	fi
	@mkdir -p $(ZHUZH_INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(ZHUZH_INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(ZHUZH_INSTALL_DIR)/$(BINARY_NAME)
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
	@rm -f $(ZHUZH_INSTALL_DIR)/$(BINARY_NAME)
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
	@echo "  install        Install the application to $(ZHUZH_INSTALL_DIR)"
	@echo "  uninstall      Remove the application from $(ZHUZH_INSTALL_DIR)"
	@echo "  clean          Remove build artifacts"
	@echo "  config         Set up configuration files only"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  ZHUZH_VERSION        Set application version (default: $(ZHUZH_VERSION))"
	@echo "  ZHUZH_INSTALL_DIR    Set installation directory (default: $(ZHUZH_INSTALL_DIR))"
	@echo "  GO                   Set Go binary path (default: $(GO))"
	@echo ""
	@echo "Example:"
	@echo "  GO=/usr/bin/go ZHUZH_INSTALL_DIR=~/bin make install"
	@echo ""
	@echo "Use 'make install' to build and install the application."
