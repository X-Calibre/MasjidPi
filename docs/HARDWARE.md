# MasjidPi Hardware Guide

This guide records supported release architectures, validated Raspberry Pi platforms and measured appliance performance.

## Supported release architectures

Official release packages are built for:

- **Linux ARM64 (`aarch64`)**
- **Linux AMD64 (`x86_64`)**

MasjidPi requires `systemd`. Listen requires an ALSA-compatible audio device. Board requires a DRM/KMS display environment supported by the Cog/WPE packages installed by MasjidPi.

## Raspberry Pi compatibility

Raspberry Pi 3B and Raspberry Pi 4 are production-validated. Other entries are expectations based on architecture and measured workloads, not support claims.

| Raspberry Pi | Listen | Listen + Board | Status |
|---|---:|---:|---|
| Zero W / Pi 1 | No | No | ARMv6 is not supported by current release packages. |
| Pi 2 Model B v1.1 | Unvalidated | Unvalidated | Older 32-bit Cortex-A7 platform. |
| Pi 2 Model B v1.2 | Expected | Expected with constraints | Cortex-A53 platform, but not physically validated. |
| Zero 2 W | Expected | Unvalidated | CPU should be sufficient; 512 MB RAM requires validation for Board. |
| Pi 3A+ | Expected | Unvalidated | Cortex-A53 with 512 MB RAM. |
| **Pi 3B** | **Validated** | **Validated** | Current 64-bit reference platform. |
| Pi 3B+ | Expected | Expected | Same memory class as the validated Pi 3B with a faster CPU. |
| **Pi 4** | **Validated** | **Validated** | Validated with simultaneous Listen and GLES-backed Board operation. |
| Pi 5 | Expected | Expected | Substantially more performance than MasjidPi requires. |
| Compute Module 3/3+ | Expected | Expected | Depends on carrier, RAM and audio/display hardware. |
| Compute Module 4/5 | Expected | Expected | Suitable performance; appliance integration remains hardware-specific. |

## Raspberry Pi 3B performance baseline

A production-style Listen + Board installation was measured on a Raspberry Pi 3B running 64-bit Raspberry Pi OS Lite/Debian 13-class software.

During simultaneous Listen playback and Cog/WPE Board rendering:

- approximately 337 MiB RAM was used and 567 MiB remained available;
- swap was unused after a fresh boot;
- the MasjidPi backend used approximately 15 MiB RSS and about 2% CPU;
- mpv used approximately 75–80 MiB RSS and about 7% CPU;
- the WPE renderer used approximately 190–195 MiB RSS and about 31% CPU;
- CPU temperature was approximately 58–60 °C; and
- no thermal or undervoltage throttling was reported.

With Board stopped and Listen still running, total memory use fell to approximately 195 MiB. The Board display stack therefore added roughly 140–150 MiB in this sample.

These measurements show comfortable headroom on the Pi 3B. They do not replace real-device and soak validation for 512 MB models.

## Display hardware

Board supports:

- a responsive landscape TV/Monitor profile; and
- a portrait 600 × 1024 Appliance Display profile with touch controls.

The local display service uses Cog/WPE directly on DRM/KMS. Display detection, orientation and touchscreen calibration are handled by the appliance runtime rather than saved as Board layout preferences.

## Audio hardware

Listen works with ALSA-compatible outputs. Hardware Master Volume is available only when the selected output exposes a controllable mixer. Masjid and Radio software volumes remain available independently.

USB audio hot-plug fallback and restoration have been validated with a Logitech H390. Individual sound cards, HDMI sinks and speaker/amplifier combinations should still be tested as complete hardware paths.

## Remaining hardware validation

- longer-duration testing of the final feature set on Raspberry Pi 3B;
- Raspberry Pi Zero 2 W and Pi 3A+ memory validation;
- HDMI disconnect/reconnect behavior;
- physical HDMI-CEC support; and
- Waveshare display audio and volume behavior.
