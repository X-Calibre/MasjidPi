#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

export MASJIDBOARD_DISPLAY_LIBRARY_ONLY=1
export MASJIDPI_USB_SYSFS_ROOT="$TMP/usb"
export MASJIDPI_DRM_SYSFS_ROOT="$TMP/drm"
mkdir -p "$MASJIDPI_USB_SYSFS_ROOT" "$MASJIDPI_DRM_SYSFS_ROOT"

# shellcheck source=../scripts/masjidboard-display.sh
source "$ROOT/scripts/masjidboard-display.sh"

assert_profile() {
    local want="$1" got
    got="$(display_profile)"
    if [[ "$got" != "$want" ]]; then
        echo "expected profile=$want, got profile=$got" >&2
        exit 1
    fi
}

assert_profile standard

usb="$MASJIDPI_USB_SYSFS_ROOT/1-1.3"
mkdir -p "$usb"
printf '0eef\n' > "$usb/idVendor"
printf '0005\n' > "$usb/idProduct"
printf 'WaveShare\n' > "$usb/manufacturer"
printf 'WS170120\n' > "$usb/product"
assert_profile standard

hdmi="$MASJIDPI_DRM_SYSFS_ROOT/card0-HDMI-A-1"
mkdir -p "$hdmi"
printf 'connected\n' > "$hdmi/status"
printf '1024x600\n1920x1080\n' > "$hdmi/modes"
assert_profile appliance

printf 'disconnected\n' > "$hdmi/status"
assert_profile standard

printf 'connected\n' > "$hdmi/status"
printf '1920x1080\n' > "$hdmi/modes"
assert_profile standard

printf '1024x600\n' > "$hdmi/modes"
printf 'OtherVendor\n' > "$usb/manufacturer"
assert_profile standard

printf 'WaveShare\n' > "$usb/manufacturer"
printf 'WS170120\n' > "$usb/product"
assert_profile appliance

if [[ "$(display_url standard)" != "$MASJIDBOARD_BASE_URL" ]]; then
    echo "standard display URL is incorrect" >&2
    exit 1
fi
if [[ "$(display_url appliance)" != "${MASJIDBOARD_BASE_URL}?profile=appliance" ]]; then
    echo "appliance display URL is incorrect" >&2
    exit 1
fi
if [[ "$(launch_url standard)" != "$MASJIDBOARD_BASE_URL" ]]; then
    echo "standard launch URL must remain the board URL" >&2
    exit 1
fi
if [[ "$(launch_url appliance)" != "${MASJIDBOARD_STARTUP_URL}?profile=appliance" ]]; then
    echo "appliance launch URL must use the startup screen" >&2
    exit 1
fi

MASJIDBOARD_BASE_URL='http://127.0.0.1:8080/masjidboard.html?date=2026-08-31'
if [[ "$(display_url appliance)" != "${MASJIDBOARD_BASE_URL}&profile=appliance" ]]; then
    echo "appliance display URL did not preserve existing query parameters" >&2
    exit 1
fi

MASJIDBOARD_URL='http://127.0.0.1:8080/custom-board.html'
if [[ "$(launch_url appliance)" != "$MASJIDBOARD_URL" ]]; then
    echo "custom display URL override must bypass the appliance startup screen" >&2
    exit 1
fi

printf 'MasjidBoard display profile tests passed\n'
