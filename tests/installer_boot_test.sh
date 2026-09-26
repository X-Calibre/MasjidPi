#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

export MASJIDFRAME_BOOT_FIRMWARE_DIR="$TMP/boot"
export MASJIDFRAME_PLYMOUTH_THEME_DIR="$TMP/plymouth-theme"
export MASJIDFRAME_PLYMOUTH_QUIT_DROPIN_DIR="$TMP/plymouth-dropin"
export MASJIDFRAME_BOOT_READONLY_SERVICE_FILE="$TMP/systemd/masjidframe-boot-readonly.service"
export MASJIDFRAME_BOOT_APT_HOOK_FILE="$TMP/apt/99-masjidframe-boot-firmware"
export MASJIDFRAME_FORCE_RASPBERRY_PI=1
export PROJECT_ROOT="$ROOT"
mkdir -p "$MASJIDFRAME_BOOT_FIRMWARE_DIR" "$TMP/bin"

reset_boot_files() {
    cat > "$MASJIDFRAME_BOOT_FIRMWARE_DIR/cmdline.txt" <<'EOF'
console=serial0,115200 console=tty1 root=PARTUUID=test-02 rootfstype=ext4 rootwait
EOF

    cat > "$MASJIDFRAME_BOOT_FIRMWARE_DIR/config.txt" <<'EOF'
dtoverlay=vc4-kms-v3d
[all]
EOF
}

cat > "$TMP/bin/plymouth-set-default-theme" <<EOF
#!/usr/bin/env bash
printf '%s\n' "\$*" >> "$TMP/plymouth-calls"
EOF
chmod +x "$TMP/bin/plymouth-set-default-theme"
export PATH="$TMP/bin:$PATH"

INSTALL_BOARD=true
SYSTEMCTL_CALLS=""
MOUNT_CALLS=""

info() { :; }
warn() { :; }
success() { :; }
systemctl() {
    SYSTEMCTL_CALLS+="$*\n"
    return 0
}
mountpoint() {
    [[ "$1" == "-q" && "$2" == "$MASJIDFRAME_BOOT_FIRMWARE_DIR" ]]
}
mount() {
    MOUNT_CALLS+="$*\n"
}
sync() {
    MOUNT_CALLS+="sync $*\n"
}

# shellcheck source=../scripts/boot.sh
source "$ROOT/scripts/boot.sh"

# Raspberry Pi Board installations use the upright, non-rotated Plymouth assets.
reset_boot_files
mkdir -p "$MASJIDFRAME_PLYMOUTH_THEME_DIR"
touch "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash-logo-appliance.png"
prepare_boot_firmware_update
configure_quiet_boot
configure_quiet_boot
configure_boot_splash
configure_boot_splash
configure_boot_firmware_protection

standard_cmdline="$(cat "$MASJIDFRAME_BOOT_FIRMWARE_DIR/cmdline.txt")"
standard_config="$(cat "$MASJIDFRAME_BOOT_FIRMWARE_DIR/config.txt")"

[[ "$standard_cmdline" == *"console=tty1"* ]]
[[ "$(grep -o '\bquiet\b' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'loglevel=3' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'systemd.show_status=false' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o '\bsplash\b' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'plymouth.ignore-serial-consoles' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -c '^disable_splash=1$' <<< "$standard_config")" -eq 1 ]]
[[ "$SYSTEMCTL_CALLS" == *"disable getty@tty1.service"* ]]
[[ -f "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe.plymouth" ]]
[[ -f "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash.script" ]]
[[ -f "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash-logo.png" ]]
[[ -f "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash-logo-portrait.png" ]]
[[ -f "$MASJIDFRAME_PLYMOUTH_QUIT_DROPIN_DIR/masjidframe.conf" ]]
[[ -f "$MASJIDFRAME_BOOT_READONLY_SERVICE_FILE" ]]
[[ -f "$MASJIDFRAME_BOOT_APT_HOOK_FILE" ]]
[[ "$SYSTEMCTL_CALLS" == *"enable masjidframe-boot-readonly.service"* ]]
[[ "$MOUNT_CALLS" == *"-o remount,rw $MASJIDFRAME_BOOT_FIRMWARE_DIR"* ]]
[[ "$MOUNT_CALLS" == *"sync -f $MASJIDFRAME_BOOT_FIRMWARE_DIR"* ]]
[[ "$MOUNT_CALLS" == *"-o remount,ro $MASJIDFRAME_BOOT_FIRMWARE_DIR"* ]]
if grep -q 'Rotate' "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash.script"; then
    echo "[FAIL] standard Plymouth theme unexpectedly rotates the splash" >&2
    exit 1
fi

# Reconfiguration remains idempotent and does not restore retired assets.
configure_quiet_boot
configure_boot_splash

cmdline="$(cat "$MASJIDFRAME_BOOT_FIRMWARE_DIR/cmdline.txt")"
config="$(cat "$MASJIDFRAME_BOOT_FIRMWARE_DIR/config.txt")"

[[ "$(grep -o '\bquiet\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'loglevel=3' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'systemd.show_status=false' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o '\bsplash\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -c '^disable_splash=1$' <<< "$config")" -eq 1 ]]
cmp -s "$ROOT/scripts/masjidframe-splash-standard.script" \
    "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash.script"
[[ ! -e "$MASJIDFRAME_PLYMOUTH_THEME_DIR/masjidframe-splash-logo-appliance.png" ]]
[[ "$(grep -c -- '^-R masjidframe$' "$TMP/plymouth-calls")" -eq 3 ]]

printf '[PASS] quiet boot and upright branded splash cover supported Raspberry Pi Board profiles\n'
