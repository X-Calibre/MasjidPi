#!/usr/bin/env bash

UPDATE_STAGING="/opt/.masjidpi-staging"
UPDATE_BACKUP="/opt/.masjidpi-backup"

prepare_update() {
    info "Preparing new MasjidPi runtime..."

    rm -rf "$UPDATE_STAGING"
    mkdir -p "$UPDATE_STAGING"

    RUNTIME_TARGET="$UPDATE_STAGING"
    install_runtime
    unset RUNTIME_TARGET

    [[ -x "$UPDATE_STAGING/bin/masjidpi" ]] || die "Staged MasjidPi binary is missing."
    [[ -f "$UPDATE_STAGING/VERSION" ]] || die "Staged MasjidPi version file is missing."
    [[ -f "$UPDATE_STAGING/frontend/index.html" ]] || die "Staged MasjidPi frontend is missing."

    success "New runtime prepared."
}

activate_update() {
    local expected_version="$1"

    rm -rf "$UPDATE_BACKUP"

    stop_service

    info "Activating MasjidPi ${expected_version}..."
    mv "$INSTALL_DIR" "$UPDATE_BACKUP"
    mv "$UPDATE_STAGING" "$INSTALL_DIR"

    if start_service && run_selftest "$expected_version"; then
        rm -rf "$UPDATE_BACKUP"
        success "MasjidPi ${expected_version} installed and validated."
        return 0
    fi

    error "MasjidPi ${expected_version} failed validation. Rolling back..."

    stop_service || true
    rm -rf "$INSTALL_DIR"
    mv "$UPDATE_BACKUP" "$INSTALL_DIR"

    if start_service && run_selftest; then
        success "Previous MasjidPi version restored successfully."
    else
        error "Automatic rollback failed. MasjidPi may require manual recovery."
        return 1
    fi

    return 1
}

cleanup_update() {
    rm -rf "$UPDATE_STAGING"
}
