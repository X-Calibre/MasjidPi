# MasjidBoard Implementation Status

**Status:** v1.3.0 integration candidate
**Branch:** `docs/masjidboard-live-data-inventory` / PR #38

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
- Cog/WPE DRM appliance display runtime on Raspberry Pi OS Lite;
- responsive TV / Monitor and dedicated 7-inch Appliance Display modes;
- Core timetable retrieval with optional Premium community-content enrichment; and
- adaptive rotating community-content cards in Landscape mode.

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
cog --platform=drm --platform-params=renderer=gles http://127.0.0.1:8080/masjidboard.html
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

## Display Layouts and Community Content

The read-only display is `/masjidboard.html`, supports one, two or three selected Masjids, and offers a responsive TV / Monitor mode plus a dedicated 7-inch Appliance Display (600 × 1024) mode.

TV / Monitor mode shows current local time/date, selected Masjid names, Fajr, Dhuhr or Friday Jumu'ah, Asr, Maghrib, Esha, per-board stale/unavailable state, per-Masjid countdowns and a full-width Daily Times footer. The 7-inch Appliance Display presents the same core timetable information for the integrated physical-appliance screen.

Jumu'ah replaces Dhuhr on Friday. Timed Jumu'ah events are displayed chronologically with provider labels preserved. The countdown rolls to the following day's first Fajr event after the final visible event of the day.

When a selected board exposes public Premium content, it can enrich the successful Core timetable with active announcements and structured cards. Supported categories include Nikah, funeral, Eid, Salaah changes, well-wishes, Taleem, Dawah/Gasht, three-day Jamaat, contributions and calculated new-moon information. Core remains authoritative for the timetable, and missing or failed Premium enrichment does not make a successful Core board stale. Landscape rotates additional card pages, suppresses duplicates, converts upstream HTML to plain text and labels each card with its source.

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
- TV / Monitor and 7-inch Appliance Display preference changes and persistence;
- complete anonymised rotating-card fixtures, including RTL and dense-content cases;
- optional Premium enrichment with Core fallback;
- Raspberry Pi 4 native-1080p display using Cog's GLES renderer;
- simultaneous Listen playback and Board rendering on Raspberry Pi 4; and
- Raspberry Pi 4 recovery after reinstall and full reboot without throttling or unexpected display restarts.

## Remaining Release Work

The automated CI gate and Raspberry Pi 4 functional acceptance pass are complete. Remaining work is to complete the documentation review, merge PR #38, prepare and verify v1.3.0 release artifacts, and publish only after those artifacts pass installation validation.

Longer-duration soak testing, practical HDMI disconnect/reconnect testing, broader hardware validation and poster/media support remain follow-up work rather than release blockers.
