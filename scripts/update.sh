#!/usr/bin/env bash

UPDATE_STAGING="/opt/.masjidpi-staging"
UPDATE_BACKUP="/opt/.masjidpi-backup"
UPDATE_MARKER="/opt/.masjidpi-update-in-progress"

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

    local staged_version
    staged_version="$(cat "$UPDATE_STAGING/VERSION")"
    [[ "$staged_version" == "$RELEASE_VERSION" ]] || \
        die "Staged version mismatch: expected $RELEASE_VERSION, found $staged_version."

    success "New runtime prepared."
}

rollback_update() {
    error "MasjidPi ${RELEASE_VERSION} failed validation. Rolling back..."

    stop_service || true

    if [[ -d "$INSTALL_DIR" ]]; then
        rm -rf "$INSTALL_DIR"
    fi

    if [[ -d "$UPDATE_BACKUP" ]]; then
        mv "$UPDATE_BACKUP" "$INSTALL_DIR"
    else
        error "Previous MasjidPi runtime is not available for rollback."
        return 1
    fi

    if start_service && run_selftest; then
        rm -f "$UPDATE_MARKER"
        success "Previous MasjidPi version restored successfully."
        return 0
    fi

    error "Automatic rollback failed. MasjidPi may require manual recovery."
    return 1
}

activate_update() {
    local expected_version="$1"

    [[ -d "$UPDATE_STAGING" ]] || die "Update staging directory is missing."
    [[ -d "$INSTALL_DIR" ]] || die "Current MasjidPi runtime is missing."

    rm -rf "$UPDATE_BACKUP"
    printf '%s\n' "$expected_version" > "$UPDATE_MARKER"

    stop_service

    info "Activating MasjidPi ${expected_version}..."

    if ! mv "$INSTALL_DIR" "$UPDATE_BACKUP"; then
        rm -f "$UPDATE_MARKER"
        die "Unable to preserve the current MasjidPi runtime. Update cancelled."
    fi

    if ! mv "$UPDATE_STAGING" "$INSTALL_DIR"; then
        error "Unable to activate the new MasjidPi runtime. Restoring previous version..."
        mv "$UPDATE_BACKUP" "$INSTALL_DIR" || die "Previous MasjidPi runtime could not be restored."
        rm -f "$UPDATE_MARKER"
        start_service || true
        die "Update cancelled before the new runtime was activated."
    fi

    if start_service && run_selftest "$expected_version"; then
        rm -rf "$UPDATE_BACKUP"
        rm -f "$UPDATE_MARKER"
        success "MasjidPi ${expected_version} installed and validated."
        return 0
    fi

    if rollback_update; then
        return 1
    fi

    return 1
}

cleanup_update() {
    rm -rf "$UPDATE_STAGING"
}
