#!/usr/bin/env bash

set -Eeuo pipefail

MASJIDBOARD_URL="${MASJIDBOARD_URL:-http://127.0.0.1:8080/masjidboard.html}"
MASJIDPI_READY_URL="${MASJIDPI_READY_URL:-http://127.0.0.1:8080/api/version}"

wait_for_masjidpi() {
    while ! curl --silent --show-error --fail --max-time 2 \
        "$MASJIDPI_READY_URL" >/dev/null 2>&1; do
        sleep 2
    done
}

main() {
    if ! command -v cog >/dev/null 2>&1; then
        echo "MasjidBoard display: Cog is not installed." >&2
        exit 1
    fi

    wait_for_masjidpi

    exec cog --platform=drm "$MASJIDBOARD_URL"
}

main "$@"
