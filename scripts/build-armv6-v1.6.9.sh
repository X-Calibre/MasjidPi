#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BASE_COMMIT=22bc292561ed42b85614108ef15c57f0710509a0
SUPPORT_COMMIT=d778d1998fe3eb215a30bb0fe4f2a885cbad1569
SUPPORT_TREE=11ba39a5de45a0c83122c5989bc75cb7a6136dd4
VERSION=v1.6.9
OUTPUT="${1:-$ROOT/dist-armv6-v1.6.9}"
mkdir -p "$OUTPUT"
OUTPUT="$(cd "$OUTPUT" && pwd)"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

# Keep the existing release tag fixed. Add only the reviewed ARMv6 support
# files, and record both immutable source revisions in the package.
git -C "$ROOT" archive "$BASE_COMMIT" | tar -x -C "$WORK"
support_files=(
    backend/internal/player/mpv.go
    backend/internal/player/startup_timeout_test.go
    backend/internal/app/app.go
    backend/internal/app/update_checks.go
    backend/internal/app/armv6_updates_test.go
    scripts/armv6-listen.sh
    scripts/install-listen-armv6.sh
    scripts/install-listen-armv6-latest.sh
    scripts/components.sh
    scripts/github.sh
    scripts/go.sh
    scripts/install.sh
    scripts/packages.sh
    scripts/selftest.sh
    tests/installer_armv6_test.sh
)
git -C "$ROOT" archive "$SUPPORT_TREE" "${support_files[@]}" | tar -x -C "$WORK"

bash "$WORK/tests/installer_armv6_test.sh"
for path in "$WORK"/scripts/*.sh; do bash -n "$path"; done
(
    cd "$WORK/backend"
    go test ./internal/player ./internal/app ./internal/api ./internal/updates
)

NAME="masjidframe-${VERSION}-linux-armv6"
PACKAGE="$OUTPUT/$NAME"
[[ ! -e "$PACKAGE" && ! -e "$OUTPUT/$NAME.tar.gz" ]] \
    || { printf 'Output package already exists; choose a fresh output directory.\n' >&2; exit 1; }
mkdir -p "$PACKAGE/scripts"
(
    cd "$WORK/backend"
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -trimpath \
        -ldflags "-X github.com/X-Calibre/MasjidFrame/backend/internal/version.Version=$VERSION" \
        -o "$PACKAGE/masjidframe" ./cmd/masjidframe
)
go version -m "$PACKAGE/masjidframe" | grep -q 'GOARM=6'
cp "$WORK/backend/configs/default.yaml" "$PACKAGE/default.yaml"
cp "$WORK/backend/data/catalogue.json" "$PACKAGE/catalogue.json"
cp -R "$WORK/frontend" "$PACKAGE/frontend"
cp "$WORK"/scripts/*.sh "$WORK"/scripts/*.service \
    "$WORK"/scripts/*.plymouth "$WORK"/scripts/*.script \
    "$WORK/scripts/99-masjidframe-boot-firmware" "$PACKAGE/scripts/"
chmod +x "$PACKAGE/scripts/install-listen-armv6.sh"
printf '%s\n' "$VERSION" > "$PACKAGE/VERSION"
cat > "$PACKAGE/BUILD_INFO.txt" <<EOF
MasjidFrame v1.6.9 ARMv6 Listen-only supplement, revision 2
Base release source: $BASE_COMMIT
ARMv6 installer and updater safeguards: $SUPPORT_COMMIT (PR #123)
Build target: CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6
Toolchain: $(go version)
Physical Pi Zero playback and soak validation: pending
Install: sudo bash scripts/install-listen-armv6.sh
Pi 3 appliance updates are disabled; rerun the dedicated installer to update.
ARMv6 MPV startup allowance: 30 seconds; HTTP startup allowance: 60 seconds.
The v1.6.9 tag and existing ARM64/AMD64 packages are unchanged.
EOF
source_epoch="$(git -C "$ROOT" show -s --format=%ct "$BASE_COMMIT")"
tar --sort=name --mtime="@$source_epoch" --owner=0 --group=0 --numeric-owner \
    -C "$OUTPUT" -cf - "$NAME" | gzip -n > "$OUTPUT/$NAME.tar.gz"
(
    cd "$OUTPUT"
    sha256sum "$NAME.tar.gz" > "$NAME.tar.gz.sha256"
    sha256sum -c "$NAME.tar.gz.sha256"
)
printf 'Built: %s/%s.tar.gz\n' "$OUTPUT" "$NAME"
