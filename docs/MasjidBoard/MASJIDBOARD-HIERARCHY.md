# MasjidBoard Discovery Hierarchy

**Status:** implementation design  
**Branch:** `research/masjidboard-live`

## Purpose

MasjidPi persists a lightweight global MasjidBoard Live location hierarchy independently from the scoped board catalogue.

```text
MasjidBoard hierarchy
    Country
        -> Province / Region
            -> Town / City

Configured scope
    one Country + Region + City

Scoped catalogue
    actual board records for that configured city only
```

The hierarchy exists so initial configuration and later location changes can be performed from local persisted data without downloading the global board catalogue.

## Persistent state

Installed appliances use:

```text
/var/lib/masjidpi/masjidboard_hierarchy.json
```

Development uses:

```text
backend/data/masjidboard_hierarchy.json
```

The hierarchy contains location names, upstream board counts, and MasjidPi-owned retrieval/validation timestamps. It does not contain timetable or individual board records.

## Upstream source

The public FindMasjid endpoint is queried with the same hierarchy operations used by its frontend:

```text
country
province
cityProvince
```

Observed examples include:

```text
South Africa -> 615 boards
North West   -> 23 boards
```

Upstream data can contain anomalies such as duplicate region labels and blank region buckets. MasjidPi preserves blank regions and normalises duplicate labels at the same hierarchy level by summing their counts.

## Refresh policy

The hierarchy should refresh:

- automatically no more than once every seven days; or
- immediately when the user explicitly requests an update.

A failed refresh must never destroy the last-known-good hierarchy. A candidate hierarchy is persisted only after complete retrieval and validation.

The hierarchy refresh policy is separate from both:

- scoped catalogue refreshes; and
- selected-board timetable refreshes.

## Counts

Upstream counts are retained because they are cheap and useful validation metadata. A later hierarchy updater should use them to detect incomplete lower-level traversal before replacing persisted state.

## User flow

```text
WebUI/API
    -> read persisted hierarchy
    -> choose country
    -> choose province/region
    -> choose town/city
    -> persist discovery scope
    -> build/refresh board catalogue for that scope
    -> choose 1-3 boards
```

All configuration remains in the API/WebUI. The MasjidBoard display remains presentation-only.
