#!/bin/bash
set -e

# better-fg installer script

REPO="Lewenhaupt/better-fg"
BINARY_NAME="better-fg"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

case $OS in
    linux) OS="Linux" ;;
    darwin) OS="Darwin" ;;
    *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

# Get latest release
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to get latest release"
    exit 1
fi

echo "Installing $BINARY_NAME $LATEST_RELEASE..."

# Download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"

# Create temp directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
echo "Downloading from $DOWNLOAD_URL"
curl -sL "$DOWNLOAD_URL" | tar xz

# Install binary
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Installing to $INSTALL_DIR (requires sudo)"
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

echo "✅ $BINARY_NAME installed successfully!"
echo ""
echo "To use better-fg, add this to your shell config (~/.bashrc, ~/.zshrc):"
echo "  eval \"\$(better-fg init)\""
echo ""
echo "Then restart your shell and use the 'bfg' command."
