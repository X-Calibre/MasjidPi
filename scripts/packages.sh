#!/usr/bin/env bash

install_packages() {

    info "Updating package lists..."

    apt-get update

    info "Installing dependencies..."

    apt-get install -y \
        git \
        curl \
        wget \
        jq \
        mpv \
        ffmpeg \
        alsa-utils \
        build-essential \
        ca-certificates \
        tar

    success "Dependencies installed."
}

install_go() {

    if command -v go >/dev/null; then

        success "Go already installed."

        return
    fi

    die "Go is not installed.

Please install Go manually for now.

Automatic Go installation will be added in the next version."
}
