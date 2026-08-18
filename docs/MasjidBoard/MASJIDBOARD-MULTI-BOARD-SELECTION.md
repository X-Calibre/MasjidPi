# MasjidBoard Multi-Board Selection

**Status:** Architecture decision / implementation in progress  
**Branch:** `research/masjidboard-live`

## Decision

MasjidPi will allow the user to select and display up to three MasjidBoard entries at the same time.

The primary use case is timetable comparison between nearby masjids. A user who cannot make a Jamaah at one masjid should be able to see the next practical nearby option without reopening the catalogue or changing configuration.

Example:

```text
1. Local masjid
2. Nearby masjid with a later Asr Jamaah
3. Nearby masjid with a later Esha Jamaah
```

Selection order is significant and must be preserved. The application must not sort the user's selected boards unless a future UI explicitly offers a separate comparison/sort view.

## Persisted Selection Model

The selection file stores an ordered list with a hard maximum of three boards:

```text
SelectionState
    Boards[0..3]

SelectedBoard
    CatalogueID
    Provider
    ExternalID
    Name
    TimeZoneOffsetMS
```

For MasjidBoard Live, a typical record is:

```text
CatalogueID       masjidboardlive:brits-jamia
Provider          masjidboardlive
ExternalID        brits-jamia
Name              Brits Jamia Masjid
TimeZoneOffsetMS  7200000
```

The selection deliberately contains the last-known display name and exact timezone offset so already-selected boards can operate without loading or refreshing the full catalogue.

The opaque MasjidBoard Live Premium `boardId` is not part of the required selection identity.

## Runtime Lifecycle

The full catalogue remains disk-first and is loaded only when required for browsing/searching or catalogue maintenance.

The much smaller selection state is different: it is needed for normal operation and should be loaded once at startup and retained in memory.

```text
startup
    -> load selection state
    -> 0 boards: MasjidBoard not configured
    -> 1 board: start one board provider
    -> 2 boards: start two board providers
    -> 3 boards: start three board providers
```

Each selected board is independent:

```text
Selected board 1 -> provider -> current/cached board data
Selected board 2 -> provider -> current/cached board data
Selected board 3 -> provider -> current/cached board data
```

A failure retrieving one selected board must not prevent the other selected boards from operating.

## Validation Rules

- Zero selected boards is valid.
- One, two, or three selected boards are valid.
- Four or more boards are invalid and must be rejected.
- Duplicate catalogue IDs are invalid.
- `CatalogueID` must match `provider + ":" + external_id`.
- Provider, external ID, and name are required.
- Selection order is preserved exactly as chosen by the user.
- Exact millisecond timezone offsets are retained; they must not be reduced to whole hours.

## Persistence

The selection should be persisted independently of the full catalogue, currently planned as:

```text
/var/lib/masjidpi/masjidboard_selection.json
```

Storage follows the existing MasjidPi appliance rules:

- load once at startup and cache in memory;
- do not rewrite unchanged selection state;
- write changed state through a temporary file;
- `fsync` the temporary file;
- atomically rename it into place; and
- keep the in-memory last-known-good selection unchanged if a save fails.

## Catalogue Independence

Normal board operation must not require the full discovery catalogue to be resident in memory.

```text
normal runtime
    -> selection state in memory
    -> current selected-board data in memory
    -> full catalogue not loaded

user opens board-selection WebUI
    -> load catalogue from local storage
    -> browse/search locally
    -> save updated ordered selection
    -> catalogue may then be discarded from memory
```

This keeps the normal runtime state small while still allowing a catalogue of hundreds of masjids to be searched when needed.

## Display Implications

The display layer is not defined by this storage work, but supporting up to three selected boards enables two useful future views:

1. individual timetable views for each selected masjid; and
2. a compact comparison view that highlights Jamaah times across the selected masjids, especially Asr and Esha.

The domain/provider layer should therefore keep each board independent rather than merging timetable values into one synthetic board.

## Cache Implications

Each selected board should eventually have its own last-known-good board-data cache. If one upstream board is temporarily unavailable, MasjidPi should be able to display its cached data while continuing to refresh and display the other selected boards normally.

## Implementation Status

The initial implementation provides:

- an ordered `SelectionState`;
- a hard maximum of three boards;
- duplicate and identity validation;
- conversion from provider-neutral catalogue records;
- load-once runtime selection storage;
- unchanged-save suppression; and
- atomic persistence.

The next integration step is to construct one provider/runtime instance per selected board at startup and later maintain one current/last-known-good data state per board.
