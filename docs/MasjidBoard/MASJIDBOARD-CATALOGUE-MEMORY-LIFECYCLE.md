# MasjidBoard Catalogue Memory Lifecycle

**Status:** Stage 3/4 architecture decision  
**Branch:** `research/masjidboard-live`

## Decision

The MasjidBoard catalogue is **disk-first, not permanently memory-resident**.

The catalogue is selection-time data. Normal MasjidBoard runtime does not need the complete catalogue after the user has selected the board or boards to display.

## Runtime model

```text
Normal runtime
    -> load persisted selected board(s)
    -> construct selected provider(s)
    -> retrieve/cache selected board data
    -> full catalogue remains on disk

User opens board-selection UI
    -> load catalogue from local storage
    -> keep/use it while browsing and searching
    -> release it when no longer required

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

Selected-board persistence has a different lifecycle. It is small runtime-critical state and should be loaded at startup and may remain in memory while MasjidBoard is operating.

The resulting memory model is:

```text
full catalogue          disk-first; load on demand
selected board(s)       load at startup; retain in memory
current board data      retain in memory while operating
last-known board data   persisted cache for recovery/offline use
```

## Consequence for WebUI search

User-facing catalogue search should operate on the locally persisted catalogue after loading it when the selection UI is opened. It should not require a remote FindMasjid request for each search operation or keystroke.

The catalogue may be released from memory after the selection/browsing workflow is complete.
