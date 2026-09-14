#!/bin/bash

set -eu

filesystem=$1
genimage_input=$2

test -d "$filesystem"

mkenvimage="$LAYER_DIR/../../u-boot/build/tools/mkenvimage"
environment_file="$LAYER_DIR/../../u-boot/masjidpi.env"
primary_environment="$genimage_input/uboot-env-primary.bin"
redundant_environment="$genimage_input/uboot-env-redundant.bin"

test -x "$mkenvimage"
test -f "$environment_file"

# Start with one valid redundant-format environment and one erased copy.
# The next environment update writes the alternate copy atomically.
"$mkenvimage" \
   -r \
   -s 0x00004000 \
   -o "$primary_environment" \
   "$environment_file"

dd \
   if=/dev/zero \
   of="$redundant_environment" \
   bs=16384 \
   count=1 \
   status=none

for required_boot_file in \
   u-boot.bin \
   kernel8.img \
   kernel8-uboot.img \
   initramfs8 \
   bcm2710-rpi-3-b.dtb \
   extlinux/system_a.conf \
   extlinux/system_b.conf \
   extlinux/extlinux.conf; do
   if [[ ! -f "$filesystem/boot/firmware/$required_boot_file" ]]; then
      echo "Missing U-Boot input: /boot/firmware/$required_boot_file" >&2
      exit 1
   fi
done

boot_dir="$filesystem/boot/firmware"

# Keep the boot artefacts for each system slot independent. An updater can
# replace the inactive set without changing the confirmed slot's boot files.
for slot in a b; do
   slot_dir="$boot_dir/slots/$slot"

   install -d -m 0755 "$slot_dir"
   install -m 0644 \
      "$boot_dir/kernel8-uboot.img" \
      "$slot_dir/kernel8-uboot.img"
   install -m 0644 \
      "$boot_dir/initramfs8" \
      "$slot_dir/initramfs8"
   install -m 0644 \
      "$boot_dir/bcm2710-rpi-3-b.dtb" \
      "$slot_dir/bcm2710-rpi-3-b.dtb"
done

rm -f \
   "$boot_dir/kernel8-uboot.img" \
   "$boot_dir/initramfs8"

for required_slot_file in \
   slots/a/kernel8-uboot.img \
   slots/a/initramfs8 \
   slots/a/bcm2710-rpi-3-b.dtb \
   slots/b/kernel8-uboot.img \
   slots/b/initramfs8 \
   slots/b/bcm2710-rpi-3-b.dtb; do
   if [[ ! -f "$boot_dir/$required_slot_file" ]]; then
      echo "Missing slot boot file: /boot/firmware/$required_slot_file" >&2
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
