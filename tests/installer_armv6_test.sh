#!/usr/bin/env bash
set -Eeuo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

source "$ROOT/scripts/components.sh"
source "$ROOT/scripts/packages.sh"
source "$ROOT/scripts/go.sh"
source "$ROOT/scripts/armv6-listen.sh"
info() { :; }
success() { :; }
die() { printf '%s\n' "$*" >&2; exit 1; }
uname() { printf '%s\n' "$MOCK_ARCH"; }
apt-get() { printf '%s\n' "$*" >> "$TMP/packages"; }
# Redirect only the system configuration directory operation during tests.
install() { :; }

MOCK_ARCH=armv6l
export MASJIDFRAME_ARMV6_LISTEN=1
MASJIDFRAME_DEVICE_MODEL_PATH="$TMP/model"
printf 'Raspberry Pi Zero W Rev 1.1\0' > "$TMP/model"
detect_arch
[[ "$ARCH" == armv6l ]]
[[ "$(get_go_arch)" == linux-armv6l ]]

COMPONENTS_FILE="$TMP/components.env"
printf 'MASJIDFRAME_COMPONENTS=listen,board\n' > "$COMPONENTS_FILE"
select_components </dev/null
[[ "$(component_profile)" == listen ]]
[[ "$PREVIOUS_COMPONENT_PROFILE" == listen,board ]]
save_components
grep -qx 'MASJIDFRAME_COMPONENTS=listen' "$COMPONENTS_FILE"
grep -qx 'MASJIDFRAME_APPLIANCE_UPDATES=disabled' "$COMPONENTS_FILE"
install_packages release
grep -q 'mpv ffmpeg alsa-utils' "$TMP/packages"
if grep -Eq 'cog|webkit|plymouth|build-essential|git' "$TMP/packages"; then
    printf 'Unexpected Board or build dependencies\n' >&2
    exit 1
fi

for MOCK_ARCH in aarch64 armv7l x86_64; do
    if (detect_arch) >/dev/null 2>&1; then
        printf 'Wrong architecture accepted: %s\n' "$MOCK_ARCH" >&2
        exit 1
    fi
done
MOCK_ARCH=armv6l
printf 'Other ARMv6 hardware\n' > "$TMP/model"
if (detect_arch) >/dev/null 2>&1; then
    printf 'Non-Pi hardware accepted\n' >&2
    exit 1
fi

# ARMv6 support must not silently broaden the general installer's scope.
source "$ROOT/scripts/os.sh"
if (detect_arch) >/dev/null 2>&1; then
    printf 'General installer unexpectedly accepted ARMv6\n' >&2
    exit 1
fi

# Exercise release download/extraction with a local archive and mocked network.
# Intentionally source function definitions without running main.
# shellcheck disable=SC1090
source <(sed '$d' "$ROOT/scripts/install-listen-armv6-latest.sh")
MOCK_ARCH=armv6l
printf 'Raspberry Pi Zero W\n' > "$TMP/model"
id() { printf '0\n'; }
fixture="$TMP/fixture/masjidframe-v9.0.0-linux-armv6"
mkdir -p "$fixture/scripts"
printf '#!/bin/sh\nexit 0\n' > "$fixture/masjidframe"
printf '#!/bin/sh\nprintf installed > "%s"\n' "$TMP/installed" > "$fixture/scripts/install-listen-armv6.sh"
touch "$fixture/scripts/armv6-listen.sh"
chmod +x "$fixture/masjidframe" "$fixture/scripts/install-listen-armv6.sh"
archive=masjidframe-v9.0.0-linux-armv6.tar.gz
tar -czf "$TMP/$archive" -C "$TMP/fixture" "$(basename "$fixture")"
(cd "$TMP" && sha256sum "$archive" > SHA256SUMS)
curl() {
    local url="" output=""
    while [[ $# -gt 0 ]]; do
        case "$1" in
            https:*) url="$1" ;;
            -o) shift; output="$1" ;;
        esac
        shift
    done
    case "$url" in
        */releases/latest) printf '{"tag_name":"v9.0.0"}\n' ;;
        */SHA256SUMS) cp "$TMP/SHA256SUMS" "$output" ;;
        *.tar.gz)
            printf '%s\n' "$(dirname "$output")" > "$TMP/download-dir"
            cp "$TMP/$archive" "$output"
            ;;
        *) return 1 ;;
    esac
}
(main) >/dev/null
[[ -f "$TMP/installed" ]]
[[ ! -d "$(cat "$TMP/download-dir")" ]]
rm "$TMP/installed"
printf '%064d  %s\n' 0 "$archive" > "$TMP/SHA256SUMS"
if (main) >/dev/null 2>&1; then
    printf 'Bad checksum accepted\n' >&2
    exit 1
fi
[[ ! -f "$TMP/installed" ]]
[[ ! -d "$(cat "$TMP/download-dir")" ]]
printf '[PASS] ARMv6 hardware gating, forced Listen profile and dependencies\n'
