#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/common.sh"
source "$SCRIPT_DIR/os.sh"
source "$SCRIPT_DIR/version.sh"
source "$SCRIPT_DIR/packages.sh"
source "$SCRIPT_DIR/go.sh"
source "$SCRIPT_DIR/github.sh"
source "$SCRIPT_DIR/runtime.sh"
source "$SCRIPT_DIR/build.sh"
source "$SCRIPT_DIR/install_mode.sh"
source "$SCRIPT_DIR/service.sh"
source "$SCRIPT_DIR/selftest.sh"

SOURCE_MODE=false

usage() {
    cat <<EOF
MasjidPi installer

Usage:
  sudo $0              Install the latest pre-built release
  sudo $0 --source     Install from the local Git source tree
EOF
}

parse_args() {
    case "${1:-}" in
        "") ;;
        --source)
            SOURCE_MODE=true
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            usage
            die "Unknown option: $1"
            ;;
    esac
}

main() {
    require_root
    parse_args "${1:-}"

    detect_os
    detect_arch
    detect_install_mode

    if [ "$INSTALL_MODE" = "install" ]; then
        info "Fresh installation detected."
    else
        info "Existing installation detected."
    fi

    if $SOURCE_MODE; then
        info "Source installation selected."
        install_packages source
        install_go
        update_repository
        build_project
    else
        info "Release installation selected."
        install_packages release
        prepare_release
    fi

    stop_service
    install_runtime
    install_service
    start_service
    run_selftest

    success "Installation completed."
    print_summary
}

main "$@"
