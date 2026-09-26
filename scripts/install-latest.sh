#!/usr/bin/env bash

set -Eeuo pipefail

REPO_OWNER="X-Calibre"
REPO_NAME="MasjidFrame"
RELEASE_API="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"

info() { printf 'MasjidFrame: %s\n' "$*"; }
die() { printf 'MasjidFrame: ERROR: %s\n' "$*" >&2; exit 1; }

require_root() {
    if [[ "$(id -u)" -ne 0 ]]; then
        exec sudo bash "$0" "$@"
    fi
}

main() {
    require_root "$@"

    command -v curl >/dev/null 2>&1 || die "curl is required."
    command -v tar >/dev/null 2>&1 || die "tar is required."

    local arch release_arch latest_tag archive_name base_url work_dir release_dir

    arch="$(uname -m)"
    case "$arch" in
        aarch64)
            release_arch="arm64"
            ;;
        x86_64)
            release_arch="amd64"
            ;;
        armv7l|armv6l)
            die "No official MasjidFrame release is available for $arch. Use the source installer instead."
            ;;
        *)
            die "Unsupported architecture: $arch"
            ;;
    esac

    latest_tag="$(curl -fsSL --retry 3 "$RELEASE_API" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
    [[ "$latest_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "Unable to determine the latest MasjidFrame release."

    archive_name="masjidframe-${latest_tag}-linux-${release_arch}.tar.gz"
    base_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${latest_tag}"
    work_dir="$(mktemp -d /tmp/masjidframe-latest.XXXXXX)"
    trap 'rm -rf "$work_dir"' EXIT

    info "Downloading MasjidFrame ${latest_tag} (${release_arch})..."
    curl -fL --retry 3 "${base_url}/${archive_name}" -o "$work_dir/$archive_name"
    curl -fsSL --retry 3 "${base_url}/SHA256SUMS" -o "$work_dir/SHA256SUMS"

    info "Verifying release checksum..."
    (
        cd "$work_dir"
        grep -F "  $archive_name" SHA256SUMS | sha256sum -c -
    ) || die "Release checksum verification failed."

    info "Extracting release..."
    tar -xzf "$work_dir/$archive_name" -C "$work_dir"
    release_dir="$work_dir/masjidframe-${latest_tag}-linux-${release_arch}"

    [[ -x "$release_dir/masjidframe" ]] || die "Release binary is missing."
    [[ -f "$release_dir/default.yaml" ]] || die "Release configuration is missing."
    [[ -f "$release_dir/VERSION" ]] || die "Release version file is missing."
    [[ -f "$release_dir/frontend/index.html" ]] || die "Release frontend is missing."
    [[ -x "$release_dir/scripts/install.sh" ]] || die "Release installer is missing."
    [[ -f "$release_dir/scripts/masjidframe.service" ]] || die "Release service file is missing."
    [[ -f "$release_dir/scripts/masjidframe-display.service" ]] || die "MasjidBoard display service file is missing."

    info "Installing MasjidFrame ${latest_tag}..."

    # The recommended install command pipes this bootstrap script into bash,
    # which means stdin is not a TTY. Reconnect the bundled installer to the
    # controlling terminal when one exists so users are prompted to choose
    # Listen, Board, or Listen + Board. Unattended installs retain the bundled
    # installer's documented non-interactive default.
    if [[ -r /dev/tty && -w /dev/tty ]]; then
        exec "$release_dir/scripts/install.sh" </dev/tty >/dev/tty
    fi

    exec "$release_dir/scripts/install.sh"
}

main "$@"
