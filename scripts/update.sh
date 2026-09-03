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
    [[ -f "$UPDATE_STAGING/frontend/index.html" ]] || die "Staged MasjidPi frontend is missing."

    if $SOURCE_MODE; then
        info "Source build staged successfully."
    else
        [[ -f "$UPDATE_STAGING/VERSION" ]] || die "Staged MasjidPi version file is missing."

        local staged_version
        staged_version="$(cat "$UPDATE_STAGING/VERSION")"
        [[ "$staged_version" == "$RELEASE_VERSION" ]] || \
            die "Staged version mismatch: expected $RELEASE_VERSION, found $staged_version."
    fi

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

    # The transaction is resolved once the backup has replaced the failed
    # runtime. Clear the marker before broader service validation so a
    # non-runtime self-test failure cannot leave an unrecoverable marker with
    # no backup directory.
    rm -f "$UPDATE_MARKER"

    if ! restore_previous_components; then
        error "Unable to restore the previous MasjidPi component profile."
        return 1
    fi

    if ! install_service; then
        error "Unable to restore the MasjidPi systemd service during rollback."
        return 1
    fi

    if ! install_component_services; then
        error "Unable to restore component services during rollback."
        return 1
    fi

    if start_service && run_selftest; then
        success "Previous MasjidPi version and component profile restored successfully."
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

        if ! mv "$UPDATE_BACKUP" "$INSTALL_DIR"; then
            die "Previous MasjidPi runtime could not be restored."
        fi

        restore_previous_components || true
        install_component_services || true
        rm -f "$UPDATE_MARKER"

        if start_service && run_selftest; then
            die "Update cancelled before the new runtime was activated. Previous version restored successfully."
        fi

        die "Update cancelled and previous MasjidPi runtime failed validation."
    fi

    # Updates may change service-level runtime requirements such as
    # RuntimeDirectory. Install the main unit before starting the new runtime;
    # otherwise a migrated configuration can depend on a directory that the
    # previous unit does not create.
    if install_service && install_component_services && start_service && run_selftest "$expected_version"; then
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
