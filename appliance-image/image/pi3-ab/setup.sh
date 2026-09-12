#!/bin/bash

set -eu

component=$1

case "$component" in
   SYSTEM)
      cat > "$IMAGEMOUNTPATH/etc/fstab" <<'EOF'
/dev/mmcblk0p2  /               ext4  rw,relatime,errors=remount-ro,commit=30  0  1
/dev/mmcblk0p1  /boot/firmware  vfat  defaults,rw,noatime,errors=remount-ro     0  2
/dev/mmcblk0p4  /persistent     ext4  defaults,rw,noatime                       0  2
EOF
      ;;
   BOOT)
      sed -i 's|root=[^ ]*|root=/dev/mmcblk0p2|' "$IMAGEMOUNTPATH/cmdline.txt"
      ;;
   *)
      echo "Unknown image component: $component" >&2
      exit 1
      ;;
esac
