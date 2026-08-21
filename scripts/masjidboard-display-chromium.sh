#!/usr/bin/env bash

set -Eeuo pipefail

MASJIDBOARD_URL="${MASJIDBOARD_URL:-http://127.0.0.1:8080/masjidboard.html}"
MASJIDPI_READY_URL="${MASJIDPI_READY_URL:-http://127.0.0.1:8080/api/version}"

find_browser() {
    if [[ -x /usr/lib/chromium/chromium ]]; then
        printf '%s\n' /usr/lib/chromium/chromium
        return 0
    fi

    local candidate
    for candidate in chromium chromium-browser; do
        if command -v "$candidate" >/dev/null 2>&1; then
            command -v "$candidate"
            return 0
        fi
    done

    echo "MasjidBoard display: Chromium is not installed." >&2
    return 1
}

wait_for_masjidpi() {
    while ! curl --silent --show-error --fail --max-time 2 \
        "$MASJIDPI_READY_URL" >/dev/null 2>&1; do
        sleep 2
    done
}

main() {
    local browser
    browser="$(find_browser)"

    wait_for_masjidpi

    rm -rf /tmp/masjidpi-display-profile /tmp/masjidpi-display-cache
    mkdir -p /tmp/masjidpi-display-profile /tmp/masjidpi-display-cache

    exec "$browser" \
        --ozone-platform=wayland \
        --use-angle=gles \
        --enable-gpu-rasterization \
        --kiosk \
        --start-maximized \
        --no-first-run \
        --no-default-browser-check \
        --noerrdialogs \
        --disable-infobars \
        --disable-session-crashed-bubble \
        --disable-notifications \
        --disable-translate \
        --disable-pinch \
        --disable-sync \
        --disable-background-networking \
        --disable-component-extensions-with-background-pages \
        --disable-component-update \
        --disable-default-apps \
        --disable-domain-reliability \
        --disable-extensions \
        --disable-features=MediaRouter,OptimizationHints,Translate \
        --password-store=basic \
        --overscroll-history-navigation=0 \
        --user-data-dir=/tmp/masjidpi-display-profile \
        --disk-cache-dir=/tmp/masjidpi-display-cache \
        "$MASJIDBOARD_URL"
}

main "$@"
