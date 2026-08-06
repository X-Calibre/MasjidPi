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

start_service() {

    info "Starting MasjidPi service..."

    systemctl restart masjidpi

    if systemctl is-active --quiet masjidpi; then
        success "MasjidPi service is running."
    else
        error "MasjidPi failed to start."
        journalctl -u masjidpi --no-pager -n 20
        exit 1
    fi

    success "MasjidPi service started."

}