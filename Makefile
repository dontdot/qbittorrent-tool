# Get version from git tag
VERSION ?= $(shell git describe --tags --always --dirty)
NAME = qbittorrent-tool
RELEASE_DIR = release

# Python executable
PYTHON ?= python

# Platform list
PLATFORM_LIST = linux-amd64 darwin-amd64 darwin-arm64
WINDOWS_ARCH_LIST = windows-amd64

# Default target
.PHONY: all
all: clean build-exe

# Clean build artifacts
.PHONY: clean
clean:
	if exist dist rmdir /s /q dist 2>nul || rm -rf dist 2>/dev/null || true
	if exist build rmdir /s /q build 2>nul || rm -rf build 2>/dev/null || true
	for /d %%i in (*.egg-info) do if exist "%%i" rmdir /s /q "%%i" 2>nul || rm -rf *.egg-info 2>/dev/null || true
	del /q *.pyc 2>nul || rm -f *.pyc 2>/dev/null || true
	for /d /r . %%j in (__pycache__) do @if exist "%%j" rmdir /s /q "%%j" 2>nul || find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true

# Install dependencies
.PHONY: deps
deps:
	@(if exist requirements.txt ( \
		(if exist uv.lock ( \
			uv pip install -r requirements.txt \
		) else ( \
			$(PYTHON) -m pip install -r requirements.txt \
		)) \
	) else ( \
		echo "No requirements.txt found, using pyproject.toml" \
	)) || \
	(if [ -f requirements.txt ]; then \
		if [ -f uv.lock ]; then \
			uv pip install -r requirements.txt; \
		else \
			$(PYTHON) -m pip install -r requirements.txt; \
		fi \
	else \
		echo "No requirements.txt found, using pyproject.toml"; \
	fi)

# Install build dependencies
.PHONY: deps-build
deps-build:
	$(PYTHON) -m pip install pyinstaller

# Platform-specific builds
.PHONY: windows-amd64
windows-amd64: deps deps-build
	if not exist $(RELEASE_DIR) mkdir $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	copy dist\$(NAME).exe $(RELEASE_DIR)\$(NAME)-windows-amd64.exe
	copy example.config.json $(RELEASE_DIR)\

.PHONY: linux-amd64
linux-amd64: deps deps-build
	mkdir -p $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-linux-amd64
	cp example.config.json $(RELEASE_DIR)/

.PHONY: darwin-amd64
darwin-amd64: deps deps-build
	mkdir -p $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-darwin-amd64
	cp example.config.json $(RELEASE_DIR)/

.PHONY: darwin-arm64
darwin-arm64: deps deps-build
	mkdir -p $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-darwin-arm64
	cp example.config.json $(RELEASE_DIR)/

# Build executable with pyinstaller for current platform
.PHONY: build-exe
build-exe: deps deps-build
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	@if exist dist ( \
		copy example.config.json dist\ \
	) else ( \
		cp example.config.json dist/ \
	)

# Install the package in development mode
.PHONY: dev-install
dev-install:
	$(PYTHON) -m pip install -e .

# Run the tool directly with Python
.PHONY: run
run:
	$(PYTHON) main.py

# Create release builds for current platform
.PHONY: release
release: clean deps deps-build
	mkdir -p $(RELEASE_DIR)
	
	# Build executable
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	
	# Copy platform-specific executable with platform suffix
	@if exist dist ( \
		if exist dist\$(NAME).exe ( \
			copy dist\$(NAME).exe $(RELEASE_DIR)\$(NAME)-windows-amd64.exe \
		) else ( \
			copy dist\$(NAME) $(RELEASE_DIR)\$(NAME)-windows-amd64 \
		) \
		copy example.config.json $(RELEASE_DIR)\ \
	) else ( \
		cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-$$(uname -s | tr '[:upper:]' '[:lower:]')-$$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/') \
		cp example.config.json $(RELEASE_DIR)/ \
	)

# Help information
.PHONY: help
help:
	@echo Available targets:
	@echo   all             - Clean and build executable (default)
	@echo   clean           - Remove build artifacts
	@echo   deps            - Install Python dependencies
	@echo   deps-build      - Install build dependencies (pyinstaller)
	@echo   windows-amd64   - Build for Windows
	@echo   linux-amd64     - Build for Linux
	@echo   darwin-amd64    - Build for macOS Intel
	@echo   darwin-arm64    - Build for macOS Apple Silicon
	@echo   build-exe       - Build executable using pyinstaller for current platform
	@echo   dev-install     - Install package in development mode
	@echo   run             - Run the tool directly with Python
	@echo   release         - Create release builds for current platform
	@echo   help            - Show this help message