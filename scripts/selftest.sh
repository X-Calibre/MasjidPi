#!/usr/bin/env bash

run_selftest() {

    info "Running self test..."

    if ! systemctl is-active --quiet masjidpi; then
        die "MasjidPi service is not running."
    fi

    success "MasjidPi service is running."

    info "Checking HTTP interface..."

    for i in {1..50}; do

        if curl -fs http://localhost:8080 >/dev/null 2>&1; then
            success "HTTP interface is responding."
            break
        fi

        sleep 0.2

    done

    curl -fs http://localhost:8080 >/dev/null 2>&1 \
        || die "HTTP interface failed to respond."

    if aplay -l >/dev/null 2>&1; then
        success "Audio device detected."
    else
        warn "No audio device found."
    fi
}