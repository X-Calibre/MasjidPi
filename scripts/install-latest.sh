#!/usr/bin/env bash

set -Eeuo pipefail

REPO_OWNER="X-Calibre"
REPO_NAME="MasjidPi"
RELEASE_API="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"

info() { printf 'MasjidPi: %s\n' "$*"; }
die() { printf 'MasjidPi: ERROR: %s\n' "$*" >&2; exit 1; }

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
            die "No official MasjidPi release is available for $arch. Use the source installer instead."
            ;;
        *)
            die "Unsupported architecture: $arch"
            ;;
    esac

    latest_tag="$(curl -fsSL --retry 3 "$RELEASE_API" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
    [[ "$latest_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "Unable to determine the latest MasjidPi release."

    archive_name="masjidpi-${latest_tag}-linux-${release_arch}.tar.gz"
    base_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${latest_tag}"
    work_dir="$(mktemp -d /tmp/masjidpi-latest.XXXXXX)"
    trap 'rm -rf "$work_dir"' EXIT

    info "Downloading MasjidPi ${latest_tag} (${release_arch})..."
    curl -fL --retry 3 "${base_url}/${archive_name}" -o "$work_dir/$archive_name"
    curl -fsSL --retry 3 "${base_url}/SHA256SUMS" -o "$work_dir/SHA256SUMS"

    info "Verifying release checksum..."
    (
        cd "$work_dir"
        grep -F "  $archive_name" SHA256SUMS | sha256sum -c -
    ) || die "Release checksum verification failed."

    info "Extracting release..."
    tar -xzf "$work_dir/$archive_name" -C "$work_dir"
    release_dir="$work_dir/masjidpi-${latest_tag}-linux-${release_arch}"

    [[ -x "$release_dir/masjidpi" ]] || die "Release binary is missing."
    [[ -f "$release_dir/default.yaml" ]] || die "Release configuration is missing."
    [[ -f "$release_dir/VERSION" ]] || die "Release version file is missing."
    [[ -f "$release_dir/frontend/index.html" ]] || die "Release frontend is missing."
    [[ -x "$release_dir/scripts/install.sh" ]] || die "Release installer is missing."

    info "Installing MasjidPi ${latest_tag}..."
    exec "$release_dir/scripts/install.sh"
}

main "$@"
