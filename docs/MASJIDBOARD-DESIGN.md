# MasjidBoard Application Design

**Status:** Design — proposed for review
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

The cache provides resilience when MasjidBoard Live is unavailable.

It should distinguish at least:

- fresh usable data
- stale but usable data
- no usable cached data

The cache should not require the renderer to understand the upstream API.

### `scheduler`

Determines what should currently appear on the display.

Inputs include:

- normalised board data
- current time
- display configuration
- available content
- media availability

The scheduler should produce presentation state rather than drawing anything itself.

### `display`

Defines the rendering boundary.

The scheduler supplies presentation state and the display implementation turns that state into pixels.

The display layer must not fetch MasjidBoard Live data.

### `display/hdmi`

Initial concrete renderer/output implementation for the Raspberry Pi HDMI display.

The exact graphics technology is intentionally not fixed yet. The key requirement is that it should be substantially lighter than running a full browser simply to display the MasjidBoard webpage.

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

Each prayer should be represented semantically, for example:

```text
PrayerTime
├── Prayer
├── Adhan
└── Jamaah
```

The exact Go representation is still to be finalised, but the model must not require callers to understand upstream row/column positions.

If a source supplies additional prayer-related information, it should be represented as an extension of the semantic model rather than as an opaque source field.

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

The exact interface and storage format are not yet fixed.

A cached board should include freshness metadata so that the application can tell the difference between current and stale data.

The cache should preserve the last known valid **core prayer schedule** even if a later refresh fails.

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

For example:

```text
Provider refresh       periodically
Media refresh          when changed/required
Scheduler evaluation   continuously/as required
Clock/countdown        approximately once per second
Slide rotation         according to presentation rules
```

The exact intervals are not yet fixed.

A slow upstream refresh interval must not make the on-screen clock or prayer countdown appear frozen.

## Media Handling

Images/posters are first-class board content.

The provider/cache side should:

1. Discover media references from MasjidBoard Live.
2. Download media when required.
3. Store it locally.
4. Associate it with the normalised `Media` model.
5. Allow the scheduler/renderer to use the local copy.

The renderer should not need a live Internet connection to display already-cached media.

Failed optional media downloads should not invalidate the core prayer schedule.

## Stale Data Behaviour

MasjidBoard must distinguish between **stale but usable** and **invalid/unavailable**.

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

The exact user-facing stale-data indicator will be decided during display design.

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

### Scheduler tests

Given a fixed board and timestamp, verify:

- current prayer
- next prayer
- countdown state
- Jumu'ah selection
- optional slide selection
- missing optional content
- stale data behaviour

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

## Open Design Questions

These remain to be decided during implementation/design review:

1. Exact Go types for the semantic model.
2. Cache storage format and location.
3. Exact renderer technology.
4. Exact display slide layouts.
5. Exact provider refresh interval.
6. Exact media refresh/cache policy.
7. User-facing stale-data indication.
8. Whether MasjidBoard should expose its own HTTP/API endpoints or share selected MasjidPi infrastructure.
9. Exact service/executable naming.

These questions do not change the core architectural boundaries already agreed.
