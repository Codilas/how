#!/bin/bash
set -e

# Installation script for 'how' AI shell assistant

REPO="tzvonimir/how"
INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.config/how"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

# Check if running on supported platform
check_platform() {
    case "$(uname -s)" in
        Linux*)     PLATFORM=linux;;
        Darwin*)    PLATFORM=darwin;;
        *)          error "Unsupported platform: $(uname -s)";;
    esac
    
    case "$(uname -m)" in
        x86_64)     ARCH=amd64;;
        arm64)      ARCH=arm64;;
        aarch64)    ARCH=arm64;;
        *)          error "Unsupported architecture: $(uname -m)";;
    esac
    
    info "Detected platform: $PLATFORM-$ARCH"
}

# Get latest release version
get_latest_version() {
    info "Fetching latest release..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi
    info "Latest version: $VERSION"
}

# Download and install binary
install_binary() {
    BINARY_NAME="how-$PLATFORM-$ARCH"
    if [ "$PLATFORM" = "windows" ]; then
        BINARY_NAME="$BINARY_NAME.exe"
    fi
    
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME"
    
    info "Downloading from: $DOWNLOAD_URL"
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        curl -L "$DOWNLOAD_URL" -o "$INSTALL_DIR/how"
    elif command -v wget >/dev/null 2>&1; then
        wget "$DOWNLOAD_URL" -O "$INSTALL_DIR/how"
    else
        error "Neither curl nor wget is available"
    fi
    
    # Make executable
    chmod +x "$INSTALL_DIR/how"
    
    success "Binary installed to $INSTALL_DIR/how"
}

# Check if ~/.local/bin is in PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warning "$INSTALL_DIR is not in your PATH"
        echo
        echo "Add this to your shell configuration file (~/.bashrc, ~/.zshrc, etc.):"
        echo "export PATH=\"$INSTALL_DIR:\$PATH\""
        echo
    fi
}

# Offer to run setup
offer_setup() {
    echo
    info "Installation complete!"
    echo
    echo "Next steps:"
    echo "1. Make sure $INSTALL_DIR is in your PATH"
    echo "2. Run 'how setup' to configure your AI provider"
    echo "3. Run 'how install' to set up shell integration"
    echo
    
    read -p "Would you like to run the setup now? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        "$INSTALL_DIR/how" setup
    fi
}

# Main installation flow
main() {
    echo "Installing 'how' AI Shell Assistant..."
    echo
    
    check_platform
    get_latest_version
    install_binary
    check_path
    offer_setup
}

main "$@"
