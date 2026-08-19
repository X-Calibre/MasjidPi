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

Astronomical values remain optional. The presentation model currently carries the normalised astronomical fields so the later screen design can choose which are actually rendered without making the display consume the full internal `Board` object.

Clock times have an explicit display JSON shape:

```json
{
  "hour": 16,
  "minute": 45
}
```

The display model owns this JSON contract independently of the internal Go domain structs.

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

Time formatting (for example 12-hour versus 24-hour presentation), visual layout, stale-warning styling and which optional astronomical fields are shown belong to the later display/UI design rather than the provider or runtime layers.

## Live Validation — 19 August 2026

Three real selected boards were returned in configured order with independent current timetable data. A deliberate failure of one board proved that the other two remained current while the failed board continued to expose its cached timetable as stale. Restoring the provider identity returned all three to current state.

The new presentation model is therefore built on already-live-validated runtime semantics rather than a speculative provider contract.

## Independence from Audio Playback

MasjidBoard remains independent from the audio-streaming subsystem. A MasjidBoard display or refresh failure must not prevent audio playback, and audio availability must not determine whether MasjidBoard information can be displayed.
