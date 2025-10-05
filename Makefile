# Get version from git tag
VERSION ?= $(shell git describe --tags --always --dirty)
NAME = qbittorrent-tool
RELEASE_DIR = release

# Python executable
PYTHON ?= python

# Platform list
PLATFORM_LIST = linux-amd64 darwin-amd64 darwin-arm64
WINDOWS_ARCH_LIST = windows-amd64

# Compressed releases - 统一使用zip格式
ZIP_RELEASES = $(addsuffix .zip, $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST))

# Default target
.PHONY: all
all: $(ZIP_RELEASES)

# Platform-specific builds
.PHONY: windows-amd64
windows-amd64: deps deps-build
	@if [ ! -d "$(RELEASE_DIR)" ]; then mkdir $(RELEASE_DIR); fi
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	@cp dist/$(NAME).exe $(RELEASE_DIR)/$(NAME)-$@.exe
	@cp example.config.json $(RELEASE_DIR)/

.PHONY: linux-amd64
linux-amd64: deps deps-build
	mkdir -p $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	@cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-$@
	@cp example.config.json $(RELEASE_DIR)/

.PHONY: darwin-amd64
darwin-amd64: deps deps-build
	mkdir -p $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	@cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-$@
	@cp example.config.json $(RELEASE_DIR)/

.PHONY: darwin-arm64
darwin-arm64: deps deps-build
	mkdir -p $(RELEASE_DIR)
	$(PYTHON) -m PyInstaller --onefile main.py --name $(NAME)
	@cp dist/$(NAME) $(RELEASE_DIR)/$(NAME)-$@
	@cp example.config.json $(RELEASE_DIR)/

# Rules for creating .zip releases for all platforms
.PHONY: $(ZIP_RELEASES)
$(ZIP_RELEASES): %.zip : %
	@if [ -f "$(RELEASE_DIR)/$(NAME)-$(basename $@).exe" ]; then \
		zip -j $(RELEASE_DIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(RELEASE_DIR)/$(NAME)-$(basename $@).exe $(RELEASE_DIR)/example.config.json; \
	else \
		chmod +x $(RELEASE_DIR)/$(NAME)-$(basename $@); \
		zip -j $(RELEASE_DIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(RELEASE_DIR)/$(NAME)-$(basename $@) $(RELEASE_DIR)/example.config.json; \
	fi

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@if [ -f "requirements.txt" ]; then \
		echo "Found requirements.txt, installing dependencies with pip..."; \
		$(PYTHON) -m pip install -r requirements.txt; \
	else \
		echo "No requirements.txt found, using pyproject.toml"; \
	fi

# Install build dependencies
.PHONY: deps-build
deps-build:
	$(PYTHON) -m pip install pyinstaller

# Create release builds for current platform
.PHONY: release
release: all

# Clean build artifacts
.PHONY: clean
clean:
	if exist dist rmdir /s /q dist 2>nul || rm -rf dist 2>/dev/null || true
	if exist build rmdir /s /q build 2>nul || rm -rf build 2>/dev/null || true
	for /d %%i in (*.egg-info) do if exist "%%i" rmdir /s /q "%%i" 2>nul || rm -rf *.egg-info 2>/dev/null || true
	del /q *.pyc 2>nul || rm -f *.pyc 2>/dev/null || true
	for /d /r . %%j in (__pycache__) do @if exist "%%j" rmdir /s /q "%%j" 2>nul || find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	if exist $(RELEASE_DIR) rmdir /s /q $(RELEASE_DIR) 2>nul || rm -rf $(RELEASE_DIR) 2>/dev/null || true
