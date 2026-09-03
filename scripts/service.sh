# shellcheck shell=bash

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

        if is_raspberry_pi_board; then
            if [[ ! -f "$PROJECT_ROOT/scripts/masjidboard-warmup.sh" || ! -f "$PROJECT_ROOT/scripts/masjidpi-display-warmup.service" ]]; then
                die "MasjidBoard Raspberry Pi warm-up runtime files are missing."
            fi

            install -m 0755 "$PROJECT_ROOT/scripts/masjidboard-warmup.sh" \
                /opt/masjidpi/bin/masjidboard-warmup
            install -m 0644 "$PROJECT_ROOT/scripts/masjidpi-display-warmup.service" \
                /etc/systemd/system/masjidpi-display-warmup.service
        else
            systemctl disable --now masjidpi-display-warmup.service >/dev/null 2>&1 || true
            rm -f /etc/systemd/system/masjidpi-display-warmup.service
            rm -f /opt/masjidpi/bin/masjidboard-warmup
        fi

        if [[ -f "$PROJECT_ROOT/scripts/99-masjidpi-appliance-touchscreen.rules" ]]; then
            install -m 0644 "$PROJECT_ROOT/scripts/99-masjidpi-appliance-touchscreen.rules" \
                /etc/udev/rules.d/99-masjidpi-appliance-touchscreen.rules
            udevadm control --reload-rules
            udevadm trigger --subsystem-match=input || true
        fi

        systemctl daemon-reload
        systemctl enable masjidpi-display.service
        if is_raspberry_pi_board; then
            systemctl enable masjidpi-display-warmup.service
            success "MasjidBoard Raspberry Pi WebKit warm-up service installed and enabled."
        fi
        success "MasjidBoard display service installed and enabled."
    else
        if systemctl list-unit-files masjidpi-display.service >/dev/null 2>&1; then
            systemctl disable --now masjidpi-display.service || true
        fi
        systemctl disable --now masjidpi-display-warmup.service >/dev/null 2>&1 || true
        rm -f /etc/systemd/system/masjidpi-display.service
        rm -f /etc/systemd/system/masjidpi-display-warmup.service
        rm -f /opt/masjidpi/bin/masjidboard-display
        rm -f /opt/masjidpi/bin/masjidboard-warmup
        rm -f /etc/udev/rules.d/99-masjidpi-appliance-touchscreen.rules
        udevadm control --reload-rules 2>/dev/null || true
        systemctl daemon-reload
        systemctl reset-failed masjidpi-display.service 2>/dev/null || true
        systemctl reset-failed masjidpi-display-warmup.service 2>/dev/null || true
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
        if ! systemctl restart masjidpi-display.service; then
            error "MasjidBoard display service failed to start."
            journalctl -u masjidpi-display.service --no-pager -n 20 || true
            return 1
        fi
    fi

    success "MasjidPi service started."
}
