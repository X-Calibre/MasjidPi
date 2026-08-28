#!/usr/bin/env bash

set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TEST_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEST_ROOT"' EXIT

mkdir -p "$TEST_ROOT/bin"
PATH="$TEST_ROOT/bin:$PATH"

cat > "$TEST_ROOT/bin/ip" <<'EOF'
#!/usr/bin/env bash
[[ -n "${TEST_ROUTE_ADDRESS:-}" ]] || exit 1
printf '1.1.1.1 via 10.0.0.1 dev wlan0 src %s uid 0\n' "$TEST_ROUTE_ADDRESS"
EOF
cat > "$TEST_ROOT/bin/hostname" <<'EOF'
#!/usr/bin/env bash
case "${1:-}" in
    --fqdn) printf '%s\n' "${TEST_FQDN:-}" ;;
    --short) printf '%s\n' "${TEST_SHORT_NAME:-}" ;;
    -I) printf '%s\n' "${TEST_HOST_ADDRESSES:-}" ;;
esac
EOF
cat > "$TEST_ROOT/bin/systemctl" <<'EOF'
#!/usr/bin/env bash
[[ "${TEST_AVAHI_ACTIVE:-false}" == true ]]
EOF
chmod +x "$TEST_ROOT/bin/ip" "$TEST_ROOT/bin/hostname" "$TEST_ROOT/bin/systemctl"

source "$ROOT/scripts/common.sh"
get_version() { printf 'v1.5.0-rc.3\n'; }

TEST_ROUTE_ADDRESS=10.78.63.4 \
TEST_FQDN=masjidpi.example.test \
TEST_SHORT_NAME=masjidpi \
TEST_AVAHI_ACTIVE=false \
    print_summary > "$TEST_ROOT/configured-summary"
grep -Fq 'IP address:  http://10.78.63.4:8080' "$TEST_ROOT/configured-summary"
grep -Fq 'Hostname:    http://masjidpi.example.test:8080' "$TEST_ROOT/configured-summary"

TEST_ROUTE_ADDRESS=10.78.63.5 \
TEST_FQDN=MasjidPi-Test \
TEST_SHORT_NAME=MasjidPi-Test \
TEST_AVAHI_ACTIVE=true \
    print_summary > "$TEST_ROOT/mdns-summary"
grep -Fq 'Hostname:    http://MasjidPi-Test.local:8080' "$TEST_ROOT/mdns-summary"

TEST_ROUTE_ADDRESS= \
TEST_HOST_ADDRESSES= \
TEST_FQDN=localhost \
TEST_SHORT_NAME=localhost \
TEST_AVAHI_ACTIVE=false \
    print_summary > "$TEST_ROOT/offline-summary"
grep -Fq 'http://localhost:8080' "$TEST_ROOT/offline-summary"

printf 'installer summary tests passed\n'
