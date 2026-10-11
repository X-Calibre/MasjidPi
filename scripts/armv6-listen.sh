#!/usr/bin/env bash

# Overrides used only by the dedicated ARMv6 Listen installer.
GO_MIN_VERSION="1.26.5"

detect_arch() {
    ARCH="$(uname -m)"
    [[ "$ARCH" == armv6l ]] || die "This installer requires an ARMv6 Raspberry Pi (detected $ARCH)."
    local model_path="${MASJIDFRAME_DEVICE_MODEL_PATH:-/proc/device-tree/model}"
    [[ -r "$model_path" ]] && grep -aq 'Raspberry Pi' "$model_path" \
        || die "This installer requires Raspberry Pi hardware."
    success "ARMv6 Raspberry Pi detected."
}

select_components() {
    if load_existing_components; then
        PREVIOUS_COMPONENT_PROFILE="$(component_profile)"
    fi
    set_component_profile listen
    info "Installing Listen only; Board is disabled."
}
