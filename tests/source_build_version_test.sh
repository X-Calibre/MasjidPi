#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

PROJECT_ROOT="$TMP/project"
mkdir -p "$PROJECT_ROOT/backend"
cp "$ROOT/version.json" "$PROJECT_ROOT/version.json"

info() { :; }
success() { :; }
die() {
    printf '%s\n' "$*" >&2
    return 1
}
go() {
    printf '%s\n' "$*" > "$TMP/go-args"
}

# shellcheck source=../scripts/build.sh
source "$ROOT/scripts/build.sh"

build_project

expected_version="$(jq -r '.version' "$ROOT/version.json")-dev"
grep -F -- "-ldflags -X github.com/X-Calibre/MasjidPi/backend/internal/version.Version=$expected_version" "$TMP/go-args" >/dev/null

printf '[PASS] source builds embed the current project version with a development suffix\n'
