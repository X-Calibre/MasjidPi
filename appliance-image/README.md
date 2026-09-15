# MasjidPi Raspberry Pi 3 A/B appliance image

This directory builds a Raspberry Pi 3 Model B appliance image using
[`rpi-image-gen`](https://github.com/raspberrypi/rpi-image-gen).

The image provides two bootable system slots, persistent shared storage,
redundant U-Boot environment storage, boot-attempt counting, automatic health
confirmation, shared Wi-Fi profiles, and rollback to the last confirmed slot.

## Disk layout

| Region | Location | Size | Purpose |
|---|---:|---:|---|
| Primary U-Boot environment | 1 MiB offset | 16 KiB | Active redundant environment record |
| Redundant U-Boot environment | 4 MiB offset | 16 KiB | Alternate environment record |
| `BOOT` | Partition 1 | 256 MiB | Firmware, U-Boot, kernel, initramfs and extlinux configuration |
| `SYSTEM_A` | Partition 2 | 4 GiB | System slot A |
| `SYSTEM_B` | Partition 3 | 4 GiB | System slot B |
| `PERSISTENT` | Partition 4 | 5 GiB | State shared between system slots |

The two environment records are stored in the unpartitioned 8 MiB gap before
`BOOT`. The four visible partitions use an MBR partition table. The resulting
image is 14,235,467,776 bytes (approximately 13.26 GiB), requires a nominal
16 GB or larger microSD card, and leaves approximately 1.5 GiB unused on the
validated 14.8 GiB card.

Both system partitions initially contain the same slot-neutral filesystem.
U-Boot selects the root partition using separate
`extlinux/system_a.conf` and `extlinux/system_b.conf` entries.

## Shared networking state

The appliance uses NetworkManager with `wpa_supplicant`, matching the networking
interface used by current Raspberry Pi OS. Both slots share NetworkManager
connection profiles and runtime state through these bind mounts:

- `/persistent/network-manager/system-connections` at
  `/etc/NetworkManager/system-connections`
- `/persistent/network-manager/state` at `/var/lib/NetworkManager`

The mounts explicitly depend on `/persistent`. Connection-profile files use
mode `0600`, and their parent directories use mode `0700`.

The tracked build configuration contains no Wi-Fi credentials or network names.
A developer may use an ignored local configuration and `nm.cmds` file to
preconfigure a test network.

During first-run setup, the selected IANA time zone is applied to the system and
saved as `/var/lib/masjidpi/timezone.json`. Because `/var/lib/masjidpi` is
shared, either system slot restores the same time zone at application startup.

## U-Boot

The image uses a reproducible custom build of Debian U-Boot source version:

```text
2025.01-3+deb13u1
```

The custom configuration provides:

- Redundant raw MMC environment storage.
- System-slot selection through `active_slot`.
- A remembered confirmed slot through `rollback_slot`.
- Persistent boot counting during trial boots.
- A two-attempt trial boot limit.
- Automatic rollback after the trial limit is exceeded.
- An uncompressed `kernel8-uboot.img` for ARM64 U-Boot booting.

A stable confirmed system has no `bootlimit`. An updater must create
`bootlimit=2` only when arming a trial boot. Successful confirmation and
automatic rollback both remove it.

The U-Boot source tree is not modified. MasjidPi's configuration and compiled
default environment are stored under `appliance-image/u-boot/`.

## A/B management command

The image installs `/usr/local/sbin/masjidpi-ab` in both system slots.

Display the running slot and boot state:

```bash
sudo masjidpi-ab status
```

Collect detailed boot diagnostics:

```bash
sudo masjidpi-ab diagnostics
```

Arm the other slot for a trial boot:

```bash
sudo masjidpi-ab trial b
sudo reboot
```

A systemd one-shot service automatically evaluates a pending trial after a
60-second settling period. It confirms the running slot only when:

- The running root partition agrees with `active_slot`.
- The root filesystem is mounted read-write.
- `PERSISTENT` is mounted read-write from `/dev/mmcblk0p4`.
- No systemd service has failed.

Successful automatic confirmation makes the running slot the new rollback
target and clears `upgrade_available`, `bootcount`, and `bootlimit`.
Ordinary stable boots exit immediately without changing the environment.

Manual confirmation remains available for diagnostics and recovery:

```bash
sudo masjidpi-ab confirm
```

These commands and the confirmation service provide the boot-control foundation
for an updater. Downloading, installing and verifying update bundles remain
separate work.

## Reproducible inputs

The image build is pinned to `rpi-image-gen` commit:

```text
d1021e82dd578b588cc3b4d45cd7b4b86e57b796
```

The build rejects another revision unless
`ALLOW_UNPINNED_RPI_IMAGE_GEN=1` is explicitly set.

The U-Boot helper verifies the Debian source version and SHA-256 checksum of
the hardware-tested default environment before building.

## Build

The default locations are:

- `rpi-image-gen`: `$HOME/rpi-image-gen`
- U-Boot source: `$HOME/masjidpi-u-boot-source/u-boot-2025.01`

From the repository root:

```bash
./appliance-image/build-pi3-ab-prototype.sh
```

Explicit paths can be supplied as follows:

```bash
./appliance-image/build-pi3-ab-prototype.sh \
    "$HOME/rpi-image-gen" \
    "$PWD/appliance-image/config/pi3-ab-prototype.yaml" \
    "$HOME/masjidpi-u-boot-source/u-boot-2025.01"
```

For a hardware test requiring a login, create the ignored file
`appliance-image/config/pi3-ab-local.yaml`:

```yaml
include:
  file: pi3-ab-prototype.yaml

device:
  user1pass: "choose-a-valid-temporary-password"

nm:
  cmds: ${@SRCROOT}/config/local/network-manager.cmds
```

The optional `nm.cmds` file contains one `nmcli --offline` command per line.
Keep that file under `appliance-image/config/local/`; the directory is ignored
by Git because connection commands can contain Wi-Fi credentials. Builds
without `nm.cmds` start with no saved Wi-Fi network and use MasjidPi's first-run
network setup.

Then build with that configuration:

```bash
./appliance-image/build-pi3-ab-prototype.sh \
    "$HOME/rpi-image-gen" \
    "$PWD/appliance-image/config/pi3-ab-local.yaml" \
    "$HOME/masjidpi-u-boot-source/u-boot-2025.01"
```

Never commit a test password.

Generated images, compressed deployment artefacts, SBOM data and IDP metadata
are written below `work/`.

## Hardware validation

The integrated image has been tested on a Raspberry Pi 3 Model B Rev 1.2.

Validated behaviour:

1. An untouched image boots system A from `/dev/mmcblk0p2`.
2. A healthy trial system boots from the alternate root partition.
3. The confirmation service waits 60 seconds and automatically confirms the
   healthy running slot.
4. Stable boots leave the confirmed environment unchanged.
5. A deliberately failed systemd service prevents automatic confirmation.
6. The unhealthy candidate receives exactly two permitted boot attempts.
7. U-Boot then automatically restores the confirmed rollback slot.
8. Rollback clears `upgrade_available`, `bootcount`, and `bootlimit`.
9. Both system slots mount `PERSISTENT` at `/persistent`.
10. No power-throttling flags were observed during successful boots.

The NetworkManager-based replacement for the earlier direct-IWD configuration
has passed layer-resolution, filesystem-image, and Raspberry Pi hardware
validation. Automatic Wi-Fi reconnection retained the same DHCP address across
both slots, and NetworkManager state, MasjidBoard selection, audio selection,
and display operation persisted after switching to and confirming system B.
A temporary ARM64 runtime deployed to system B also validated the revised
touchscreen first-run flow: an unconfigured MasjidBoard opened location setup
instead of a blank display, South Africa automatically selected
`Africa/Johannesburg`, the time zone was applied and persisted, the board clock
updated correctly, and the original selection was restored successfully.
The clean expanded-layout image subsequently validated the complete first-run
flow and automatic time-zone selection. After switching to system B, Wi-Fi,
MasjidBoard selection, time zone, touch operation, and USB audio persisted;
system B automatically confirmed and restored `Africa/Johannesburg` from shared
state.

A reversible system-B boot test first validated quieter console settings. The
appliance integration was then extended to use MasjidPi's branded Plymouth
theme throughout early Linux startup. Hardware validation confirmed that no
firmware or kernel text was exposed, the logo was correctly scaled and
oriented, the touchscreen login prompt was suppressed, and the WebKit startup
image and Board followed normally. Touch, Wi-Fi, USB audio, and all required
services remained healthy. A filesystem build confirmed that the theme,
script, and logo are embedded in `initramfs8`, all A/B extlinux entries enable
the splash, and serial recovery remains available.
