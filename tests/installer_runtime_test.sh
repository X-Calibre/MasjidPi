#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

info() { :; }
success() { :; }

# shellcheck source=../scripts/runtime.sh
source "$ROOT/scripts/runtime.sh"

config="$TMP/config.yaml"

cat > "$config" <<'EOF'
player:
  socket: "/tmp/masjidframe.sock"
EOF

migrate_runtime_socket_path "$config"
grep -qx '  socket: "/run/masjidframe/mpv.sock"' "$config"

cat > "$config" <<'EOF'
player:
  socket: "/run/masjidpi/mpv.sock"
EOF

migrate_runtime_socket_path "$config"
grep -qx '  socket: "/run/masjidframe/mpv.sock"' "$config"

cat > "$config" <<'EOF'
player:
  socket: "/srv/custom/mpv.sock"
EOF

migrate_runtime_socket_path "$config"
grep -qx '  socket: "/srv/custom/mpv.sock"' "$config"

printf '[PASS] runtime socket migration moves only previous packaged defaults\n'
