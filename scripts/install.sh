#!/usr/bin/env bash
set -euo pipefail
# Install linkding from GitHub Releases (draft or published).
# Usage: curl -fsSL ... | bash   or   bash scripts/install.sh [version]
set -u
REPO="swtormy-agents/linkding-cli"
BIN="linkding"
VERSION="${1:-latest}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac
case "$os" in
  linux|darwin) ;;
  mingw*|msys*|cygwin*|windows*) os="windows" ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

if [[ "$VERSION" == "latest" ]]; then
  url="https://github.com/${REPO}/releases/latest/download/${BIN}_${os}_${arch}.tar.gz"
  if [[ "$os" == "windows" ]]; then
    url="https://github.com/${REPO}/releases/latest/download/${BIN}_${os}_${arch}.zip"
  fi
else
  url="https://github.com/${REPO}/releases/download/${VERSION}/${BIN}_${VERSION#v}_${os}_${arch}.tar.gz"
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT
echo "Downloading $url"
curl -fsSL "$url" -o "$tmpdir/archive"
# shellcheck disable=SC2164
cd "$tmpdir"
if [[ "$url" == *.zip ]]; then
  unzip -q archive
else
  tar -xzf archive
fi
install -m 755 "$BIN" "${INSTALL_DIR:-/usr/local/bin}/$BIN"
echo "Installed $(command -v "$BIN" 2>/dev/null || echo "${INSTALL_DIR:-/usr/local/bin}/$BIN")"
