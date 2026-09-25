#!/usr/bin/env bash
# These single-quoted patterns intentionally match literal shell expressions.
# shellcheck disable=SC2016

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
UPDATER="$ROOT/appliance-image/update/masjidpi-update"

bash -n "$UPDATER"

grep -Fqx 'extraction_parent=$(dirname "$bundle")' "$UPDATER"
grep -Fqx 'extraction_dir=$(mktemp -d "$extraction_parent/.verify.XXXXXX")' "$UPDATER"
grep -Fqx '    staging_parent=/persistent/updates' "$UPDATER"

if grep -Fqx 'extraction_dir=$(mktemp -d)' "$UPDATER"; then
    echo '[FAIL] update verification still uses the RAM-backed default temporary directory' >&2
    exit 1
fi

if grep -Fq '\\nextraction_parent=' "$UPDATER"; then
    echo '[FAIL] update workspace statements contain literal newline escapes' >&2
    exit 1
fi

grep -Fqx '    exec 9>/run/lock/masjidpi-update.lock' "$UPDATER"
grep -Fqx '    if ! flock -n 9; then' "$UPDATER"
grep -Fqx '    stale_install_dirs=("$staging_parent"/install.*)' "$UPDATER"
grep -Fqx '    stale_target_roots=("$staging_parent"/target-root.*)' "$UPDATER"
grep -Fqx '        if mountpoint -q "$stale_target_root"; then' "$UPDATER"
grep -Fqx '        rm -rf -- "$stale_target_root"' "$UPDATER"
grep -Fqx '        rm -rf -- "$stale_install_dir"' "$UPDATER"
grep -Fqx '        rm -f -- "$boot_dir/.$boot_file.masjidpi-update"' "$UPDATER"

lock_line=$(grep -nF '    exec 9>/run/lock/masjidpi-update.lock' "$UPDATER" | cut -d: -f1)
cleanup_line=$(grep -nF '    stale_install_dirs=("$staging_parent"/install.*)' "$UPDATER" | cut -d: -f1)
plan_line=$(grep -nF '    "$0" plan "$bundle" "$signature"' "$UPDATER" | head -n 1 | cut -d: -f1)

if (( lock_line >= cleanup_line || cleanup_line >= plan_line )); then
    echo '[FAIL] interrupted-update cleanup must run under the lock before verification' >&2
    exit 1
fi

printf '[PASS] appliance update verification and installation use persistent storage\n'
