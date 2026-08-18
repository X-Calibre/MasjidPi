#!/usr/bin/env bash

preflight_install() {
    info "Checking installation environment..."

    command_exists apt-get || die "apt-get is required on supported MasjidPi systems."
    command_exists systemctl || die "systemctl is required. MasjidPi must run as a system service."
    command_exists curl || die "curl is required."
    command_exists tar || die "tar is required."
    command_exists sha256sum || die "sha256sum is required."

    if [[ ! -d /run/systemd/system ]]; then
        die "systemd is not running. MasjidPi production installation requires systemd."
    fi

    local system_state
    system_state="$(systemctl is-system-running 2>/dev/null || true)"

    case "$system_state" in
        running|degraded|starting)
            ;;
        *)
            die "systemd is not ready for installation (state: ${system_state:-unknown})."
            ;;
    esac

    success "Installation environment is ready."
}
