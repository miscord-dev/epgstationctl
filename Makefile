# Makefile for epgstationctl

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOLINT=golangci-lint

# Binary name
BINARY_NAME=epgstationctl
BINARY_PATH=./bin/$(BINARY_NAME)

# Source files
CMD_PATH=./cmd/epgstationctl

.PHONY: all build install lint test clean generate

all: lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_PATH) $(CMD_PATH)

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	@$(GOINSTALL) $(CMD_PATH)

# Run the linter
lint:
	@echo "Running linter..."
	@$(GOLINT) run

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) ./...

# Clean up build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(BINARY_PATH)

# Generate API client
generate:
	@echo "Generating API client..."
	@cd api && $(GOGENERATE) .