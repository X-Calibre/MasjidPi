#!/usr/bin/env bash

install_runtime() {

    local target_dir="${RUNTIME_TARGET:-$INSTALL_DIR}"

    info "Installing runtime..."

    mkdir -p "$target_dir/bin"

    if [[ -n "${RUNTIME_TARGET:-}" ]]; then
        # During an update, configuration and persistent runtime data remain
        # outside the replaceable application runtime.
        if $SOURCE_MODE; then
            cp "$PROJECT_ROOT/backend/build/masjidpi" \
                "$target_dir/bin/"
            rm -rf "$target_dir/frontend"
            cp -R "$PROJECT_ROOT/frontend" "$target_dir/"
            rm -f "$target_dir/VERSION"
        else
            cp "$RELEASE_DIR/masjidpi" "$target_dir/bin/"
            rm -rf "$target_dir/frontend"
            cp -R "$RELEASE_DIR/frontend" "$target_dir/"
            cp "$RELEASE_DIR/VERSION" "$target_dir/VERSION"
        fi

        chmod +x "$target_dir/bin/masjidpi"
        success "Runtime staged."
        return
    fi

    mkdir -p "$INSTALL_DIR/bin"
    mkdir -p /etc/masjidpi
    mkdir -p /var/lib/masjidpi

    if $SOURCE_MODE; then
        cp "$PROJECT_ROOT/backend/build/masjidpi" \
            "$INSTALL_DIR/bin/"

        if [[ ! -f /etc/masjidpi/config.yaml ]]; then
            cp "$PROJECT_ROOT/backend/configs/default.yaml" \
                "/etc/masjidpi/config.yaml"
        else
            info "Keeping existing configuration."
        fi

        if [[ ! -f /var/lib/masjidpi/catalogue.json ]]; then
            cp "$PROJECT_ROOT/backend/data/catalogue.json" \
                "/var/lib/masjidpi/catalogue.json"
        else
            info "Keeping existing catalogue."
        fi

        rm -rf "$INSTALL_DIR/frontend"
        cp -R "$PROJECT_ROOT/frontend" "$INSTALL_DIR/"
        rm -f "$INSTALL_DIR/VERSION"
    else
        cp "$RELEASE_DIR/masjidpi" "$INSTALL_DIR/bin/"

        if [[ ! -f /etc/masjidpi/config.yaml ]]; then
            cp "$RELEASE_DIR/default.yaml" "/etc/masjidpi/config.yaml"
        else
            info "Keeping existing configuration."
        fi

        if [[ ! -f /var/lib/masjidpi/catalogue.json ]]; then
            cp "$RELEASE_DIR/catalogue.json" "/var/lib/masjidpi/catalogue.json"
        else
            info "Keeping existing catalogue."
        fi

        rm -rf "$INSTALL_DIR/frontend"
        cp -R "$RELEASE_DIR/frontend" "$INSTALL_DIR/"

        cp "$RELEASE_DIR/VERSION" "$INSTALL_DIR/VERSION"
    fi

    chmod +x "$INSTALL_DIR/bin/masjidpi"

    success "Runtime installed."
}
