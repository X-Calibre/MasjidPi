#!/usr/bin/env bash

run_selftest() {

    local expected_version="${1:-}"

    info "Running self test..."

    if ! systemctl is-active --quiet masjidpi; then
        error "MasjidPi service is not running."
        return 1
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

    if ! curl -fs http://localhost:8080 >/dev/null 2>&1; then
        error "HTTP interface failed to respond."
        return 1
    fi

    info "Checking application version..."

    local version_response
    version_response="$(curl -fsS http://localhost:8080/api/version)" || {
        error "Version endpoint failed to respond."
        return 1
    }

    if [[ -n "$expected_version" ]] && \
       ! grep -Eq '"version"[[:space:]]*:[[:space:]]*"'"$expected_version"'"' <<<"$version_response"; then
        error "Running application version does not match expected version ${expected_version}."
        return 1
    fi

    success "Application version verified."

    info "Checking player status endpoint..."

    if ! curl -fs http://localhost:8080/api/player/status >/dev/null 2>&1; then
        error "Player status endpoint failed."
        return 1
    fi

    success "Player status endpoint is responding."

    if aplay -l >/dev/null 2>&1; then
        success "Audio device detected."
    else
        warn "No audio device found."
    fi

    success "Self test passed."
}
