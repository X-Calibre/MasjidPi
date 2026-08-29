#!/usr/bin/env bash

configure_raspberry_pi_wifi() {
    local model_path="${MASJIDPI_DEVICE_MODEL_PATH:-/proc/device-tree/model}"
    local config_dir="${MASJIDPI_NETWORKMANAGER_CONF_DIR:-/etc/NetworkManager/conf.d}"
    local config_path="$config_dir/masjidpi-wifi-powersave.conf"
    local model

    if [[ ! -r "$model_path" ]]; then
        return
    fi

    model="$(tr -d '\0' < "$model_path")"
    if [[ "$model" != *"Raspberry Pi"* ]]; then
        return
    fi

    if ! command_exists nmcli || ! systemctl is-active --quiet NetworkManager.service; then
        info "Raspberry Pi detected; NetworkManager Wi-Fi tuning is unavailable."
        return
    fi

    local -a wifi_interfaces=()
    mapfile -t wifi_interfaces < <(
        nmcli -t -f DEVICE,TYPE device status 2>/dev/null \
            | awk -F: '$2 == "wifi" && $1 != "" {print $1}'
    )
    if [[ ${#wifi_interfaces[@]} -eq 0 ]]; then
        info "Raspberry Pi detected; no NetworkManager Wi-Fi interface found."
        return
    fi

    info "Configuring Raspberry Pi Wi-Fi for appliance operation..."

    mkdir -p "$config_dir"
    local temporary_config
    temporary_config="$(mktemp)"
    printf '%s\n' '[connection]' 'wifi.powersave=2' > "$temporary_config"
    if [[ ! -f "$config_path" ]] || ! cmp -s "$temporary_config" "$config_path"; then
        install -m 0644 "$temporary_config" "$config_path"
    fi
    rm -f "$temporary_config"

    local interface connection_uuid
    for interface in "${wifi_interfaces[@]}"; do
        connection_uuid="$(nmcli -g GENERAL.CON-UUID device show "$interface" 2>/dev/null || true)"
        if [[ -n "$connection_uuid" && "$connection_uuid" != "--" ]]; then
            nmcli connection modify "$connection_uuid" 802-11-wireless.powersave 2 \
                >/dev/null 2>&1 || warn "Could not persist Wi-Fi power saving policy for $interface."
        fi

        if command_exists iw; then
            iw dev "$interface" set power_save off >/dev/null 2>&1 \
                || warn "Could not disable Wi-Fi power saving immediately for $interface."
        fi
    done

    success "Wi-Fi power saving disabled."
}
