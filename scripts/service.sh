# shellcheck shell=bash

install_service() {

    if [[ ! -f "$PROJECT_ROOT/scripts/masjidframe.service" ]]; then
        die "MasjidFrame systemd service template is missing."
    fi

    cp "$PROJECT_ROOT/scripts/masjidframe.service" \
       /etc/systemd/system/

    systemctl daemon-reload

    systemctl enable masjidframe

    success "Systemd service installed and enabled."
}

install_component_services() {
    if $INSTALL_BOARD; then
        if [[ ! -f "$PROJECT_ROOT/scripts/masjidboard-display.sh" || ! -f "$PROJECT_ROOT/scripts/masjidframe-display.service" ]]; then
            die "MasjidBoard display runtime files are missing."
        fi

        install -m 0755 "$PROJECT_ROOT/scripts/masjidboard-display.sh" \
            /opt/masjidframe/bin/masjidboard-display
        install -m 0644 "$PROJECT_ROOT/scripts/masjidframe-display.service" \
            /etc/systemd/system/masjidframe-display.service

        if is_raspberry_pi_board; then
            if [[ ! -f "$PROJECT_ROOT/scripts/masjidboard-warmup.sh" || ! -f "$PROJECT_ROOT/scripts/masjidframe-display-warmup.service" ]]; then
                die "MasjidBoard Raspberry Pi warm-up runtime files are missing."
            fi

            install -m 0755 "$PROJECT_ROOT/scripts/masjidboard-warmup.sh" \
                /opt/masjidframe/bin/masjidboard-warmup
            install -m 0644 "$PROJECT_ROOT/scripts/masjidframe-display-warmup.service" \
                /etc/systemd/system/masjidframe-display-warmup.service
        else
            systemctl disable --now masjidframe-display-warmup.service >/dev/null 2>&1 || true
            rm -f /etc/systemd/system/masjidframe-display-warmup.service
            rm -f /opt/masjidframe/bin/masjidboard-warmup
        fi

        # Remove the calibration rule used by the retired rotated Waveshare profile.
        rm -f /etc/udev/rules.d/99-masjidframe-appliance-touchscreen.rules
        udevadm control --reload-rules 2>/dev/null || true

        systemctl daemon-reload
        systemctl enable masjidframe-display.service
        if is_raspberry_pi_board; then
            systemctl enable masjidframe-display-warmup.service
            success "MasjidBoard Raspberry Pi WebKit warm-up service installed and enabled."
        fi
        success "MasjidBoard display service installed and enabled."
    else
        if systemctl list-unit-files masjidframe-display.service >/dev/null 2>&1; then
            systemctl disable --now masjidframe-display.service || true
        fi
        systemctl disable --now masjidframe-display-warmup.service >/dev/null 2>&1 || true
        rm -f /etc/systemd/system/masjidframe-display.service
        rm -f /etc/systemd/system/masjidframe-display-warmup.service
        rm -f /opt/masjidframe/bin/masjidboard-display
        rm -f /opt/masjidframe/bin/masjidboard-warmup
        rm -f /etc/udev/rules.d/99-masjidframe-appliance-touchscreen.rules
        udevadm control --reload-rules 2>/dev/null || true
        systemctl daemon-reload
        systemctl reset-failed masjidframe-display.service 2>/dev/null || true
        systemctl reset-failed masjidframe-display-warmup.service 2>/dev/null || true
    fi
}

stop_service() {

    if systemctl is-active --quiet masjidframe; then
        info "Stopping MasjidFrame service..."
        systemctl stop masjidframe
        success "MasjidFrame service stopped."
    fi
}

start_service() {

    info "Starting MasjidFrame service..."

    systemctl start masjidframe

    if systemctl is-active --quiet masjidframe; then
        success "MasjidFrame service is running."
    else
        error "MasjidFrame failed to start."
        journalctl -u masjidframe --no-pager -n 20
        return 1
    fi

    if $INSTALL_BOARD; then
        info "Starting MasjidBoard display service..."
        if ! systemctl restart masjidframe-display.service; then
            error "MasjidBoard display service failed to start."
            journalctl -u masjidframe-display.service --no-pager -n 20 || true
            return 1
        fi
    fi

    success "MasjidFrame service started."
}
