#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/common.sh"
source "$SCRIPT_DIR/os.sh"
source "$SCRIPT_DIR/packages.sh"
source "$SCRIPT_DIR/go.sh"
source "$SCRIPT_DIR/github.sh"
source "$SCRIPT_DIR/runtime.sh"
source "$SCRIPT_DIR/build.sh"
source "$SCRIPT_DIR/install_mode.sh"
source "$SCRIPT_DIR/service.sh"
source "$SCRIPT_DIR/selftest.sh"

main() {

    require_root

    detect_os
    detect_arch

    detect_install_mode

    if [ "$INSTALL_MODE" = "install" ]; then
        info "Fresh installation detected."
    else
        info "Existing installation detected."
    fi

    install_packages
    install_go

    update_repository

    build_project

    install_runtime

    install_service

    start_service

    run_selftest

    success "Installation completed."

    print_summary
    }

main "$@"
