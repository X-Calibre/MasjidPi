#!/usr/bin/env bash
set -Eeuo pipefail

REPO_OWNER="X-Calibre"
REPO_NAME="MasjidPi"
RELEASE_API="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"

die() { printf 'MasjidFrame Listen: ERROR: %s\n' "$*" >&2; exit 1; }

main() {
    [[ $# -eq 0 ]] || die "This bootstrap accepts no arguments. Use the local installer for --source."
    [[ "$(id -u)" -eq 0 ]] || die "Run this installer with sudo bash."
    [[ "$(uname -m)" == armv6l ]] || die "This installer requires an ARMv6 Raspberry Pi."
    local model_path="${MASJIDFRAME_DEVICE_MODEL_PATH:-/proc/device-tree/model}"
    [[ -r "$model_path" ]] && grep -aq 'Raspberry Pi' "$model_path" \
        || die "This installer requires Raspberry Pi hardware."
    local command
    for command in curl tar sha256sum; do
        command -v "$command" >/dev/null || die "$command is required."
    done

    local latest_tag archive_name base_url work_dir release_dir
    latest_tag="$(curl -fsSL --retry 3 "$RELEASE_API" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
    [[ "$latest_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "Unable to determine the latest stable release."
    archive_name="masjidframe-${latest_tag}-linux-armv6.tar.gz"
    base_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${latest_tag}"
    work_dir="$(mktemp -d /tmp/masjidframe-armv6.XXXXXX)"
    local cleanup_path
    printf -v cleanup_path '%q' "$work_dir"
    trap "rm -rf -- $cleanup_path" EXIT

    printf 'MasjidFrame Listen: Downloading %s for ARMv6...\n' "$latest_tag"
    curl -fL --retry 3 "${base_url}/${archive_name}" -o "$work_dir/$archive_name" \
        || die "This release has no downloadable ARMv6 package. Use a release with ARMv6 support or the local --source installer."
    curl -fsSL --retry 3 "${base_url}/SHA256SUMS" -o "$work_dir/SHA256SUMS"
    (
        cd "$work_dir"
        grep -F "  $archive_name" SHA256SUMS | sha256sum -c -
    ) || die "Release checksum verification failed."
    tar -xzf "$work_dir/$archive_name" -C "$work_dir"
    release_dir="$work_dir/masjidframe-${latest_tag}-linux-armv6"
    [[ -x "$release_dir/masjidframe" ]] || die "Release binary is missing."
    [[ -f "$release_dir/scripts/armv6-listen.sh" ]] || die "ARMv6 installer support is missing."
    [[ -x "$release_dir/scripts/install-listen-armv6.sh" ]] || die "Listen installer is missing."

    # Run as a child so the EXIT trap removes the downloaded package.
    bash "$release_dir/scripts/install-listen-armv6.sh"
}

main "$@"
