# MasjidBoard Last-Known-Good Board Cache

**Status:** Implemented and live-validated  
**Branch:** `research/masjidboard-live`

## Behavioural Rule

Each selected MasjidBoard keeps its own persisted last-known-good timetable. A failed refresh must never overwrite or remove usable timetable data from the previous successful refresh.

The runtime states are:

```text
current
    latest refresh succeeded
    using_cached_data = false
    update_failed = false

stale
    latest refresh failed
    last-known-good timetable exists and is displayed
    using_cached_data = true
    update_failed = true

unavailable
    latest refresh failed
    no successful cached timetable exists yet
```

`stale` means usable data remains available, but it came from the last successful refresh rather than the latest attempt. The display must indicate that the latest update failed.

## Per-Board Independence

For one to three configured boards:

```text
Board 1 -> own provider -> own runtime state -> own cache
Board 2 -> own provider -> own runtime state -> own cache
Board 3 -> own provider -> own runtime state -> own cache
```

A failure for one board does not prevent the others from refreshing or displaying normally.

## Persisted Cache

Each successful cache entry contains:

```text
CatalogueID
SuccessfulAt
Board
```

Cache filenames are derived from a SHA-256 hash of the stable catalogue ID. The catalogue ID remains inside the persisted record for validation.

Installed cache storage is:

```text
/var/lib/masjidpi/masjidboard_cache/
```

Development cache storage lives under `backend/data/`.

Failure metadata is not written into the last-known-good cache. Failed attempts leave the successful cache unchanged. The runtime owns latest-attempt/error state and exposes it through the API.

## Refresh Sequence

```text
attempt live refresh
        |
        +--> success
        |       -> validate result
        |       -> persist new last-known-good data
        |       -> display new data
        |       -> status = current
        |
        +--> failure
                -> existing last-known-good data?
                        |
                        +--> yes
                        |       -> display cached data
                        |       -> status = stale
                        |       -> using_cached_data = true
                        |       -> update_failed = true
                        |
                        +--> no
                                -> status = unavailable
```

The cache itself does not decide `current`, `stale` or `unavailable`; those are runtime states derived from the latest refresh attempt and cache availability.

## Persistence / SD-Card Rules

- write only after a successful validated board refresh;
- never write on a failed refresh;
- suppress an identical successful write;
- use temporary-file + sync + atomic rename for changed data;
- reject invalid candidates before storage; and
- persist normalised board data, not raw upstream HTML/provider payloads.

## Live Failure and Recovery Validation — 19 August 2026

The fallback path was deliberately exercised with three selected real boards. Brits Jamia Masjid was temporarily reconfigured to an invalid MasjidBoard Live external ID while its previous successful cache was retained under the test identity.

Observed result:

```text
Jamiah Yusuf Darul Uloom Brits
    status = current

Brits Jamia test entry
    live request failed with 404
    status = stale
    using_cached_data = true
    update_failed = true
    last_successful_update preserved
    previously cached timetable still returned

Masjid Taqwa
    status = current
```

After restoring the genuine selection, the next successful live refresh returned Brits Jamia Masjid to:

```text
status = current
using_cached_data = false
update_failed = false
```

This validates both failure isolation and recovery. The last-known-good requirement is therefore considered implemented and live-validated.
