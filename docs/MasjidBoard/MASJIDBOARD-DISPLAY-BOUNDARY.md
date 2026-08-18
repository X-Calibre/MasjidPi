# MasjidBoard Display and Configuration Boundary

**Status:** Architecture decision  
**Branch:** `research/masjidboard-live`

## Decision

The MasjidBoard display is a **read-only presentation surface**.

It exists only to display MasjidBoard information for the user's configured one to three selected masjids. It must not contain configuration controls, discovery/search workflows, board selection controls, preferences, or administrative actions.

All configuration and preferences are handled through the MasjidPi API and/or the administrative WebUI.

## Responsibility Split

```text
MasjidPi WebUI / configuration API
    |
    +-- discover/search MasjidBoards
    +-- select 1-3 boards
    +-- reorder/replace selected boards
    +-- configure display preferences
    +-- inspect diagnostics / administrative status
    |
    v
MasjidBoard backend service
    |
    +-- persisted selection
    +-- provider construction
    +-- live refresh
    +-- last-known-good cache
    +-- current / stale / unavailable runtime status
    |
    v
MasjidBoard display API
    |
    v
MasjidBoard display
    +-- READ ONLY
    +-- render selected board information
    +-- show stale/update-failed warning when required
```

The display does not need to know how a board was discovered, selected, or configured.

## Display-Facing API

The first display-facing boundary is:

```text
GET /api/masjidboard/status
```

This endpoint is read-only and provides the current runtime view of the configured boards in user-selected order.

It exposes, per selected board:

- stable catalogue/provider identity;
- display name and timezone offset;
- runtime status (`current`, `stale`, or `unavailable`);
- whether cached data is currently being displayed;
- whether the last live update failed;
- last attempt timestamp;
- last successful update timestamp;
- update/persistence diagnostic messages where applicable; and
- the currently displayable normalised board data when available.

The `stale` case is especially important:

```text
latest refresh fails
    + last-known-good timetable exists
        -> display cached timetable
        -> using_cached_data = true
        -> update_failed = true
        -> status = stale
```

A failed update therefore does not blank a previously usable display.

## Configuration-Facing APIs

Configuration endpoints are separate from the display-facing endpoint. Planned responsibilities include concepts such as:

```text
GET/refresh catalogue
GET/PUT selected boards
GET/PUT display preferences
```

Exact routes are implementation details to be defined as the configuration API is built.

The display itself must not call configuration-changing endpoints.

## UI Consequence

The MasjidBoard display frontend should remain deliberately simple. Its normal operation should be:

```text
load display page
    -> GET /api/masjidboard/status
    -> render configured board(s)
    -> periodically read refreshed status
```

No configuration state machine is required in the display frontend.

The administrative WebUI remains the place where users search for masjids, choose and order up to three boards, and change display preferences.

## Independence from Audio Playback

This boundary does not change the existing architectural rule that MasjidBoard remains independent from the audio-streaming subsystem.

A MasjidBoard display or refresh failure must not prevent audio playback from operating, and audio availability must not determine whether MasjidBoard information can be displayed.
