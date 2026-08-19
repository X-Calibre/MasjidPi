# MasjidBoard Catalogue and Persistence Design

**Status:** Implemented core architecture; live integration validated  
**Branch:** `research/masjidboard-live`

## Purpose

Define the stable MasjidPi catalogue identity, geographic discovery scope, hierarchy persistence, multi-board selection, refresh behaviour and last-known-good strategy for MasjidBoard discovery/runtime.

```text
FindMasjid hierarchy
        |
        v
persisted available geography
        |
        v
configured 1-3 discovery locations
        |
        v
partitioned scoped local catalogue
        |
        +--> merged WebUI/API browse/search view
        |
        +--> persisted ordered 1-3 selected boards
        |
        v
independent Core providers/runtimes
        |
        v
normalised current/stale board results
```

Discovery hierarchy, configured scope, catalogue persistence, board selection, individual timetable retrieval and display are separate concerns.

## Discovery Scope

MasjidPi does **not** maintain a worldwide mirror of all MasjidBoard timetable records.

The user may configure **one to three locations**. Each location is expressed as:

```text
Country
  -> Province / Region
      -> Town / City
```

Country and city are required. Region may be blank because the upstream hierarchy contains entries without a province/region value.

Multiple locations support real-world border cases where a nearby masjid may belong to an adjacent town/city. The maximum is three locations.

The persisted global hierarchy is maintained separately so the WebUI can know which countries, regions and towns/cities are currently available. Board catalogue data itself remains scoped to the configured locations.

## Stable Catalogue Identity

For MasjidBoard Live, the working external identity is the public `web_url` slug:

```text
provider    = masjidboardlive
external_id = brits-jamia
catalogue_id = masjidboardlive:brits-jamia
```

The opaque Premium `boardId` is not the catalogue key. `MBL_ID` remains provider metadata.

## Catalogue Record

```text
CatalogueRecord
    ID
    Provider
    ExternalID
    Name
    City
    Country
    Region
    TimeZoneOffsetMS
    ProviderMetadata
    DiscoveredAt
    LastSeenAt
    Status
```

FindMasjid timetable summary fields are not authoritative selected-board data and do not need to drive the runtime timetable. Selected boards retrieve their normalised timetable through the Core provider.

## Partitioned Catalogue

Each configured discovery location owns an independent persisted catalogue partition. The local browse/search view merges the one to three partitions and deduplicates records by stable catalogue ID.

A successful refresh replaces only the relevant partition. A failed location refresh leaves that partition's last-known-good state untouched and does not affect other locations.

This design provides both multi-location discovery and failure isolation without keeping a worldwide board mirror.

## Multi-Board Selection

Selection is persisted independently from discovery scope/catalogue:

```text
unconfigured -> no persisted board selection
configured   -> exactly 1-3 ordered selected boards
```

Each selected board persists:

```text
CatalogueID
Provider
ExternalID
Name
TimeZoneOffsetMS
```

Selection order is preserved and maps to display order. The selection can be reconfigured through API/WebUI while the service is running.

## Startup and Runtime

An already-configured appliance does not depend on discovery availability:

```text
load selected-board state
        |
        +--> configured
        |       -> construct 1-3 independent runtimes
        |       -> refresh each independently
        |       -> use per-board last-known-good cache on failure
        |
        +--> unconfigured
                -> wait for API/WebUI configuration
```

## Persistent State

Installed appliance state includes conceptually:

```text
/var/lib/masjidpi/masjidboard_scope.json
/var/lib/masjidpi/masjidboard_hierarchy.json
/var/lib/masjidpi/masjidboard_catalogue.json
/var/lib/masjidpi/masjidboard_selection.json
/var/lib/masjidpi/masjidboard_cache/
```

Development equivalents live under `backend/data/`.

These remain separate because their lifecycles and write frequencies differ.

## Refresh Policy

### Hierarchy

The hierarchy is persisted independently so available countries, regions and towns/cities can be refreshed without downloading every board timetable globally.

### Scoped catalogue

Automatic scoped catalogue refresh occurs at most once every seven days based on persisted MasjidPi freshness timestamps. The user may explicitly request an immediate refresh through API/WebUI. Opening configuration does not itself force a refresh.

### Selected timetables

Selected-board timetable refresh is independent of catalogue maintenance and occurs on the provider/runtime schedule.

## Last-Known-Good Behaviour

Catalogue partition failure:

```text
retrieve location
    -> validate/reconcile candidate
        -> success: atomically replace that partition if changed
        -> failure: retain existing partition
```

Selected-board failure:

```text
live refresh succeeds
    -> current
    -> persist accepted timetable as last-known-good

live refresh fails + cache exists
    -> stale
    -> using_cached_data = true
    -> update_failed = true
    -> keep displaying previous timetable

live refresh fails + no cache
    -> unavailable
```

Failures are isolated per selected board.

## Timezone Handling

`time_zone_milli` is preserved exactly and is not rounded to whole hours. MasjidPi does not invent an IANA timezone where upstream only provides a fixed offset.

## Configuration / Display Boundary

All configuration belongs to API/WebUI:

- configure one to three discovery locations;
- refresh/browse the local catalogue;
- select and order one to three boards; and
- manage future display preferences.

The physical MasjidBoard display is read-only. It consumes presentation data only and never performs discovery, catalogue maintenance or configuration.

`GET /api/masjidboard/status` currently exposes detailed runtime/diagnostic board state. The next implementation boundary is a smaller stable presentation-oriented display model/API rather than coupling the display directly to the complete raw normalised board structure.

## Independence from Audio

MasjidBoard remains independent from the LiveMasjid audio-stream catalogue and playback subsystem. Failure in either subsystem must not prevent the other from operating.

## Live Validation — 19 August 2026

Development exercised the live FindMasjid hierarchy, including alternate upstream response shapes and blank region values.

The complete selected-board path was then tested with three real Brits boards. The tests confirmed:

- ordered three-board selection and persistence;
- independent live provider retrieval and normalisation;
- simultaneous `current` state for all three boards;
- optional missing timetable fields do not invalidate a board;
- intentional failure of one board leaves the other two current;
- the failed board returns its persisted timetable as `stale` with cached/update-failed flags; and
- successful recovery returns it to `current`.

The catalogue/selection/runtime reliability path is therefore considered sufficiently validated to proceed to the display presentation model.

## Current Decisions

1. `provider + external_id` is the stable catalogue identity; MasjidBoard Live uses the public `web_url` slug.
2. Discovery hierarchy is persisted separately from scoped board catalogue data.
3. The user may configure one to three discovery locations.
4. Catalogue data is partitioned per configured location and merged/deduplicated for browsing.
5. A configured selection contains exactly one to three ordered boards.
6. Selected boards operate independently of catalogue availability after configuration.
7. Catalogue and board caches use last-known-good transactional persistence.
8. Automatic scoped catalogue refresh is weekly; manual refresh may run immediately.
9. The catalogue is disk-first and loaded on demand.
10. Exact millisecond timezone offsets are preserved.
11. Per-board runtime states are `current`, `stale` or `unavailable`.
12. Display is read-only; all configuration belongs to WebUI/API.
13. MasjidBoard remains independent from audio playback.

## Next Implementation Boundary

Define and test the provider-neutral **display presentation model/API** that transforms the validated runtime state into the small stable data contract required by the physical display.
