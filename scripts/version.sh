#!/usr/bin/env bash

get_version() {

    if [[ -n "${RELEASE_VERSION:-}" ]]; then
        echo "$RELEASE_VERSION"
        return
    fi

    if [[ -f "$INSTALL_DIR/VERSION" ]]; then
        cat "$INSTALL_DIR/VERSION"
        return
    fi

    local version_file="$PROJECT_ROOT/version.json"

    if [[ -f "$version_file" ]]; then
        jq -r '.version' "$version_file"
    else
        echo "development"
    fi
}
