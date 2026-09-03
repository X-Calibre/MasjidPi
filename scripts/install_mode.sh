#!/usr/bin/env bash

detect_install_mode() {

    local update_backup="/opt/.masjidpi-backup"
    local update_marker="/opt/.masjidpi-update-in-progress"

    # Recover an interrupted update before deciding whether this is a fresh
    # installation. The marker means the previous runtime swap did not finish.
    if [[ -f "$update_marker" ]]; then
        if [[ -d "$update_backup" ]]; then
            warn "Incomplete previous update detected. Restoring previous runtime..."

            stop_service || true

            if [[ -d "$INSTALL_DIR" ]]; then
                rm -rf "$INSTALL_DIR"
            fi

            if ! mv "$update_backup" "$INSTALL_DIR"; then
                error "Previous MasjidPi runtime could not be restored."
                return 1
            fi

            # The previous runtime is now unambiguously active and the backup
            # has been consumed. Remove the transaction marker before service
            # validation so a separate self-test failure cannot block retries.
            rm -f "$update_marker"

            if start_service && run_selftest; then
                success "Previous MasjidPi runtime restored and validated."
            else
                error "Previous MasjidPi runtime was restored but failed validation."
                return 1
            fi
        else
            warn "Found an incomplete update marker without a backup runtime."
            error "MasjidPi cannot safely determine which runtime should be active."
            return 1
        fi
    elif [[ ! -d "$INSTALL_DIR" && -d "$update_backup" ]]; then
        # Backward-compatible recovery for an interrupted update created before
        # the transaction marker was present.
        warn "Incomplete previous update detected. Restoring previous runtime..."

        if ! mv "$update_backup" "$INSTALL_DIR"; then
            error "Previous MasjidPi runtime could not be restored."
            return 1
        fi

        if start_service && run_selftest; then
            success "Previous MasjidPi runtime restored and validated."
        else
            error "Previous MasjidPi runtime was restored but failed validation."
            return 1
        fi
    fi

    if [ -d "$INSTALL_DIR" ]; then
        INSTALL_MODE="update"
    else
        INSTALL_MODE="install"
    fi

}
