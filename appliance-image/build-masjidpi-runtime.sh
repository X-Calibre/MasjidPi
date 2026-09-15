#!/bin/bash

set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
project_root=$(cd "$script_dir/.." && pwd)
output_dir=${1:-"$project_root/work/masjidpi-runtime-arm64"}

for command_name in \
    file \
    go \
    jq \
    qemu-aarch64
do
    if ! command -v "$command_name" >/dev/null; then
        echo "Required command not found: $command_name" >&2
        exit 1
    fi
done

source_version=$(
    jq -er '
        .version |
        select(
            type == "string" and
            test("^v[0-9]+\\.[0-9]+\\.[0-9]+(-rc\\.[0-9]+)?$")
        )
    ' "$project_root/version.json"
)

runtime_version=${MASJIDPI_IMAGE_VERSION:-"${source_version}-image"}
staging_dir=$(mktemp -d "$project_root/work/masjidpi-runtime.XXXXXX")

cleanup()
{
    rm -rf "$staging_dir"
}
trap cleanup EXIT

mkdir -p \
    "$staging_dir/bin" \
    "$staging_dir/scripts"

echo "Building MasjidPi ARM64 runtime..."

(
    cd "$project_root/backend"

    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64 \
    go build \
        -trimpath \
        -ldflags \
        "-X github.com/X-Calibre/MasjidPi/backend/internal/version.Version=$runtime_version" \
        -o "$staging_dir/bin/masjidpi" \
        ./cmd/masjidpi
)

cp -R \
    "$project_root/frontend" \
    "$staging_dir/frontend"

install -m 0644 \
    "$project_root/backend/configs/default.yaml" \
    "$staging_dir/default.yaml"

install -m 0644 \
    "$project_root/backend/data/catalogue.json" \
    "$staging_dir/catalogue.json"

printf '%s\n' "$runtime_version" > "$staging_dir/VERSION"

for runtime_script in \
    masjidboard-display.sh \
    masjidboard-warmup.sh
do
    install -m 0755 \
        "$project_root/scripts/$runtime_script" \
        "$staging_dir/scripts/$runtime_script"
done

for service_file in \
    masjidpi.service \
    masjidpi-display.service \
    masjidpi-display-warmup.service
do
    install -m 0644 \
        "$project_root/scripts/$service_file" \
        "$staging_dir/scripts/$service_file"
done

if ! file "$staging_dir/bin/masjidpi" |
    grep -Fq 'ARM aarch64'; then
    echo "Built backend is not an ARM64 executable." >&2
    exit 1
fi

if ! file "$staging_dir/bin/masjidpi" |
    grep -Fq 'statically linked'; then
    echo "Built backend is not statically linked." >&2
    exit 1
fi

reported_version=$(
    qemu-aarch64 "$staging_dir/bin/masjidpi" --version
)

if [[ "$reported_version" != "$runtime_version" ]]; then
    echo "Built backend reported an unexpected version." >&2
    echo "Expected: $runtime_version" >&2
    echo "Found:    $reported_version" >&2
    exit 1
fi

for required_file in \
    bin/masjidpi \
    frontend/index.html \
    frontend/masjidboard.html \
    frontend/masjidboard-startup.html \
    default.yaml \
    catalogue.json \
    VERSION \
    scripts/masjidboard-display.sh \
    scripts/masjidboard-warmup.sh \
    scripts/masjidpi.service \
    scripts/masjidpi-display.service \
    scripts/masjidpi-display-warmup.service
do
    if [[ ! -f "$staging_dir/$required_file" ]]; then
        echo "Missing runtime file: $required_file" >&2
        exit 1
    fi
done

case "$output_dir" in
    "$project_root"/work/*)
        ;;
    *)
        echo "Refusing output directory outside project work/: $output_dir" >&2
        exit 1
        ;;
esac

rm -rf "$output_dir"
mv "$staging_dir" "$output_dir"
trap - EXIT

echo
echo "MasjidPi ARM64 runtime complete"
echo "Version: $runtime_version"
echo "Output:  $output_dir"
du -sh "$output_dir"
