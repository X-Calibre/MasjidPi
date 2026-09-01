#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

export MASJIDPI_BOOT_FIRMWARE_DIR="$TMP/boot"
export MASJIDPI_PLYMOUTH_THEME_DIR="$TMP/plymouth-theme"
export MASJIDPI_PLYMOUTH_QUIT_DROPIN_DIR="$TMP/plymouth-dropin"
export MASJIDPI_USB_SYSFS_ROOT="$TMP/usb"
export MASJIDPI_DRM_SYSFS_ROOT="$TMP/drm"
export MASJIDPI_FORCE_RASPBERRY_PI=1
export PROJECT_ROOT="$ROOT"
mkdir -p "$MASJIDPI_BOOT_FIRMWARE_DIR" "$TMP/bin" "$MASJIDPI_USB_SYSFS_ROOT" "$MASJIDPI_DRM_SYSFS_ROOT"

reset_boot_files() {
    cat > "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt" <<'EOF'
console=serial0,115200 console=tty1 root=PARTUUID=test-02 rootfstype=ext4 rootwait
EOF

    cat > "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt" <<'EOF'
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

info() { :; }
warn() { :; }
success() { :; }
systemctl() {
    SYSTEMCTL_CALLS+="$*\n"
    return 0
}

# shellcheck source=../scripts/boot.sh
source "$ROOT/scripts/boot.sh"

# A Raspberry Pi Board installation with ordinary HDMI hardware gets the same
# quiet branded boot, using the unrotated landscape Plymouth script.
reset_boot_files
configure_quiet_boot
configure_quiet_boot
configure_boot_splash
configure_boot_splash

standard_cmdline="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt")"
standard_config="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt")"

[[ "$standard_cmdline" == *"console=tty1"* ]]
[[ "$(grep -o '\bquiet\b' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'loglevel=3' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'systemd.show_status=false' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o '\bsplash\b' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'plymouth.ignore-serial-consoles' <<< "$standard_cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -c '^disable_splash=1$' <<< "$standard_config")" -eq 1 ]]
[[ "$SYSTEMCTL_CALLS" == *"disable getty@tty1.service"* ]]
[[ -f "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi.plymouth" ]]
[[ -f "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi-splash.script" ]]
[[ -f "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi-splash-logo.png" ]]
[[ -f "$MASJIDPI_PLYMOUTH_QUIT_DROPIN_DIR/masjidpi.conf" ]]
! grep -q 'Rotate' "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi-splash.script"

# Add the exact Waveshare USB identity plus connected 1024x600 HDMI mode used by
# the appliance profile. Re-running configuration must switch Plymouth to the
# validated portrait rotation without duplicating boot parameters.
usb="$MASJIDPI_USB_SYSFS_ROOT/1-1.3"
mkdir -p "$usb"
printf '0eef\n' > "$usb/idVendor"
printf '0005\n' > "$usb/idProduct"
printf 'WaveShare\n' > "$usb/manufacturer"
printf 'WS170120\n' > "$usb/product"

hdmi="$MASJIDPI_DRM_SYSFS_ROOT/card0-HDMI-A-1"
mkdir -p "$hdmi"
printf 'connected\n' > "$hdmi/status"
printf '1024x600\n1920x1080\n' > "$hdmi/modes"

configure_quiet_boot
configure_boot_splash

cmdline="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/cmdline.txt")"
config="$(cat "$MASJIDPI_BOOT_FIRMWARE_DIR/config.txt")"

[[ "$(grep -o '\bquiet\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'loglevel=3' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'systemd.show_status=false' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o '\bsplash\b' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -o 'plymouth.ignore-serial-consoles' <<< "$cmdline" | wc -l)" -eq 1 ]]
[[ "$(grep -c '^disable_splash=1$' <<< "$config")" -eq 1 ]]
grep -q 'Rotate(-Math.Pi / 2)' "$MASJIDPI_PLYMOUTH_THEME_DIR/masjidpi-splash.script"
[[ "$(grep -c -- '^-R masjidpi$' "$TMP/plymouth-calls")" -eq 3 ]]

printf '[PASS] quiet boot and branded splash cover standard and appliance Raspberry Pi Board profiles\n'
