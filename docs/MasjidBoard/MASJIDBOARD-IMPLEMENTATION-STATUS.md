# MasjidBoard Implementation Status

**Status:** Implemented in MasjidPi v1.5.2

This document summarises the current production implementation. Provider research and superseded design proposals are retained separately as historical records.

## Architecture

MasjidPi supports three component profiles:

- Listen
- Board
- Listen + Board

The profile is stored in `/etc/masjidpi/components.env` and controls backend startup, API registration, dependencies, systemd services and installer self-tests. Listen and Board remain independently operable within one application.

Current Board packages cover:

- MasjidBoard Live hierarchy and scoped discovery;
- selected-board persistence and ordering;
- Core timetable retrieval and optional Premium enrichment;
- normalized models;
- per-board last-known-good caches;
- shared daily Islamic content and economic data;
- runtime refresh/status handling; and
- a stable display API consumed by browser renderers.

## Installation and display runtime

Board installations use Cog/WPE directly on DRM/KMS. The systemd-managed runtime includes display warm-up, Plymouth-to-Cog handoff, display restart behavior and Raspberry Pi boot-firmware write protection.

The launcher selects:

- the standard responsive landscape profile for ordinary displays; or
- the portrait Appliance profile for the validated Waveshare display/touchscreen combination.

Display profile is detected from hardware and is not a saved Board preference.

## Discovery and selection

The Web UI can:

- refresh the global location hierarchy;
- select and persist one to three locations;
- build and refresh a scoped masjid catalogue;
- select and order up to three masjids;
- enable or disable detailed Jumu'ah cards per masjid; and
- refresh selected timetables.

Hierarchy, catalogue, selection and timetable refresh remain separate operations so one upstream failure does not erase other usable state.

## Timetables and recovery

Selected boards refresh independently. Each board is reported as:

- `current` when fresh data is available;
- `stale` when its last-known-good cache is being used; or
- `unavailable` when neither current nor cached data can be displayed.

Identical successful data is not rewritten unnecessarily. Core remains authoritative for the timetable; failed optional enrichment does not make a successful Core timetable stale.

## Presentation

Both profiles support:

- one to three ordered masjids;
- five daily prayers with Jumu'ah replacing Dhuhr during the Islamic-Friday interval;
- chronological Jumu'ah events and optional detailed schedule cards;
- next-event countdowns with overnight Fajr rollover;
- Gregorian and masjid-adjusted Islamic dates;
- primary-masjid Daily Times;
- provider-supplied special Dhuhr when distinct from normal Dhuhr;
- Zawaal/Istiwaa clock/date warning;
- stale/unavailable state; and
- ten colour themes.

The Appliance profile uses one timetable slide per masjid and includes a touch-control sheet for selected Listen and theme actions.

## Community and shared content

Supported content includes:

- general and urgent announcements;
- Salaah and class-time changes;
- weekly and Ramadan programmes;
- funeral, Nikah and Eid notices;
- well-wishes;
- Taleem, Dawah/Gasht and three-day Jamaat information;
- contribution appeals;
- new-moon information;
- Daily Ayah, Hadith and Sunnah; and
- Islamic Economic Indicators.

Community cards complete one selected masjid at a time in priority groups. Shared daily Islamic content and economic data are shown once after masjid-specific content.

Arabic and mixed-direction text is sanitized and rendered with automatic direction. Static cellphone-reminder chrome and poster/media content are deliberately excluded.

The optional built-in Dua after Adhan is disabled by default. When enabled, it takes over the Appliance slide or complete landscape notice column for five minutes beginning at a primary-masjid Adhan time.

## Saved preferences

Saved Board preferences include:

- theme;
- slide duration;
- Islamic Economic Indicators visibility;
- Daily Ayah, Hadith and Sunnah visibility;
- Dua after Adhan visibility; and
- per-masjid detailed Jumu'ah visibility.

The display profile itself is not persisted.

## Validation

v1.5.2 validation includes automated Go, race, shell, installer and JavaScript checks plus Raspberry Pi 4 source-install testing. Portrait and landscape rendering, detailed Jumu'ah, structured community content, Dua priority, Zawaal warning, special Dhuhr suppression, 11-item Daily Times layout, persistence and service health were verified.

See [../VALIDATION_CHECKLIST.md](../VALIDATION_CHECKLIST.md) and [../RELEASE_CANDIDATE_v1.5.2.md](../RELEASE_CANDIDATE_v1.5.2.md).

## Remaining work

Non-blocking follow-up work includes:

- poster/media transport and caching;
- broader real-world Arabic/RTL validation;
- Ramadan Maghrib/Iftar source validation;
- HDMI reconnect and wider hardware testing; and
- longer-duration soak monitoring of the final feature set.
