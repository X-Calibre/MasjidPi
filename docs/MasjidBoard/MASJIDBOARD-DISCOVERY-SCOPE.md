# MasjidBoard Discovery Scope

**Status:** Stage 3 architecture  
**Branch:** `research/masjidboard-live`

## Decision

MasjidPi will not maintain a worldwide local copy of the MasjidBoard Live board directory.

During MasjidBoard configuration, the user chooses **one to three geographic discovery locations** through the WebUI/API:

```text
Location 1
  Country
    -> Province / Region
        -> Town / City

Location 2 (optional)
  Country
    -> Province / Region
        -> Town / City

Location 3 (optional)
  Country
    -> Province / Region
        -> Town / City
```

This supports users who live near a town, municipal, provincial or national boundary. A mosque in an adjacent location may be physically closer than another mosque in the user's nominal town.

MasjidPi persists the ordered location set and maintains a local catalogue containing the union of boards returned for those locations.

The selected 1–3 displayed boards are persisted separately from the discovery scope and catalogue. The limit of three discovery locations and the limit of three displayed boards are independent product rules.

## Why

The upstream FindMasjid service exposes a hierarchical directory and can contain hundreds of boards in a single country. A worldwide board mirror would add network traffic, update complexity, validation work, storage churn, and data the appliance will normally never use.

A small multi-location catalogue retains the important product behaviour while keeping the appliance lightweight:

- configuration can browse every board in up to three nearby towns/cities;
- users near boundaries can include adjacent locations;
- normal catalogue browsing is local and disk-backed;
- only small upstream result sets need refreshing;
- selected boards remain independent of catalogue availability;
- locations may cross province/region or country boundaries when useful.

## Persisted Scope

Logical model:

```text
DiscoveryScope
    Locations[1..3]
        Country
        Region
        City
```

Every location requires a country and city. Region may be blank because the upstream hierarchy can contain entries without a province/region value.

Duplicate locations are not allowed. Locations are compared after whitespace normalization and case-insensitively for duplicate detection. User order is preserved.

Example:

```json
{
  "locations": [
    {
      "country": "South Africa",
      "region": "North West",
      "city": "Brits"
    },
    {
      "country": "South Africa",
      "region": "Gauteng",
      "city": "Akasia"
    }
  ]
}
```

Persistent paths:

```text
development:
backend/data/masjidboard_scope.json

installed appliance:
/var/lib/masjidpi/masjidboard_scope.json
```

The scope file is separate from:

```text
masjidboard_hierarchy.json
masjidboard_catalogue.json
masjidboard_selection.json
masjidboard_cache/
```

because each has a different lifecycle.

## Configuration Flow

The configuration API/WebUI reads the persisted global hierarchy and allows the user to add up to three locations:

```text
choose country
    -> choose region
    -> choose town/city
    -> add location

optionally repeat for location 2 / 3
    -> persist discovery scope
    -> retrieve catalogue data for each location
    -> merge/deduplicate catalogue for browsing
    -> user selects 1–3 boards
```

The display remains presentation-only and does not participate in configuration.

## Combined Catalogue

The WebUI/API presents a single combined catalogue made from all configured locations.

Records that appear through more than one location are deduplicated using the provider-neutral catalogue identity derived from provider + external ID. For MasjidBoard Live, the public `MBL_ID` / provider identity remains the primary source identity and `web_url` can be used as a consistency check.

The persistence layer must still retain location provenance internally so refresh health can be tracked separately for each configured location.

## Catalogue Refresh Policy

Each configured location has two refresh paths only.

### Automatic

Refresh at most once every seven days.

The due decision must use MasjidPi-owned persisted successful-refresh timestamps, not process uptime and not upstream per-board `last_updated` values. A reboot therefore does not restart the seven-day interval.

### Manual

The user can explicitly request a catalogue refresh through the configuration API/WebUI. A manual refresh attempts all configured locations immediately regardless of age.

## Failure Isolation

Catalogue refresh is **location-isolated**.

For example:

```text
Location A refresh succeeds
    -> replace Location A candidate data

Location B refresh fails
    -> keep Location B last-known-good data

Location C refresh succeeds
    -> replace Location C candidate data
```

A failure in one configured location must not discard successfully refreshed data from another location, and must not erase the failed location's previous last-known-good data.

The combined catalogue can therefore continue to be served even when one upstream location is temporarily unavailable. The API should expose refresh health sufficiently for the WebUI to indicate that part of the catalogue is stale.

Catalogue refresh failure does not affect the selected 1–3 boards. Their timetable providers and per-board last-known-good caches operate independently.

## Scope Changes

Adding or removing a discovery location changes only the relevant location partition of the catalogue.

The implementation must not silently treat a board as missing merely because one of several source locations was removed if that board is still present through another configured location.

Existing selected boards are a separate configuration concern. Selection validity is not implicitly tied to current catalogue membership; a selected board can continue operating from persisted identity/cache even when discovery data changes.

## Memory Lifecycle

The discovery scope is tiny persisted configuration. The catalogue remains disk-first and does not need to remain resident in memory during normal appliance operation.

The expected normal runtime is therefore:

```text
selected-board state      -> resident as needed by runtime
selected timetable data   -> resident / cached independently
catalogue                  -> disk, loaded while configuration UI browses it
discovery scope            -> disk, loaded for catalogue maintenance/configuration
global hierarchy           -> disk, loaded by configuration UI / weekly maintenance
```
