# MasjidBoard Multi-Board Selection

**Status:** Implemented and live-validated

## Decision

MasjidPi allows the user to select and display between one and three MasjidBoard entries at the same time once MasjidBoard has been configured.

The primary use case is timetable comparison between nearby masjids. Selection order is significant and is preserved exactly as chosen by the user.

## Configured State

```text
no selection file / zero-value state -> unconfigured
1 selected board                  -> valid
2 selected boards                 -> valid
3 selected boards                 -> valid
0 saved boards                    -> invalid
4+ selected boards                -> invalid
```

An explicit future "disable MasjidBoard" function must be modelled separately rather than as an empty configured selection.

## Persisted Selection

```text
SelectionState
    Boards[1..3]

SelectedBoard
    CatalogueID
    Provider
    ExternalID
    Name
    TimeZoneOffsetMS
    ShowDetailedJumuah
```

The selection contains enough identity, display-name and timezone information for already-selected boards to operate without loading or refreshing the discovery catalogue.

Validation rejects duplicate catalogue IDs and inconsistent `CatalogueID`, provider and external-ID combinations. Exact millisecond timezone offsets are preserved.

## Runtime Model

Each selected board has an independent provider, runtime state and last-known-good cache:

```text
Selected board 1 -> provider -> runtime -> cache
Selected board 2 -> provider -> runtime -> cache
Selected board 3 -> provider -> runtime -> cache
```

A failure retrieving one board does not prevent the others from operating.

The ordered selection can be changed through the configuration API while MasjidPi is running. The service reconfigures the selected-board runtimes without requiring the display to manage configuration.

## Persistence

Installed state is persisted independently of the discovery catalogue at:

```text
/var/lib/masjidpi/masjidboard_selection.json
```

Development state lives under `backend/data/`.

Persistence follows the MasjidPi appliance rules: validate first, suppress unchanged writes, write through a temporary file, sync, atomically rename, and retain the previous in-memory last-known-good state if persistence fails.

## Catalogue Independence

Normal selected-board operation does not require the full discovery catalogue in memory:

```text
normal runtime
    -> ordered selection retained in memory
    -> one current runtime state per selected board
    -> discovery catalogue remains disk-first
```

The catalogue is loaded only for browsing, selection and maintenance.

## Display Boundary

The timetable presentation is read-only. It receives configured boards in selection order and renders their timetable/status data. Searching, selecting, reordering and replacing boards are API/WebUI responsibilities. The Appliance touch overlay is limited to everyday Listen controls and theme selection.

Missing optional timetable values are not update failures. A board may legitimately omit individual Jamaah or Jumu'ah event times while the board itself remains `current`.

## Live Validation — 19 August 2026

The implementation was exercised against live MasjidBoard Live data using three real boards in Brits:

1. Jamiah Yusuf Darul Uloom Brits
2. Brits Jamia Masjid
3. Masjid Taqwa

The test confirmed:

- three boards can be selected and persisted together;
- selection order is preserved by the selection and status APIs;
- all three providers retrieve and normalise live timetable data independently;
- all three can simultaneously report `status = current`;
- one board can fail independently while the other two remain current;
- the failed board continues to expose its last-known-good timetable as `stale`; and
- a later successful refresh returns that board to `current`.

The 1–3 selected-board runtime path is therefore considered implemented and live-validated.
