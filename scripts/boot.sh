#!/usr/bin/env bash

BOOT_FIRMWARE_DIR="${MASJIDPI_BOOT_FIRMWARE_DIR:-/boot/firmware}"
CMDLINE_FILE="${MASJIDPI_CMDLINE_FILE:-$BOOT_FIRMWARE_DIR/cmdline.txt}"
CONFIG_FILE="${MASJIDPI_CONFIG_FILE:-$BOOT_FIRMWARE_DIR/config.txt}"
RPI_MODEL_FILE="${MASJIDPI_RPI_MODEL_FILE:-/proc/device-tree/model}"

is_raspberry_pi() {
    if [[ "${MASJIDPI_FORCE_RASPBERRY_PI:-0}" == "1" ]]; then
        return 0
    fi

    [[ -r "$RPI_MODEL_FILE" ]] || return 1
    grep -aqi 'Raspberry Pi' "$RPI_MODEL_FILE"
}

append_cmdline_parameter() {
    local parameter="$1"

    grep -qw -- "$parameter" "$CMDLINE_FILE" && return 0

    sed -i "s/$/ $parameter/" "$CMDLINE_FILE"
}

configure_quiet_boot() {
    $INSTALL_BOARD || return 0
    is_raspberry_pi || return 0

    if [[ ! -f "$CMDLINE_FILE" || ! -f "$CONFIG_FILE" ]]; then
        warn "Raspberry Pi boot configuration not found; skipping quiet boot configuration."
        return 0
    fi

    info "Configuring quiet Raspberry Pi appliance boot..."

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
    # KMS startup on validated Raspberry Pi hardware. Disable only the login getty
    # so the console remains available to the kernel without exposing a prompt.
    systemctl disable getty@tty1.service >/dev/null 2>&1 || true

    success "Quiet Raspberry Pi boot configured."
}
