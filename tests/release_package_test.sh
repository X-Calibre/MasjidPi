#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/release.yml"

# shellcheck disable=SC2016
grep -Fq 'cp scripts/99-masjidframe-boot-firmware "$package_dir/scripts/"' "$WORKFLOW"
# shellcheck disable=SC2016
grep -Fq 'require_file "$package_dir/scripts/masjidframe-boot-readonly.service"' "$WORKFLOW"
# shellcheck disable=SC2016
grep -Fq 'require_file "$package_dir/scripts/99-masjidframe-boot-firmware"' "$WORKFLOW"
# shellcheck disable=SC2016
grep -Fq 'require_file "$package_dir/frontend/updates.html"' "$WORKFLOW"
# shellcheck disable=SC2016
grep -Fq 'require_file "$package_dir/frontend/update-status.js"' "$WORKFLOW"

printf '[PASS] release workflow validates boot protection and update Web UI assets\n'
