#!/bin/bash

set -euo pipefail
umask 022

usage()
{
    echo "Usage: $0 VERSION [IMAGE_OUTPUT_DIR] [PRIVATE_KEY]" >&2
    exit 1
}

version=${1:-}
image_output_dir=${2:-"$PWD/work/image-masjidpi-pi3-ab-prototype"}
private_key=${3:-"$PWD/appliance-image/update/masjidpi-update-private.key"}

if [[ -z "$version" ||
      ! "$version" =~ ^[A-Za-z0-9][A-Za-z0-9._+-]*$ ]]; then
    usage
fi

readonly rootfs_image="$image_output_dir/system_a.ext4"
readonly boot_image="$image_output_dir/boot.vfat"
readonly image_metadata="$image_output_dir/image.json"
readonly output_dir="$PWD/work/update-bundles"
readonly bundle_name="masjidpi-update-${version}-pi3.tar.zst"
readonly bundle="$output_dir/$bundle_name"
readonly signature="$bundle.minisig"
readonly current_source_commit=$(git rev-parse HEAD)
readonly current_build_version=$(git describe --tags --always --dirty)
readonly source_date_epoch=$(
    git show -s --format=%ct "$current_source_commit"
)
readonly created_at=$(
    date -u \
        --date="@$source_date_epoch" \
        '+%Y-%m-%dT%H:%M:%SZ'
)

for command_name in \
    date \
    debugfs \
    git \
    jq \
    mcopy \
    minisign \
    sha256sum \
    stat \
    tar \
    zstd
do
    if ! command -v "$command_name" >/dev/null; then
        echo "Required command not found: $command_name" >&2
        exit 1
    fi
done

for required_file in \
    "$rootfs_image" \
    "$boot_image" \
    "$image_metadata" \
    "$private_key"
do
    if [[ ! -f "$required_file" ]]; then
        echo "Required input not found: $required_file" >&2
        exit 1
    fi
done

if [[ -n "$(git status --porcelain)" ]]; then
    echo "Refusing to sign from a dirty source tree." >&2
    git status --short >&2
    exit 1
fi

embedded_release=$(
    debugfs \
        -R 'cat /usr/share/masjidpi/update/release.json' \
        "$rootfs_image" \
        2>/dev/null
)

if ! jq -e '
    .schema_version == 1 and
    .product == "masjidpi" and
    .device_class == "pi3" and
    (.release_version | type == "string") and
    (.build_version | type == "string") and
    (.runtime_version | type == "string") and
    (.source_commit | type == "string" and
        test("^[0-9a-f]{40}$")) and
    (.created_at | type == "string")
' <<<"$embedded_release" >/dev/null; then
    echo "The root filesystem has no valid embedded release record." >&2
    echo "Rebuild the image before signing an update." >&2
    exit 1
fi

embedded_source_commit=$(
    jq -r '.source_commit' <<<"$embedded_release"
)
embedded_build_version=$(
    jq -r '.build_version' <<<"$embedded_release"
)
embedded_runtime_version=$(
    jq -r '.runtime_version' <<<"$embedded_release"
)
image_build_version=$(
    jq -er '.IGmeta.IGconf_image_version' "$image_metadata"
)

if [[ "$embedded_source_commit" != "$current_source_commit" ]]; then
    echo "Root filesystem source commit does not match HEAD." >&2
    echo "Embedded: $embedded_source_commit" >&2
    echo "HEAD:     $current_source_commit" >&2
    exit 1
fi

if [[ "$embedded_build_version" != "$current_build_version" ]]; then
    echo "Root filesystem build version does not match the source tree." >&2
    echo "Embedded: $embedded_build_version" >&2
    echo "Source:   $current_build_version" >&2
    exit 1
fi

if [[ "$image_build_version" != "$embedded_build_version" ]]; then
    echo "Image metadata does not match the root filesystem." >&2
    echo "Image:    $image_build_version" >&2
    echo "Rootfs:   $embedded_build_version" >&2
    exit 1
fi

embedded_version_file=$(
    debugfs \
        -R 'cat /opt/masjidpi/VERSION' \
        "$rootfs_image" \
        2>/dev/null |
    tr -d '\r\n'
)

if [[ "$embedded_version_file" != "$embedded_runtime_version" ]]; then
    echo "Embedded runtime VERSION does not match release provenance." >&2
    echo "VERSION: $embedded_version_file" >&2
    echo "Record:  $embedded_runtime_version" >&2
    exit 1
fi

readonly source_commit="$embedded_source_commit"
readonly build_version="$embedded_build_version"
readonly runtime_version="$embedded_runtime_version"

echo "=== VERIFIED IMAGE PROVENANCE ==="
echo "Build version:   $build_version"
echo "Runtime version: $runtime_version"
echo "Source commit:   $source_commit"

mkdir -p "$output_dir"
staging_dir=$(mktemp -d)

cleanup()
{
    rm -rf "$staging_dir"
}
trap cleanup EXIT

install -d -m 0755 "$staging_dir/boot"

echo "Extracting slot boot artifacts..."

for boot_file in \
    kernel8-uboot.img \
    initramfs8 \
    bcm2710-rpi-3-b.dtb
do
    mcopy \
        -i "$boot_image" \
        "::/slots/a/$boot_file" \
        "$staging_dir/boot/$boot_file"
done

echo "Compressing root filesystem..."

zstd \
    --threads=0 \
    -19 \
    --force \
    "$rootfs_image" \
    -o "$staging_dir/rootfs.ext4.zst"

rootfs_size=$(stat -c '%s' "$rootfs_image")
rootfs_sha256=$(sha256sum "$rootfs_image" | awk '{print $1}')
rootfs_archive_size=$(
    stat -c '%s' "$staging_dir/rootfs.ext4.zst"
)
rootfs_archive_sha256=$(
    sha256sum "$staging_dir/rootfs.ext4.zst" |
    awk '{print $1}'
)

kernel_size=$(
    stat -c '%s' "$staging_dir/boot/kernel8-uboot.img"
)
kernel_sha256=$(
    sha256sum "$staging_dir/boot/kernel8-uboot.img" |
    awk '{print $1}'
)

initramfs_size=$(
    stat -c '%s' "$staging_dir/boot/initramfs8"
)
initramfs_sha256=$(
    sha256sum "$staging_dir/boot/initramfs8" |
    awk '{print $1}'
)

dtb_size=$(
    stat -c '%s' "$staging_dir/boot/bcm2710-rpi-3-b.dtb"
)
dtb_sha256=$(
    sha256sum "$staging_dir/boot/bcm2710-rpi-3-b.dtb" |
    awk '{print $1}'
)

jq -n \
    --arg version "$version" \
    --arg build_version "$build_version" \
    --arg runtime_version "$runtime_version" \
    --arg commit "$source_commit" \
    --arg created "$created_at" \
    --arg rootfs_sha "$rootfs_sha256" \
    --arg rootfs_archive_sha "$rootfs_archive_sha256" \
    --arg kernel_sha "$kernel_sha256" \
    --arg initramfs_sha "$initramfs_sha256" \
    --arg dtb_sha "$dtb_sha256" \
    --argjson rootfs_size "$rootfs_size" \
    --argjson rootfs_archive_size "$rootfs_archive_size" \
    --argjson kernel_size "$kernel_size" \
    --argjson initramfs_size "$initramfs_size" \
    --argjson dtb_size "$dtb_size" \
    '{
        schema_version: 1,
        product: "masjidpi",
        device_class: "pi3",
        release_version: $version,
        build_version: $build_version,
        runtime_version: $runtime_version,
        source_commit: $commit,
        created_at: $created,
        rootfs: {
            file: "rootfs.ext4.zst",
            compression: "zstd",
            uncompressed_size: $rootfs_size,
            uncompressed_sha256: $rootfs_sha,
            archive_size: $rootfs_archive_size,
            archive_sha256: $rootfs_archive_sha,
            source_label: "SYSTEM_A"
        },
        boot: {
            kernel: {
                file: "boot/kernel8-uboot.img",
                size: $kernel_size,
                sha256: $kernel_sha
            },
            initramfs: {
                file: "boot/initramfs8",
                size: $initramfs_size,
                sha256: $initramfs_sha
            },
            dtb: {
                file: "boot/bcm2710-rpi-3-b.dtb",
                size: $dtb_size,
                sha256: $dtb_sha
            }
        }
    }' > "$staging_dir/manifest.json"

rm -f "$bundle" "$signature"

tar \
    --sort=name \
    --mtime="@$source_date_epoch" \
    --owner=0 \
    --group=0 \
    --numeric-owner \
    -C "$staging_dir" \
    -cf - \
    manifest.json \
    rootfs.ext4.zst \
    boot |
zstd \
    --threads=0 \
    -19 \
    --force \
    -o "$bundle"

echo
echo "Signing update bundle..."
minisign \
    -S \
    -s "$private_key" \
    -m "$bundle"

echo
echo "MasjidPi update bundle complete"
echo "Bundle:    $bundle"
echo "Signature: $signature"
sha256sum "$bundle" "$signature"
