# MasjidBoard Catalogue and Persistence Design

**Status:** Stage 3 research / architecture  
**Branch:** `research/masjidboard-live`

## Purpose

Define the stable MasjidPi catalogue record, selected-board persistence model, refresh/reconciliation behaviour, and last-known-good strategy for MasjidBoard Live discovery.

This design sits between the already validated FindMasjid discovery client and the already validated Core provider path.

```text
FindMasjid discovery
        |
        v
normalised catalogue
        |
        +--> user search / selection
        |
        +--> persisted selected board
        |
        v
Core provider
        |
        v
normalised Board
```

Discovery, catalogue persistence, selected-board persistence, and individual-board retrieval remain separate concerns.

## Design Goals

The catalogue subsystem should:

- use stable public MasjidBoard Live identifiers rather than opaque implementation IDs;
- allow a selected board to continue working even if discovery is temporarily unavailable;
- tolerate board renames and metadata changes;
- tolerate boards disappearing temporarily or permanently;
- avoid unnecessary SD-card writes;
- persist only changed data and use atomic replacement;
- distinguish catalogue freshness from individual-board timetable freshness;
- avoid treating upstream `last_updated` as catalogue freshness; and
- remain provider-specific at the discovery boundary while exposing stable MasjidPi catalogue records to the rest of the application.

## Stable Catalogue Identity

### Primary key

The current working key is the public MasjidBoard Live `web_url` slug.

Example:

```text
brits-jamia
```

This is the same public identifier used by:

```text
https://masjidboardlive.com/boards/?brits-jamia
https://premium.masjidboardlive.com/v2/?mid=brits-jamia
```

The opaque Premium `boardId` must not become the catalogue key.

`MBL_ID` should also remain upstream metadata rather than the MasjidPi primary key until authoritative evidence establishes stronger semantics or stability guarantees.

### MasjidPi catalogue ID

The provider-neutral catalogue ID should be namespaced so another provider can coexist later:

```text
masjidboardlive:brits-jamia
```

Conceptually:

```text
catalogue_id = provider + ":" + external_id
```

For MasjidBoard Live:

```text
provider    = "masjidboardlive"
external_id = web_url
```

This prevents accidental collisions if another timetable provider uses the same public slug.

## Proposed Catalogue Record

The persisted MasjidPi catalogue should contain only identity, geography, selection/search metadata, provider metadata needed for retrieval, and discovery-state information.

Suggested logical model:

```text
CatalogueRecord
    ID                  masjidboardlive:<web_url>
    Provider            masjidboardlive
    ExternalID          <web_url>
    Name                upstream masjid name
    City                upstream city
    Country             discovery hierarchy value
    Province            discovery hierarchy value
    TimeZoneOffsetMS    exact upstream offset
    MBLID               optional opaque upstream metadata
    LadiesFacility      optional discovery metadata
    DiscoveredAt        MasjidPi timestamp
    LastSeenAt          MasjidPi timestamp
    Status              active | missing | unavailable
```

The first implementation does not need to persist timetable summary fields such as Fajr Jamaah or sunset in the catalogue because the selected board's authoritative timetable is obtained from Core. Those fields can remain available on provider-level discovery results where useful for validation or UI previews.

Likewise, upstream `last_updated` may be retained as optional provider metadata but must not drive catalogue freshness or removal decisions.

## Selected-Board Persistence

The user's selected MasjidBoard should be persisted independently from the catalogue itself.

Suggested state:

```text
SelectedBoard
    CatalogueID         masjidboardlive:<web_url>
    Provider            masjidboardlive
    ExternalID          <web_url>
    Name                last-known display name
    TimeZoneOffsetMS    last-known exact offset
```

The important fields are the provider and external ID. The name and timezone are last-known metadata that allow already-selected board retrieval to continue if catalogue refresh is unavailable.

The selected-board state must not contain the opaque Premium `boardId` as a required identifier.

### Startup behaviour

On startup:

```text
load selected-board state
        |
        +--> valid selection exists
        |       |
        |       v
        |   construct provider from persisted selection
        |       |
        |       v
        |   retrieve selected Core board
        |
        +--> no selection
                -> wait for user selection
```

A successful catalogue refresh is therefore not a prerequisite for an already-configured appliance to start displaying its selected masjid.

## Catalogue Persistence

Proposed persistent files:

```text
/var/lib/masjidpi/masjidboard_catalogue.json
/var/lib/masjidpi/masjidboard_selection.json
```

These names are provisional until implementation, but the catalogue and selection should remain separate files because they have different lifecycles and write frequencies.

Persistence should follow the existing MasjidPi SD-card policy:

- load once and cache in memory;
- do not rewrite unchanged state;
- compare newly generated catalogue content against current state;
- write changed state through a temporary file;
- atomically rename the temporary file into place; and
- avoid persisting raw discovery HTTP responses.

## Catalogue Freshness

MasjidPi-owned freshness metadata must be separate from MasjidBoard Live's `last_updated` field.

Recommended catalogue-level timestamps:

```text
retrieved_at
    when a discovery refresh request successfully completed

validated_at
    when the returned catalogue parsed and passed validation

last_seen_at
    per-record timestamp recording the most recent successful discovery
```

`last_updated` remains optional upstream metadata only.

A catalogue may therefore be considered operationally fresh even when a board's upstream `last_updated` value is old or blank.

## Refresh Strategy

The first implementation should favour conservative refresh behaviour over frequent polling.

Recommended initial model:

```text
startup
    -> load last-known-good catalogue immediately
    -> do not block selected-board operation on refresh

manual refresh
    -> available to user / administrative API

periodic refresh
    -> low frequency, exact interval to be decided during Stage 4
```

A daily or weekly catalogue refresh is more appropriate than polling every few minutes because board discovery changes much less frequently than individual prayer times.

The individual Core board provider may refresh independently at a shorter interval.

## Last-Known-Good Behaviour

A failed discovery refresh must never replace a valid persisted catalogue with an empty or partially parsed result.

Refresh should conceptually be transactional:

```text
retrieve full candidate catalogue
        |
        v
parse / validate candidate
        |
        +--> failure
        |       -> keep current catalogue unchanged
        |
        +--> success
                -> reconcile candidate with current catalogue
                -> persist only if changed
```

The in-memory catalogue should also remain the current last-known-good version until a complete replacement has passed validation.

## Reconciliation Rules

### Existing record seen again

If the same `provider + external_id` appears again:

- retain the same MasjidPi catalogue ID;
- update mutable metadata such as name, city, timezone, `MBL_ID`, and facilities;
- update `last_seen_at`; and
- do not treat a name change as a new masjid.

### Rename

A board whose `web_url` remains unchanged but whose display name changes is considered the same catalogue record.

The selected-board state may update its cached display name after a successful catalogue refresh, but the selection remains valid regardless.

### `web_url` change

If the upstream `web_url` changes, it appears as a new external identity under the current evidence.

MasjidPi must not automatically infer that the new slug is the same masjid based only on matching name/city or `MBL_ID`.

Possible migration heuristics can be researched later, but automatic identity merging would risk selecting the wrong board.

### Temporarily missing record

A board absent from one refresh should not be deleted immediately.

Recommended state transition:

```text
active
  |
  | absent from successful full refresh
  v
missing
```

The record remains available as last-known metadata, especially if it is currently selected.

Permanent removal should require a conservative policy such as repeated successful refreshes over time before pruning. Exact thresholds belong to Stage 4.

### Selected board disappears from catalogue

If the selected board is missing from discovery but its Core URL still works:

- continue operating the selected board;
- retain the persisted selection;
- mark its catalogue record missing; and
- do not force the user to reselect.

If both discovery and Core retrieval fail, MasjidPi should fall back to the last-known-good board data cache once that layer is implemented.

## Availability Status

Initial catalogue status values should be deliberately simple:

```text
active
missing
unavailable
```

`active` means present in the latest successful discovery result.

`missing` means absent from one or more successful discovery results but retained locally.

`unavailable` should only be used when explicit validation proves that the board cannot currently be retrieved; a transient network error must not automatically mark a board unavailable.

Premium capability should be represented separately from availability if/when it becomes part of catalogue enrichment.

## Discovery and Provider Cross-Validation

The live Brits validation established that:

```text
FindMasjid MBL_ID = MBL11517PRP
Core mbl_number    = MBL11517PRP
```

This provides a useful optional consistency check after selecting a catalogue record.

However, `MBL_ID` remains opaque. A mismatch should initially be treated as a warning/diagnostic signal rather than as proof that either source is invalid until broader behaviour is understood.

The authoritative retrieval key remains `web_url`.

## Timezone Handling

`time_zone_milli` must be persisted exactly as supplied and converted without reducing it to whole hours.

The catalogue-to-provider handoff already supports exact offsets such as:

```text
GMT+02:00
GMT+05:30
GMT+05:45
GMT+09:30
GMT-03:30
```

This exact offset should be preserved in selected-board state so Core retrieval remains correctly contextualised without requiring a catalogue refresh first.

Longer term, an IANA timezone name would be preferable where MasjidBoard Live exposes one reliably, because fixed offsets do not encode daylight-saving rules. The current discovery contract only gives a millisecond offset, so MasjidPi must not invent an IANA zone.

## Search and User Selection

The catalogue should support user-facing search over at least:

- masjid name;
- city;
- province/region; and
- country.

Search should operate against the local catalogue after it has been generated rather than making a remote request for every keystroke.

User selection should persist the stable catalogue ID/provider/external ID immediately after confirmation.

## Separation from Audio Catalogue

The MasjidBoard catalogue is separate from the existing LiveMasjid audio-stream catalogue.

They may eventually be correlated in the application UI, but neither catalogue should use the other's identity scheme as its storage key.

MasjidBoard retrieval must continue to operate independently from audio playback availability.

## Stage 3 Decisions

The following are now the working Stage 3 decisions:

1. `provider + web_url` is the stable MasjidPi catalogue identity.
2. `MBL_ID` remains opaque metadata.
3. Catalogue and selected-board state are persisted independently.
4. An already-selected board must work without a successful discovery refresh.
5. Discovery refresh is transactional and last-known-good.
6. Unchanged catalogue/state produces no persistent write.
7. Rename with unchanged `web_url` preserves identity.
8. A changed `web_url` is not automatically merged with an older record.
9. Missing records are retained conservatively rather than immediately deleted.
10. Catalogue freshness uses MasjidPi-owned timestamps, not upstream `last_updated`.
11. Exact millisecond timezone offsets are preserved.
12. The MasjidBoard catalogue remains separate from the existing audio-stream catalogue.

## Stage 4 Implementation Boundary

With the Stage 3 design established, implementation should proceed in small independently tested pieces:

```text
1. provider-neutral catalogue model
2. catalogue reconciliation logic
3. atomic catalogue storage
4. selected-board storage
5. discovery -> catalogue builder
6. last-known-good refresh service
7. selected-board loader/provider construction
8. board-data cache and refresh
9. API/UI search and selection
```

The first implementation step should be the provider-neutral catalogue model plus reconciliation tests. Persistence should be added only after the reconciliation behaviour is proven in memory.
