#!/usr/bin/env bash

build_project() {

    info "Building MasjidPi..."

    cd "$PROJECT_ROOT/backend"

    mkdir -p build

    go build -o build/masjidpi ./cmd/masjidpi

    success "Build complete."
}