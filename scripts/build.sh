#!/usr/bin/env bash

build_project() {

    info "Building MasjidPi..."

    cd "$PROJECT_ROOT/backend"

    go build -o masjidpi ./cmd/masjidpi

    success "Build complete."
}