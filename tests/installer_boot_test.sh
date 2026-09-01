#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

export MASJIDPI_BOOT_FIRMWARE_DIR="$TMP/boot"
export MASJIDPI_FORCE_RASPBERRY_PI=1
mkdir -p "$MASJIDPI_BOOT_FIRMWARE_DIR"

cat > "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt" <<'EOF'
console=serial0,115200 console=tty1 root=PARTUUID=test-02 rootfstype=ext4 rootwait
EOF

cat > "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt" <<'EOF'
dtoverlay=vc4-kms-v3d
[all]
EOF

INSTALL_BOARD=true
SYSTEMCTL_CALLS=""

info() { :; }
warn() { :; }
success() { :; }
systemctl() {
    SYSTEMCTL_CALLS+="$*\n"
    return 0
}

# shellcheck source=../scripts/boot.sh
source "$ROOT/scripts/boot.sh"

configure_quiet_boot
configure_quiet_boot

cmdline="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt")"
config="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt")"

[[ "$cmdline" == *"console=tty1"* ]]
[[ "$(grep -o '\bquiet\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'loglevel=3' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'systemd.show_status=false' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -c '^disable_splash=1$' <<< "$config")" -eq 1 ]]
[[ "$SYSTEMCTL_CALLS" == *"disable getty@tty1.service"* ]]

printf '[PASS] quiet boot configuration is idempotent and preserves console=tty1\n'
