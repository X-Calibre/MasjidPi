#!/usr/bin/env bash

install_runtime() {

    info "Installing runtime..."

    mkdir -p "$INSTALL_DIR/bin"
    mkdir -p /etc/masjidpi
    mkdir -p /var/lib/masjidpi

    cp "$PROJECT_ROOT/backend/build/masjidpi" \
        "$INSTALL_DIR/bin/"

    cp "$PROJECT_ROOT/backend/configs/default.yaml" \
        "/etc/masjidpi/config.yaml"

    cp "$PROJECT_ROOT/backend/data/catalogue.json" \
        "/var/lib/masjidpi/catalogue.json"

    rm -rf "$INSTALL_DIR/frontend"

    cp -R "$PROJECT_ROOT/frontend" \
        "$INSTALL_DIR/"

    success "Runtime installed."
}