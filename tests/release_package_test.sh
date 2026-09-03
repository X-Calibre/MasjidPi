#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/release.yml"

grep -Fq 'cp scripts/99-masjidpi-boot-firmware "'\
grep -Fq 'test -f "'\
grep -Fq 'test -f "'\

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/masjidpi-boot-readonly.service"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/masjidpi-boot-readonly.service"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/masjidpi-boot-readonly.service"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/masjidpi-boot-readonly.service"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/masjidpi-boot-readonly.service"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
package_dir/scripts/"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/masjidpi-boot-readonly.service"' "$WORKFLOW"
grep -Fq 'test -f "$package_dir/scripts/99-masjidpi-boot-firmware"' "$WORKFLOW"

printf '[PASS] release workflow packages and validates boot protection assets\n'
