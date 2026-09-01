#!/usr/bin/env bash

set -u

STARTUP_FILE="${MASJIDBOARD_STARTUP_FILE:-/opt/masjidpi/frontend/masjidboard-startup.html}"
WARMUP_TIMEOUT="${MASJIDBOARD_WARMUP_TIMEOUT:-20}"
WARMUP_PROFILE="${MASJIDBOARD_WARMUP_PROFILE:-standard}"

if [[ "$WARMUP_PROFILE" != "appliance" ]]; then
    WARMUP_PROFILE="standard"
fi

WARMUP_URL="file://${STARTUP_FILE}?profile=${WARMUP_PROFILE}"

if ! command -v cog >/dev/null 2>&1; then
    echo "MasjidBoard warm-up: Cog is not installed; skipping." >&2
    exit 0
fi

if [[ ! -f "$STARTUP_FILE" ]]; then
    echo "MasjidBoard warm-up: startup page is unavailable; skipping." >&2
    exit 0
fi

# The headless Cog backend initializes WPE/WebKit and decodes the shared splash
# artwork without touching DRM. This lets Plymouth remain visible while the
# expensive cold-start work is brought into memory. The page itself is profile-
# aware, but either profile warms the same WebKit path and static assets.
coproc COG_WARMUP { cog --platform=headless "$WARMUP_URL" 2>&1; }
warmup_pid="$COG_WARMUP_PID"
loaded=false

(
    sleep "$WARMUP_TIMEOUT"
    kill -TERM "$warmup_pid" 2>/dev/null || true
) &
watchdog_pid=$!

while IFS= read -r line <&"${COG_WARMUP[0]}"; do
    printf '%s\n' "$line" >&2
    if [[ "$line" == *"Loaded successfully."* ]]; then
        loaded=true
        break
    fi
done

kill "$watchdog_pid" 2>/dev/null || true
wait "$watchdog_pid" 2>/dev/null || true
kill -TERM "$warmup_pid" 2>/dev/null || true
wait "$warmup_pid" 2>/dev/null || true

if $loaded; then
    echo "MasjidBoard warm-up: WPE/WebKit cold start completed." >&2
else
    echo "MasjidBoard warm-up: warm-up did not complete before timeout; continuing boot." >&2
fi

exit 0
