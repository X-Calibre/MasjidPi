#!/usr/bin/env bash

set -Eeuo pipefail

MASJIDBOARD_BASE_URL="${MASJIDBOARD_BASE_URL:-http://127.0.0.1:8080/masjidboard.html}"
MASJIDBOARD_LAYOUT_URL="${MASJIDBOARD_LAYOUT_URL:-http://127.0.0.1:8080/api/masjidboard/layout}"
MASJIDPI_READY_URL="${MASJIDPI_READY_URL:-http://127.0.0.1:8080/api/version}"

wait_for_masjidpi() {
    while ! curl --silent --show-error --fail --max-time 2 \
        "$MASJIDPI_READY_URL" >/dev/null 2>&1; do
        sleep 2
    done
}

saved_layout() {
    curl --silent --show-error --fail --max-time 3 "$MASJIDBOARD_LAYOUT_URL" 2>/dev/null \
        | sed -n 's/.*"layout"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p'
}

display_url() {
    # An explicit MASJIDBOARD_URL remains an advanced override for testing or
    # custom deployments. Normal appliance HDMI output follows the saved WebUI
    # layout preference.
    if [[ -n "${MASJIDBOARD_URL:-}" ]]; then
        printf '%s\n' "$MASJIDBOARD_URL"
        return
    fi

    local layout
    layout="$(saved_layout || true)"
    if [[ "$layout" == "detailed" || "$layout" == "portrait" ]]; then
        if [[ "$MASJIDBOARD_BASE_URL" == *\?* ]]; then
            printf '%s&layout=%s\n' "$MASJIDBOARD_BASE_URL" "$layout"
        else
            printf '%s?layout=%s\n' "$MASJIDBOARD_BASE_URL" "$layout"
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

    exec cog --platform=drm "$(display_url)"
}

main "$@"
