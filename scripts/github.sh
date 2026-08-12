#!/usr/bin/env bash

REPO_OWNER="X-Calibre"
REPO_NAME="MasjidPi"
REPO_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}.git"
RELEASE_API="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"

update_repository() {

    if [[ ! -d "$PROJECT_ROOT/.git" ]]; then
        info "Cloning MasjidPi..."
        git clone "$REPO_URL" "$PROJECT_ROOT"
        success "Repository cloned."
        return
    fi

    info "Using local development repository."
}

prepare_release() {

    local local_version_file="$PROJECT_ROOT/VERSION"

    if [[ -x "$PROJECT_ROOT/masjidpi" && -f "$local_version_file" ]]; then
        RELEASE_DIR="$PROJECT_ROOT"
        RELEASE_VERSION="$(cat "$local_version_file")"
        info "Using bundled MasjidPi release $RELEASE_VERSION."
        return
    fi

    local latest_tag
    latest_tag="$(curl -fsSL "$RELEASE_API" | jq -r '.tag_name')"

    [[ "$latest_tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] \
        || die "Unable to determine a valid latest MasjidPi release."

    RELEASE_VERSION="$latest_tag"

    case "$ARCH" in
        x86_64)
            RELEASE_ARCH="amd64"
            ;;
        aarch64)
            RELEASE_ARCH="arm64"
            ;;
        armv7l)
            die "No official MasjidPi release is available for armv7l. Use --source instead."
            ;;
        *)
            die "Unsupported release architecture: $ARCH"
            ;;
    esac

    local archive_name="masjidpi-${RELEASE_VERSION}-linux-${RELEASE_ARCH}.tar.gz"
    local base_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${RELEASE_VERSION}"
    local url="${base_url}/${archive_name}"
    local checksums_url="${base_url}/SHA256SUMS"
    local download_dir="$(mktemp -d /tmp/masjidpi-release.XXXXXX)"

    RELEASE_DIR="$download_dir/masjidpi-${RELEASE_VERSION}-linux-${RELEASE_ARCH}"

    info "Downloading MasjidPi ${RELEASE_VERSION} (${RELEASE_ARCH})..."
    curl -fL --retry 3 "$url" -o "$download_dir/$archive_name"
    curl -fsSL --retry 3 "$checksums_url" -o "$download_dir/SHA256SUMS"

    info "Verifying release checksum..."
    (
        cd "$download_dir"
        grep -F "  $archive_name" SHA256SUMS | sha256sum -c -
    ) || die "Release checksum verification failed."

    info "Extracting release..."
    tar -xzf "$download_dir/$archive_name" -C "$download_dir"

    [[ -x "$RELEASE_DIR/masjidpi" ]] || die "Release binary is missing."
    [[ -f "$RELEASE_DIR/default.yaml" ]] || die "Release configuration is missing."
    [[ -f "$RELEASE_DIR/VERSION" ]] || die "Release version file is missing."
    [[ -f "$RELEASE_DIR/frontend/index.html" ]] || die "Release frontend is missing."

    if [[ ! -f "$RELEASE_DIR/catalogue.json" ]]; then
        warn "Release does not contain a bundled catalogue. Existing runtime catalogue will be preserved."
    fi

    success "MasjidPi ${RELEASE_VERSION} downloaded and verified."
}
