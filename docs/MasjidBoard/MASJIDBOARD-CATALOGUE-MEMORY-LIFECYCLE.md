# MasjidBoard Catalogue Memory Lifecycle

**Status:** Implemented architecture

## Decision

The MasjidBoard catalogue is **disk-first, not permanently memory-resident**.

MasjidPi may be configured with **one to three discovery locations**. Each location has an independently persisted last-known-good catalogue partition. The WebUI/API sees a merged, deduplicated catalogue view across those partitions.

The selected **one to three boards** are separate runtime-critical state.

## Runtime Model

```text
Normal runtime
    -> load persisted ordered board selection
    -> construct one provider/runtime per selected board
    -> retrieve/cache each selected board independently
    -> catalogue partitions remain on disk

User opens board-selection UI
    -> load catalogue partitions
    -> merge + deduplicate in memory
    -> browse/search locally
    -> select/reorder 1-3 boards
    -> persist selection
    -> release merged catalogue when no longer required

Catalogue maintenance
    -> inspect configured 1-3 discovery locations
    -> refresh each due/requested location independently
    -> validate + reconcile only that partition
    -> replace only successful partitions
    -> retain another location's last-known-good partition on failure
```

## Discovery Hierarchy vs Scoped Catalogue

The global FindMasjid hierarchy and the scoped board catalogue have different purposes.

The persisted hierarchy tracks which countries, regions/provinces and towns/cities are available for configuration. The scoped catalogue does **not** mirror all boards globally. It contains board records only for the user's configured one to three locations.

This keeps board data small while still allowing configuration choices to follow changes in the upstream hierarchy.

## Storage Behaviour

`catalogue.Store.Load()` reads and validates the partitioned catalogue on demand.

`catalogue.Store.SavePartition()` validates a location partition, preserves all other partitions, suppresses identical writes and atomically replaces changed state. A failed candidate does not replace the persisted last-known-good partition.

The merged catalogue deduplicates by stable provider-neutral catalogue ID. If a board appears in more than one selected discovery location, it appears once to the WebUI/API while partition provenance remains persisted.

## Selected-Board State

Selected-board persistence is loaded at startup and retained in memory because it is required for normal operation.

```text
catalogue partitions       disk-first; load on demand
merged catalogue view      temporary for configuration
selected boards (1-3)      retained in memory
current board data         independent state per selected board
last-known-good board data persisted independently per board
```

The selected-board runtime does not depend on discovery being available after configuration.

## Refresh Policy

Scoped catalogue refresh is deliberately low-frequency:

- automatic refresh when due, at most once every seven days based on persisted freshness timestamps; or
- immediate refresh when explicitly requested through API/WebUI.

Opening the configuration UI does not itself force a remote refresh. Individual selected-board timetable refreshes are independent and run on their own provider schedule.

## Failure Isolation

Failure is isolated at both catalogue-partition and selected-board levels:

- failure refreshing one discovery location does not discard other location partitions;
- failure refreshing a selected board does not affect other selected boards; and
- a selected board with last-known-good timetable data remains displayable as stale after a live refresh failure.

## Live Validation

The provider hierarchy was exercised against the live FindMasjid service across its reported countries and region/city structures during development, including upstream irregularities such as blank region values and alternate response shapes.

The selected-board path was then live-tested with three real boards from the persisted local catalogue. All three were retrieved concurrently and returned in configured order. A deliberate failure of one board demonstrated independent last-known-good fallback while the other two remained current.

## Consequence for WebUI Search

Ordinary user-facing board search/filtering operates against the locally persisted merged catalogue. It does not perform a remote FindMasjid request for each search operation or keystroke.

The WebUI/API owns discovery-location configuration and the ordered one-to-three board selection. The physical MasjidBoard display remains presentation-only.
