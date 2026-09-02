#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

# shellcheck source=../scripts/update.sh
source "$ROOT/scripts/update.sh"

INSTALL_DIR="$TMP/install"
UPDATE_STAGING="$TMP/staging"
UPDATE_BACKUP="$TMP/backup"
UPDATE_MARKER="$TMP/update-in-progress"
PREVIOUS_COMPONENT_PROFILE="listen,board"
RELEASE_VERSION="v-test"
EVENTS=""

record() {
    EVENTS+="$1\n"
}

info() { :; }
warn() { :; }
error() { :; }
success() { :; }
stop_service() { record stop_service; }
install_service() { record install_service; }
install_component_services() { record install_component_services; }
start_service() { record start_service; }
run_selftest() { record run_selftest; }
restore_previous_components() { record restore_previous_components; }

mkdir -p "$INSTALL_DIR" "$UPDATE_STAGING"
printf 'old\n' > "$INSTALL_DIR/runtime"
printf 'new\n' > "$UPDATE_STAGING/runtime"

activate_update v-test

expected_success='stop_service\ninstall_service\ninstall_component_services\nstart_service\nrun_selftest\n'
[[ "$EVENTS" == "$expected_success" ]]
[[ "$(<"$INSTALL_DIR/runtime")" == "new" ]]
[[ ! -e "$UPDATE_BACKUP" ]]
[[ ! -e "$UPDATE_MARKER" ]]

EVENTS=""
rm -rf "$INSTALL_DIR" "$UPDATE_BACKUP"
mkdir -p "$INSTALL_DIR" "$UPDATE_BACKUP"
printf 'failed-new\n' > "$INSTALL_DIR/runtime"
printf 'restored-old\n' > "$UPDATE_BACKUP/runtime"

rollback_update

expected_rollback='stop_service\nrestore_previous_components\ninstall_service\ninstall_component_services\nstart_service\nrun_selftest\n'
[[ "$EVENTS" == "$expected_rollback" ]]
[[ "$(<"$INSTALL_DIR/runtime")" == "restored-old" ]]

printf '[PASS] update activation and rollback install the main service before startup\n'
