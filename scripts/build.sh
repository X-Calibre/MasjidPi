#!/usr/bin/env bash

build_project() {

    info "Building MasjidPi..."

    local source_version
    source_version="$(jq -er '
        .version |
        select(type == "string" and test("^v[0-9]+\\.[0-9]+\\.[0-9]+$"))
    ' "$PROJECT_ROOT/version.json")" \
        || die "Unable to determine a valid source version."

    local development_version="${source_version}-dev"

    cd "$PROJECT_ROOT/backend" || die "Unable to enter backend source directory."

    mkdir -p build

    go build \
        -ldflags "-X github.com/X-Calibre/MasjidPi/backend/internal/version.Version=${development_version}" \
        -o build/masjidpi \
        ./cmd/masjidpi

    # Keep build artifacts writable by the user who invoked sudo, if any.
    if [[ -n "${SUDO_UID:-}" && -n "${SUDO_GID:-}" ]]; then
        chown -R "$SUDO_UID:$SUDO_GID" build
    fi

    success "Build complete."
}
