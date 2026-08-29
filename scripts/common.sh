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

primary_ipv4_address() {
    local address=""

    if command_exists ip; then
        address="$(
            ip -4 route get 1.1.1.1 2>/dev/null \
                | awk '{for (i = 1; i <= NF; i++) if ($i == "src") {print $(i + 1); exit}}'
        )"
    fi

    if [[ -z "$address" ]] && command_exists hostname; then
        address="$(hostname -I 2>/dev/null | awk '{for (i = 1; i <= NF; i++) if ($i !~ /^127\./ && $i !~ /^169\.254\./) {print $i; exit}}')"
    fi

    printf '%s\n' "$address"
}

web_hostname() {
    command_exists hostname || return

    local fqdn short_name
    fqdn="$(hostname --fqdn 2>/dev/null || true)"
    short_name="$(hostname --short 2>/dev/null || true)"

    if [[ "$fqdn" == *.* && "$fqdn" != "localhost.localdomain" ]]; then
        printf '%s\n' "$fqdn"
        return
    fi

    if [[ -n "$short_name" ]] \
        && command_exists systemctl \
        && systemctl is-active --quiet avahi-daemon.service 2>/dev/null; then
        printf '%s.local\n' "$short_name"
    fi
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

    local ip_address host_name
    ip_address="$(primary_ipv4_address)"
    host_name="$(web_hostname)"

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
    if [[ -n "$ip_address" ]]; then
        echo "    IP address:  http://${ip_address}:8080"
    fi
    if [[ -n "$host_name" ]]; then
        echo "    Hostname:    http://${host_name}:8080"
    fi
    if [[ -z "$ip_address" && -z "$host_name" ]]; then
        echo "    http://localhost:8080"
    fi
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
