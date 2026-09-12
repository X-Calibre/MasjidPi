# MasjidPi Raspberry Pi 3 appliance image

This directory contains the Raspberry Pi 3 image prototype built with
[`rpi-image-gen`](https://github.com/raspberrypi/rpi-image-gen).

The first milestone deliberately boots only system slot A. It proves that a
Raspberry Pi 3 can boot from a four-partition MBR image and that the proposed
layout fits comfortably on a nominal 16 GB microSD card before U-Boot, RAUC,
and automatic rollback are introduced.

## Prototype layout

| Partition | Label | Prototype size | Purpose |
|---|---|---:|---|
| 1 | `BOOT` | 256 MiB | Raspberry Pi firmware, kernel, and later U-Boot |
| 2 | `SYSTEM_A` | 3 GiB | Initially active operating-system slot |
| 3 | `SYSTEM_B` | 3 GiB | Inactive operating-system slot |
| 4 | `PERSISTENT` | 2 GiB | State shared across system slots |

Both system partitions initially contain the same filesystem image. The boot
command line explicitly selects partition 2. Do not treat this milestone as a
working updater or rollback implementation.

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

The compressed image and its associated metadata are written under the
`rpi-image-gen/work/deploy-*` directory reported at the end of the build.

## Acceptance criteria

Before adding slot selection, the image must:

1. Build without modifying the `rpi-image-gen` checkout.
2. Contain exactly four MBR partitions with the labels above.
3. Boot a Raspberry Pi 3 Model B from system slot A.
4. Mount `PERSISTENT` at `/persistent`.
5. Leave enough unused card capacity to support the nominal 16 GB target.

The prototype has no production provisioning flow and creates no login
password. Serial console and local boot diagnostics should be used for this
first smoke test.
