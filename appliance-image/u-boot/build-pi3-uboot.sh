#!/bin/bash

set -euo pipefail

readonly expected_version="2025.01-3+deb13u1"
readonly expected_environment_sha256="934ba277e205dade87d445c33e92c93c15d567bd84974c80f4b91274aec8eadf"

readonly script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
readonly source_dir=${1:-"$HOME/masjidframe-u-boot-source/u-boot-2025.01"}
readonly build_dir="$script_dir/build"
readonly environment_file="$script_dir/masjidframe.env"
readonly config_fragment="$script_dir/pi3-ab.config"
readonly environment_relative=$(realpath --relative-to="$source_dir" "$environment_file")

for command_name in \
    aarch64-linux-gnu-gcc \
    date \
    dpkg-parsechangelog \
    realpath \
    make \
    sha256sum \
    strings
do
    if ! command -v "$command_name" >/dev/null; then
        echo "Required command not found: $command_name" >&2
        exit 1
    fi
done

for required_file in \
    "$source_dir/Makefile" \
    "$source_dir/configs/rpi_3_defconfig" \
    "$source_dir/scripts/config" \
    "$source_dir/scripts/kconfig/merge_config.sh" \
    "$source_dir/debian/changelog" \
    "$environment_file" \
    "$config_fragment"
do
    if [[ ! -f "$required_file" ]]; then
        echo "Required file not found: $required_file" >&2
        exit 1
    fi
done

actual_version=$(
    cd "$source_dir"
    dpkg-parsechangelog -S Version
)

if [[ "$actual_version" != "$expected_version" ]]; then
    echo "Expected U-Boot source version: $expected_version" >&2
    echo "Found U-Boot source version:    $actual_version" >&2
    exit 1
fi

actual_environment_sha256=$(
    sha256sum "$environment_file" |
    awk '{print $1}'
)

if [[ "$actual_environment_sha256" != "$expected_environment_sha256" ]]; then
    echo "The MasjidFrame U-Boot environment differs from the hardware-tested input." >&2
    echo "Expected: $expected_environment_sha256" >&2
    echo "Found:    $actual_environment_sha256" >&2
    exit 1
fi

source_date=$(
    cd "$source_dir"
    dpkg-parsechangelog -S Date
)
export SOURCE_DATE_EPOCH
SOURCE_DATE_EPOCH=$(date -d "$source_date" +%s)

mkdir -p "$build_dir"

make \
    -C "$source_dir" \
    O="$build_dir" \
    ARCH=arm \
    CROSS_COMPILE=aarch64-linux-gnu- \
    mrproper

make \
    -C "$source_dir" \
    O="$build_dir" \
    ARCH=arm \
    CROSS_COMPILE=aarch64-linux-gnu- \
    rpi_3_defconfig

(
    cd "$source_dir"

    KCONFIG_CONFIG="$build_dir/.config" \
        scripts/kconfig/merge_config.sh \
        -m \
        -O "$build_dir" \
        "$build_dir/.config" \
        "$config_fragment"
)

"$source_dir/scripts/config" \
    --file "$build_dir/.config" \
    --set-str DEFAULT_ENV_FILE "$environment_relative"

make \
    -C "$source_dir" \
    O="$build_dir" \
    ARCH=arm \
    CROSS_COMPILE=aarch64-linux-gnu- \
    olddefconfig

make \
    -C "$source_dir" \
    O="$build_dir" \
    ARCH=arm \
    CROSS_COMPILE=aarch64-linux-gnu- \
    -j"$(nproc)"

readonly binary="$build_dir/u-boot.bin"

if [[ ! -s "$binary" ]]; then
    echo "Custom U-Boot binary was not produced: $binary" >&2
    exit 1
fi

for required_variable in \
    'active_slot=a' \
    'rollback_slot=a' \
    'upgrade_available=0' \
    'bootcount=0' \
    'bootcmd=run boot_selected_slot'
do
    if ! strings "$binary" | grep -Fx "$required_variable" >/dev/null; then
        echo "Compiled environment variable missing: $required_variable" >&2
        exit 1
    fi
done

echo
echo "MasjidFrame Pi 3 U-Boot build complete"
echo "Version: $actual_version"
echo "Source date epoch: $SOURCE_DATE_EPOCH"
echo "Binary: $binary"
sha256sum "$binary"
