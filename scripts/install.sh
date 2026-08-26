#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

source "$SCRIPT_DIR/common.sh"
source "$SCRIPT_DIR/os.sh"
source "$SCRIPT_DIR/version.sh"
source "$SCRIPT_DIR/components.sh"
source "$SCRIPT_DIR/packages.sh"
source "$SCRIPT_DIR/go.sh"
source "$SCRIPT_DIR/github.sh"
source "$SCRIPT_DIR/runtime.sh"
source "$SCRIPT_DIR/build.sh"
source "$SCRIPT_DIR/install_mode.sh"
source "$SCRIPT_DIR/service.sh"
source "$SCRIPT_DIR/selftest.sh"
source "$SCRIPT_DIR/update.sh"
source "$SCRIPT_DIR/preflight.sh"

SOURCE_MODE=false
INSTALL_SUCCEEDED=false

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

cleanup_failed_install() {
    local status=$?

    if [[ "$status" -ne 0 && "${INSTALL_MODE:-}" == "install" && "$INSTALL_SUCCEEDED" != true ]]; then
        warn "Installation did not complete. Cleaning up the application runtime..."

        if systemctl is-active --quiet masjidpi-display.service 2>/dev/null; then
            systemctl stop masjidpi-display.service || true
        fi

        stop_service || true
        rm -rf "$INSTALL_DIR"
        rm -rf "$UPDATE_STAGING" "$UPDATE_BACKUP"
        rm -f "$UPDATE_MARKER"
    fi

    return "$status"
}

main() {
    require_root
    parse_args "${1:-}"
    preflight_install

    detect_os
    detect_arch
    detect_install_mode
    select_components
    migrate_legacy_preferences
    migrate_catalogue_refresh_interval

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
        RELEASE_VERSION="$("$PROJECT_ROOT/backend/build/masjidpi" --version)"
    else
        info "Release installation selected."
        install_packages release
        prepare_release
    fi

    save_components

    if [ "$INSTALL_MODE" = "update" ]; then
        trap cleanup_update EXIT
        prepare_update
        activate_update "$RELEASE_VERSION"
    else
        stop_service
        install_runtime
        install_service
        install_component_services
        start_service
        run_selftest "$RELEASE_VERSION"
    fi

    INSTALL_SUCCEEDED=true
    success "Installation completed."
    print_summary
}

trap cleanup_failed_install EXIT

main "$@"
