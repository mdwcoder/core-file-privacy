#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="cfp"
REPO_URL="https://github.com/mdwcoder/core-file-privacy"

detect_arch() {
  local arch
  arch="$(uname -m)"
  case "$arch" in
    x86_64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *) echo "$arch" ;;
  esac
}

detect_os() {
  local os
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux) echo "linux" ;;
    darwin) echo "darwin" ;;
    msys*|mingw*|cygwin*) echo "windows" ;;
    *) echo "$os" ;;
  esac
}

download_binary() {
  local os="$1"
  local arch="$2"
  local target="$INSTALL_DIR/$BINARY_NAME"

  local asset_name="cfp_${os}_${arch}.tar.gz"
  if [ "$os" = "windows" ]; then
    asset_name="cfp_${os}_${arch}.zip"
  fi

  local url="${REPO_URL}/releases/latest/download/${asset_name}"

  echo "Downloading precompiled binary: $url"
  mkdir -p "$INSTALL_DIR"

  local tmp_dir
  tmp_dir="$(mktemp -d)"
  trap 'rm -rf "$tmp_dir"' EXIT

  if ! curl -fsSL -o "$tmp_dir/$asset_name" "$url"; then
    echo "Failed to download binary. You can build from source if Go is installed." >&2
    return 1
  fi

  if [ "$os" = "windows" ]; then
    unzip -q "$tmp_dir/$asset_name" -d "$tmp_dir"
  else
    tar -xzf "$tmp_dir/$asset_name" -C "$tmp_dir"
  fi

  if [ -f "$tmp_dir/$BINARY_NAME" ]; then
    cp "$tmp_dir/$BINARY_NAME" "$target"
  elif [ -f "$tmp_dir/${BINARY_NAME}.exe" ]; then
    cp "$tmp_dir/${BINARY_NAME}.exe" "${target}.exe"
  else
    echo "Binary not found in downloaded archive." >&2
    return 1
  fi

  chmod +x "$target" 2>/dev/null || true
  echo "Installed: $target"
}

build_from_source() {
  local target="$INSTALL_DIR/$BINARY_NAME"
  echo "Building from source with Go..."
  cd "$REPO_DIR"
  go build -buildvcs=false -ldflags="-s -w" -o "$BINARY_NAME" ./cmd/cfp
  mkdir -p "$INSTALL_DIR"
  cp "$BINARY_NAME" "$target"
  chmod +x "$target"
  echo "Built and installed: $target"
}

add_to_path() {
  local shell_rc=""
  case "${SHELL:-}" in
    */zsh) shell_rc="${HOME}/.zshrc" ;;
    */bash) shell_rc="${HOME}/.bashrc" ;;
    */fish) shell_rc="${HOME}/.config/fish/config.fish" ;;
  esac

  if [ -n "$shell_rc" ] && [ -f "$shell_rc" ]; then
    if ! grep -q "$INSTALL_DIR" "$shell_rc" 2>/dev/null; then
      echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$shell_rc"
      echo "Added $INSTALL_DIR to PATH in $shell_rc"
      echo "Run 'source $shell_rc' to update your current session."
    fi
  fi

  if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    export PATH="$INSTALL_DIR:$PATH"
  fi
}

main() {
  local os arch
  os="$(detect_os)"
  arch="$(detect_arch)"

  if [ "$os" = "windows" ]; then
    INSTALL_DIR="${USERPROFILE:-${HOME}}/bin"
    BINARY_NAME="cfp.exe"
  fi

  if command -v go >/dev/null 2>&1; then
    build_from_source
  else
    echo "Go not found. Attempting to download precompiled binary..."
    if ! download_binary "$os" "$arch"; then
      echo "ERROR: Could not install core-file-privacy." >&2
      echo "Please install Go (https://go.dev/dl/) or download the binary manually from:" >&2
      echo "  $REPO_URL/releases" >&2
      exit 1
    fi
  fi

  add_to_path
  echo ""
  echo "core-file-privacy installed successfully!"
  echo "Run 'cfp --help' to get started."
}

main
