#!/usr/bin/env bash

BOOT_FIRMWARE_DIR="${MASJIDPI_BOOT_FIRMWARE_DIR:-/boot/firmware}"
CMDLINE_FILE="${MASJIDPI_CMDLINE_FILE:-$BOOT_FIRMWARE_DIR/cmdline.txt}"
CONFIG_FILE="${MASJIDPI_CONFIG_FILE:-$BOOT_FIRMWARE_DIR/config.txt}"
RPI_MODEL_FILE="${MASJIDPI_RPI_MODEL_FILE:-/proc/device-tree/model}"
PLYMOUTH_THEME_DIR="${MASJIDPI_PLYMOUTH_THEME_DIR:-/usr/share/plymouth/themes/masjidpi}"
PLYMOUTH_QUIT_DROPIN_DIR="${MASJIDPI_PLYMOUTH_QUIT_DROPIN_DIR:-/etc/systemd/system/plymouth-quit.service.d}"
BOOT_READONLY_SERVICE_FILE="${MASJIDPI_BOOT_READONLY_SERVICE_FILE:-/etc/systemd/system/masjidpi-boot-readonly.service}"
BOOT_APT_HOOK_FILE="${MASJIDPI_BOOT_APT_HOOK_FILE:-/etc/apt/apt.conf.d/99-masjidpi-boot-firmware}"
APPLIANCE_USB_SYSFS_ROOT="${MASJIDPI_USB_SYSFS_ROOT:-/sys/bus/usb/devices}"
APPLIANCE_DRM_SYSFS_ROOT="${MASJIDPI_DRM_SYSFS_ROOT:-/sys/class/drm}"
BOOT_FIRMWARE_UPDATE_ACTIVE=false

is_raspberry_pi() {
    if [[ "${MASJIDPI_FORCE_RASPBERRY_PI:-0}" == "1" ]]; then
        return 0
    fi

    [[ -r "$RPI_MODEL_FILE" ]] || return 1
    grep -aqi 'Raspberry Pi' "$RPI_MODEL_FILE"
}

is_raspberry_pi_board() {
    $INSTALL_BOARD || return 1
    is_raspberry_pi
}

waveshare_appliance_touch_present() {
    local device

    for device in "$APPLIANCE_USB_SYSFS_ROOT"/*; do
        [[ -r "$device/idVendor" && -r "$device/idProduct" ]] || continue
        [[ "$(<"$device/idVendor")" == "0eef" ]] || continue
        [[ "$(<"$device/idProduct")" == "0005" ]] || continue

        if [[ -r "$device/manufacturer" && -r "$device/product" ]]; then
            [[ "$(<"$device/manufacturer")" == "WaveShare" ]] || continue
            [[ "$(<"$device/product")" == "WS170120" ]] || continue
        fi

        return 0
    done

    return 1
}

appliance_hdmi_mode_present() {
    local connector

    for connector in "$APPLIANCE_DRM_SYSFS_ROOT"/card*-HDMI-A-*; do
        [[ -r "$connector/status" && -r "$connector/modes" ]] || continue
        [[ "$(<"$connector/status")" == "connected" ]] || continue
        grep -Fxq '1024x600' "$connector/modes" && return 0
    done

    return 1
}

is_appliance_display_hardware() {
    is_raspberry_pi_board || return 1
    waveshare_appliance_touch_present && appliance_hdmi_mode_present
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

    local service_file="$PROJECT_ROOT/scripts/masjidpi-boot-readonly.service"
    local apt_hook_file="$PROJECT_ROOT/scripts/99-masjidpi-boot-firmware"

    if [[ ! -f "$service_file" || ! -f "$apt_hook_file" ]]; then
        warn "Raspberry Pi boot firmware protection assets are missing; persistent protection was not installed."
        finish_boot_firmware_update
        return 0
    fi

    install -D -m 0644 "$service_file" "$BOOT_READONLY_SERVICE_FILE"
    install -D -m 0644 "$apt_hook_file" "$BOOT_APT_HOOK_FILE"
    systemctl daemon-reload
    systemctl enable masjidpi-boot-readonly.service >/dev/null

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

    local theme_file="$PROJECT_ROOT/scripts/masjidpi-splash.plymouth"
    local script_file="$PROJECT_ROOT/scripts/masjidpi-splash-standard.script"
    local logo_file="$PROJECT_ROOT/frontend/masjidpi-splash-logo.png"
    local appliance_logo_file="$PROJECT_ROOT/frontend/masjidpi-splash-logo-appliance.png"

    if is_appliance_display_hardware; then
        script_file="$PROJECT_ROOT/scripts/masjidpi-splash.script"
        logo_file="$appliance_logo_file"
    fi

    if [[ ! -f "$theme_file" || ! -f "$script_file" || ! -f "$logo_file" ]]; then
        warn "MasjidPi Plymouth splash assets are missing; skipping branded boot splash."
        return 0
    fi

    if ! command -v plymouth-set-default-theme >/dev/null 2>&1; then
        warn "Plymouth is unavailable; skipping branded boot splash."
        return 0
    fi

    info "Installing MasjidPi Raspberry Pi Board boot splash..."

    install -d -m 0755 "$PLYMOUTH_THEME_DIR"
    install -m 0644 "$theme_file" "$PLYMOUTH_THEME_DIR/masjidpi.plymouth"
    install -m 0644 "$script_file" "$PLYMOUTH_THEME_DIR/masjidpi-splash.script"
    install -m 0644 "$logo_file" "$PLYMOUTH_THEME_DIR/$(basename "$logo_file")"

    append_cmdline_parameter splash
    append_cmdline_parameter plymouth.ignore-serial-consoles

    # Debian normally clears the Plymouth framebuffer when plymouth-quit runs.
    # Keep the final splash frame on screen after Plymouth releases DRM so the
    # quiet console cannot flash between Plymouth and Cog startup.
    install -d -m 0755 "$PLYMOUTH_QUIT_DROPIN_DIR"
    cat > "$PLYMOUTH_QUIT_DROPIN_DIR/masjidpi.conf" <<'EOF'
[Service]
ExecStart=
ExecStart=-/usr/bin/plymouth quit --retain-splash
EOF

    systemctl daemon-reload
    plymouth-set-default-theme -R masjidpi

    success "MasjidPi Raspberry Pi Board boot splash installed."
}
