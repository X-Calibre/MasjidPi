#!/usr/bin/env bash

run_selftest() {

    info "Running self test..."

    command -v mpv >/dev/null || die "mpv not found"

    command -v go >/dev/null || die "Go not found"

    [[ -f "$PROJECT_ROOT/backend/data/catalogue.json" ]] \
        || die "Catalogue missing"

    cd "$PROJECT_ROOT/backend"

    timeout 5 go run ./cmd/masjidpi >/tmp/masjidpi.log 2>&1 || true

    grep -q "Starting HTTP server" /tmp/masjidpi.log \
        || die "Application failed to start."

    success "Application started successfully."

    if aplay -l >/dev/null 2>&1; then
        success "Audio device detected."
    else
        warn "No audio device found."
    fi
}