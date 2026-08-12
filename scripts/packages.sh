#!/usr/bin/env bash

install_packages() {

    local mode="${1:-release}"

    info "Updating package lists..."
    apt-get update

    info "Installing runtime dependencies..."

    apt-get install -y \
        curl \
        jq \
        mpv \
        ffmpeg \
        alsa-utils \
        ca-certificates \
        tar

    if [[ "$mode" == "source" ]]; then
        info "Installing source-build dependencies..."
        apt-get install -y \
            git \
            build-essential
    fi

    success "Dependencies installed."
}
