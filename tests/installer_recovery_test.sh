#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

# shellcheck source=../scripts/install_mode.sh
source "$ROOT/scripts/install_mode.sh"

INSTALL_DIR="$TMP/install"
UPDATE_BACKUP="$TMP/backup"
UPDATE_MARKER="$TMP/update-in-progress"
INSTALL_MODE=""
SELFTEST_SUCCEEDS=false

info() { :; }
warn() { :; }
error() { :; }
success() { :; }
stop_service() { :; }
start_service() { :; }
run_selftest() { $SELFTEST_SUCCEEDS; }

mkdir -p "$INSTALL_DIR" "$UPDATE_BACKUP"
printf 'failed-new\n' > "$INSTALL_DIR/runtime"
printf 'restored-old\n' > "$UPDATE_BACKUP/runtime"
printf 'v-test\n' > "$UPDATE_MARKER"

if detect_install_mode; then
    printf '[FAIL] interrupted recovery unexpectedly passed validation\n' >&2
    exit 1
fi

[[ "$(<"$INSTALL_DIR/runtime")" == "restored-old" ]]
[[ ! -e "$UPDATE_BACKUP" ]]
[[ ! -e "$UPDATE_MARKER" ]]

printf '[PASS] interrupted recovery clears a resolved transaction marker\n'
