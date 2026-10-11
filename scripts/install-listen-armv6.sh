#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
export MASJIDFRAME_ARMV6_LISTEN=1
export GOARM=6
export CGO_ENABLED=0
# Keep source compilation within the original Zero's limited memory.
export GOMAXPROCS=1
export GOFLAGS="${GOFLAGS:-} -p=1"

# Run this bundled/local entry point for release or --source installations.
exec bash "$SCRIPT_DIR/install.sh" "$@"
