#!/usr/bin/env bash

run_selftest() {

    info "Running self test..."

    command -v mpv >/dev/null || die "mpv not found"
    command -v go >/dev/null || die "Go not found"

    [[ -f "$PROJECT_ROOT/backend/data/catalogue.json" ]] \
        || die "Catalogue missing"

    cd "$PROJECT_ROOT/backend"

    LOG="/tmp/masjidpi-selftest.log"

    info "Starting MasjidPi..."

    ./masjidpi >"$LOG" 2>&1 &
    PID=$!

    # Wait up to 5 seconds for the HTTP server to start.
    for i in {1..50}; do
        if grep -q "Starting HTTP server" "$LOG"; then
            break
        fi

        sleep 0.1
    done

    grep -q "Starting HTTP server" "$LOG" \
        || {
            kill "$PID" 2>/dev/null || true
            wait "$PID" 2>/dev/null || true
            die "Application failed to start."
        }

    success "Application started successfully."

    info "Stopping MasjidPi..."

    kill "$PID" 2>/dev/null || true
    wait "$PID" 2>/dev/null || true

    if aplay -l >/dev/null 2>&1; then
        success "Audio device detected."
    else
        warn "No audio device found."
    fi
}