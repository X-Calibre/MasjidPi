#!/usr/bin/env bash
# These single-quoted patterns intentionally match literal shell expressions.
# shellcheck disable=SC2016

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
UPDATER="$ROOT/appliance-image/update/masjidpi-update"

bash -n "$UPDATER"

grep -Fq 'extraction_parent=$(dirname "$bundle")' "$UPDATER"
grep -Fq 'extraction_dir=$(mktemp -d "$extraction_parent/.verify.XXXXXX")' "$UPDATER"
grep -Fq 'staging_parent=/persistent/updates' "$UPDATER"

if grep -Fqx 'extraction_dir=$(mktemp -d)' "$UPDATER"; then
    echo '[FAIL] update verification still uses the RAM-backed default temporary directory' >&2
    exit 1
fi

printf '[PASS] appliance update verification and installation use persistent storage\n'
