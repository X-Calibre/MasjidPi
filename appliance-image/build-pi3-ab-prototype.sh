#!/bin/bash

set -euo pipefail
export PATH="/usr/local/sbin:/usr/sbin:/sbin:$PATH"

readonly expected_revision=d1021e82dd578b588cc3b4d45cd7b4b86e57b796
readonly script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
readonly image_gen_dir=${1:-"$HOME/rpi-image-gen"}
readonly config_file=${2:-"$script_dir/config/pi3-ab-prototype.yaml"}

if [[ ! -x "$image_gen_dir/rpi-image-gen" ]]; then
   echo "rpi-image-gen was not found at: $image_gen_dir" >&2
   echo "Pass its checkout directory as the first argument." >&2
   exit 1
fi

if [[ ! -f "$config_file" ]]; then
   echo "Image configuration was not found at: $config_file" >&2
   exit 1
fi

actual_revision=$(git -C "$image_gen_dir" rev-parse HEAD)
if [[ "$actual_revision" != "$expected_revision" && \
      "${ALLOW_UNPINNED_RPI_IMAGE_GEN:-0}" != 1 ]]; then
   echo "Expected rpi-image-gen revision: $expected_revision" >&2
   echo "Found revision:                  $actual_revision" >&2
   echo "Set ALLOW_UNPINNED_RPI_IMAGE_GEN=1 only for an intentional test." >&2
   exit 1
fi

"$image_gen_dir/rpi-image-gen" build \
   -S "$script_dir" \
   -c "$config_file"
