#!/usr/bin/env bash

migrate_legacy_preferences() {
    local legacy_path="$INSTALL_DIR/backend/data/preferences.json"
    local persistent_path="/var/lib/masjidframe/preferences.json"

    if [[ ! -f "$legacy_path" || -f "$persistent_path" ]]; then
        return
    fi

    mkdir -p /var/lib/masjidframe
    cp "$legacy_path" "$persistent_path"
    chmod 0600 "$persistent_path"
    info "Migrated saved Web UI preferences to persistent storage."
}

migrate_catalogue_refresh_interval() {
    local config_path="${1:-/etc/masjidframe/config.yaml}"

    if [[ ! -f "$config_path" ]]; then
        return
    fi

    # Change only the previous packaged default. Other explicitly configured
    # refresh intervals remain untouched.
    if grep -qx '  refresh_interval: "168h"' "$config_path"; then
        sed -i 's/^  refresh_interval: "168h"$/  refresh_interval: "672h"/' \
            "$config_path"
        info "Updated automatic Listen catalogue refresh interval to 28 days."
    fi
}

migrate_runtime_socket_path() {
    local config_path="${1:-/etc/masjidframe/config.yaml}"

    if [[ ! -f "$config_path" ]]; then
        return
    fi

    # The mpv IPC socket is disposable runtime state. Move only MasjidFrame's
    # previous packaged default to tmpfs-backed /run; preserve custom paths.
    if grep -qx '  socket: "/tmp/masjidframe.sock"' "$config_path"; then
        sed -i 's|^  socket: "/tmp/masjidframe.sock"$|  socket: "/run/masjidframe/mpv.sock"|' \
            "$config_path"
        info "Moved the mpv runtime socket from /tmp to /run."
    fi
}

install_runtime() {

    local target_dir="${RUNTIME_TARGET:-$INSTALL_DIR}"

    info "Installing runtime..."

    mkdir -p "$target_dir/bin"

    if [[ -n "${RUNTIME_TARGET:-}" ]]; then
        # During an update, configuration and persistent runtime data remain
        # outside the replaceable application runtime.
        if $SOURCE_MODE; then
            cp "$PROJECT_ROOT/backend/build/masjidframe" \
                "$target_dir/bin/"
            rm -rf "$target_dir/frontend"
            cp -R "$PROJECT_ROOT/frontend" "$target_dir/"
            rm -f "$target_dir/VERSION"
        else
            cp "$RELEASE_DIR/masjidframe" "$target_dir/bin/"
            rm -rf "$target_dir/frontend"
            cp -R "$RELEASE_DIR/frontend" "$target_dir/"
            cp "$RELEASE_DIR/VERSION" "$target_dir/VERSION"
        fi

        chmod +x "$target_dir/bin/masjidframe"
        success "Runtime staged."
        return
    fi

    mkdir -p "$INSTALL_DIR/bin"
    mkdir -p /etc/masjidframe
    mkdir -p /var/lib/masjidframe

    if $SOURCE_MODE; then
        cp "$PROJECT_ROOT/backend/build/masjidframe" \
            "$INSTALL_DIR/bin/"

        if [[ ! -f /etc/masjidframe/config.yaml ]]; then
            cp "$PROJECT_ROOT/backend/configs/default.yaml" \
                "/etc/masjidframe/config.yaml"
        else
            info "Keeping existing configuration."
        fi

        if [[ ! -f /var/lib/masjidframe/catalogue.json ]]; then
            cp "$PROJECT_ROOT/backend/data/catalogue.json" \
                "/var/lib/masjidframe/catalogue.json"
        else
            info "Keeping existing catalogue."
        fi

        rm -rf "$INSTALL_DIR/frontend"
        cp -R "$PROJECT_ROOT/frontend" "$INSTALL_DIR/"
        rm -f "$INSTALL_DIR/VERSION"
    else
        cp "$RELEASE_DIR/masjidframe" "$INSTALL_DIR/bin/"

        if [[ ! -f /etc/masjidframe/config.yaml ]]; then
            cp "$RELEASE_DIR/default.yaml" "/etc/masjidframe/config.yaml"
        else
            info "Keeping existing configuration."
        fi

        if [[ ! -f /var/lib/masjidframe/catalogue.json ]]; then
            cp "$RELEASE_DIR/catalogue.json" "/var/lib/masjidframe/catalogue.json"
        else
            info "Keeping existing catalogue."
        fi

        rm -rf "$INSTALL_DIR/frontend"
        cp -R "$RELEASE_DIR/frontend" "$INSTALL_DIR/"

        cp "$RELEASE_DIR/VERSION" "$INSTALL_DIR/VERSION"
    fi

    chmod +x "$INSTALL_DIR/bin/masjidframe"

    success "Runtime installed."
}
