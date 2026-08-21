# MasjidBoard Implementation Status

**Status:** Functional appliance implementation / final pre-integration work  
**Branch:** `research/masjidboard-live`

## Purpose

Record the current implemented and validated state of MasjidBoard on the research branch. Older research documents may describe approaches that have since been replaced; this document records the current working behaviour.

## Current Architecture

MasjidPi now supports three installation profiles:

```text
Listen
Board
Listen + Board
```

The selected profile is persisted in `/etc/masjidpi/components.env` and is used by the backend, API registration, installer, runtime dependencies, systemd services and installer self-test.

Listen and Board remain independent capabilities built on the common MasjidPi application runtime. A Board provider/cache/display failure must not prevent Listen from operating, and a Board-only appliance does not start MPV.

Current MasjidBoard components include:

- MasjidBoard Live discovery and hierarchy retrieval;
- persisted location scope;
- scoped Masjid catalogue;
- persisted selection and ordering of up to three Masjids;
- MasjidBoard Live Core timetable provider;
- normalized board/prayer model;
- per-board last-known-good cache;
- per-board runtime status and recovery;
- display presentation API;
- dedicated MasjidBoard configuration WebUI;
- read-only MasjidBoard display page;
- Cog/WPE DRM appliance display runtime on Raspberry Pi OS Lite.

## Component-Aware Appliance Installation

The installer prompts for Listen, Board or Listen + Board on fresh installation and reinstall/update flows. Existing installations show the current profile as the default while still allowing the profile to be changed.

Runtime dependencies are profile-aware:

- Listen installs MPV/FFmpeg/ALSA dependencies;
- Board installs Cog/WPE display dependencies;
- combined installs receive both sets.

Backend subsystem startup and API registration are also profile-aware. A disabled component does not expose its API endpoints.

Validated runtime behaviour is:

```text
Listen only
    -> masjidpi + mpv
    -> no Cog
    -> player API available
    -> Board API returns 404

Board only
    -> masjidpi + Cog/WPE
    -> no mpv
    -> Board API available
    -> player API returns 404

Listen + Board
    -> masjidpi + mpv + Cog/WPE
    -> both APIs available
```

The configuration WebUI also hides configuration/navigation for components that are not installed.

Component-profile changes participate in the safe update workflow. Profile state is transactional and can be restored with the previous runtime during rollback. Installer self-tests validate only the components selected for the candidate installation.

Profile transitions between Listen-only, Board-only and Listen + Board have been exercised on Raspberry Pi hardware, including adding and removing each component. Removing Board stops/disables the display runtime, removes the display unit/launcher and clears stale systemd failed-unit state.

## Display Runtime

The production Raspberry Pi display direction is now Cog with WPE WebKit rendering directly through DRM/KMS:

```text
cog --platform=drm http://127.0.0.1:8080/masjidboard.html
```

This has been validated on Raspberry Pi OS Lite 64-bit with an HDMI display. It does not require Chromium, a desktop environment, X11, a Wayland compositor or graphical login session.

`masjidpi-display.service` starts the display automatically when Board is installed and restarts it after unexpected Cog termination. Crash/restart behaviour has been tested successfully.

See `MASJIDBOARD-DISPLAY-RUNTIME.md` for the current runtime design and validation details.

## Discovery and Selection

The configuration UI supports refreshing the location hierarchy, selecting/persisting locations, building and refreshing the scoped Masjid catalogue, selecting up to three Masjids, and preserving/reordering the selected display order.

Refresh actions remain intentionally distinct:

- **Refresh Location List** — refresh the upstream hierarchy;
- **Refresh Masjid List** — refresh Masjids available within the selected locations;
- **Refresh Timetables** — refresh live timetable data for selected Masjids.

## Timetable Retrieval, Cache and Recovery

The MasjidBoard Live Core board page is the working primary timetable source. Selected boards refresh independently and a failure for one board does not invalidate the others.

Automatic timetable refresh runs every 30 minutes, with manual refresh available from the configuration UI. Startup refresh is asynchronous.

Each selected Masjid has an independent last-known-good timetable cache. Verified behaviour is:

```text
provider available
    -> current live timetable

provider unavailable with cache
    -> stale
    -> cached timetable remains displayed
    -> update error recorded

provider available again
    -> next successful refresh returns board to current
```

This lifecycle has been verified in automated tests and a real provider-outage runtime test. Identical timetable data is not rewritten on every refresh.

## Default Display

The read-only display is `/masjidboard.html` and supports one, two or three selected Masjids.

The default layout shows current local time/date, selected Masjid names, Fajr, Dhuhr or Friday Jumu'ah, Asr, Maghrib, Esha, per-board stale/unavailable state, and a per-Masjid countdown to the next visible timetable event.

Jumu'ah replaces Dhuhr on Friday. Timed Jumu'ah events are displayed chronologically with provider labels preserved. The countdown rolls to the following day's first Fajr event after the final visible event of the day.

Alternative layouts remain future presentation work and should consume the same normalized display API.

## Validation Completed

Automated Go coverage includes API endpoints, hierarchy/discovery, scoped catalogue and reconciliation, selection persistence, provider parsing/normalization, Jumu'ah handling, cache persistence/write suppression, display presentation, runtime behaviour and stale-cache recovery.

Manual/runtime validation now includes:

- location and Masjid selection;
- three-board display;
- Friday/Jumu'ah behaviour;
- next-event and overnight countdowns;
- provider outage/cache/recovery;
- Raspberry Pi OS Lite HDMI display using Cog/WPE DRM;
- display process restart recovery;
- reboot/service restart appliance behaviour;
- Listen-only appliance operation;
- Board-only appliance operation;
- combined Listen + Board operation;
- transitions among all three component profiles;
- component-aware API exposure and process startup;
- component-aware installer dependencies and self-tests;
- transactional profile/update handling and rollback;
- Board display-service cleanup when Board is removed.

## Remaining Pre-Integration Work

The major installer/component-profile and Raspberry Pi display-runtime work is now complete. Remaining work before branch integration should be focused:

1. Verify production ownership/permissions and preservation of hierarchy, scope, catalogue, selection and per-board cache files across fresh install/update paths.
2. Perform first-run validation with Board installed but no MasjidBoard configuration/cache.
3. Review user-facing API/frontend error messages while retaining diagnostic detail in logs/status data.
4. Run a final Listen/audio regression pass after the complete branch changes.
5. Perform longer-duration Raspberry Pi stability/resource testing and practical HDMI power-cycle/reconnect testing.
6. Complete final documentation/integration review and prepare the branch for merge/release planning.

Additional layouts and richer optional Board content are not blockers for this integration milestone.
