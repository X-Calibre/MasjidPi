#!/usr/bin/env bash

REPO_URL="https://github.com/X-Calibre/MasjidPi.git"

update_repository() {

    if [[ ! -d "$PROJECT_ROOT/.git" ]]; then

        info "Cloning MasjidPi..."

        git clone "$REPO_URL" "$PROJECT_ROOT"

        success "Repository cloned."

        return
    fi

    cd "$PROJECT_ROOT"

    info "Updating repository..."

    git fetch origin

    STATUS="$(git status --porcelain)"

    if [[ -n "$STATUS" ]]; then

        warn "Repository contains local changes."

        warn "Skipping automatic update."

        return
    fi

    CURRENT_BRANCH="$(git branch --show-current)"

    git pull --ff-only origin "$CURRENT_BRANCH"

    success "Repository updated."
}