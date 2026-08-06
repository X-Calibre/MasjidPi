#!/usr/bin/env bash

detect_install_mode() {

    if [ -d "$INSTALL_DIR" ]; then
        INSTALL_MODE="update"
    else
        INSTALL_MODE="install"
    fi

}