# MasjidBoard Application Design

**Status:** Design — agreed
**Branch:** `research/masjidboard-live`

## Purpose

This document turns the architectural decisions in `MASJIDBOARD-ARCHITECTURE.md` into a concrete application/component design.

MasjidBoard is a standalone application within the MasjidPi repository. It uses MasjidBoard Live as its primary data source, normalises the upstream data into a semantic model, caches the result and media locally, schedules board content, and renders the resulting presentation to an HDMI display.

The design deliberately keeps MasjidBoard independent from the MasjidPi audio application.

## Component Structure

The initial Go package structure should be approximately:

```text
backend/
├── cmd/
│   ├── masjidpi/
│   └── masjidboard/
│
└── internal/
    └── masjidboard/
        ├── app/
        ├── model/
        ├── provider/
        │   └── masjidboardlive/
        ├── cache/
        ├── scheduler/
        └── display/
            └── hdmi/
```

This is a starting structure, not a requirement that every package become a large abstraction. Packages should remain small until there is a real reason to split them further.

## Responsibilities

### `cmd/masjidboard`

Application entry point.

Responsible for:

- loading configuration
- creating dependencies
- starting the MasjidBoard application
- handling process signals and graceful shutdown
- reporting startup/runtime errors

It should contain very little business logic.

### `app`

Owns the MasjidBoard application lifecycle and coordinates the major components.

It should manage:

- initial cache loading
- provider refresh lifecycle
- model updates
- scheduler lifecycle
- renderer lifecycle
- shutdown

It should not contain MasjidBoard Live parsing logic or detailed rendering logic.

### `model`

Contains the normalised semantic MasjidBoard domain model.

The model must not know that MasjidBoard Live uses a positional 29-row response.

### `provider`

Defines the boundary between the application and external board data sources.

The provider abstraction exists so that the rest of the application consumes normalised `Board` data rather than source-specific structures.

### `provider/masjidboardlive`

Implements the MasjidBoard Live provider.

Responsibilities:

- request MasjidBoard Live data
- parse the positional upstream response
- interpret source-specific fields
- validate core data
- normalise it into the internal model
- identify/download required media
- expose source metadata needed for diagnostics/cache refresh

No MasjidBoard Live-specific field names or positional indexes should leak into the rest of the application.

### `cache`

Persists the normalised board data and downloaded media.

The agreed initial cache location is:

```text
/var/lib/masjidpi/masjidboard/
├── board.json
├── metadata.json
└── media/
```

The cache provides resilience when MasjidBoard Live is unavailable.

It should distinguish at least:

- fresh usable data
- stale but usable data
- no usable cached data

The cache should not require the renderer to understand the upstream API.

The normalised model is the primary runtime cache representation. The raw upstream response may optionally be retained for diagnostics, but it is not the application's primary runtime data format.

### `scheduler`

Determines what should currently appear on the display.

Inputs include:

- normalised board data
- current time
- display configuration
- available content
- media availability

The scheduler should produce presentation state rather than drawing anything itself.

Prayer times form the primary scheduler context. Optional content is scheduled around the core prayer-time presentation rather than being a prerequisite for board operation.

### `display`

Defines the rendering boundary.

The scheduler supplies presentation state and the display implementation turns that state into pixels.

The display layer must not fetch MasjidBoard Live data.

### `display/hdmi`

Initial concrete renderer/output implementation for the Raspberry Pi HDMI display.

The renderer will be native rather than browser-based. A lightweight graphics implementation will be used for the initial prototype; **SDL2 is the preferred candidate**, but the core display interface must not depend on SDL2 so that it can be replaced if target-hardware testing identifies a better option.

The display target is initially **1920×1080 landscape**.

The renderer must be substantially lighter than running a full browser simply to display the MasjidBoard webpage.

## Domain Model

The model is divided into **required core data** and **optional enrichment**.

### Required core data

```text
Board
├── Identity
└── PrayerTimes
    ├── Fajr
    ├── Zuhr
    ├── Asr
    ├── Maghrib
    └── Esha
```

A usable board dataset requires:

- a valid masjid identity
- a valid date/context
- a usable schedule for the five daily prayers

The five prayers are the primary purpose of the board.

### Prayer model

Each prayer is represented semantically:

```text
PrayerTime
├── Prayer
├── Adhan
└── Jamaah
```

The five daily prayers use a fixed semantic structure rather than a generic map so that the application's core purpose is explicit and cannot be accidentally omitted.

Adhan and Jamaah are individually optional values because the upstream source may not provide both.

Actual prayer times will use Go `time.Time` values internally with explicit timezone/date context. The provider converts MasjidBoard Live's source strings into these values; the rest of the application must not manipulate prayer times as arbitrary strings.

### Date context

Prayer schedules are date-specific. The board model therefore includes explicit date context containing at least:

```text
DateContext
├── Gregorian date
├── Islamic date
└── Location / timezone
```

A cached schedule for a previous Gregorian date must not silently be presented as the current day's prayer schedule after a date boundary. If the current date has no valid prayer schedule, the application must enter an explicit no-valid-data/recovery state rather than masquerading yesterday's schedule as today's.

### Optional enrichment

The initial model should support:

```text
Board
├── Identity                 REQUIRED
├── PrayerTimes              REQUIRED
├── AstronomicalTimes        OPTIONAL
├── JumuahServices[]         OPTIONAL
├── Announcements[]          OPTIONAL
├── Programmes[]             OPTIONAL
├── Notices[]                OPTIONAL
├── Media[]                  OPTIONAL
├── Banking                  OPTIONAL
├── ContributionInformation  OPTIONAL
├── NewMoon                  OPTIONAL
├── DisplayConfiguration     OPTIONAL
└── DeferredContent          OPTIONAL / FUTURE
```

Ayah, Hadith and Sunnah are deliberately deferred from the initial implementation.

## Provider Interface

The application should depend on a small provider contract rather than the concrete MasjidBoard Live implementation.

Conceptually:

```go
type Provider interface {
    Fetch(ctx context.Context) (Board, error)
}
```

The exact interface can be refined when implementation begins.

The important rule is that `Provider` returns the normalised model.

It must not return the 29-row MasjidBoard Live structure.

### Provider outcomes

The provider should allow the application to distinguish between:

1. Successful refresh with valid core data.
2. Successful refresh with valid core data and missing optional content.
3. Upstream failure where cached core data can continue to be used.
4. Invalid upstream response where cached core data can continue to be used.
5. No valid data available at all.

Missing optional content is not a provider failure.

## Cache Interface

Conceptually:

```go
type Cache interface {
    Load(ctx context.Context) (CachedBoard, error)
    Save(ctx context.Context, board Board) error
}
```

The exact implementation is not yet fixed, but the agreed initial storage is normalised JSON under `/var/lib/masjidpi/masjidboard/` plus locally cached media.

A cached board should include freshness metadata so that the application can tell the difference between current and stale data.

The cache should preserve the last known valid **core prayer schedule** even if a later refresh fails.

Media should be cached independently enough that a failed optional media download cannot invalidate the core board.

## Scheduler Interface

Conceptually:

```go
type Scheduler interface {
    Current(board Board, now time.Time) DisplayState
}
```

The scheduler should be deterministic for a given board, configuration and timestamp wherever possible.

It should be responsible for deciding:

- current/next prayer state
- countdown information
- which content is eligible to display
- slide ordering
- display duration
- special/seasonal content
- behaviour when optional content is absent
- behaviour when data is stale

It should not perform network requests or render graphics.

## Display State

The scheduler should not pass the entire domain model directly to the renderer.

Instead:

```text
Board data
    ↓
Scheduler
    ↓
DisplayState
    ↓
Renderer
```

`DisplayState` represents what the renderer needs for the current presentation.

It may contain concepts such as:

```text
DisplayState
├── Current slide/content type
├── Masjid identity
├── Date information
├── Prayer display
├── Current/next prayer state
├── Optional content for the current slide
├── Media reference
├── Display duration
└── Status/staleness information
```

The renderer should not need to understand provider failures, API refreshes or cache storage.

## Refresh Lifecycle

The application should separate data refresh from display timing.

Conceptually:

```text
Startup
  │
  ├── Load cached board
  │
  ├── Start display
  │
  ├── Start scheduler
  │
  └── Start provider refresh loop
             │
             ▼
       MasjidBoard Live
             │
       success / failure
          /       \
     success     failure
        │            │
        ▼            ▼
    normalise     keep cache
        │
        ▼
      cache
        │
        ▼
    update board
```

The display should not have to wait for the first network refresh if valid cached data exists.

If no cache exists and the provider is unavailable, the application should enter an explicit no-data/recovery state rather than silently pretending the board is current.

## Refresh versus Display Timing

These are independent concerns.

The initial provider refresh interval is **10 minutes**, but it should be configurable internally rather than hard-coded throughout the provider.

The local display clock and countdown must not depend on this refresh interval.

Conceptually:

```text
Provider refresh       every ~10 minutes initially
Media refresh          when changed/required
Scheduler evaluation   continuously/as required
Clock/countdown        approximately once per second
Slide rotation         according to presentation rules
```

The exact user-configurable refresh options can be considered later.

A slow upstream refresh interval must not make the on-screen clock or prayer countdown appear frozen.

## Media Handling

Images/posters are first-class board content.

The provider/cache side should:

1. Discover media references from MasjidBoard Live.
2. Determine whether the referenced content is already cached.
3. Download media when required.
4. Store it locally.
5. Associate it with the normalised `Media` model.
6. Allow the scheduler/renderer to use the local copy.

Media identity should use available upstream identifiers and/or content hashes so unchanged media is not unnecessarily downloaded again.

The cache should retain enough metadata to support later cleanup of obsolete media. Exact garbage-collection rules are deferred until implementation/testing.

The renderer should not need a live Internet connection to display already-cached media.

Failed optional media downloads should not invalidate the core prayer schedule.

## Stale Data Behaviour

MasjidBoard must distinguish between **fresh**, **stale but usable**, and **no valid data**.

If MasjidBoard Live is temporarily unreachable:

```text
Network unavailable
       ↓
Last known valid Board
       ↓
Display continues
       ↓
Status indicates stale data internally / where appropriate
```

The application should never discard valid cached prayer data merely because optional content could not be refreshed.

Stale cached prayer data may continue to be displayed while it remains trustworthy, but a cached schedule from a previous date must not be used as though it were today's schedule.

The initial implementation should keep the stale indication subtle rather than placing a large warning over the primary prayer display. The exact user-facing indicator will be decided during display design.

## Failure Isolation

The MasjidBoard process must not depend on the audio process.

Likewise, a MasjidBoard failure must not stop audio playback.

Within MasjidBoard:

```text
Provider failure
    → use cache

Media failure
    → omit affected media
    → continue board operation

Scheduler failure
    → application recovery/restart

Renderer failure
    → application recovery/restart
```

The exact service-level restart strategy will follow the existing MasjidPi service/recovery conventions.

## Display Modes

The initial implementation should support the information required for functional/content parity with MasjidBoard Live without requiring pixel-perfect reproduction of its webpage.

The renderer should be capable of presenting at least:

1. Core daily prayer information.
2. Jumu'ah information when available.
3. Optional astronomical/perpetual times when available.
4. Announcements/programmes/notices when available.
5. Posters and other media when available.
6. Seasonal/special content when available.

The precise visual design and slide layouts are a separate design task.

## Configuration Boundary

MasjidBoard configuration should eventually cover at least:

- MasjidBoard Live board/masjid identifier.
- Enable/disable MasjidBoard.
- Display configuration.
- Refresh behaviour where user-configurable.
- Any required language/display settings.

Configuration should not be embedded in the provider implementation.

## API Boundary

MasjidBoard will have its own application/API boundary rather than being coupled to the existing audio API.

The initial API surface should remain deliberately small and should only expose functionality needed by the MasjidBoard application and MasjidPi UI. Candidate endpoints include:

```text
/api/masjidboard/status
/api/masjidboard/board
/api/masjidboard/display
```

A refresh endpoint may be added later if required:

```text
/api/masjidboard/refresh
```

These endpoints are not required before the core provider/model/cache/scheduler path is working.

## Executable and Service

The agreed conceptual executable and service names are:

```text
masjidboard
masjidboard.service
```

alongside the existing MasjidPi audio application/service:

```text
masjidpi
masjidpi.service
```

Both must be independently startable, stoppable and recoverable.

## Testing Strategy

The design should allow most MasjidBoard logic to be tested without a Raspberry Pi or HDMI display.

### Provider tests

Use captured MasjidBoard Live responses to verify:

- correct parsing
- correct normalisation
- required prayer validation
- optional-field handling
- malformed/incomplete upstream data

### Model/cache tests

Verify:

- persistence and reload
- freshness metadata
- preservation of last known valid prayer data
- optional content absence
- date-boundary behaviour

### Scheduler tests

Given a fixed board and timestamp, verify:

- current prayer
- next prayer
- countdown state
- Jumu'ah selection
- optional slide selection
- missing optional content
- stale data behaviour
- previous-date cached data is not treated as today's schedule

### Renderer tests

Keep renderer tests separate from provider/network tests.

Where practical, presentation logic should be testable without physical HDMI hardware.

## Implementation Order

The initial implementation should proceed in this order:

### 1. Model

Define the normalised core board/prayer model and optional enrichment structures.

### 2. Provider contract

Define the provider boundary without exposing MasjidBoard Live's positional schema.

### 3. MasjidBoard Live provider

Implement parsing and normalisation using captured upstream responses.

### 4. Cache

Persist the normalised board and media references/assets.

### 5. Scheduler

Implement prayer-aware board state and initial content scheduling.

### 6. Display state

Define the presentation contract between scheduler and renderer.

### 7. Renderer prototype

Implement a minimal 1920×1080 renderer showing the core prayer board.

### 8. Optional content

Add Jumu'ah, announcements, programmes, notices, media and other enrichment incrementally.

### 9. Raspberry Pi HDMI integration

Validate actual display output and resource usage on the target hardware.

### 10. Service integration

Add independent startup, shutdown, recovery and configuration integration with MasjidPi.

## Deliberate Non-Goals for Initial Implementation

The first implementation will not attempt to:

- reproduce the MasjidBoard Live webpage itself
- reproduce its HTML/CSS/JavaScript architecture
- require a full web browser
- block implementation on Ayah/Hadith/Sunnah
- support every possible upstream field before the core board works
- couple MasjidBoard to MPV/audio playback

## Settled Design Decisions

The following questions have been reviewed and settled for the initial implementation:

1. **Exact Go model:** explicit semantic structs and enums; no generic maps for the core domain model.
2. **Prayer representation:** fixed five-prayer structure with `PrayerTime` values containing Adhan and Jamaah where available.
3. **Time representation:** Go `time.Time` internally with explicit date/timezone context.
4. **Cache location:** `/var/lib/masjidpi/masjidboard/`.
5. **Cache format:** normalised JSON plus locally cached media; raw upstream data is optional diagnostics only.
6. **Renderer:** native graphics; SDL2 is the preferred initial candidate, behind a renderer interface.
7. **Display target:** 1920×1080 landscape initially.
8. **Browser:** no full browser dependency.
9. **Provider refresh:** 10 minutes initially, with the implementation keeping the interval configurable internally.
10. **Media caching:** local cache using available source identity and/or content hash to avoid unnecessary downloads; exact garbage collection deferred.
11. **Stale data:** last known valid prayer data may continue to display while trustworthy; stale status is tracked; previous-date data cannot masquerade as today's schedule.
12. **API boundary:** MasjidBoard has its own small API boundary; exact endpoints can evolve as implementation requires.
13. **Executable/service:** `masjidboard` / `masjidboard.service`.

These decisions are now the working baseline for implementation. Changes should be recorded explicitly if later testing demonstrates that a decision needs to change.
