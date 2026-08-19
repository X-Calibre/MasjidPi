# MasjidBoard Catalogue Memory Lifecycle

**Status:** Stage 3/4 architecture decision  
**Branch:** `research/masjidboard-live`

## Decision

The MasjidBoard catalogue is **disk-first, not permanently memory-resident**.

MasjidPi may be configured with **1–3 discovery locations**. Each location has an independently persisted last-known-good catalogue partition. The WebUI/API sees a merged, deduplicated catalogue view across those partitions.

The selected **1–3 boards** are separate runtime state and are therefore treated differently from the discovery catalogue.

## Runtime model

```text
Normal runtime
    -> load 1-3 persisted selected boards when configured
    -> construct one provider/runtime instance per selected board
    -> retrieve/cache each selected board independently
    -> catalogue partitions remain on disk

User opens board-selection UI
    -> load catalogue partitions from local storage
    -> merge + deduplicate records in memory
    -> browse/search the combined view
    -> select/reorder 1-3 boards
    -> persist the ordered selection
    -> release the merged catalogue when no longer required

Catalogue refresh
    -> load persisted partitions
    -> refresh each configured location independently
    -> validate + reconcile that location only
    -> replace only that successful partition
    -> keep another location's last-known-good partition on failure
    -> do not retain a process-lifetime catalogue cache
```

## Storage behaviour

`catalogue.Store.Load()` reads and validates the partitioned catalogue file on demand. It does not keep a process-lifetime in-memory copy.

`catalogue.Store.SavePartition()` is the normal refresh persistence primitive. It:

- validates the proposed location partition;
- preserves all other location partitions;
- avoids rewriting identical state;
- writes changed state through a temporary file;
- syncs the temporary file;
- atomically renames it into place; and
- leaves the existing persisted last-known-good state unchanged if validation or replacement fails.

`catalogue.Store.Save()` remains available for replacing the complete persisted partition set when configuration maintenance explicitly requires it.

This design favours low steady-state memory use and avoids unnecessary SD-card writes while providing failure isolation between configured discovery locations.

## Merged catalogue view

The persisted file retains provenance by location:

```text
catalogue state
    partition: Town A
        retrieved_at
        validated_at
        records...

    partition: Town B
        retrieved_at
        validated_at
        records...
```

The WebUI/API does not need separate copies of a mosque that appears in more than one location result. `catalogue.Merge()` deduplicates by the stable provider-neutral catalogue ID.

If the same board appears in multiple partitions:

- it appears once in the merged view;
- an active copy keeps the merged record active;
- the most recently seen copy supplies mutable metadata; and
- the earliest discovery timestamp is retained.

The merged `retrieved_at` and `validated_at` values are conservative: they represent the oldest partition timestamps. Per-location freshness remains authoritative for refresh scheduling.

## Selected-board state

Selected-board persistence has a different lifecycle. It is small runtime-critical state and is loaded at startup and retained in memory while MasjidBoard is operating.

The persisted selection is an ordered list of one to three boards once configured. Selection order is significant and may be used directly by the display layer.

Each selected board is independent. A retrieval failure for one selected board must not prevent the other selected boards from continuing to operate.

The resulting memory model is:

```text
catalogue partitions       disk-first; load on demand
merged catalogue view      temporary while browsing/configuring
selected boards (1-3)      load at startup; retain in memory
current board data         one independent in-memory state per selected board
last-known board data      one persisted cache per selected board
```

## Consequence for WebUI search

User-facing catalogue search operates on the locally persisted merged catalogue after loading it when the selection UI is opened. It does not require a remote FindMasjid request for each search operation or keystroke.

The user may choose and order up to three boards from the merged catalogue. The ordered selection is persisted separately and remains available at normal startup even when discovery is unavailable.
