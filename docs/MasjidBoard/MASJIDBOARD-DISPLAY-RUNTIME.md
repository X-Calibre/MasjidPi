# MasjidBoard Raspberry Pi Display Runtime

**Status:** Initial appliance-display implementation / hardware validation required  
**Branch:** `research/masjidboard-live`

## Purpose

Define how the existing browser-based MasjidBoard display becomes an automatic Raspberry Pi appliance display rather than requiring a user to browse manually to `/masjidboard.html`.

The current default display remains:

```text
http://127.0.0.1:8080/masjidboard.html
```

The appliance runtime is responsible for launching that page full-screen on an attached TV/monitor after boot.

## Current Direction

The first production candidate is a Chromium kiosk running in the Raspberry Pi OS graphical session.

Raspberry Pi OS uses Wayland with the `labwc` compositor by default from Bookworm onward. Raspberry Pi's current kiosk guidance also uses Chromium launched from `labwc` autostart. MasjidPi should follow that supported platform direction rather than build a separate X11-only kiosk stack.

The intended runtime path is:

```text
Raspberry Pi boots
        ↓
MasjidPi backend starts
        ↓
graphical Wayland/labwc session starts
        ↓
MasjidBoard kiosk launcher waits for MasjidPi HTTP API
        ↓
Chromium opens /masjidboard.html in kiosk mode
        ↓
TV/monitor shows MasjidBoard automatically
```

## Browser-Based Display Decision

The working HTML/CSS/JavaScript display is retained as the renderer for the initial appliance implementation.

This avoids duplicating the presentation in SDL/framebuffer code and preserves:

- the existing responsive layout;
- Friday/Jumu'ah behaviour;
- per-Masjid countdowns;
- stale/current presentation;
- one-, two- and three-board layouts;
- future selectable layout support through the same display API.

A native renderer is not required unless later Raspberry Pi measurements show that the browser approach is unsuitable on supported hardware.

## Initial Launcher

The research branch contains:

```text
scripts/masjidboard-display.sh
```

The launcher:

- locates `chromium` or `chromium-browser`;
- waits for the local MasjidPi API before opening the kiosk;
- opens `http://127.0.0.1:8080/masjidboard.html` full-screen;
- uses Wayland explicitly;
- suppresses first-run/crash/notification browser UI;
- stores Chromium profile/cache data under `/tmp` to avoid unnecessary persistent SD-card writes.

The target URL and readiness URL can be overridden for development with:

```text
MASJIDBOARD_URL
MASJIDPI_READY_URL
```

## Session Startup

For Raspberry Pi OS Desktop/labwc, the intended automatic startup mechanism is the user's:

```text
~/.config/labwc/autostart
```

The final installer is expected to add a MasjidPi-managed kiosk launch entry there, preferably through `lwrespawn` where available so the browser is relaunched after an unexpected exit.

Conceptually:

```text
/usr/bin/lwrespawn /opt/masjidpi/bin/masjidboard-display &
```

The exact installed path and installer behaviour will be finalized after hardware testing.

## Raspberry Pi OS Lite

Existing MasjidPi audio deployments have used Raspberry Pi OS Lite. Lite does not include the complete graphical desktop/session stack required by a browser kiosk.

Before production integration, MasjidPi must choose and validate one of these approaches:

1. use Raspberry Pi OS Desktop for appliances that enable MasjidBoard display; or
2. install a minimal Wayland/labwc/Chromium session on top of Lite.

The initial hardware test should determine whether the Desktop image is sufficiently lightweight on a Raspberry Pi 3 B running Listen + Board + Chromium. If resource use is acceptable, using the supported Raspberry Pi OS graphical environment is preferred over maintaining a custom minimal display stack.

## Failure Isolation

The display is optional relative to Listen and the MasjidBoard backend.

Expected behaviour:

```text
No monitor / graphical display unavailable
    -> Listen continues
    -> MasjidBoard backend continues

Chromium exits unexpectedly
    -> display launcher/session restarts Chromium
    -> Listen remains unaffected

MasjidPi backend restarts
    -> existing display may show connection interruption
    -> page resumes when backend returns

MasjidBoard Live unavailable
    -> cached timetable remains displayed
    -> browser remains running
    -> Listen remains unaffected
```

## Hardware Support Boundary

Raspberry Pi's current kiosk guidance requires a Raspberry Pi 3 or newer with at least 1 GB RAM for Chromium/Firefox kiosk use. The Raspberry Pi 3 B therefore meets the documented minimum, but combined MasjidPi Listen + Board + browser resource use still needs direct measurement.

## Hardware Validation Plan

The first full test should be on the Raspberry Pi 3 B and an HDMI TV/monitor.

Validate:

```text
1. boot with monitor connected
2. automatic graphical login/session
3. automatic MasjidBoard kiosk launch
4. correct 1920x1080 or monitor-native layout
5. no visible browser chrome
6. Listen playback operates concurrently
7. CPU and RAM usage during idle display
8. Chromium crash/relaunch
9. MasjidPi service restart while kiosk remains active
10. provider/network outage with cached display
11. reboot recovery
12. monitor power off/on
13. HDMI disconnect/reconnect where practical
```

## Production Integration Still Required

After the hardware prototype is validated:

- install Chromium/display dependencies when Board display is enabled;
- install the launcher into the runtime rather than running it from the source tree;
- configure graphical auto-login/session startup;
- add the labwc autostart entry safely and idempotently;
- preserve user/session configuration on upgrades;
- decide how a Listen-only appliance opts out of browser dependencies;
- include display checks in installation/self-test where appropriate;
- document supported Raspberry Pi OS image requirements.

## Future Layouts

Automatic display startup is independent of layout selection. Future user-selectable MasjidBoard layouts should still be served through the same browser-based display runtime and normalized display API unless hardware testing establishes a reason to change renderer technology.
