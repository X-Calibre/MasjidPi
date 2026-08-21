# MasjidBoard Display and Configuration Boundary

**Status:** Architecture decision; presentation API implemented  
**Branch:** `research/masjidboard-live`

## Decision

The MasjidBoard display is a **read-only presentation surface**.

It exists only to display MasjidBoard information for the configured one to three selected masjids. It must not contain configuration controls, discovery/search workflows, board-selection controls, catalogue maintenance actions, ordering controls or administrative preferences.

All configuration and preferences belong to the MasjidPi API and/or administrative WebUI.

## Responsibility Split

```text
MasjidPi WebUI / configuration API
    +-- configure 1-3 discovery locations
    +-- browse locally persisted catalogue
    +-- select 1-3 boards
    +-- reorder/replace selected boards
    +-- request catalogue refresh
    +-- configure display preferences
    +-- inspect diagnostics
             |
             v
MasjidBoard backend service
    +-- persisted discovery scope/catalogue/selection
    +-- provider construction
    +-- independent live refresh
    +-- independent last-known-good caches
    +-- current / stale / unavailable runtime state
             |
             v
GET /api/masjidboard/display
             |
             v
MasjidBoard display
    +-- READ ONLY
    +-- render selected board information
    +-- preserve configured board order
    +-- show per-board stale indication
```

The display does not need to know how a board was discovered, selected, refreshed or configured.

## Presentation API

The display-facing endpoint is:

```text
GET /api/masjidboard/display
```

It is derived from the selected-board runtime state but exposes a deliberately smaller provider-neutral contract than the diagnostic status API.

The top-level model is:

```text
DisplayView
    configured
    boards[]
```

Each board contains presentation information only:

```text
catalogue_id
name
alternate_name (optional)
time_zone (optional)
status              current | stale | unavailable
stale               boolean
last_successful_update (optional)
date
prayers[]
jumuah[] (optional)
astronomical (optional)
```

The five daily prayers are always presented in stable order:

```text
Fajr
Dhuhr
Asr
Maghrib
Esha
```

Each row contains a stable key, display label and optional Adhan/Jamaah local wall-clock times. Missing individual times are omitted rather than fabricated.

Jumu'ah presentation retains available source headings/events but also exposes `effective_salaah`, which uses a supplied Jumu'ah Jamaah time when available and otherwise the normalised Khutbah fallback already defined by the domain model.

Astronomical values remain optional. The presentation model currently carries the normalised astronomical fields so later screen layouts/preferences can use them without making the display consume the full internal `Board` object.

Clock times have an explicit display JSON shape:

```json
{
  "hour": 16,
  "minute": 45
}
```

The display model owns this JSON contract independently of the internal Go domain structs.

## Initial Display Content Scope

The initial household display should focus on the information required to decide **where and when to pray** rather than displaying every available timetable field.

The default layout should show:

- masjid name;
- current local date/time;
- Fajr, Dhuhr, Asr, Maghrib and Esha;
- Adhan where supplied;
- Jamaah where supplied;
- Jumu'ah information where supplied; and
- a concise per-board stale/unavailable indication.

The following data may remain available through the presentation API but is **not part of the default initial layout**:

```text
Suhur
Fajr Start
Sunrise
Ishraaq
Duha
Istiwa / Zawaal
Asr Shafi'i
Asr Hanafi
Sunset
Esha Start
```

These values are reserved for later layouts/preferences such as an extended timetable or Ramadan-oriented view.

## Prayer Time Presentation Rules

The information order for each daily prayer is fixed:

```text
Prayer name
Adhan
Jamaah
```

Adhan must therefore appear before Jamaah in reading order.

Visual emphasis is different from reading order:

- when both Adhan and Jamaah are present, **Jamaah is visually dominant** (for example larger and/or bolder);
- Adhan remains visible first but with secondary emphasis;
- when only one of Adhan or Jamaah is supplied, that available time takes the dominant styling; and
- the display must never render a fabricated placeholder such as `--:--` for a missing value.

Example when both values exist:

```text
Asr
Adhan    16:30
Jamaah   16:45   <- dominant
```

Example when only Adhan exists:

```text
Maghrib
Adhan    17:54   <- dominant
```

This omission rule is general, not Maghrib-specific. Missing optional timetable values do not indicate a failed board update.

## What the Display API Deliberately Omits

The presentation endpoint does **not** expose:

- provider name or provider-specific external IDs;
- discovery scope or catalogue partitions;
- upstream provider metadata;
- update error strings;
- persistence error strings;
- configuration operations; or
- refresh controls.

Those remain available through administrative/configuration APIs where appropriate.

## Runtime/Diagnostic API

The separate diagnostic endpoint remains:

```text
GET /api/masjidboard/status
```

It exposes richer runtime information such as provider identity, update-failed flags, attempt timestamps and diagnostic error messages. It is useful for WebUI administration and troubleshooting but should not be the physical display's normal data source.

## Failure Behaviour

Per-board failure state is reduced to the display concepts that matter visually:

```text
latest refresh succeeds
    -> status = current
    -> stale = false

latest refresh fails + last-known-good exists
    -> cached timetable remains visible
    -> status = stale
    -> stale = true
    -> last_successful_update retained

latest refresh fails + no cache exists
    -> status = unavailable
    -> timetable omitted
```

A failed update must not blank a previously usable timetable. Raw error text is intentionally excluded from the display endpoint.

## Unconfigured / Initial Refresh Behaviour

When MasjidBoard has not been configured:

```json
{
  "configured": false,
  "boards": []
}
```

When boards are configured but a particular board has not yet produced live or cached data, its display slot remains present in selection order with `status = unavailable`. This prevents the display layout from silently changing while an initial refresh is pending or unsuccessful.

## Display Lifecycle

The physical display frontend should remain simple:

```text
load display
    -> GET /api/masjidboard/display
    -> render configured board(s)
    -> periodically GET the same endpoint
```

It must not call configuration-changing endpoints.

Time formatting (for example 12-hour versus 24-hour presentation), exact typography, stale-warning styling and optional future layouts belong to the display/UI layer rather than the provider or runtime layers.

## Live Validation — 19 August 2026

Three real selected boards were returned in configured order with independent current timetable data. A deliberate failure of one board proved that the other two remained current while the failed board continued to expose its cached timetable as stale. Restoring the provider identity returned all three to current state.

The presentation API was subsequently live-tested with the same three real boards and returned the expected ordered compact display data.

## Independence from Audio Playback

MasjidBoard remains independent from the audio-streaming subsystem. A MasjidBoard display or refresh failure must not prevent audio playback, and audio availability must not determine whether MasjidBoard information can be displayed.
