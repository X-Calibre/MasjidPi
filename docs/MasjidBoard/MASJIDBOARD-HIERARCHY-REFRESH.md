# MasjidBoard Hierarchy Refresh

**Status:** Implemented

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

A candidate hierarchy must be retrieved and structurally validated before it can replace the persisted hierarchy.

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
retrieve available cities for every region
        |
        v
normalise + validate counts / unresolved coverage
        |
        +--> structural or over-count failure
        |       -> retain previous hierarchy unchanged
        |
        +--> success
                -> atomically persist candidate
```

Network errors while retrieving countries, regions or otherwise resolvable city lists, malformed responses, count over-runs and persistence failures must not erase the previous hierarchy.

## Upstream Counts and Unresolved Coverage

MasjidBoard Live exposes board counts at each hierarchy level. These are retained and used as completeness checks.

The live all-country validation on 2026-08-18 detected 24 countries and 721 boards. It also established that the upstream hierarchy is not perfectly reversible at city level. South Africa exposed two boards that are counted by higher hierarchy levels but cannot currently be represented by a unique returned city row:

- one board in a blank region bucket for which the city endpoint returns HTTP 500; and
- one extra board included in a duplicate `Limpopo` region count, while the merged city response accounts for only 25 of the 26 boards.

The hierarchy therefore distinguishes **resolved** city coverage from **unresolved** upstream coverage. `Region.UnresolvedCount` records the number of boards that are included in the authoritative upstream region count but are not represented by usable city rows.

For a valid candidate:

```text
sum(region counts) == country count
sum(city counts) + unresolved_count == region count
```

A city total greater than its region count is invalid and rejects the candidate. A city total below its region count is retained with the difference recorded as unresolved rather than causing all otherwise-valid countries and cities to remain stale.

A blank region is handled according to context:

- if it is the country's only region, it represents a country without a province layer and its cities are resolved normally;
- if it appears alongside named regions, it is preserved as an unresolved bucket instead of calling an upstream city request known to be non-resolvable for this shape.

Duplicate labels returned by upstream are merged before lower-level traversal and their counts are summed. This accommodates observed duplicate province labels while retaining the upstream total for validation.

Blank country and city names are not accepted.

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

Unresolved hierarchy counts are diagnostic metadata; they do not create synthetic city choices in the WebUI. Users select from the actual country / region / city rows that upstream exposes.

Refreshing the hierarchy does not refresh selected board timetables. Those are separate runtime operations with their own last-known-good caches.
