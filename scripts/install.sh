#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/common.sh"
source "$SCRIPT_DIR/os.sh"
source "$SCRIPT_DIR/packages.sh"
source "$SCRIPT_DIR/go.sh"
source "$SCRIPT_DIR/github.sh"
source "$SCRIPT_DIR/build.sh"
source "$SCRIPT_DIR/service.sh"
source "$SCRIPT_DIR/selftest.sh"

main() {

    require_root

    detect_os
    detect_arch

    install_packages
    install_go

    build_project

    run_selftest

    success ""
    success "Installation completed."
}

main "$@"