# MasjidBoard Last-Known-Good Board Cache

**Status:** Stage 4 implementation decision  
**Branch:** `research/masjidboard-live`

## Behavioural rule

Each selected MasjidBoard keeps its own persisted last-known-good timetable.

A failed refresh must never overwrite or remove usable timetable data from the previous successful refresh.

The user-facing distinction is:

```text
current
    latest refresh succeeded

stale
    latest refresh failed
    last-known-good timetable exists and is displayed
    UI must indicate that the last update failed

unavailable
    latest refresh failed
    no successful cached timetable exists yet
```

`stale` therefore does not mean that there is no timetable to display. It means that the displayed timetable is the last successfully validated data rather than the result of the latest refresh attempt.

## Per-board independence

For one to three configured boards, cache and update state remain independent:

```text
Board 1 -> own provider -> own cache -> own update status
Board 2 -> own provider -> own cache -> own update status
Board 3 -> own provider -> own cache -> own update status
```

A failed update for one board must not prevent the other selected boards from updating or displaying normally.

## Persisted cache contents

Each successful cache entry contains:

```text
CatalogueID
SuccessfulAt
Board
```

`SuccessfulAt` is a MasjidPi-owned timestamp recording when the provider result was successfully retrieved and accepted.

Failure metadata is deliberately not written into the last-known-good cache. A failed attempt must leave the successful cache unchanged. The runtime coordinator will hold the latest attempt/error status and expose it to the API/UI.

## Storage layout

One cache file is maintained per selected board. Cache filenames are derived from a SHA-256 hash of the stable catalogue ID rather than embedding the ID directly in the filename. This avoids filesystem-character and cross-platform naming problems while retaining the catalogue ID inside the persisted record for validation.

The cache directory is expected to live under MasjidPi persistent data storage, for example:

```text
/var/lib/masjidpi/masjidboard_cache/
```

The exact runtime path will be wired during application integration.

## Persistence rules

A cache entry is written only after a successful validated board refresh.

Changed entries are persisted by:

```text
validate entry
    -> encode JSON
    -> write temporary file
    -> fsync temporary file
    -> close
    -> atomic rename over previous entry
```

Identical entries are not rewritten.

A candidate that fails validation is rejected before storage and cannot replace the previous cache.

## Refresh/runtime sequence

The future runtime coordinator should implement:

```text
attempt provider refresh
        |
        +--> success
        |       -> validate result
        |       -> persist as new last-known-good entry
        |       -> display new data
        |       -> status = current
        |
        +--> failure
                -> load existing cache
                        |
                        +--> cache exists
                        |       -> display cached data
                        |       -> status = stale
                        |       -> expose "last update failed" flag/error
                        |
                        +--> no cache
                                -> status = unavailable
```

The cache itself does not decide `current`, `stale`, or `unavailable`; those are runtime states derived from the latest refresh attempt plus cache availability.

## SD-card behaviour

This cache follows the existing MasjidPi appliance write policy:

- no write on failed refresh;
- no write for an identical successful entry;
- one atomic replacement only when last-known-good data actually changes; and
- no raw upstream HTML or provider payload is persisted.
