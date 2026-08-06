#!/usr/bin/env bash

REPO_URL="https://github.com/X-Calibre/MasjidPi.git"

update_repository() {

    if [[ ! -d "$PROJECT_ROOT/.git" ]]; then

        info "Cloning MasjidPi..."

        git clone "$REPO_URL" "$PROJECT_ROOT"

        success "Repository cloned."

        return
    fi

    info "Using local development repository."

}