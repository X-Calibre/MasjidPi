#!/usr/bin/env bash
# The single-quoted grep patterns below intentionally match literal shell
# expressions in the updater and must not expand in this test process.
# shellcheck disable=SC2016

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
UPDATER="$ROOT/appliance-image/update/masjidframe-update"

bash -n "$UPDATER"

grep -Fq 'readonly appliance_user=masjidframe' "$UPDATER"
grep -Fq 'getent shadow "$appliance_user"' "$UPDATER"
grep -Fq 'Running $appliance_user account has no usable password identity.' "$UPDATER"
grep -Fq 'Target system has an incompatible $appliance_user account.' "$UPDATER"
grep -Fq "printf '%s:%s\\n'" "$UPDATER"
grep -Fq '"$running_user_password_hash" |' "$UPDATER"
grep -Fq '/usr/sbin/chpasswd' "$UPDATER"
grep -Fq -- '--encrypted' "$UPDATER"
grep -Fq 'Unable to preserve the $appliance_user account password.' "$UPDATER"
grep -Fq 'echo "PASS  $appliance_user account password"' "$UPDATER"

if grep -Eq 'echo .*password_hash|printf .*password_hash.*>&2' "$UPDATER"; then
    echo '[FAIL] updater may print the encrypted appliance password hash' >&2
    exit 1
fi

printf '[PASS] appliance updates preserve the masjidframe password hash without logging it\n'
