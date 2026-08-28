#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT

export MASJIDPI_DEVICE_MODEL_PATH="$TEST_ROOT/model"
export MASJIDPI_NETWORKMANAGER_CONF_DIR="$TEST_ROOT/conf"
export TEST_CALLS="$TEST_ROOT/calls"

mkdir -p "$TEST_ROOT/bin"
PATH="$TEST_ROOT/bin:$PATH"

info() { :; }
success() { :; }
warn() { printf 'WARN %s\n' "$*" >> "$TEST_CALLS"; }
command_exists() { command -v "$1" >/dev/null 2>&1; }

source "$ROOT/scripts/network.sh"

cat > "$TEST_ROOT/bin/systemctl" <<'EOF'
#!/usr/bin/env bash
[[ "${TEST_NETWORKMANAGER_ACTIVE:-false}" == true ]]
EOF
cat > "$TEST_ROOT/bin/nmcli" <<'EOF'
#!/usr/bin/env bash
printf 'nmcli %s\n' "$*" >> "$TEST_CALLS"
case "$*" in
    '-t -f DEVICE,TYPE device status') printf 'wlan0:wifi\n' ;;
    '-g GENERAL.CON-UUID device show wlan0') printf 'test-uuid\n' ;;
esac
EOF
cat > "$TEST_ROOT/bin/iw" <<'EOF'
#!/usr/bin/env bash
printf 'iw %s\n' "$*" >> "$TEST_CALLS"
EOF
chmod +x "$TEST_ROOT/bin/systemctl" "$TEST_ROOT/bin/nmcli" "$TEST_ROOT/bin/iw"

printf 'Generic Linux computer\0' > "$MASJIDPI_DEVICE_MODEL_PATH"
TEST_NETWORKMANAGER_ACTIVE=true configure_raspberry_pi_wifi
[[ ! -e "$MASJIDPI_NETWORKMANAGER_CONF_DIR/masjidpi-wifi-powersave.conf" ]]
[[ ! -e "$TEST_CALLS" ]]

printf 'Raspberry Pi 4 Model B Rev 1.5\0' > "$MASJIDPI_DEVICE_MODEL_PATH"
TEST_NETWORKMANAGER_ACTIVE=false configure_raspberry_pi_wifi
[[ ! -e "$MASJIDPI_NETWORKMANAGER_CONF_DIR/masjidpi-wifi-powersave.conf" ]]

TEST_NETWORKMANAGER_ACTIVE=true configure_raspberry_pi_wifi
config="$MASJIDPI_NETWORKMANAGER_CONF_DIR/masjidpi-wifi-powersave.conf"
[[ "$(cat "$config")" == $'[connection]\nwifi.powersave=2' ]]
grep -qx 'nmcli connection modify test-uuid 802-11-wireless.powersave 2' "$TEST_CALLS"
grep -qx 'iw dev wlan0 set power_save off' "$TEST_CALLS"

mtime="$(stat -c %Y "$config")"
sleep 1
TEST_NETWORKMANAGER_ACTIVE=true configure_raspberry_pi_wifi
[[ "$(stat -c %Y "$config")" == "$mtime" ]]

printf 'installer network tests passed\n'
