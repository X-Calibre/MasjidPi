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

    # The backend currently initializes the playback subsystem for every
    # profile. Keep its runtime dependencies installed until Listen is fully
    # split from the common backend lifecycle.
    apt-get install -y \
        mpv \
        ffmpeg \
        alsa-utils

    if $INSTALL_BOARD; then
        info "Installing MasjidBoard display dependencies..."
        apt-get install -y \
            cog \
            libwpewebkit-2.0-1
    fi

    if [[ "$mode" == "source" ]]; then
        info "Installing source-build dependencies..."
        apt-get install -y \
            git \
            build-essential
    fi

    success "Dependencies installed."
}
