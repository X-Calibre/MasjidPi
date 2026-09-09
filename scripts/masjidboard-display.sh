#!/usr/bin/env bash

set -Eeuo pipefail

MASJIDBOARD_BASE_URL="${MASJIDBOARD_BASE_URL:-http://127.0.0.1:8080/masjidboard.html}"
MASJIDBOARD_STARTUP_FILE="${MASJIDBOARD_STARTUP_FILE:-/opt/masjidpi/frontend/masjidboard-startup.html}"
MASJIDPI_READY_URL="${MASJIDPI_READY_URL:-http://127.0.0.1:8080/api/version}"
MASJIDPI_DRM_SYSFS_ROOT="${MASJIDPI_DRM_SYSFS_ROOT:-/sys/class/drm}"
MASJIDPI_RPI_MODEL_FILE="${MASJIDPI_RPI_MODEL_FILE:-/proc/device-tree/model}"

wait_for_masjidpi() {
    while ! curl --silent --show-error --fail --max-time 2 \
        "$MASJIDPI_READY_URL" >/dev/null 2>&1; do
        sleep 2
    done
}

is_raspberry_pi_runtime() {
    if [[ "${MASJIDPI_FORCE_RASPBERRY_PI:-0}" == "1" ]]; then
        return 0
    fi

    [[ -r "$MASJIDPI_RPI_MODEL_FILE" ]] || return 1
    grep -aqi 'Raspberry Pi' "$MASJIDPI_RPI_MODEL_FILE"
}

dsi_720x1280_present() {
    local connector
    for connector in "$MASJIDPI_DRM_SYSFS_ROOT"/card*-DSI-*; do
        [[ -r "$connector/status" && -r "$connector/modes" ]] || continue
        [[ "$(<"$connector/status")" == "connected" ]] || continue
        grep -Fxq '720x1280' "$connector/modes" && return 0
    done
    return 1
}

display_profile() {
    if dsi_720x1280_present; then
        printf 'appliance-720\n'
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

    if [[ "$profile" == "appliance-720" ]]; then
        if [[ "$MASJIDBOARD_BASE_URL" == *\?* ]]; then
            printf '%s&profile=%s\n' "$MASJIDBOARD_BASE_URL" "$profile"
        else
            printf '%s?profile=%s\n' "$MASJIDBOARD_BASE_URL" "$profile"
        fi
        return
    fi
    printf '%s\n' "$MASJIDBOARD_BASE_URL"
}

uses_startup_screen() {
    local profile="$1"

    [[ -z "${MASJIDBOARD_URL:-}" ]] || return 1
    [[ "$profile" == "appliance-720" ]] && return 0
    is_raspberry_pi_runtime
}

launch_url() {
    local profile="$1"

    if ! uses_startup_screen "$profile"; then
        display_url "$profile"
        return
    fi

    # Raspberry Pi Board installations load the startup page directly from the
    # installed runtime, so Cog can paint branding before the HTTP server is up.
    printf 'file://%s?profile=%s\n' "$MASJIDBOARD_STARTUP_FILE" "$profile"
}

main() {
    if ! command -v cog >/dev/null 2>&1; then
        echo "MasjidBoard display: Cog is not installed." >&2
        exit 1
    fi

    local profile platform_params
    local -a cog_args
    profile="$(display_profile)"

    # Raspberry Pi Board installations launch a local startup page immediately;
    # that page performs its own backend readiness poll. Other standard installs
    # retain the established wait-before-launch behaviour.
    if ! uses_startup_screen "$profile"; then
        wait_for_masjidpi
    fi

    platform_params="renderer=gles"
    cog_args=(--platform=drm --bg-color='#050f0d')

    cog_args+=(--platform-params="$platform_params")

    echo "MasjidBoard display: detected profile=$profile" >&2

    # GLES is compatible with Raspberry Pi DRM/KMS devices whose preferred
    # framebuffer format cannot be used by Cog's direct modeset renderer.
    exec cog "${cog_args[@]}" "$(launch_url "$profile")"
}

if [[ "${MASJIDBOARD_DISPLAY_LIBRARY_ONLY:-0}" != "1" ]]; then
    main "$@"
fi
