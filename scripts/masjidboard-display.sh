#!/usr/bin/env bash

set -Eeuo pipefail

MASJIDBOARD_BASE_URL="${MASJIDBOARD_BASE_URL:-http://127.0.0.1:8080/masjidboard.html}"
MASJIDPI_READY_URL="${MASJIDPI_READY_URL:-http://127.0.0.1:8080/api/version}"
MASJIDPI_USB_SYSFS_ROOT="${MASJIDPI_USB_SYSFS_ROOT:-/sys/bus/usb/devices}"
MASJIDPI_DRM_SYSFS_ROOT="${MASJIDPI_DRM_SYSFS_ROOT:-/sys/class/drm}"

wait_for_masjidpi() {
    while ! curl --silent --show-error --fail --max-time 2 \
        "$MASJIDPI_READY_URL" >/dev/null 2>&1; do
        sleep 2
    done
}

waveshare_touch_present() {
    local device
    for device in "$MASJIDPI_USB_SYSFS_ROOT"/*; do
        [[ -r "$device/idVendor" && -r "$device/idProduct" ]] || continue
        [[ "$(<"$device/idVendor")" == "0eef" ]] || continue
        [[ "$(<"$device/idProduct")" == "0005" ]] || continue

        # Some kernels/udev timing windows may expose VID/PID before optional
        # descriptive strings. When the strings are present, require the
        # validated Waveshare identity as an additional guard.
        if [[ -r "$device/manufacturer" && -r "$device/product" ]]; then
            [[ "$(<"$device/manufacturer")" == "WaveShare" ]] || continue
            [[ "$(<"$device/product")" == "WS170120" ]] || continue
        fi
        return 0
    done
    return 1
}

hdmi_1024x600_present() {
    local connector
    for connector in "$MASJIDPI_DRM_SYSFS_ROOT"/card*-HDMI-A-*; do
        [[ -r "$connector/status" && -r "$connector/modes" ]] || continue
        [[ "$(<"$connector/status")" == "connected" ]] || continue
        grep -Fxq '1024x600' "$connector/modes" && return 0
    done
    return 1
}

display_profile() {
    if waveshare_touch_present && hdmi_1024x600_present; then
        printf 'appliance\n'
    else
        printf 'standard\n'
    fi
}

display_url() {
    local profile="$1"

    # An explicit MASJIDBOARD_URL remains an advanced override for testing or
    # custom deployments. Normal HDMI output follows the detected hardware
    # profile and never a persisted user preference.
    if [[ -n "${MASJIDBOARD_URL:-}" ]]; then
        printf '%s\n' "$MASJIDBOARD_URL"
        return
    fi

    if [[ "$profile" == "appliance" ]]; then
        if [[ "$MASJIDBOARD_BASE_URL" == *\?* ]]; then
            printf '%s&profile=appliance\n' "$MASJIDBOARD_BASE_URL"
        else
            printf '%s?profile=appliance\n' "$MASJIDBOARD_BASE_URL"
        fi
        return
    fi
    printf '%s\n' "$MASJIDBOARD_BASE_URL"
}

main() {
    if ! command -v cog >/dev/null 2>&1; then
        echo "MasjidBoard display: Cog is not installed." >&2
        exit 1
    fi

    wait_for_masjidpi

    local profile platform_params
    profile="$(display_profile)"
    platform_params="renderer=gles"
    if [[ "$profile" == "appliance" ]]; then
        platform_params+=",rotation=1"
    fi

    echo "MasjidBoard display: detected profile=$profile" >&2

    # GLES is compatible with Raspberry Pi DRM/KMS devices whose preferred
    # framebuffer format cannot be used by Cog's direct modeset renderer.
    exec cog --platform=drm --platform-params="$platform_params" "$(display_url "$profile")"
}

if [[ "${MASJIDBOARD_DISPLAY_LIBRARY_ONLY:-0}" != "1" ]]; then
    main "$@"
fi
