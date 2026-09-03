# MasjidBoard Application Architecture

**Status:** Architecture / design — agreed
**Branch:** `research/masjidboard-live`

## Decision

MasjidBoard will be designed as a **standalone application within the MasjidPi repository**, rather than as a tightly coupled feature of the audio-streaming application.

An end user will be able to choose whether to run:

- Audio streaming only
- MasjidBoard display only
- Both audio streaming and MasjidBoard display

The two applications/subsystems must remain operationally independent.

## Intended Architecture

```text
                         MasjidPi
                            |
              +-------------+-------------+
              |                           |
              v                           v
      Audio Application            MasjidBoard Application
              |                           |
             MPV                    MasjidBoard Live
                                          |
                              +-----------+-----------+
                              |                       |
                         Data / Cache            Display
                              |                       |
                              +-----------+-----------+
                                          |
                                         HDMI
```

The MasjidBoard application should be independently executable and should not require the audio application to be running.

## Application Boundaries

### Audio application

Responsible for the existing MasjidPi audio-streaming functionality, including stream selection, playback, volume and related audio controls.

### MasjidBoard application

Responsible for:

- MasjidBoard Live data retrieval
- Normalisation of upstream data
- Persistent local caching
- Image/media caching
- Board state
- Content scheduling
- Display rendering
- HDMI output
- MasjidBoard-specific configuration and recovery

MasjidBoard must not depend on MPV or the audio playback subsystem.

## Repository Strategy

MasjidBoard remains in the **MasjidPi repository** so that shared infrastructure, configuration, documentation, testing and release management remain together.

It should nevertheless be implemented as a distinct application/package boundary rather than being embedded directly into the audio playback code.

The eventual implementation may expose separate executables/services, for example:

```text
/opt/masjidpi/bin/masjidpi
/opt/masjidpi/bin/masjidboard
```

and corresponding independently managed services:

```text
masjidpi.service
masjidboard.service
```

Exact executable/service names are implementation details and are not yet fixed.

## Configuration

The MasjidPi configuration/UI should allow the end user to enable either application independently or both together.

Conceptually:

```text
Audio Streaming:   [ enabled / disabled ]
MasjidBoard:       [ enabled / disabled ]
```

This configuration must not cause one subsystem to fail simply because the other is disabled or unavailable.

## Failure Isolation

The architecture must support independent recovery:

```text
MasjidBoard failure
    -> Audio continues
    -> MasjidBoard restarts/recover independently

MasjidBoard Live unavailable
    -> Audio continues
    -> MasjidBoard displays cached data where possible

Audio failure
    -> MasjidBoard continues
    -> Audio restarts/recover independently
```

## Display Strategy

MasjidBoard will **not** simply embed or mirror the MasjidBoard Live webpage.

MasjidBoard Live is the primary data/content provider. MasjidPi will implement its own presentation layer using the structured upstream data discovered during the research phase.

The initial goal is **functional and content parity** with MasjidBoard Live: all available information and relevant display behaviour should be represented. The implementation should not unnecessarily reproduce MasjidBoard Live's HTML, CSS or JavaScript architecture.

This gives MasjidPi control over:

- HDMI output
- 1920×1080 display layout
- Content scheduling
- Offline behaviour
- Local caching
- Image caching
- Rendering performance
- Recovery behaviour
- Future MasjidPi-specific display improvements

## Data Flow

```text
MasjidBoard Live
        |
        v
MasjidBoard Provider
        |
        v
Normalised MasjidBoard Model
        |
        +----> Persistent Cache
        |
        v
Board State / Scheduler
        |
        v
MasjidBoard Renderer
        |
        v
HDMI Display
```

The upstream MasjidBoard Live positional API schema must remain isolated inside the provider/normalisation layer.

## Initial Data Model

The internal model should be **semantic and collection-oriented**, rather than reproducing MasjidBoard Live's positional spreadsheet/API rows. The upstream provider is responsible for translating the opaque upstream structure into this model.

The model deliberately distinguishes **core board data** from **optional enrichment**. Masjid identity and the five daily prayer times are required for a usable board dataset. All other content is optional and may legitimately be absent for a particular masjid or day.

### Board / Masjid

```text
Board
├── Masjid identity              REQUIRED
├── English name
├── Arabic name
├── Location
├── Timezone
├── Language / display settings
├── Gregorian date
└── Islamic date
```

The model should retain the source board identifier and public masjid identifier where useful for diagnostics and refresh operations.

**Current upstream finding:** the captured MasjidBoard Live response provides the Gregorian date directly in row 2 and provides an Islamic calendar year (`١٤٤٨`) in row 6, but it does **not yet provide a verified complete Islamic date** that can safely populate `DateContext.IslamicDate`. The row-2 values labelled around “Islamic Time” are display/configuration settings rather than a complete Islamic date. We therefore leave `IslamicDate` empty until its upstream source and semantics are verified. We must not derive or guess the Islamic date from the Gregorian date merely to populate this field.

### Daily prayer times — core board data

Prayer times are the **primary purpose of MasjidBoard** and are therefore required core data rather than optional content.

```text
PrayerTimes                    REQUIRED
├── Fajr                         REQUIRED
├── Zuhr                         REQUIRED
├── Asr                          REQUIRED
├── Maghrib                      REQUIRED
└── Esha                         REQUIRED
```

Each prayer should be represented semantically rather than as separate upstream fields, with Adhan and Jamaah values where supplied.

A board dataset should be considered usable only when it has a valid masjid identity and a usable schedule for the five daily prayers. Temporary upstream failures should not invalidate an existing cached schedule; the cache may continue to provide the last known valid prayer data while marking it as stale.

### Perpetual / astronomical times — optional

```text
AstronomicalTimes
├── Suhur
├── Fajr start
├── Sunrise
├── Ishraaq
├── Duha
├── Istiwā / solar noon
├── Zuhr start
├── Asr Shafi'i
├── Asr Hanafi
├── Sunset
└── Esha start
```

Where MasjidBoard Live supplies additional related times, the model should be extensible without changing the basic prayer model. These times enrich the board but are not required for the board to be considered usable.

### Jumu'ah — optional

```text
JumuahServices[]
├── Title / service label
├── Adhan
├── Lecture / Sunan
├── Khutbah
├── Salah
└── Khateeb
```

Multiple services must be supported. The model must not assume a single Jumu'ah.

### Announcements and programmes — optional

```text
Announcements[]
Programmes[]
Notices[]
```

Each item should be able to carry its content plus relevant display/visibility/scheduling metadata where the upstream source supplies it.

This covers, as applicable:

- Configurable announcements
- Weekly programmes
- Taleem
- Dawah / Gasht
- Community notices
- Other activities

### Notices — optional

The initial model should support distinct semantic notice types where this improves rendering or scheduling, including:

- Nikah
- Funeral
- Sickness / well-wishes
- Community notices
- Eid / Eidgah information

The provider should not force every upstream notice into a rigid subtype if the source does not provide enough information to distinguish it reliably.

### Media — optional

```text
Media[]
├── Source/media identifier
├── Media type
├── Cached/local reference
├── Visibility
├── Display duration (if supplied)
└── Presentation metadata
```

Posters and images are first-class board content. Media should be cached locally so that display operation is not dependent on repeatedly downloading remote assets.

### Banking / contribution information — optional

```text
Banking
ContributionInformation
```

These should remain semantic data rather than being tied to a particular upstream field position.

### New Moon / lunar information — optional

```text
NewMoon
```

The model should allow the upstream source's New Moon information to be represented without assuming that every board has populated it.

### Shared daily Islamic content — optional and implemented

Ayah, Hadith and Sunnah were not required for the initial implementation and
remain optional enrichment. They are now obtained from MasjidBoard Live's
shared public translations feed through a dedicated client, semantic model
and cache outside the per-board provider path.

The service loads last-known-good content at startup and refreshes it at most
once per Africa/Johannesburg calendar day. The display API applies the user's
independent Ayah, Hadith and Sunnah preferences. This boundary preserves the
important rule that generic content is available regardless of which masjids
are selected or which Premium features those masjids enable.

## Data Model Principles

1. **Core versus enrichment** — masjid identity and the five daily prayer times are required core data; all other content is optional enrichment.
2. **Semantic over positional** — the application must not expose MasjidBoard Live's row/column structure outside the provider.
3. **Graceful degradation** — optional content may be absent without making the board invalid; cached core prayer data may be used when the upstream source is temporarily unavailable.
4. **Collections for repeated content** — announcements, programmes, Jumu'ah services and media must support multiple entries.
5. **Source metadata retained where useful** — source IDs and identifiers should be retained where they help refresh, cache, diagnose or reproduce content.
6. **Display metadata is separate from content** — content should describe what exists; scheduling/rendering determines how and when it appears.
7. **Extensible** — future MasjidBoard Live fields and deferred content must be addable without redesigning the entire model.
8. **Do not infer unsupported upstream semantics** — if a value is not demonstrably a complete semantic field, leave the corresponding optional model field absent rather than deriving or guessing it.

## Data Provider Boundary

The MasjidBoard Live provider should have three conceptual responsibilities:

```text
MasjidBoard Live
      |
      v
Fetch upstream data/media
      |
      v
Parse positional schema
      |
      v
Validate / normalise
      |
      v
Normalised MasjidBoard Model
```

The rest of the application should not need to know that the source uses a 29-row positional response.

The provider should also expose enough information for the cache layer to distinguish:

- successfully refreshed data
- stale but usable cached data
- missing optional content
- failed media downloads
- invalid/unusable upstream responses

The provider must treat missing optional enrichment differently from missing required prayer data. Missing optional content is normal; missing a usable five-prayer schedule is a data/recovery condition.

## Cache Boundary

The cache should persist the **normalised model and required media**, rather than only the raw upstream response.

This allows the display system to operate without MasjidBoard Live being reachable and keeps the display layer independent from upstream schema changes.

The raw upstream response may optionally be retained for diagnostics, but it should not be the application's primary runtime data representation.

The last known valid core prayer schedule should remain available in the cache when possible, even when a subsequent refresh fails. Cached data should carry enough metadata to determine its freshness/staleness.

## Scheduler / Board State

The scheduler should operate on the normalised model and determine the current presentation state.

Conceptually:

```text
Normalised board data
        |
        v
Board scheduler
        |
        +--> current time / prayer state
        +--> configured slide order
        +--> content availability
        +--> display durations
        +--> special/seasonal content
        |
        v
Current slide/state
```

The scheduler should not need to know how the upstream API is structured.

Prayer times should form the primary scheduler context. Optional content is scheduled around the core prayer-time presentation rather than being a prerequisite for board operation.

## Renderer

The renderer receives a presentation state and renders it for the HDMI display.

The initial target is a **1920×1080 Full HD display**.

The renderer should be lightweight and should not require a full web browser merely to reproduce MasjidBoard content.

The initial implementation should aim for functional/content parity rather than pixel-perfect reproduction of MasjidBoard Live's webpage.

## Future Extensibility

Although MasjidBoard Live is the primary data source, the provider boundary should allow another data source to be introduced later without redesigning the display system.

The architecture should therefore distinguish between:

```text
Data provider
        -> Normalised board data
        -> Display system
```

rather than coupling the display directly to MasjidBoard Live.

## Consequences

### Benefits

- Audio-only installations remain lightweight and unaffected by HDMI requirements.
- Display-only installations are possible.
- Audio and display can be enabled together.
- Failures are isolated between subsystems.
- MasjidBoard can have its own caching and recovery behaviour.
- The display is not dependent on the MasjidBoard Live website implementation.
- The MasjidBoard application can evolve independently while remaining part of the MasjidPi project.
- The internal data model remains independent of MasjidBoard Live's positional schema.
- Prayer times remain available as the core purpose of the board even when optional enrichment is absent.
- Shared Ayah/Hadith/Sunnah content was added without redesigning the per-board provider/display architecture.

### Trade-offs

- There is additional application/service architecture compared with simply displaying the website.
- MasjidPi must maintain its own renderer and display behaviour.
- Changes to the upstream MasjidBoard Live data model may require provider updates.
- A normalisation layer introduces some implementation work before rendering can begin.

These trade-offs are accepted because reliability, independence and appliance-style operation are more important than minimising implementation effort.

## Implementation Guardrail

The architecture is now sufficiently defined to begin implementation planning, but production code should still be introduced in stages.

The first implementation should establish the provider and normalised model before building the complete renderer. The initial scope does not need to block on Ayah, Hadith or Sunnah.
