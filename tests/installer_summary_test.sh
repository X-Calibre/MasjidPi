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
    -I) printf '%s\n' "${TEST_HOST_ADDRESSES:-}" ;;
esac
EOF
cat > "$TEST_ROOT/bin/curl" <<'EOF'
#!/usr/bin/env bash
[[ "${TEST_DEVICE_ACCESS_AVAILABLE:-true}" == true ]] || exit 22
printf '{"ip_address":"%s","fqdn":"%s"}\n' \
    "${TEST_ROUTE_ADDRESS:-}" "${TEST_FQDN:-}"
EOF
chmod +x "$TEST_ROOT/bin/ip" "$TEST_ROOT/bin/hostname" "$TEST_ROOT/bin/curl"

source "$ROOT/scripts/common.sh"
get_version() { printf 'v1.5.0-rc.3\n'; }

TEST_ROUTE_ADDRESS=10.78.63.4 \
TEST_FQDN=masjidpi.example.test \
    print_summary > "$TEST_ROOT/configured-summary"
grep -Fq 'IP address:  http://10.78.63.4:8080' "$TEST_ROOT/configured-summary"
grep -Fq 'Hostname:    http://masjidpi.example.test:8080' "$TEST_ROOT/configured-summary"

TEST_ROUTE_ADDRESS=10.78.63.5 \
TEST_FQDN='' \
    print_summary > "$TEST_ROOT/ip-only-summary"
grep -Fq 'IP address:  http://10.78.63.5:8080' "$TEST_ROOT/ip-only-summary"
if grep -Fq 'Hostname:' "$TEST_ROOT/ip-only-summary"; then
    printf 'installer summary invented a hostname without a network-issued FQDN\n' >&2
    exit 1
fi

TEST_ROUTE_ADDRESS=10.78.63.6 \
TEST_FQDN=MasjidPi-Test.local \
    print_summary > "$TEST_ROOT/network-local-summary"
grep -Fq 'Hostname:    http://MasjidPi-Test.local:8080' "$TEST_ROOT/network-local-summary"

TEST_ROUTE_ADDRESS='' \
TEST_HOST_ADDRESSES='' \
TEST_DEVICE_ACCESS_AVAILABLE=false \
    print_summary > "$TEST_ROOT/offline-summary"
grep -Fq 'http://localhost:8080' "$TEST_ROOT/offline-summary"

printf 'installer summary tests passed\n'
