# MasjidBoard Display and Configuration Boundary

**Status:** Implemented in MasjidPi v1.5.2

## Principle

The timetable/content presentation is read-only. Discovery, selected-board administration and persistent settings belong to the MasjidPi configuration Web UI and APIs.

The Appliance profile includes a deliberately narrow touch-control overlay for everyday Listen actions and theme selection. It does not turn the HDMI presentation into an administrative interface.

## Display responsibilities

The display may:

- render selected timetable and community data;
- format board-local times and dates;
- calculate Friday, countdown and Zawaal warning state;
- rotate or prioritize cards;
- show stale/unavailable state;
- show Listen source notifications;
- expose the limited Appliance touch overlay; and
- retry presentation retrieval while retaining usable rendered data.

The display must not:

- search or alter the MasjidBoard catalogue;
- change geographic scope or selected boards;
- edit schedules;
- select audio devices;
- expose provider/cache diagnostics; or
- write upstream content.

## Configuration responsibilities

The Web UI and administrative APIs own:

- location hierarchy and scope;
- catalogue refresh/search;
- selected-board add/remove/order;
- per-board detailed Jumu'ah preference;
- theme and slide duration;
- Daily Ayah/Hadith/Sunnah preferences;
- Dua-after-Adhan preference;
- Islamic Economic Indicator preference; and
- timetable refresh/status.

## Presentation API

```text
GET /api/masjidboard/display
```

The response contains only the stable presentation state required by the renderers. It excludes raw upstream payloads, cache paths and detailed errors.

Administrative state uses separate endpoints including:

```text
GET  /api/masjidboard/status
GET  /api/masjidboard/hierarchy
POST /api/masjidboard/hierarchy/refresh
GET  /api/masjidboard/scope
PUT  /api/masjidboard/scope
GET  /api/masjidboard/catalogue
POST /api/masjidboard/catalogue/refresh
GET  /api/masjidboard/selection
PUT  /api/masjidboard/selection
GET  /api/masjidboard/layout
PUT  /api/masjidboard/layout
POST /api/masjidboard/boards/refresh
```

## Profiles

The local launcher chooses standard or appliance presentation from attached hardware. Profile is not persisted through the layout API.

A remote browser may open:

```text
/masjidboard.html
/masjidboard.html?profile=appliance
```

The query parameter changes only that browser view.

## Touch-control exception

The Appliance overlay may:

- select a favourited Listen masjid;
- select a Radio station;
- choose scheduled/immediate/stopped Radio mode;
- start or stop Listen;
- adjust source and supported Master volumes; and
- change the saved Board theme.

Catalogue search, Radio schedule editing, Board selection and audio-device configuration remain in the full Web UI.

## Failure behavior

An upstream or optional-content failure must not blank already usable prayer data. Per-board state is current, stale or unavailable, and selected positions remain stable. A temporary display-API failure leaves the last rendered view on screen while retrying.

Board failures do not stop Listen, and Listen failures do not stop Board timetable presentation.
