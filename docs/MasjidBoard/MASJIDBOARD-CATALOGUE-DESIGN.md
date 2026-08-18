# MasjidBoard Catalogue and Persistence Design

**Status:** Stage 3 research / architecture  
**Branch:** `research/masjidboard-live`

## Purpose

Define the stable MasjidPi catalogue record, geographic discovery scope, multi-board selection model, refresh/reconciliation behaviour, and last-known-good strategy for MasjidBoard Live discovery.

```text
FindMasjid hierarchy
        |
        v
configured discovery scope
        |
        v
scoped local catalogue
        |
        +--> WebUI/API search and selection
        |
        +--> persisted 1–3 selected boards
        |
        v
independent Core providers
        |
        v
normalised Board results
```

Discovery scope, catalogue persistence, selected-board persistence, individual-board retrieval, and display remain separate concerns.

## Discovery Scope

MasjidPi does **not** maintain a worldwide mirror of the MasjidBoard Live directory.

During initial configuration the user chooses:

```text
Country
  -> Province / Region
      -> Town / City
```

That scope is persisted independently in `masjidboard_scope.json`. The local catalogue contains the boards for that configured town/city only.

Country and city are required for a configured scope. Region may be blank because the upstream FindMasjid hierarchy contains entries without a province/region value.

Changing scope rebuilds the local catalogue for the new location. Records from two scopes must not be silently mixed.

See `MASJIDBOARD-DISCOVERY-SCOPE.md` for the lifecycle decision.

## Design Goals

The catalogue subsystem should:

- use stable public MasjidBoard Live identifiers rather than opaque implementation IDs;
- keep already-selected boards operational if discovery is unavailable;
- keep the local catalogue small and relevant to the user's configured area;
- tolerate board renames and metadata changes;
- tolerate boards disappearing temporarily or permanently;
- avoid unnecessary SD-card writes;
- persist only changed data and use atomic replacement;
- distinguish catalogue freshness from individual-board timetable freshness;
- avoid treating upstream `last_updated` as catalogue freshness; and
- remain provider-specific at the discovery boundary while exposing provider-neutral records internally.

## Stable Catalogue Identity

The working external key is the public MasjidBoard Live `web_url` slug, for example:

```text
brits-jamia
```

MasjidPi namespaces that identity:

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

The opaque Premium `boardId` must not become the catalogue key. `MBL_ID` remains upstream metadata unless stronger stability guarantees are established.

## Catalogue Record

The persisted catalogue contains identity, geography, provider metadata needed for retrieval, and discovery-state information.

```text
CatalogueRecord
    ID                  masjidboardlive:<web_url>
    Provider            masjidboardlive
    ExternalID          <web_url>
    Name                upstream masjid name
    City                upstream city
    Country             configured country
    Region              configured province/region
    TimeZoneOffsetMS    exact upstream offset
    ProviderMetadata    optional upstream metadata such as MBL_ID
    DiscoveredAt        MasjidPi timestamp
    LastSeenAt          MasjidPi timestamp
    Status              active | missing | unavailable
```

Timetable summary fields from FindMasjid do not need to be persisted in the catalogue because selected boards retrieve their authoritative timetable through the Core provider.

Upstream `last_updated` may be retained as provider metadata but must not drive catalogue freshness or removal decisions.

## Multi-Board Selection

Selection is persisted independently from discovery scope and catalogue.

The product model is:

```text
unconfigured
    no persisted board selection

configured
    exactly 1–3 ordered selected boards
```

An empty configured selection and a fourth board are invalid.

Each selected board persists enough stable metadata to start without a catalogue refresh:

```text
SelectedBoard
    CatalogueID
    Provider
    ExternalID
    Name
    TimeZoneOffsetMS
```

Selection order is preserved and later maps directly to display order.

## Startup Behaviour

An already-configured appliance must not depend on catalogue availability:

```text
load selected-board state
        |
        +--> configured
        |       -> construct 1–3 independent providers
        |       -> refresh each independently
        |       -> use per-board last-known-good cache on failure
        |
        +--> unconfigured
                -> wait for WebUI/API configuration
```

The catalogue and discovery scope are not required for normal selected-board display once configuration is complete.

## Persistent Files

Installed appliance:

```text
/var/lib/masjidpi/masjidboard_scope.json
/var/lib/masjidpi/masjidboard_catalogue.json
/var/lib/masjidpi/masjidboard_selection.json
/var/lib/masjidpi/masjidboard_cache/
```

Development equivalents live under `backend/data/`.

These are separate because their lifecycles and write frequencies differ.

## Catalogue Memory Lifecycle

The catalogue is disk-first. It does **not** need to remain in memory throughout normal appliance operation.

It is loaded when the WebUI/API browses the configured location catalogue or when catalogue maintenance requires it.

The selected 1–3 board state and current timetable results are the runtime-relevant data.

## Catalogue Freshness

MasjidPi-owned timestamps are separate from MasjidBoard Live's per-board `last_updated` field:

```text
retrieved_at
    successful upstream scoped discovery retrieval

validated_at
    candidate parsed and passed validation

last_seen_at
    per-record most recent successful discovery
```

A catalogue can therefore be fresh even when an individual board has a blank or old upstream `last_updated` value.

## Refresh Policy

The scoped catalogue refreshes only in two circumstances.

### Automatic

At most once every seven days.

The due decision is based on persisted `validated_at` / `retrieved_at`, not process uptime. Reboots therefore do not reset the weekly interval.

### Manual

The user may explicitly request an immediate refresh through the configuration API/WebUI regardless of catalogue age.

Opening the configuration UI does not itself force a refresh.

Selected-board timetable refreshes are independent and occur on their own shorter provider schedule.

## Last-Known-Good Catalogue Behaviour

A failed discovery refresh must never replace a valid persisted catalogue with an empty, malformed, or partial result.

```text
retrieve configured scope
        |
        v
parse / normalise / validate candidate
        |
        +--> failure
        |       -> keep current catalogue unchanged
        |
        +--> success
                -> reconcile with current scoped catalogue
                -> atomically persist if changed
```

## Reconciliation Rules

### Existing record seen again

For the same `provider + external_id`:

- retain the same catalogue ID;
- update mutable metadata such as name, city, timezone and provider metadata;
- update `last_seen_at`;
- keep status active.

### Rename

An unchanged `web_url` with a changed display name is the same board.

### `web_url` change

A changed slug is treated as a new external identity. MasjidPi must not automatically merge records based only on similar names, city or `MBL_ID`.

### Temporarily missing

A board absent from a successful refresh is retained conservatively as `missing` rather than deleted immediately.

### Selected board disappears from catalogue

Selection remains valid. If Core retrieval still works, the board continues normally. If live retrieval fails, the runtime falls back to that board's last-known-good timetable cache.

## Availability Status

Catalogue status remains deliberately simple:

```text
active
missing
unavailable
```

A transient network failure must not automatically mark a board unavailable.

Runtime timetable state is separate:

```text
current
stale
unavailable
```

where `stale` means the latest board refresh failed but persisted last-known-good timetable data is still displayed with an update-failed indication.

## Timezone Handling

`time_zone_milli` is persisted exactly and must not be rounded to whole hours. The catalogue-to-provider handoff supports offsets such as:

```text
GMT+02:00
GMT+05:30
GMT+05:45
GMT+09:30
GMT-03:30
```

MasjidPi must not invent an IANA timezone where upstream only exposes a fixed millisecond offset.

## Search and Configuration

The WebUI/API owns all configuration. The MasjidBoard display is presentation-only.

Initial configuration walks the upstream hierarchy:

```text
countries
    -> regions/provinces
        -> towns/cities
            -> persist scope
            -> fetch local scoped catalogue
            -> select 1–3 boards
```

After generation, ordinary board search/filtering operates against the persisted local catalogue rather than making a remote request for every keystroke.

## Separation from Audio Catalogue

The MasjidBoard catalogue is separate from the existing LiveMasjid audio-stream catalogue. Neither should use the other's identity scheme as its storage key, and MasjidBoard failure must not prevent audio playback from starting or operating.

## Current Decisions

1. `provider + web_url` is the stable MasjidPi catalogue identity.
2. `MBL_ID` remains opaque provider metadata.
3. Discovery scope, catalogue, selection and board caches are persisted independently.
4. Discovery is scoped to the configured country/region/city rather than mirrored globally.
5. A configured selection contains exactly 1–3 ordered boards.
6. Selected boards operate without a successful catalogue refresh.
7. Catalogue refresh is transactional and last-known-good.
8. Automatic catalogue refresh is weekly; manual refresh may run immediately.
9. Catalogue freshness uses MasjidPi-owned timestamps, not upstream `last_updated`.
10. The catalogue is disk-first and loaded on demand for configuration/browsing.
11. Exact millisecond timezone offsets are preserved.
12. Display is read-only; all configuration belongs to WebUI/API.
13. MasjidBoard remains independent from audio playback.

## Next Implementation Boundary

With discovery scope persistence now defined, the next implementation pieces are:

```text
1. hierarchy client methods for country / region / city discovery
2. scoped catalogue builder for the persisted location
3. weekly-due decision based on persisted catalogue timestamps
4. manual catalogue update operation
5. configuration API for reading/updating discovery scope
6. configuration API for ordered 1–3 board selection
```
