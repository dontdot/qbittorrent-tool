# Get version from git tag
VERSION ?= $(shell git describe --tags --always --dirty)

# Python executable
PYTHON ?= python3

# Default target
.PHONY: all
all: clean build

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf dist/
	rm -rf build/
	rm -rf *.egg-info/
	find . -type f -name "*.pyc" -delete
	find . -type d -name "__pycache__" -delete

# Install dependencies
.PHONY: deps
deps:
	$(PYTHON) -m pip install -r requirements.txt

# Install pyinstaller for building executable
.PHONY: deps-build
deps-build:
	$(PYTHON) -m pip install pyinstaller

# Build executable using pyinstaller
.PHONY: build
build: deps deps-build
	$(PYTHON) -m PyInstaller pyinstaller.spec

# Install the package in development mode
.PHONY: dev-install
dev-install:
	$(PYTHON) -m pip install -e .

# Run the tool directly with Python
.PHONY: run
run:
	$(PYTHON) main.py

# Help information
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build executable (default)"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install Python dependencies"
	@echo "  deps-build   - Install build dependencies (pyinstaller)"
	@echo "  build        - Build executable using pyinstaller"
	@echo "  dev-install  - Install package in development mode"
	@echo "  run          - Run the tool directly with Python"
	@echo "  help         - Show this help message"