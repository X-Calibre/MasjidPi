# MasjidBoard Implementation Status

**Status:** Active research implementation / pre-integration  
**Branch:** `research/masjidboard-live`

## Purpose

Record the current implemented state of MasjidBoard on the research branch and distinguish it from earlier design/research documents that describe ideas which have since evolved.

Where older MasjidBoard documents describe a different implementation choice, this document records the current working behaviour.

## Current Architecture

MasjidBoard is implemented as an independent subsystem inside the existing MasjidPi application process.

The audio/Listen subsystem and MasjidBoard remain functionally independent. A MasjidBoard provider, cache or display failure must not prevent audio playback from starting or continuing.

Current major components include:

- MasjidBoard Live discovery and hierarchy retrieval;
- persisted location scope;
- scoped Masjid catalogue;
- persisted selection of up to three Masjids;
- MasjidBoard Live Core timetable provider;
- normalized board/prayer model;
- per-board last-known-good cache;
- per-board runtime status and recovery;
- display presentation API;
- separate MasjidBoard configuration WebUI;
- separate read-only MasjidBoard display page.

The current implementation does **not** use a separate native SDL renderer or a separate `masjidboard.service`. The working display is browser-based and served by MasjidPi. Those older ideas remain historical design exploration rather than current implementation requirements.

## Discovery and Selection

The configuration UI supports:

- refreshing the MasjidBoard Live location hierarchy;
- selecting and persisting up to three locations;
- building a scoped Masjid catalogue from those locations;
- refreshing the scoped Masjid list;
- selecting up to three Masjids;
- preserving user-defined display column order;
- reordering/removing selected Masjids.

The audio streaming WebUI and MasjidBoard configuration UI are separate pages with navigation between Listen, configuration and display.

## Timetable Retrieval

The MasjidBoard Live Core board page is the working primary timetable source.

Selected boards are refreshed independently. A failure for one board does not invalidate the others.

Automatic timetable refresh currently runs every **30 minutes**. A manual **Refresh Timetables** action is also available from the configuration UI.

A configured board is also refreshed asynchronously at application startup so provider access never blocks the audio appliance startup path.

## Cache and Recovery

Each selected Masjid has an independent last-known-good timetable cache.

Verified runtime behaviour is:

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

This outage/recovery lifecycle has been verified both by automated tests and by a real runtime test in which `masjidboardlive.com` was temporarily made unreachable.

Unchanged timetable data is not rewritten to disk every 30 minutes. Cache persistence suppresses identical writes, with a periodic freshness checkpoint so persisted successful-retrieval metadata does not become indefinitely old.

## Configuration UI

The current configuration workflow is:

```text
Locations
    -> scoped Masjid list
    -> Selected Masjids
    -> timetable refresh/status
    -> display
```

Refresh actions are intentionally distinct:

- **Refresh Location List** — refresh the upstream hierarchy;
- **Refresh Masjid List** — refresh Masjids available within the selected locations;
- **Refresh Timetables** — refresh live timetable data for selected Masjids.

Per-board status exposes current/stale/unavailable state, cached-data use, last successful update, last attempt and update errors.

## Default Display Layout

The working read-only display is:

```text
/masjidboard.html
```

It supports one, two or three selected Masjids and preserves responsive sizing as the viewport changes.

The current layout is the **default layout**, not the intended permanent single-layout limitation. Alternative user-selectable layouts are deferred until the overall MasjidBoard appliance path is more production-ready.

The default layout shows:

- current local time and date;
- selected Masjid names;
- Fajr;
- Dhuhr on Saturday–Thursday;
- Jumu'ah replacing Dhuhr on Friday;
- Asr;
- Maghrib;
- Esha;
- per-board stale/unavailable indication;
- a per-Masjid countdown to the next visible timetable event.

Extended astronomical/calculation values remain available in the normalized data but are intentionally omitted from this initial layout.

## Prayer Presentation Rules

For normal prayers:

```text
Adhan
Jamaah
```

Adhan appears first. When Jamaah is supplied it receives stronger visual emphasis. If only one value exists, that available value receives the dominant styling. Missing values are omitted rather than rendered as placeholders.

## Friday / Jumu'ah Behaviour

Jumu'ah is hidden on non-Friday days and replaces Dhuhr on Friday.

Timed Jumu'ah items are displayed chronologically and source labels are preserved.

Examples include:

- Adhan;
- Lecture;
- Sunan;
- Khutbah;
- explicit Salaah/Jamaah where supplied.

Khutbah is **not** semantically relabelled as Salaah when an explicit Salaah value is absent. In that case Khutbah may receive the dominant visual treatment so the final supplied Friday event remains prominent, while retaining the correct `Khutbah` label.

## Next-Event Countdown

Each selected Masjid independently identifies its next visible timetable event.

The countdown appears immediately beneath that event's time using concise wording such as:

```text
in 15 min
in 1 hr
in 1 hr 15 min
now
```

The next event may therefore be Adhan, Jamaah or a visible Jumu'ah event such as Sunan, Lecture, Khutbah or explicit Salaah.

When the final visible event of the day has passed, the countdown rolls forward to the following day's first Fajr event rather than disappearing overnight.

Friday countdown behaviour has been manually validated across Jumu'ah event transitions.

For development/testing, the display accepts date/time overrides:

```text
/masjidboard.html?date=2026-08-21&time=12:10
```

These exist to test Friday and countdown transitions without changing the host clock.

## Current Validation

The branch currently has broad Go test coverage across:

- API endpoints;
- hierarchy/discovery;
- scoped catalogue and reconciliation;
- selection persistence;
- provider parsing and normalization;
- Jumu'ah extraction/fallback label recovery;
- cache persistence/write suppression;
- display presentation model;
- per-board runtime behaviour;
- stale-cache recovery.

The full `go test ./...` suite has been repeatedly validated during development.

Manual browser/runtime validation has covered:

- location and Masjid selection;
- three-board comparison display;
- Friday replacement of Dhuhr with Jumu'ah;
- Jumu'ah chronological event display;
- next-event countdown transitions;
- overnight countdown rollover;
- provider outage -> stale cached display -> recovery to current.

## Deferred / Future Display Work

Not required for the current default layout or branch integration:

- additional selectable display layouts;
- optional astronomical/calculation-time layouts;
- announcements/programmes/notices/media presentation;
- OLED or other small-display variants;
- layout preferences and theme variants beyond the current page behaviour.

These should consume the normalized MasjidBoard model/display API rather than duplicating provider logic.

## Remaining Pre-Integration Work

Before merging the research branch into the main development line, the focus should be production readiness rather than additional features:

1. Review installation/runtime paths so all MasjidBoard state directories/files are created and preserved correctly by fresh installs and upgrades.
2. Validate the combined Listen + Board application on Raspberry Pi hardware, including restart/reboot behaviour and resource usage.
3. Validate first-run behaviour when MasjidBoard has never been configured and when no cache exists.
4. Verify permissions/ownership for hierarchy, scope, catalogue, selection and per-board cache files under the production service account.
5. Review API/frontend error handling for user-friendly messages while retaining diagnostic detail in logs/status data.
6. Run a final full regression pass covering Listen/audio functionality to ensure MasjidBoard changes have not affected the stable playback appliance path.
7. Perform a final documentation cleanup and integration review, then prepare the branch for merge/release planning.

Additional display layouts and richer optional board content are explicitly **not blockers** for this integration milestone.
