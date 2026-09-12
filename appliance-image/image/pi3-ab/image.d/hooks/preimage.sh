#!/bin/bash

set -eu

filesystem=$1
genimage_input=$2

test -d "$filesystem"

for required_boot_file in \
   u-boot.bin \
   kernel8.img \
   kernel8-uboot.img \
   initramfs8 \
   bcm2710-rpi-3-b.dtb \
   extlinux/extlinux.conf; do
   if [[ ! -f "$filesystem/boot/firmware/$required_boot_file" ]]; then
      echo "Missing U-Boot input: /boot/firmware/$required_boot_file" >&2
      exit 1
   fi
done

# Later built-in customization hooks may alter config.txt after the U-Boot
# files are installed. Normalize the final firmware settings immediately
# before genimage copies the boot filesystem.
boot_config="$filesystem/boot/firmware/config.txt"
sed -i \
   -e '/^auto_initramfs=/d' \
   -e '/^arm_64bit=/d' \
   -e '/^kernel=/d' \
   -e '/^enable_uart=/d' \
   "$boot_config"

cat >> "$boot_config" <<'EOF'

# MasjidPi Pi 3 A/B boot prototype
auto_initramfs=0
arm_64bit=1
kernel=u-boot.bin
enable_uart=1
EOF

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
