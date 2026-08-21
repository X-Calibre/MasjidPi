#!/usr/bin/env bash

COMPONENTS_FILE="/etc/masjidpi/components.env"
INSTALL_LISTEN=true
INSTALL_BOARD=true

load_existing_components() {
    if [[ ! -f "$COMPONENTS_FILE" ]]; then
        return 1
    fi

    local value
    value="$(sed -n 's/^MASJIDPI_COMPONENTS=//p' "$COMPONENTS_FILE" | tail -1)"
    case "$value" in
        listen)
            INSTALL_LISTEN=true
            INSTALL_BOARD=false
            ;;
        board)
            INSTALL_LISTEN=false
            INSTALL_BOARD=true
            ;;
        listen,board|board,listen|both)
            INSTALL_LISTEN=true
            INSTALL_BOARD=true
            ;;
        *)
            return 1
            ;;
    esac
}

select_components() {
    if load_existing_components; then
        info "Preserving installed MasjidPi component profile: $(component_profile)"
        return 0
    fi

    if [[ ! -t 0 ]]; then
        info "Non-interactive installation detected; installing Listen + Board."
        INSTALL_LISTEN=true
        INSTALL_BOARD=true
        return 0
    fi

    cat <<'EOF'

Select MasjidPi components to install:
  1) Listen
  2) Board
  3) Listen + Board
EOF

    local choice
    while true; do
        read -r -p "Choice [3]: " choice
        choice="${choice:-3}"
        case "$choice" in
            1)
                INSTALL_LISTEN=true
                INSTALL_BOARD=false
                break
                ;;
            2)
                INSTALL_LISTEN=false
                INSTALL_BOARD=true
                break
                ;;
            3)
                INSTALL_LISTEN=true
                INSTALL_BOARD=true
                break
                ;;
            *)
                warn "Enter 1, 2, or 3."
                ;;
        esac
    done

    info "Selected MasjidPi component profile: $(component_profile)"
}

component_profile() {
    if $INSTALL_LISTEN && $INSTALL_BOARD; then
        printf '%s\n' "listen,board"
    elif $INSTALL_LISTEN; then
        printf '%s\n' "listen"
    else
        printf '%s\n' "board"
    fi
}

save_components() {
    install -d -m 0755 /etc/masjidpi
    printf 'MASJIDPI_COMPONENTS=%s\n' "$(component_profile)" > "$COMPONENTS_FILE"
    chmod 0644 "$COMPONENTS_FILE"
    success "Installed component profile saved to $COMPONENTS_FILE."
}
