#!/usr/bin/env bash

build_project() {

    info "Building MasjidPi..."

    cd "$PROJECT_ROOT/backend" || die "Unable to enter backend source directory."

    mkdir -p build

    go build -o build/masjidpi ./cmd/masjidpi

    # Keep build artifacts writable by the user who invoked sudo, if any.
    if [[ -n "${SUDO_UID:-}" && -n "${SUDO_GID:-}" ]]; then
        chown -R "$SUDO_UID:$SUDO_GID" build
    fi

    success "Build complete."
}
