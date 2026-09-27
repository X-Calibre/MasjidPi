#!/usr/bin/env bash
# These single-quoted patterns intentionally match literal shell expressions.
# shellcheck disable=SC2016

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SETUP="$ROOT/appliance-image/image/pi3-ab/setup.sh"
CUSTOMIZE="$ROOT/appliance-image/image/pi3-ab/image.d/hooks/customize82-masjidframe"

bash -n "$SETUP"
bash -n "$CUSTOMIZE"

grep -Fqx '/persistent/masjidpi/etc /etc/masjidframe none bind,x-systemd.requires-mounts-for=/persistent 0 0' "$SETUP"
grep -Fqx '/persistent/masjidpi/var /var/lib/masjidframe none bind,x-systemd.requires-mounts-for=/persistent 0 0' "$SETUP"

grep -Fq '"$filesystem/persistent/masjidpi/etc"' "$CUSTOMIZE"
grep -Fq '"$filesystem/persistent/masjidpi/var"' "$CUSTOMIZE"
grep -Fq '"$filesystem/persistent/masjidpi/etc/config.yaml"' "$CUSTOMIZE"
grep -Fq '"$filesystem/persistent/masjidpi/var/catalogue.json"' "$CUSTOMIZE"

if grep -Fq '/persistent/masjidframe/' "$SETUP"; then
    echo '[FAIL] system image still mounts an empty renamed persistent data root' >&2
    exit 1
fi

if grep -Fq '$filesystem/persistent/masjidframe/' "$CUSTOMIZE"; then
    echo '[FAIL] fresh image still initializes the incompatible persistent data root' >&2
    exit 1
fi

printf '[PASS] MasjidFrame runtime uses the v1.6.0-compatible persistent data root\n'
