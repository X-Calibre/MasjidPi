#!/usr/bin/env bash

GO_MIN_VERSION="1.26.0"

version_ge() {
    [ "$(printf '%s\n' "$2" "$1" | sort -V | head -n1)" = "$2" ]
}

get_go_arch() {

    case "$(uname -m)" in

        x86_64)
            echo linux-amd64
            ;;

        aarch64)
            echo linux-arm64
            ;;

        armv7l)
            echo linux-armv6l
            ;;

        *)
            die "Unsupported CPU architecture."
            ;;
    esac
}

install_go() {

    if command_exists go; then

        CURRENT="$(go version | awk '{print $3}' | sed 's/go//')"

        if version_ge "$CURRENT" "$GO_MIN_VERSION"; then
            success "Go $CURRENT already installed."
            return
        fi

        warn "Go $CURRENT is too old. Upgrading..."
    else
        info "Go not installed."
    fi

    GO_ARCH="$(get_go_arch)"

    info "Finding latest Go release..."

    GO_VERSION="$(curl -fsSL https://go.dev/VERSION?m=text | head -n1)"

    [[ -n "$GO_VERSION" ]] || die "Unable to determine latest Go version."

    GO_TARBALL="${GO_VERSION}.${GO_ARCH}.tar.gz"

    URL="https://go.dev/dl/${GO_TARBALL}"

    info "Downloading ${GO_VERSION}..."

    wget -q --show-progress "$URL" -O /tmp/go.tar.gz

    info "Installing Go..."

    rm -rf /usr/local/go

    tar -C /usr/local -xzf /tmp/go.tar.gz

    ln -sf /usr/local/go/bin/go /usr/local/bin/go
    ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

    rm -f /tmp/go.tar.gz

    command_exists go || die "Go installation failed."

    success "$(go version)"
}