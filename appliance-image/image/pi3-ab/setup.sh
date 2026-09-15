#!/bin/bash

set -eu

component=$1

case "$component" in
   SYSTEM)
      cat > "$IMAGEMOUNTPATH/etc/fstab" <<'EOF_FSTAB'
/dev/mmcblk0p1  /boot/firmware  vfat  defaults,rw,noatime,errors=remount-ro     0  2
/dev/mmcblk0p4  /persistent     ext4  defaults,rw,noatime                       0  2
/persistent/iwd /var/lib/iwd    none  bind,x-systemd.requires-mounts-for=/persistent 0 0
/persistent/masjidpi/etc /etc/masjidpi none bind,x-systemd.requires-mounts-for=/persistent 0 0
/persistent/masjidpi/var /var/lib/masjidpi none bind,x-systemd.requires-mounts-for=/persistent 0 0
EOF_FSTAB
      ;;
   BOOT)
      sed -i 's|root=[^ ]*|root=/dev/mmcblk0p2|' "$IMAGEMOUNTPATH/cmdline.txt"
      ;;
   *)
      echo "Unknown image component: $component" >&2
      exit 1
      ;;
esac
