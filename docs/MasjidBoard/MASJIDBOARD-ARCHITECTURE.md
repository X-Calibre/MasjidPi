# MasjidBoard Architecture

**Status:** Current architecture in MasjidPi v1.5.2

## Boundary

MasjidBoard is the Board capability inside the MasjidPi repository. It shares configuration, HTTP serving, persistence helpers and platform integration with Listen, while its provider, cache, refresh and display failures remain isolated from audio playback.

The installer supports:

- Listen
- Board
- Listen + Board

Only installed capabilities start their backend subsystems, expose their APIs and install their runtime dependencies/services.

## Data flow

```text
MasjidBoard Live hierarchy
    -> persisted global hierarchy
    -> selected location scope
    -> partitioned scoped catalogue
    -> ordered selected boards
    -> Core timetable retrieval
    -> optional Premium enrichment
    -> normalized Board model
    -> per-board last-known-good cache
    -> display presentation model
    -> standard or appliance browser renderer
```

Shared Daily Ayah/Hadith/Sunnah and Islamic Economic Indicators use independent clients and caches, then join the display presentation model.

## Packages

| Package | Responsibility |
|---|---|
| `hierarchy` | Global MasjidBoard Live location hierarchy |
| `scope` | Persisted selected geographic scope |
| `catalogue` | Scoped, partitioned and deduplicated board catalogue |
| `selection` | Ordered selected boards and saved display preferences |
| `provider` | Provider-neutral retrieval interface |
| `provider/masjidboardlive` | Core parsing, Premium enrichment and fallback extraction |
| `model` | Normalized provider-independent board data |
| `cache` | Per-board last-known-good persistence |
| `runtime` | Current/stale/unavailable result state |
| `service` | Refresh orchestration and presentation assembly |
| `display` | Stable JSON-facing presentation types |
| `dailycontent` | Shared Daily Ayah, Hadith and Sunnah retrieval/cache |
| `economic` | Jamiat Islamic Economic Indicator retrieval/cache |
| `maintenance` | Board persistent-state maintenance |

There is no separate `cmd/masjidboard` executable or scheduler service. Board runs inside the main MasjidPi backend; slideshow timing and content rotation belong to the frontend renderers.

## Persistence

Board state is stored under `/var/lib/masjidpi/masjidboard/`. Persistent concerns are separated:

- hierarchy;
- scoped catalogue partitions;
- selected boards and preferences;
- per-board normalized caches;
- shared daily Islamic content; and
- economic indicators.

Writes use atomic replacement and avoid replacing unchanged cache content. A failed refresh never replaces a usable last-known-good value.

## Refresh behavior

- Hierarchy and scoped catalogue refresh are administrative discovery operations.
- Selected timetables refresh independently on startup, periodically and on demand.
- Optional Premium failure does not invalidate a successful Core timetable.
- Daily Islamic content and economic indicators have independent refresh/cache policies.
- One failed board does not prevent other selected boards from remaining current.

## Display API

`GET /api/masjidboard/display` returns a read-only presentation model. It excludes discovery internals, cache paths and diagnostic error strings.

It contains:

- configured state and saved display preferences;
- selected boards in display order;
- per-board current/stale/unavailable state;
- prayers, Jumu'ah, special Dhuhr and astronomical times;
- supported announcements, programmes, notices, banking and new-moon data;
- optional daily Islamic content; and
- optional Islamic Economic Indicators.

The administrative APIs remain separate from the display response.

## Renderers

The common HTML page loads two presentation implementations:

- **standard** — responsive landscape TV/Monitor comparison;
- **appliance** — portrait 600 × 1024 slideshow with touch controls.

The system launcher detects the local hardware profile. Remote browsers can explicitly request the appliance preview for testing.

The frontend:

- calculates visible Friday state and next events from board-local time;
- formats times and dates;
- sanitizes upstream text before DOM insertion;
- orders and packs community cards;
- updates countdowns without rebuilding unchanged grids; and
- retains already-rendered data during temporary API failures.

## Failure isolation

- Listen operates without Board.
- Board operates without Listen, except that touch controls and source notifications are naturally unavailable.
- Core timetable data remains usable when optional enrichment fails.
- Cached boards remain visible as stale during provider failure.
- An unavailable board retains its selected position.
- Display-process failure is recovered by systemd without restarting the backend.

## Security and privacy

Upstream HTML is treated as untrusted data and converted to safe text. Named-person content is rendered only where its upstream purpose is established. Poster/media content remains excluded until retrieval, caching and presentation rules are defined.

## Related documentation

- [Implementation status](MASJIDBOARD-IMPLEMENTATION-STATUS.md)
- [Domain model](MASJIDBOARD-MODEL.md)
- [Display boundary](MASJIDBOARD-DISPLAY-BOUNDARY.md)
- [Display presentation](MASJIDBOARD-DISPLAY-LAYOUT.md)
- [Display runtime](MASJIDBOARD-DISPLAY-RUNTIME.md)
- [MasjidBoard Live contract](MASJIDBOARD-LIVE.md)
