#!/bin/bash

set -euo pipefail

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
readonly output_dir="$PWD/work/update-bundles"
readonly bundle_name="masjidpi-update-${version}-pi3.tar.zst"
readonly bundle="$output_dir/$bundle_name"
readonly signature="$bundle.minisig"
readonly source_commit=$(git rev-parse HEAD)
readonly source_date_epoch=$(
    git show -s --format=%ct "$source_commit"
)
readonly created_at=$(
    date -u \
        --date="@$source_date_epoch" \
        '+%Y-%m-%dT%H:%M:%SZ'
)

for command_name in \
    date \
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
    "$private_key"
do
    if [[ ! -f "$required_file" ]]; then
        echo "Required input not found: $required_file" >&2
        exit 1
    fi
done

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
