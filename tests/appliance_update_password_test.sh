#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
UPDATER="$ROOT/appliance-image/update/masjidpi-update"

bash -n "$UPDATER"

grep -Fq 'mapfile -t running_pi_shadow_records < <(getent shadow pi)' "$UPDATER"
grep -Fq 'Running pi account has no usable password identity.' "$UPDATER"
grep -Fq 'Target system has an incompatible pi account.' "$UPDATER"
grep -Fq "printf 'pi:%s\\n' \"\$running_pi_password_hash\"" "$UPDATER"
grep -Fq '/usr/sbin/chpasswd' "$UPDATER"
grep -Fq -- '--encrypted' "$UPDATER"
grep -Fq 'Unable to preserve the pi account password.' "$UPDATER"
grep -Fq 'echo "PASS  pi account password"' "$UPDATER"

if grep -Eq 'echo .*password_hash|printf .*password_hash.*>&2' "$UPDATER"; then
    echo '[FAIL] updater may print the encrypted pi password hash' >&2
    exit 1
fi

printf '[PASS] appliance updates preserve the pi password hash without logging it\n'
