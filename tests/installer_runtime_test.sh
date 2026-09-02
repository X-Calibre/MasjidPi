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
  socket: "/tmp/masjidpi.sock"
EOF

migrate_runtime_socket_path "$config"
grep -qx '  socket: "/run/masjidpi/mpv.sock"' "$config"

cat > "$config" <<'EOF'
player:
  socket: "/srv/custom/mpv.sock"
EOF

migrate_runtime_socket_path "$config"
grep -qx '  socket: "/srv/custom/mpv.sock"' "$config"

printf '[PASS] runtime socket migration moves only the previous packaged default\n'
