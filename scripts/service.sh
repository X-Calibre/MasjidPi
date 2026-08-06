#!/usr/bin/env bash

install_service() {

    if [[ ! -f "$PROJECT_ROOT/scripts/masjidpi.service" ]]; then
        warn "No systemd service template found."
        return
    fi

    cp "$PROJECT_ROOT/scripts/masjidpi.service" \
       /etc/systemd/system/

    systemctl daemon-reload

    systemctl enable masjidpi
    systemctl restart masjidpi

    success "Systemd service installed and started."
}