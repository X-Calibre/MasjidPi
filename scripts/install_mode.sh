#!/usr/bin/env bash

detect_install_mode() {

    local update_backup="/opt/.masjidpi-backup"
    local update_marker="/opt/.masjidpi-update-in-progress"

    # Recover an interrupted update before deciding whether this is a fresh
    # installation. A marker means the previous runtime was not committed.
    if [[ -f "$update_marker" ]]; then
        if [[ -d "$update_backup" ]]; then
            warn "Incomplete previous update detected. Restoring previous runtime..."

            if systemctl is-active --quiet masjidpi; then
                systemctl stop masjidpi
            fi

            rm -rf "$INSTALL_DIR"
            mv "$update_backup" "$INSTALL_DIR"
            rm -f "$update_marker"

            success "Previous MasjidPi runtime restored."
        else
            warn "Found an incomplete update marker without a backup runtime. Continuing with the current installation."
            rm -f "$update_marker"
        fi
    elif [[ ! -d "$INSTALL_DIR" && -d "$update_backup" ]]; then
        # Backward-compatible recovery for an interrupted update created before
        # the transaction marker was present.
        warn "Incomplete previous update detected. Restoring previous runtime..."
        mv "$update_backup" "$INSTALL_DIR"
        success "Previous MasjidPi runtime restored."
    fi

    if [ -d "$INSTALL_DIR" ]; then
        INSTALL_MODE="update"
    else
        INSTALL_MODE="install"
    fi

}
