# MasjidBoard Display and Configuration Boundary

**Status:** Architecture decision; backend status boundary implemented  
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
presentation-oriented display API
             |
             v
MasjidBoard display
    +-- READ ONLY
    +-- render selected board information
    +-- preserve configured board order
    +-- show per-board stale/update-failed indication
```

The display does not need to know how a board was discovered, selected, refreshed or configured.

## Current Runtime/Diagnostic API

The implemented read-only runtime endpoint is:

```text
GET /api/masjidboard/status
```

It returns the configured boards in user-selected order and exposes detailed runtime/provider-neutral board state including:

- stable catalogue/provider identity;
- display name and timezone offset;
- `current`, `stale` or `unavailable` status;
- cached-data and update-failed flags;
- attempt/success timestamps;
- diagnostic errors where applicable; and
- the currently displayable normalised board data.

This endpoint has been live-validated with three simultaneous boards and with an intentional single-board failure.

## Next Boundary: Presentation Model

`GET /api/masjidboard/status` is deliberately useful for runtime inspection and diagnostics, but the physical display should not be tightly coupled to its large raw normalised `Board` structure.

The next implementation step is to define a **small, stable, presentation-oriented display model/API** derived from runtime state.

That model should expose only what the display needs, for example:

- selected boards in configured order;
- board display name;
- the prayer/Jamaah times required by the chosen presentation;
- relevant Jumu'ah information;
- any deliberately selected astronomical times;
- a concise per-board stale/update-failed indicator; and
- enough freshness information to communicate stale data appropriately.

It should gracefully represent missing optional values. A missing individual Jamaah/Jumu'ah field is not itself an update failure.

The display model must remain provider-neutral. MasjidBoard Live-specific payload structure must not leak into the display frontend.

## Configuration APIs

Configuration and maintenance endpoints are separate from the display boundary. The backend now includes configuration operations for catalogue access/refresh and ordered selected-board state. The display must never call configuration-changing endpoints.

The display's normal lifecycle should remain conceptually simple:

```text
load display
    -> read presentation data
    -> render configured board(s)
    -> periodically read refreshed presentation data
```

No discovery or configuration state machine belongs in the display frontend.

## Failure Behaviour

Per-board failure state is part of presentation data:

```text
latest refresh succeeds
    -> current timetable

latest refresh fails + last-known-good exists
    -> cached timetable remains visible
    -> stale/update-failed indication shown for that board only

latest refresh fails + no cache exists
    -> board unavailable indication
```

A failed update must not blank a previously usable timetable.

## Live Validation — 19 August 2026

Three real selected boards were returned in configured order with independent current timetable data. A deliberate failure of one board subsequently proved that the other two remained current while the failed board continued to expose its cached timetable as stale. Restoring the provider identity returned all three to current state.

This validates the backend information needed by the future read-only display.

## Independence from Audio Playback

MasjidBoard remains independent from the audio-streaming subsystem. A MasjidBoard display or refresh failure must not prevent audio playback, and audio availability must not determine whether MasjidBoard information can be displayed.
