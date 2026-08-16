#!/usr/bin/env bash

detect_install_mode() {

    local update_backup="/opt/.masjidpi-backup"

    # Recover an interrupted update before deciding whether this is a fresh
    # installation. Persistent configuration and data live outside the
    # replaceable application runtime, so restoring this directory is safe.
    if [[ ! -d "$INSTALL_DIR" && -d "$update_backup" ]]; then
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
