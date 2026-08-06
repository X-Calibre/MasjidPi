#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

INSTALL_DIR="/opt/masjidpi"
SERVICE_NAME="masjidpi"

RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
RESET="\033[0m"

info() {
    echo -e "${BLUE}[INFO]${RESET} $*"
}

success() {
    echo -e "${GREEN}[ OK ]${RESET} $*"
}

warn() {
    echo -e "${YELLOW}[WARN]${RESET} $*"
}

error() {
    echo -e "${RED}[FAIL]${RESET} $*" >&2
}

die() {
    error "$*"
    exit 1
}

require_root() {
    [[ $EUID -eq 0 ]] || die "Run installer with sudo."
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

detect_os() {

    source /etc/os-release

    DISTRO="$ID"
    VERSION="$VERSION_ID"

    info "Detected $PRETTY_NAME"
}

detect_arch() {

    ARCH="$(uname -m)"

    case "$ARCH" in
        x86_64|aarch64|armv7l)
            ;;
        *)
            die "Unsupported architecture: $ARCH"
            ;;
    esac
}

print_summary() {

    VERSION="$(get_version)"

    echo
    echo "========================================="
    echo " MasjidPi installation complete"
    echo "========================================="
    echo

    echo "Version"
    echo
    echo "    $VERSION"
    echo

    echo "Web Interface"
    echo
    echo "    http://localhost:8080"
    echo

    echo "Service"
    echo
    echo "    Active and running"
    echo

    echo "Useful commands"
    echo
    echo "    sudo systemctl status masjidpi"
    echo "    sudo systemctl restart masjidpi"
    echo "    sudo systemctl stop masjidpi"
    echo "    journalctl -u masjidpi -f"
    echo
    echo "========================================="
}