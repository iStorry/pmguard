#!/bin/sh
set -e

REPO="istorry/pmguard"
BINARY="pmguard"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux|darwin) ;;
  *)
    echo "❌ Unsupported OS: $OS"
    echo "   Windows users: download from https://github.com/$REPO/releases"
    exit 1
    ;;
esac

# Get latest version tag from GitHub
VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')"

if [ -z "$VERSION" ]; then
  echo "❌ Could not fetch latest version from GitHub"
  exit 1
fi

FILENAME="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "⬇️  Downloading pmguard $VERSION for $OS/$ARCH..."
curl -fsSL "$URL" -o "/tmp/$FILENAME"

echo "📦 Extracting..."
tar -xzf "/tmp/$FILENAME" -C /tmp

echo "🔧 Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
  mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
else
  sudo mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
fi

chmod +x "$INSTALL_DIR/$BINARY"
rm -f "/tmp/$FILENAME"

echo ""
echo "✅ pmguard installed!"
echo ""
echo "👉 Add this to your ~/.zshrc or ~/.bashrc to activate:"
echo '   eval "$(pmguard install-hooks)"'
echo ""
echo "   Then restart your shell or run: source ~/.zshrc"