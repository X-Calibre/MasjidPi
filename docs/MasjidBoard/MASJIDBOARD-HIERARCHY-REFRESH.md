# MasjidBoard Hierarchy Refresh

**Status:** Stage 4 implementation  
**Branch:** `research/masjidboard-live`

## Purpose

Maintain a lightweight, global MasjidBoard location hierarchy independently of the configured city-level board catalogue.

```text
countries
    -> regions / provinces
        -> towns / cities
```

The hierarchy contains location labels and upstream board counts only. It does not contain individual board timetable records.

## Persistence

The last-known-good hierarchy is stored separately from scope, catalogue, selection and board-data cache:

```text
/var/lib/masjidpi/masjidboard_hierarchy.json
```

Development uses:

```text
backend/data/masjidboard_hierarchy.json
```

## Refresh Policy

The hierarchy is refreshed:

- at most once every seven days by scheduled maintenance; or
- immediately when the user explicitly requests a hierarchy update.

A reboot does not reset the seven-day interval. `validated_at` in the persisted hierarchy determines whether scheduled refresh is due.

A manual refresh bypasses the age check.

## Transactional Behaviour

A candidate hierarchy must be completely retrieved and validated before it can replace the persisted hierarchy.

```text
load last-known-good hierarchy
        |
        v
retrieve countries
        |
        v
retrieve regions for every country
        |
        v
retrieve cities for every region
        |
        v
normalise + validate counts
        |
        +--> failure
        |       -> retain previous hierarchy unchanged
        |
        +--> success
                -> atomically persist candidate
```

Network errors, malformed responses, incomplete traversal and persistence failures must not erase the previous hierarchy.

## Upstream Counts

MasjidBoard Live exposes board counts at each hierarchy level. These are retained and used as completeness checks.

For a valid candidate:

```text
sum(region counts) == country count
sum(city counts)   == region count
```

Duplicate labels returned by upstream are merged before lower-level traversal and their counts are summed. This specifically accommodates observed upstream anomalies such as duplicate province labels.

Blank region names are preserved because MasjidBoard Live has been observed to expose a legitimate blank region bucket. Blank country and city names are not accepted.

## Relationship to Scoped Catalogue

The hierarchy is global and lightweight. The board catalogue remains local to the user's configured discovery scope.

```text
global hierarchy
        |
        v
user chooses country / region / city
        |
        v
persisted scope
        |
        v
city-level board catalogue
        |
        v
selected 1-3 boards
```

Refreshing the hierarchy does not refresh selected board timetables. Those are separate runtime operations with their own last-known-good caches.
