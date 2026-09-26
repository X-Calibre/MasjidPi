#!/usr/bin/env bash

BOOT_FIRMWARE_DIR="${MASJIDFRAME_BOOT_FIRMWARE_DIR:-/boot/firmware}"
CMDLINE_FILE="${MASJIDFRAME_CMDLINE_FILE:-$BOOT_FIRMWARE_DIR/cmdline.txt}"
CONFIG_FILE="${MASJIDFRAME_CONFIG_FILE:-$BOOT_FIRMWARE_DIR/config.txt}"
RPI_MODEL_FILE="${MASJIDFRAME_RPI_MODEL_FILE:-/proc/device-tree/model}"
PLYMOUTH_THEME_DIR="${MASJIDFRAME_PLYMOUTH_THEME_DIR:-/usr/share/plymouth/themes/masjidframe}"
PLYMOUTH_QUIT_DROPIN_DIR="${MASJIDFRAME_PLYMOUTH_QUIT_DROPIN_DIR:-/etc/systemd/system/plymouth-quit.service.d}"
BOOT_READONLY_SERVICE_FILE="${MASJIDFRAME_BOOT_READONLY_SERVICE_FILE:-/etc/systemd/system/masjidframe-boot-readonly.service}"
BOOT_APT_HOOK_FILE="${MASJIDFRAME_BOOT_APT_HOOK_FILE:-/etc/apt/apt.conf.d/99-masjidframe-boot-firmware}"
BOOT_FIRMWARE_UPDATE_ACTIVE=false

is_raspberry_pi() {
    if [[ "${MASJIDFRAME_FORCE_RASPBERRY_PI:-0}" == "1" ]]; then
        return 0
    fi

    [[ -r "$RPI_MODEL_FILE" ]] || return 1
    grep -aqi 'Raspberry Pi' "$RPI_MODEL_FILE"
}

is_raspberry_pi_board() {
    $INSTALL_BOARD || return 1
    is_raspberry_pi
}

prepare_boot_firmware_update() {
    is_raspberry_pi_board || return 0
    mountpoint -q "$BOOT_FIRMWARE_DIR" || return 0

    info "Temporarily enabling Raspberry Pi boot firmware writes..."
    mount -o remount,rw "$BOOT_FIRMWARE_DIR"
    BOOT_FIRMWARE_UPDATE_ACTIVE=true
}

finish_boot_firmware_update() {
    $BOOT_FIRMWARE_UPDATE_ACTIVE || return 0

    protect_boot_firmware
}

protect_boot_firmware() {
    is_raspberry_pi || return 0
    mountpoint -q "$BOOT_FIRMWARE_DIR" || return 0

    sync -f "$BOOT_FIRMWARE_DIR"
    mount -o remount,ro "$BOOT_FIRMWARE_DIR"
    BOOT_FIRMWARE_UPDATE_ACTIVE=false
    success "Raspberry Pi boot firmware protected read-only."
}

configure_boot_firmware_protection() {
    is_raspberry_pi || return 0

    local service_file="$PROJECT_ROOT/scripts/masjidframe-boot-readonly.service"
    local apt_hook_file="$PROJECT_ROOT/scripts/99-masjidframe-boot-firmware"

    if [[ ! -f "$service_file" || ! -f "$apt_hook_file" ]]; then
        warn "Raspberry Pi boot firmware protection assets are missing; persistent protection was not installed."
        finish_boot_firmware_update
        return 0
    fi

    install -D -m 0644 "$service_file" "$BOOT_READONLY_SERVICE_FILE"
    install -D -m 0644 "$apt_hook_file" "$BOOT_APT_HOOK_FILE"
    systemctl daemon-reload
    systemctl enable masjidframe-boot-readonly.service >/dev/null

    protect_boot_firmware
}

append_cmdline_parameter() {
    local parameter="$1"

    grep -qw -- "$parameter" "$CMDLINE_FILE" && return 0

    sed -i "s/$/ $parameter/" "$CMDLINE_FILE"
}

configure_quiet_boot() {
    is_raspberry_pi_board || return 0

    if [[ ! -f "$CMDLINE_FILE" || ! -f "$CONFIG_FILE" ]]; then
        warn "Raspberry Pi boot configuration not found; skipping quiet boot configuration."
        return 0
    fi

    info "Configuring quiet Raspberry Pi Board boot..."

    append_cmdline_parameter quiet
    append_cmdline_parameter loglevel=3
    append_cmdline_parameter systemd.show_status=false

    if ! grep -Eq '^[[:space:]]*disable_splash=1[[:space:]]*$' "$CONFIG_FILE"; then
        if [[ "$(awk '/^\[[^]]+\][[:space:]]*$/ { section=$0 } END { print section }' "$CONFIG_FILE")" == "[all]" ]]; then
            printf 'disable_splash=1\n' >> "$CONFIG_FILE"
        else
            printf '\n[all]\ndisable_splash=1\n' >> "$CONFIG_FILE"
        fi
    fi

    # Keep console=tty1 in the kernel command line. Removing it caused unreliable
    # KMS startup during physical Raspberry Pi validation. Disable only the login
    # getty so the kernel keeps its console without exposing a prompt at boot.
    systemctl disable getty@tty1.service >/dev/null 2>&1 || true

    success "Quiet Raspberry Pi Board boot configured."
}

configure_boot_splash() {
    is_raspberry_pi_board || return 0

    local theme_file="$PROJECT_ROOT/scripts/masjidframe-splash.plymouth"
    local script_file="$PROJECT_ROOT/scripts/masjidframe-splash-standard.script"
    local logo_file="$PROJECT_ROOT/frontend/masjidframe-splash-logo.png"
    local portrait_logo_file="$PROJECT_ROOT/frontend/masjidframe-splash-logo-portrait.png"

    if [[ ! -f "$theme_file" || ! -f "$script_file" || ! -f "$logo_file" || ! -f "$portrait_logo_file" ]]; then
        warn "MasjidFrame Plymouth splash assets are missing; skipping branded boot splash."
        return 0
    fi

    if ! command -v plymouth-set-default-theme >/dev/null 2>&1; then
        warn "Plymouth is unavailable; skipping branded boot splash."
        return 0
    fi

    info "Installing MasjidFrame Raspberry Pi Board boot splash..."

    install -d -m 0755 "$PLYMOUTH_THEME_DIR"
    install -m 0644 "$theme_file" "$PLYMOUTH_THEME_DIR/masjidframe.plymouth"
    install -m 0644 "$script_file" "$PLYMOUTH_THEME_DIR/masjidframe-splash.script"
    install -m 0644 "$logo_file" "$PLYMOUTH_THEME_DIR/$(basename "$logo_file")"
    install -m 0644 "$portrait_logo_file" "$PLYMOUTH_THEME_DIR/$(basename "$portrait_logo_file")"
    rm -f "$PLYMOUTH_THEME_DIR/masjidframe-splash-logo-appliance.png"

    append_cmdline_parameter splash
    append_cmdline_parameter plymouth.ignore-serial-consoles

    # Debian normally clears the Plymouth framebuffer when plymouth-quit runs.
    # Keep the final splash frame on screen after Plymouth releases DRM so the
    # quiet console cannot flash between Plymouth and Cog startup.
    install -d -m 0755 "$PLYMOUTH_QUIT_DROPIN_DIR"
    cat > "$PLYMOUTH_QUIT_DROPIN_DIR/masjidframe.conf" <<'EOF'
[Service]
ExecStart=
ExecStart=-/usr/bin/plymouth quit --retain-splash
EOF

    systemctl daemon-reload
    plymouth-set-default-theme -R masjidframe

    success "MasjidFrame Raspberry Pi Board boot splash installed."
}
