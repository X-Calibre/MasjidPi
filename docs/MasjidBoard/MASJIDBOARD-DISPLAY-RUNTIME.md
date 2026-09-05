# MasjidBoard Raspberry Pi Display Runtime

**Status:** Cog/DRM appliance runtime implemented and validated on Raspberry Pi OS Lite 64-bit

## Purpose

Document the working MasjidBoard HDMI appliance display runtime.

MasjidPi is intended to be a home appliance. When Board is installed, the attached HDMI display is therefore the normal MasjidBoard presentation path rather than a user manually opening the page in another browser.

The display page remains:

```text
http://127.0.0.1:8080/masjidboard.html
```

## Validated Runtime

The working appliance renderer is **Cog with WPE WebKit using the DRM platform directly**:

```text
cog --platform=drm --platform-params=renderer=gles http://127.0.0.1:8080/masjidboard.html
```

This was validated on Raspberry Pi OS Lite 64-bit. A full desktop environment, Chromium, X11, Wayland compositor and graphical login session are not required.

The appliance path is:

```text
Raspberry Pi OS Lite boots
        ↓
MasjidPi backend starts
        ↓
masjidpi-display.service starts
        ↓
Cog/WPE renders directly through DRM/KMS
        ↓
TV/monitor shows MasjidBoard
```

The production boot path also uses `masjidpi-display-warmup.service` to prepare WPE while Plymouth retains the display, then hands DRM ownership to the main display service.

This is substantially better aligned with the MasjidPi appliance goal than the earlier Chromium/labwc prototype direction.

## Why Cog / WPE

The existing HTML/CSS/JavaScript display remains the renderer, while Cog provides a lightweight embedded browser shell designed for this kind of appliance use.

This preserves:

- the existing responsive layout;
- Friday/Jumu'ah behaviour;
- per-Masjid countdowns;
- stale/current presentation;
- one-, two- and three-board layouts;
- both current display profiles through the same display API.

It also avoids maintaining a second native framebuffer/SDL presentation implementation.

## Installation

Board installations require:

```text
cog
libwpewebkit-2.0-1
```

The installer installs these only when Board is selected. Listen-only installations do not require Cog/WPE.

The installed component profile is stored in:

```text
/etc/masjidpi/components.env
```

Supported profiles are:

```text
listen
board
listen,board
```

## Display Service

Board installs use the dedicated systemd unit:

```text
masjidpi-display.service
```

The service launches the installed display wrapper, which starts Cog against the local MasjidBoard page using the DRM platform.

The service is enabled and started when Board is installed. When Board is removed from the component profile, the installer stops and disables the display service, removes its unit/runtime launcher and clears any stale failed-unit state.

The display service is configured to restart after an unexpected Cog exit. This recovery behaviour has been tested by terminating Cog and confirming that systemd launches a replacement process and the Board returns to the HDMI display.

## Backend Independence

`masjidpi.service` remains the common application service because it supplies the shared HTTP/API runtime and whichever installed subsystems are enabled.

Component startup is profile-aware:

```text
Listen only
    -> MasjidPi backend + MPV
    -> no Board subsystem
    -> no Cog display

Board only
    -> MasjidPi backend + Board subsystem + Cog display
    -> no MPV

Listen + Board
    -> MasjidPi backend + MPV + Board subsystem + Cog display
```

Listen and Board remain functionally independent. Board does not control playback and Listen does not require Board.

## Validated Raspberry Pi Behaviour

The Cog/DRM runtime has been tested successfully on Raspberry Pi hardware with an HDMI Samsung display.

Validated behaviour includes:

- correct MasjidBoard rendering over HDMI;
- automatic startup under systemd;
- Cog crash/restart recovery;
- MasjidPi service restart while the display runtime remains operational;
- reboot recovery;
- Board-only operation with no MPV process;
- Listen-only operation with no Cog process or Board API;
- combined Listen + Board operation;
- component-profile transitions between Listen, Board and Listen + Board;
- component-aware installer self-tests;
- transactional update/profile handling and rollback;
- cleanup of the Board display service when Board is removed.

During combined Listen + Board testing on a roughly 1 GB Raspberry Pi, the system retained useful available memory and ran without swap pressure during the measured test. Cog/WPE is still the largest application memory consumer, so longer-duration appliance testing remains useful, but the runtime has demonstrated that Raspberry Pi OS Lite is viable for this display approach.

### Raspberry Pi 4 validation — 23 August 2026

The notice-card and default Landscape display work was validated on a Raspberry Pi 4 running 64-bit Raspberry Pi OS Lite with kernel `6.18.39+rpt-rpi-v8` and a native 1920×1080 HDMI display.

The default Cog `modeset` renderer loaded the page but failed to create a framebuffer with `Invalid argument`. The supported GLES renderer initialized successfully on the Raspberry Pi VC4 DRM device and is therefore used by the production launcher:

```text
cog --platform=drm --platform-params=renderer=gles http://127.0.0.1:8080/masjidboard.html
```

Validation passed for:

- installation and automatic HDMI display startup;
- recovery after source update/reinstallation;
- recovery after a full reboot;
- complete rotating fixtures, including one-, two- and three-card arrangements;
- Dawah/Gasht, three-day Jamaat and contribution cards;
- live theme changes and persistence of the saved Slate theme;
- persisted masjid configuration;
- simultaneous Listen playback and Board rendering;
- stable display operation without unexpected service restarts;
- no swap use and no CPU throttling.

During combined playback/display operation, approximately 3.3 GiB of 3.7 GiB RAM remained available, system load was low, temperature was 64.2°C, and `vcgencmd get_throttled` reported `0x0`.

### v1.4.0 hardware validation — 24 August 2026

The published v1.4.0 release was installed successfully on both the Raspberry Pi 4 test unit and Raspberry Pi 3B production unit. Listen and Board operated correctly on both devices, including the optional Islamic Economic Indicators content sourced from Jamiatul Ulama South Africa. The Nisaab, Krugerrand, gold, silver and Mahr values rendered correctly on the target displays.

On this platform Cog may abort during an intentional systemd stop/restart while releasing EGL resources. The replacement process starts normally and the behaviour was not observed during steady-state operation.

## Expected Failure Isolation

```text
Display process exits
    -> systemd restarts Cog
    -> Listen remains independent

MasjidPi backend restarts
    -> display may temporarily lose the local page/API
    -> display recovers when the backend returns

MasjidBoard Live unavailable
    -> cached timetable remains available where possible
    -> display remains running
    -> Listen remains unaffected

Board not installed
    -> no Cog/WPE display service
    -> Board API is not registered
```

## Remaining Production Work

The display runtime itself is now established. Remaining work should focus on production polish rather than replacing the renderer:

- longer-duration stability/resource testing on supported Raspberry Pi hardware;
- monitor/HDMI power-cycle and disconnect/reconnect behaviour where practical;
- final user-facing error presentation;
- continued release/integration validation.

A Chromium/labwc desktop kiosk is no longer the preferred production direction for Raspberry Pi appliances. It may remain useful for development or alternative platforms, but the Raspberry Pi appliance target should use Cog/WPE DRM unless later testing establishes a concrete reason to change it.

## Future Layouts

Automatic display startup is independent of presentation profile. Any future presentation should continue to use the same browser-based display runtime and normalized display API unless hardware testing establishes a reason to change renderer technology.
