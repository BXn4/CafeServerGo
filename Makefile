
# Variables
APP_NAME := cafego
BUILD_DIR := build
GO_FILES := $(wildcard *.go)
VERSION := 1.0.0

# Targets
.PHONY: all clean build-linux build-windows

# Default target
all: clean build-linux build-windows

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-$(VERSION) $(GO_FILES)

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-$(VERSION).exe $(GO_FILES)

