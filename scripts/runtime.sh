#!/usr/bin/env bash

install_runtime() {

    info "Installing runtime..."

    mkdir -p "$INSTALL_DIR"
    mkdir -p "$INSTALL_DIR/bin"
    mkdir -p "$INSTALL_DIR/configs"
    mkdir -p "$INSTALL_DIR/data"

    cp "$PROJECT_ROOT/backend/build/masjidpi" \
        "$INSTALL_DIR/bin/"

    cp "$PROJECT_ROOT/backend/configs/default.yaml" \
        "$INSTALL_DIR/configs/"

    cp "$PROJECT_ROOT/backend/data/catalogue.json" \
        "$INSTALL_DIR/data/"

    rm -rf "$INSTALL_DIR/frontend"

    cp -R "$PROJECT_ROOT/frontend" \
        "$INSTALL_DIR/"

    success "Runtime installed."

}