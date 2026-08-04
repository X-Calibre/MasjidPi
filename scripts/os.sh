#!/usr/bin/env bash

detect_os() {

    source /etc/os-release

    DISTRO="$ID"
    VERSION="$VERSION_ID"

    case "$DISTRO" in
        debian|ubuntu|linuxmint|raspbian)
            ;;
        *)
            die "Unsupported distribution: $DISTRO"
            ;;
    esac

    success "$PRETTY_NAME"
}

detect_arch() {

    ARCH="$(uname -m)"

    case "$ARCH" in
        x86_64|aarch64|armv7l)
            success "$ARCH"
            ;;
        *)
            die "Unsupported architecture: $ARCH"
            ;;
    esac
}