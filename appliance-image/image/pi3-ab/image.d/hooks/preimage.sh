#!/bin/bash

set -eu

filesystem=$1
genimage_input=$2

test -d "$filesystem"

# shellcheck disable=SC1090
source "${IGconf_image_outputdir}/img_uuids"

system_a_extraargs="-U $SYSTEM_A_UUID ${IGconf_fs_ext4_mkfs_args:-}"
system_b_extraargs="-U $SYSTEM_B_UUID ${IGconf_fs_ext4_mkfs_args:-}"
data_extraargs="-U $DATA_UUID ${IGconf_fs_ext4_mkfs_args:-}"
vfat_extraargs="-S $IGconf_device_sector_size -i $BOOT_LABEL ${IGconf_fs_vfat_mkfs_args:-}"

sed \
   -e "s|<IMAGE_NAME>|$IGconf_image_name|g" \
   -e "s|<IMAGE_SUFFIX>|$IGconf_image_suffix|g" \
   -e "s|<BOOT_SIZE>|$IGconf_image_boot_part_size|g" \
   -e "s|<SYSTEM_SIZE>|$IGconf_image_system_part_size|g" \
   -e "s|<DATA_SIZE>|$IGconf_image_data_part_size|g" \
   -e "s|<SETUP>|'$(readlink -ef "$LAYER_DIR/setup.sh")'|g" \
   -e "s|<MKE2FS_CONF>|'$(readlink -ef "$LAYER_DIR/mke2fs.conf")'|g" \
   -e "s|<SYSTEM_A_EXTRAARGS>|$system_a_extraargs|g" \
   -e "s|<SYSTEM_B_EXTRAARGS>|$system_b_extraargs|g" \
   -e "s|<DATA_EXTRAARGS>|$data_extraargs|g" \
   -e "s|<VFAT_EXTRAARGS>|$vfat_extraargs|g" \
   -e "s|<DISK_SIGNATURE>|$IGconf_image_disksig|g" \
   "$LAYER_DIR/genimage.cfg.in" > "$genimage_input/genimage.cfg"
