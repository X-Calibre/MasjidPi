#!/usr/bin/env bash

install_packages() {

    local mode="${1:-release}"

    info "Updating package lists..."
    apt-get update

    info "Installing shared runtime dependencies..."

    apt-get install -y \
        curl \
        jq \
        ca-certificates \
        tar

    if $INSTALL_LISTEN; then
        info "Installing Listen runtime dependencies..."
        apt-get install -y \
            mpv \
            ffmpeg \
            alsa-utils
    fi

    if $INSTALL_BOARD; then
        info "Installing MasjidBoard display dependencies..."
        apt-get install -y \
            cog \
            libwpewebkit-2.0-1

        if is_raspberry_pi; then
            info "Installing Raspberry Pi splash dependencies..."
            apt-get install -y \
                plymouth \
                plymouth-themes
        fi
    fi

    if [[ "$mode" == "source" ]]; then
        info "Installing source-build dependencies..."
        apt-get install -y \
            git \
            build-essential
    fi

    success "Dependencies installed."
}
