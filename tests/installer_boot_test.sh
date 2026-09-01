#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

export MASJIDPI_BOOT_FIRMWARE_DIR="$TMP/boot"
export MASJIDPI_PLYMOUTH_THEME_DIR="$TMP/plymouth-theme"
export MASJIDPI_FORCE_RASPBERRY_PI=1
export PROJECT_ROOT="$ROOT"
mkdir -p "$MASJIDPI_BOOT_FIRMWARE_DIR" "$TMP/bin"

cat > "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt" <<'EOF'
console=serial0,115200 console=tty1 root=PARTUUID=test-02 rootfstype=ext4 rootwait
EOF

cat > "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt" <<'EOF'
dtoverlay=vc4-kms-v3d
[all]
EOF

cat > "$TMP/bin/plymouth-set-default-theme" <<EOF
#!/usr/bin/env bash
printf '%s\n' "\$*" >> "$TMP/plymouth-calls"
EOF
chmod +x "$TMP/bin/plymouth-set-default-theme"
export PATH="$TMP/bin:$PATH"

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
configure_boot_splash
configure_boot_splash

cmdline="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt")"
config="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt")"

[[ "$cmdline" == *"console=tty1"* ]]
[[ "$(grep -o '\bquiet\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'loglevel=3' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'systemd.show_status=false' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o '\bsplash\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'plymouth.ignore-serial-consoles' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -c '^disable_splash=1$' <<< "$config")" -eq 1 ]]
[[ "$SYSTEMCTL_CALLS" == *"disable getty@tty1.service"* ]]
[[ -f "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi.plymouth" ]]
[[ -f "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi-splash.script" ]]
[[ "$(grep -c -- '^-R masjidpi$' "$TMP/plymouth-calls")" -eq 2 ]]

printf '[PASS] quiet boot and branded splash configuration are idempotent and preserve console=tty1\n'
