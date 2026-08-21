#!/usr/bin/env bash

install_service() {

    if [[ ! -f "$PROJECT_ROOT/scripts/masjidpi.service" ]]; then
        die "MasjidPi systemd service template is missing."
    fi

    cp "$PROJECT_ROOT/scripts/masjidpi.service" \
       /etc/systemd/system/

    systemctl daemon-reload

    systemctl enable masjidpi

    success "Systemd service installed and enabled."
}

install_component_services() {
    if $INSTALL_BOARD; then
        if [[ ! -f "$PROJECT_ROOT/scripts/masjidboard-display.sh" || ! -f "$PROJECT_ROOT/scripts/masjidpi-display.service" ]]; then
            die "MasjidBoard display runtime files are missing."
        fi

        install -m 0755 "$PROJECT_ROOT/scripts/masjidboard-display.sh" \
            /opt/masjidpi/bin/masjidboard-display
        install -m 0644 "$PROJECT_ROOT/scripts/masjidpi-display.service" \
            /etc/systemd/system/masjidpi-display.service

        systemctl daemon-reload
        systemctl enable masjidpi-display.service
        success "MasjidBoard display service installed and enabled."
    else
        if systemctl list-unit-files masjidpi-display.service >/dev/null 2>&1; then
            systemctl disable --now masjidpi-display.service || true
        fi
        rm -f /etc/systemd/system/masjidpi-display.service
        rm -f /opt/masjidpi/bin/masjidboard-display
        systemctl daemon-reload
    fi
}

stop_service() {

    if systemctl is-active --quiet masjidpi; then
        info "Stopping MasjidPi service..."
        systemctl stop masjidpi
        success "MasjidPi service stopped."
    fi
}

start_service() {

    info "Starting MasjidPi service..."

    systemctl start masjidpi

    if systemctl is-active --quiet masjidpi; then
        success "MasjidPi service is running."
    else
        error "MasjidPi failed to start."
        journalctl -u masjidpi --no-pager -n 20
        return 1
    fi

    if $INSTALL_BOARD; then
        info "Starting MasjidBoard display service..."
        systemctl restart masjidpi-display.service
    fi

    success "MasjidPi service started."

}
