# MasjidBoard Discovery Scope

**Status:** Stage 3 architecture  
**Branch:** `research/masjidboard-live`

## Decision

MasjidPi will not maintain a worldwide local copy of the MasjidBoard Live directory.

During initial MasjidBoard configuration, the user chooses a geographic discovery scope through the WebUI/API:

```text
Country
  -> Province / Region
      -> Town / City
```

MasjidPi persists that scope and maintains a local catalogue containing the boards returned for that town/city only.

The selected 1–3 boards are persisted separately from the discovery scope and catalogue.

## Why

The upstream FindMasjid service exposes a hierarchical directory and can contain hundreds of boards in a single country. A worldwide mirror would add network traffic, update complexity, validation work, storage churn, and data the appliance will normally never use.

A configured-location catalogue retains the important product behaviour while keeping the appliance lightweight:

- configuration can browse every board in the user's chosen town/city;
- normal catalogue browsing is local and disk-backed;
- only a small upstream result set needs refreshing;
- selected boards remain independent of catalogue availability;
- changing location can rebuild the catalogue for the new scope.

## Persisted Scope

Logical model:

```text
DiscoveryScope
    Country
    Region
    City
```

The configured state requires a country and city. Region may be blank because the upstream hierarchy can contain entries without a province/region value.

Persistent paths:

```text
development:
backend/data/masjidboard_scope.json

installed appliance:
/var/lib/masjidpi/masjidboard_scope.json
```

The scope file is separate from:

```text
masjidboard_catalogue.json
masjidboard_selection.json
masjidboard_cache/
```

because each has a different lifecycle.

## Initial Configuration

The configuration API/WebUI should walk the upstream hierarchy rather than requiring a global catalogue:

```text
request countries
    -> user chooses country
request regions for country
    -> user chooses region
request towns/cities for region
    -> user chooses town/city
persist discovery scope
    -> retrieve catalogue for that scope
    -> user selects 1–3 boards
```

The display remains presentation-only and does not participate in this flow.

## Catalogue Refresh Policy

The local scoped catalogue has two refresh paths only:

### Automatic

Refresh at most once every seven days.

The due decision must be based on the persisted catalogue's MasjidPi-owned `validated_at` / `retrieved_at` timestamps, not process uptime and not upstream per-board `last_updated` values.

A reboot therefore does not restart the seven-day interval.

### Manual

The user can explicitly request a catalogue refresh through the configuration API/WebUI. A manual refresh attempts immediately regardless of catalogue age.

## Failure Behaviour

A refresh is candidate-based:

```text
retrieve configured location
    -> parse
    -> normalize
    -> validate
    -> reconcile
    -> atomically persist
```

If retrieval, parsing, validation, reconciliation, or persistence fails, the previous last-known-good catalogue remains untouched.

Catalogue refresh failure does not affect the selected 1–3 boards. Their timetable providers and per-board last-known-good caches operate independently.

## Scope Changes

When the user changes country/region/city, MasjidPi should treat the existing catalogue as belonging to the previous scope and rebuild the catalogue for the new scope before presenting board choices.

The implementation should not silently mix records from two geographic scopes.

Existing selected boards are a separate configuration concern. A later configuration API step must decide whether changing the discovery scope should preserve, replace, or explicitly reconfirm selections; this document does not couple selection validity to catalogue membership.

## Memory Lifecycle

The discovery scope is tiny persisted configuration. The catalogue remains disk-first and does not need to remain resident in memory during normal appliance operation.

The expected normal runtime is therefore:

```text
selected-board state      -> resident as needed by runtime
selected timetable data   -> resident / cached independently
catalogue                  -> disk, loaded while configuration UI browses it
discovery scope            -> disk, loaded for catalogue maintenance/configuration
```
