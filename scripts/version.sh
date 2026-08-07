#!/usr/bin/env bash

get_version() {

    local version_file="$PROJECT_ROOT/version.json"

    if [[ -f "$version_file" ]]; then
        jq -r '.version' "$version_file"
    else
        echo "development"
    fi
}