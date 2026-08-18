# MasjidBoard Catalogue Memory Lifecycle

**Status:** Stage 3/4 architecture decision  
**Branch:** `research/masjidboard-live`

## Decision

The MasjidBoard catalogue is **disk-first, not permanently memory-resident**.

The catalogue is selection-time data. Normal MasjidBoard runtime does not need the complete catalogue after the user has selected the board or boards to display.

MasjidPi supports an ordered selection of up to three boards. Those selected-board records are runtime state and are therefore treated differently from the full catalogue.

## Runtime model

```text
Normal runtime
    -> load 0-3 persisted selected boards
    -> construct one provider/runtime instance per selected board
    -> retrieve/cache each selected board independently
    -> full catalogue remains on disk

User opens board-selection UI
    -> load catalogue from local storage
    -> keep/use it while browsing and searching
    -> select/reorder up to three boards
    -> persist the ordered selection
    -> release the catalogue when no longer required

Catalogue refresh
    -> load current persisted catalogue
    -> build candidate catalogue in memory
    -> validate and reconcile candidate
    -> atomically persist changed catalogue
    -> do not retain a process-lifetime catalogue cache
```

## Storage behaviour

`catalogue.Store.Load()` reads and validates the catalogue file on demand. It does not keep a process-lifetime in-memory copy.

`catalogue.Store.Save()`:

- validates the proposed catalogue;
- reads the currently persisted catalogue for comparison;
- avoids rewriting identical content;
- writes changed data through a temporary file;
- syncs the temporary file;
- atomically renames it into place; and
- leaves the existing persisted last-known-good catalogue unchanged if validation or the replacement operation fails.

This design favours low steady-state memory use without compromising the more important SD-card requirement of avoiding unnecessary writes.

## Selected-board state

Selected-board persistence has a different lifecycle. It is small runtime-critical state and is loaded once at startup and retained in memory while MasjidBoard is operating.

The persisted selection is an ordered list of zero to three boards. Selection order is significant and may be used directly by the future display/UI layer.

Each selected board is independent. A retrieval failure for one selected board must not prevent the other selected boards from continuing to operate.

The resulting memory model is:

```text
full catalogue          disk-first; load on demand
selected boards (0-3)   load at startup; retain in memory
current board data      one independent in-memory state per selected board
last-known board data   one persisted cache per selected board for recovery/offline use
```

## Consequence for WebUI search

User-facing catalogue search should operate on the locally persisted catalogue after loading it when the selection UI is opened. It should not require a remote FindMasjid request for each search operation or keystroke.

The user may choose and order up to three boards from that local catalogue. The ordered selection is persisted separately from the catalogue and remains available at normal startup even when discovery is unavailable.

The catalogue may be released from memory after the selection/browsing workflow is complete.
