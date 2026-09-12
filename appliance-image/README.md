# MasjidPi Raspberry Pi 3 appliance image

This directory contains the Raspberry Pi 3 image prototype built with
[`rpi-image-gen`](https://github.com/raspberrypi/rpi-image-gen).

This prototype uses a four-partition MBR layout that fits comfortably on a
nominal 16 GB microSD card. Both system slots have been booted successfully on
a Raspberry Pi 3 Model B through U-Boot. System A remains the default;
automatic slot selection and rollback remain separate milestones.

## Prototype layout

| Partition | Label | Prototype size | Purpose |
|---|---|---:|---|
| 1 | `BOOT` | 256 MiB | Raspberry Pi firmware, U-Boot, kernels, initramfs, and boot configuration |
| 2 | `SYSTEM_A` | 3 GiB | Default operating-system slot |
| 3 | `SYSTEM_B` | 3 GiB | Secondary operating-system slot |
| 4 | `PERSISTENT` | 2 GiB | State shared across system slots |

Both system partitions initially contain the same slot-neutral filesystem
image. The shared `extlinux.conf` contains entries for both slots: system A
uses partition 2 and system B uses partition 3. Debian's packaged U-Boot
bootflow uses the `DEFAULT` entry without presenting an interactive slot menu
on the tested Pi 3, so manual testing selected B by temporarily changing the
default. Do not treat this milestone as a working updater or rollback
implementation.

The image deliberately does not create a `uboot.env` file. The packaged kernel
is gzip-compressed, so the build also creates an uncompressed
`kernel8-uboot.img` for U-Boot to load through `extlinux.conf`. This avoids
depending on U-Boot's single, non-redundant FAT environment. Reliable automatic
slot selection and boot-attempt tracking remain separate work.

## Reproducible input

The prototype was developed against `rpi-image-gen` commit:

```text
d1021e82dd578b588cc3b4d45cd7b4b86e57b796
```

The build helper rejects a different revision unless
`ALLOW_UNPINNED_RPI_IMAGE_GEN=1` is set.

## Build

From the MasjidPi repository root:

```bash
./appliance-image/build-pi3-ab-prototype.sh "$HOME/rpi-image-gen"
```

For a local hardware test that needs a login, create the ignored file
`config/pi3-ab-local.yaml` with the production configuration as its base and
pass it as the second argument. Never commit a test password:

```yaml
include:
  file: pi3-ab-prototype.yaml

device:
  user1pass: "choose-a-valid-temporary-password"
```

```bash
./appliance-image/build-pi3-ab-prototype.sh \
  "$HOME/rpi-image-gen" \
  "$PWD/appliance-image/config/pi3-ab-local.yaml"
```

The compressed image and its associated metadata are written under the
`MasjidPi/work/deploy-*` directory reported at the end of the build.

## Acceptance criteria

Before adding automatic slot selection, the image must:

1. Build without modifying the `rpi-image-gen` checkout.
2. Contain exactly four MBR partitions with the labels above.
3. Boot a Raspberry Pi 3 Model B from both `SYSTEM_A` and `SYSTEM_B`.
4. Use a slot-neutral root filesystem and keep `SYSTEM_A` as the default.
5. Mount `PERSISTENT` at `/persistent` from either system slot.
6. Leave enough unused card capacity to support the nominal 16 GB target.

The prototype has no production provisioning flow and creates no login
password. Serial console and local boot diagnostics should be used for this
smoke test.
