# MasjidBoard Domain Model

**Status:** Design — proposed for review  
**Branch:** `research/masjidboard-live`

## Purpose

This document defines the semantic domain model that will sit between the MasjidBoard Live provider and the MasjidBoard scheduler/display system.

The model deliberately does **not** mirror MasjidBoard Live's positional 29-row response. The provider is responsible for translating the upstream structure into these semantic types.

## Core Principle

MasjidBoard is fundamentally a prayer-times display.

The model therefore has two levels of importance:

1. **Required core data** — masjid identity, date/context, and the five daily prayer times.
2. **Optional enrichment** — everything else supplied by MasjidBoard Live.

Missing optional enrichment is normal and must not make a board invalid.

A missing or unusable five-prayer schedule is a genuine data/recovery condition.

## Model Overview

```text
Board
├── Identity                  REQUIRED
├── DateContext               REQUIRED
├── PrayerTimes               REQUIRED
│   ├── Fajr
│   ├── Zuhr
│   ├── Asr
│   ├── Maghrib
│   └── Esha
│
├── AstronomicalTimes         OPTIONAL
├── JumuahServices[]          OPTIONAL
├── Announcements[]           OPTIONAL
├── Programmes[]              OPTIONAL
├── Notices[]                 OPTIONAL
├── Media[]                   OPTIONAL
├── Banking                   OPTIONAL
├── ContributionInformation   OPTIONAL
├── NewMoon                   OPTIONAL
└── DisplayConfiguration      OPTIONAL

SharedDailyIslamicContent     OPTIONAL / SERVICE-LEVEL
```

## Board

`Board` is the top-level normalised representation of one masjid's board data for a particular date/context.

Conceptually:

```text
Board
├── Identity
├── DateContext
├── PrayerTimes
├── AstronomicalTimes?
├── JumuahServices[]
├── Announcements[]
├── Programmes[]
├── Notices[]
├── Media[]
├── Banking?
├── ContributionInformation?
├── NewMoon?
└── DisplayConfiguration?
```

The exact Go representation is intentionally left open until implementation.

## Identity

Identity represents the masjid itself rather than the display configuration.

It should be capable of carrying:

- Source board/masjid identifier.
- Public masjid identifier where supplied.
- Primary masjid name.
- Alternate/local-language name where supplied.
- Arabic name where supplied.
- Location information where supplied.
- Timezone.

The source identifiers are useful for diagnostics, refreshes and cache association and should not be discarded during normalisation.

## DateContext

Prayer times are date-specific. The model must make the date and timezone context explicit.

It should carry enough information to establish:

- Gregorian date represented by the board.
- Islamic/Hijri date where supplied.
- Timezone used for local prayer times.
- Whether the board data represents the current configured date or another explicit date.

The model should use Go's standard time handling rather than inventing a custom clock representation unless a source-specific format must be retained for diagnostics.

## PrayerTimes — Required Core

The five daily prayers are the core of the model.

```text
PrayerTimes
├── Fajr
├── Zuhr
├── Asr
├── Maghrib
└── Esha
```

Each prayer should contain semantic timing information rather than source row/column positions.

Conceptually:

```text
PrayerTime
├── Prayer
├── Adhan
└── Jamaah
```

### Prayer identity

The model should use a finite semantic prayer identity for the five daily prayers rather than arbitrary strings. This prevents the provider and scheduler from having to compare display labels.

### Adhan and Jamaah

Adhan and Jamaah are separate values because MasjidBoard Live can supply both.

The model should permit a value to be absent at the individual prayer level if the upstream source genuinely does not provide it, while still requiring the prayer itself to exist in the core schedule.

### Time representation

Prayer times should be represented as local wall-clock times associated with the `Board` timezone/date context, rather than as UTC timestamps detached from the masjid's local date.

The implementation should avoid encoding prayer times as arbitrary formatted strings in the domain model. Formatting belongs to the display layer.

## Prayer Schedule Validity

A normalised board is considered to have usable core prayer data when:

- Identity is valid.
- Date/context is valid enough to establish the represented day and timezone.
- Fajr, Zuhr, Asr, Maghrib and Esha each have a usable local prayer time.

Optional Adhan/Jamaah values do not independently invalidate the prayer schedule if the upstream source does not provide them.

A provider refresh with missing optional enrichment remains successful.

A provider refresh without a usable five-prayer schedule should be treated as invalid for replacement of the current cached core schedule.

## AstronomicalTimes — Optional

These are supporting daily times supplied by MasjidBoard Live.

```text
AstronomicalTimes
├── Suhur
├── FajrStart
├── Sunrise
├── Ishraaq
├── Duha
├── Istiwa / SolarNoon
├── ZuhrStart
├── AsrShafii
├── AsrHanafi
├── Sunset
└── EshaStart
```

The model should allow individual values to be absent.

These times are enrichment. They are not required for the board to be valid because the five daily prayer times are the primary purpose of the board.

## JumuahService — Optional

A board may have zero, one or several Jumu'ah services.

```text
JumuahService
├── Label / title
├── Adhan
├── Lecture / Sunan
├── Khutbah
├── Salah
└── Khateeb
```

The model must not assume that every service has every field populated.

Different MasjidBoard Live boards use different combinations of Jumu'ah information, so optional fields are preferable to a rigid single format.

## Announcements — Optional

Announcements are repeatable board content.

Conceptually:

```text
Announcement
├── Content
├── Title (if supplied)
├── Language (if supplied)
└── Presentation metadata (if supplied)
```

The model should use a collection rather than fixed fields such as `Announcement1`, `Announcement2`, etc.

The provider is responsible for translating MasjidBoard Live's numbered slots into this collection.

## Programmes — Optional

Programmes cover recurring or scheduled masjid activities such as Taleem, Dawah/Gasht and other programmes.

Conceptually:

```text
Programme
├── Title
├── Description/content
├── Day/date information
├── Time information
└── Presentation metadata
```

Not every field will be populated for every programme.

## Notices — Optional

Notices represent informational items that do not need to be modelled as announcements or programmes.

The initial semantic categories may include:

- Nikah
- Funeral
- Sickness / well-wishes
- Community notice
- Eid / Eidgah

The category should remain extensible. The provider should not invent a category when the source does not support a reliable distinction.

Conceptually:

```text
Notice
├── Type
├── Title
├── Content
├── Date/time (if supplied)
└── Presentation metadata
```

## Media — Optional

Images and posters are first-class board content.

```text
Media
├── Source identifier/reference
├── Media type
├── Remote source/reference
├── Local cached reference
├── Visibility/scheduling metadata
├── Display duration (if supplied)
└── Presentation metadata
```

The domain model should represent the media independently from the renderer. The cache is responsible for ensuring a usable local representation when possible.

The model should support at least:

- Standard posters
- Large posters
- Community posters
- Eid/Eidgah material
- Other image-based board content

## Banking and Contribution Information — Optional

Banking/contribution information should be represented semantically rather than tied to upstream positional fields.

The model should allow the source to provide structured information where available without requiring every board to have banking information.

## NewMoon — Optional

New Moon/lunar information is represented as an optional semantic object.

The exact fields should only be added when they are supported by the verified upstream data we have chosen to implement.

We should avoid inventing fields simply because a possible future implementation might use them.

## DisplayConfiguration — Optional

Display settings supplied by MasjidBoard Live should be kept separate from content data where practical.

Examples include:

- Display language.
- RTL behaviour.
- Theme/style selection.
- Clock configuration.
- Slide/display duration settings.
- Other board presentation preferences supplied by the source.

The normalised model may preserve source display configuration, but the scheduler/renderer remains responsible for interpreting it.

## Shared Daily Islamic Content — Optional

Ayah, Hadith and Sunnah were excluded from the required initial board model.
Their upstream source and behaviour are now verified and implemented as
service-level enrichment rather than fields on `Board`:

```text
DailyIslamicContent
├── Ayah
│   ├── Surah
│   ├── AyahNumber
│   └── Text
├── Hadith
│   ├── Heading
│   ├── Text
│   └── Reference
├── Sunnah
│   ├── Heading
│   ├── Text
│   └── Reference
├── Language
├── Source
├── SourceURL
├── ContentDate
└── FetchedAt
```

This content is shared by MasjidBoard Live and is not attributed to, enabled
by or cached inside any selected masjid. The display layer independently
filters its three categories using local user preferences.

## Source Metadata

The normalised model should retain source metadata where it has operational value, but source-specific implementation details should not leak into normal domain objects.

Useful metadata may include:

- Provider/source name.
- Source board identifier.
- Retrieval timestamp.
- Data date.
- Last successful refresh.
- Data freshness/staleness state.

The exact metadata placement may ultimately belong in a cache envelope rather than the `Board` itself.

## Cache Envelope

The runtime domain model and the persisted cache representation should be considered separate concerns.

Conceptually:

```text
CachedBoard
├── Board
├── RetrievedAt
├── LastSuccessfulRefresh
└── Freshness/validation metadata
```

The cache must preserve the last known valid core prayer schedule even when a subsequent provider refresh fails.

Optional content may be stale, missing or partially unavailable without invalidating the core board.

## Normalisation Rules

The provider must perform the following translation:

```text
MasjidBoard Live positional response
              ↓
       Provider-specific parsing
              ↓
       Validation of core data
              ↓
       Semantic normalisation
              ↓
          Board model
```

The following must **not** happen outside the provider:

- indexing into the 29-row response
- interpreting source-specific field positions
- comparing source-specific labels to identify prayers
- reconstructing numbered announcement slots
- interpreting MasjidBoard Live-specific media fields

## Optionality Rules

The following distinction is intentional:

```text
Missing Fajr/Zuhr/Asr/Maghrib/Esha time
    → core data problem

Missing Jumu'ah
    → normal

Missing announcements
    → normal

Missing posters
    → normal

Missing astronomical times
    → normal

Missing banking information
    → normal

Missing New Moon information
    → normal
```

This distinction should be reflected in validation and cache replacement logic.

## What the Model Does Not Do

The domain model should not:

- fetch data
- download media
- render graphics
- decide slide order
- calculate the current display state
- know about HDMI
- know about MPV
- expose the MasjidBoard Live 29-row schema
- format times specifically for a particular screen layout

Those responsibilities belong to the provider, cache, scheduler and display layers respectively.

## Design Review Questions

Before implementing the Go types, confirm:

1. Whether `PrayerTime` should store Adhan/Jamaah as nullable values or optional typed fields.
2. Whether `DateContext` belongs directly on `Board` or partly in the cache envelope.
3. Whether source identifiers belong in `Identity`, a source metadata object, or both.
4. Whether display configuration belongs in the domain model or a separate board configuration object.
5. Whether `Banking` and `ContributionInformation` should be separate types or one contribution/information structure.
6. The exact representation of local wall-clock prayer times in Go.

These are implementation-level decisions and do not alter the agreed core principle: **the five daily prayer times are required; everything else is optional enrichment.**
