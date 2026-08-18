# MasjidBoard Live Capability Matrix Research

**Status:** Research / Stage 2  
**Branch:** `research/masjidboard-live`

## Purpose

Record the verified capabilities of the three public MasjidBoard Live access paths used in the research:

1. FindMasjid discovery/catalogue data;
2. the public Core board page; and
3. the Premium board page.

This document is separate from `MASJIDBOARD-DISCOVERY.md`, which covers catalogue enumeration, and `MASJIDBOARD-LIVE.md`, which covers detailed board semantics and Premium positional parsing.

## Reference Boards

The current research has used:

- Core-only candidate: `erasmia-abu-bakr`
- Core-only candidate: `fawkner-rahman`
- Premium-capable: `erasmia-aaisha`
- Premium-capable: `brits-jamia`
- Premium-capable: `fawkner-masjid`

Core-board schema validation has been performed across South Africa and Australia.

Direct Core/Premium same-board timetable comparisons have now been performed for:

- `brits-jamia`
- `erasmia-aaisha`

## Three Distinct Public Data Paths

```text
FindMasjid endpoint
    -> discovery and catalogue metadata

https://masjidboardlive.com/boards/?<web_url>
    -> Core board data
    -> embedded `let data = {...}` object
    -> standard timetable and astronomical values

https://premium.masjidboardlive.com/v2/?mid=<web_url>
    -> Premium board data when available
    -> embedded `boardId` and `theInfo`
    -> richer board/community/display content
```

These are distinct responsibilities. A board does not need Premium capability to provide a useful full timetable to MasjidPi.

## FindMasjid Discovery Data

The public FindMasjid endpoint exposes records containing fields such as:

```text
masjid
fajr_jamaat
zuhr_jamaat
asr_jamaat
maghrib_adhan
esha_jamaat
last_updated
MBL_ID
city
sunset
time_zone_milli
web_url
jumuah_khutbah
ramadhaanactive
date_adjust
moon_seen
ladies_facility
```

This is the preferred discovery/catalogue source because it provides the public `web_url` slug and useful location/summary metadata without retrieving each individual board first.

## Verified Core Board Interface

The public Core board page:

```text
https://masjidboardlive.com/boards/?<web_url>
```

embeds a JavaScript object:

```javascript
let data = {
    ...
}
```

The generated HTML itself contains the timetable object. `/boards/script.js` consumes it for display; no additional timetable API call was required in the observed boards.

### Verified Core schema

The exact same field set was observed for `erasmia-abu-bakr` and `fawkner-rahman`:

```text
lang
islamicDateAdjust
forceDate30
mbl_number
sehriEnds
fajrStarts
fajrAthan
fajrJamaah
sunrise
ishraaq
duha
istiwaCaution
istiwa
zawaalEnd
dhuhrAthan
dhuhrJamaah
jumuahTime1
jumuahTime2
jumuahTime3
asrShafi
asrHanafi
asrAthan
asrJamaah
sunset
maghribAthan
eshaStarts
eshaAthan
eshaJamaah
last_updated
jumuahHeadings
twentyfourhrtime
liveStreamP
liveStreamURL
ramadaanHide
customcode
sunday_zuhr_text
```

There were no key differences between these two Core-only candidates despite being in different countries.

This is strong evidence for a dedicated Core-board parser/provider. The eventual parser must nevertheless remain defensive because the research sample is not yet large enough to prove the schema is immutable across all boards.

### Core placeholder values

Core fields can contain placeholders rather than usable clock times.

For example, `fawkner-rahman` exposed:

```text
jumuahTime1 = "~~~~"
jumuahTime2 = "~~~~"
jumuahTime3 = "~~~~"
```

Values such as the following must be normalised as unavailable where semantically appropriate:

```text
""
~~~~
Hide
-
–
—
```

Other placeholder values already observed elsewhere in MasjidBoard Live research should be handled consistently by the provider rather than allowed to reach the domain model as clock times.

## Premium Capability Probe

A discovered `web_url` can be probed using:

```text
https://premium.masjidboardlive.com/v2/?mid=<web_url>
```

### Premium available

A valid Premium page exposes both:

```javascript
let boardId = "...";
let theInfo = [...];
```

Verified examples include:

- `erasmia-aaisha`
- `brits-jamia`
- `fawkner-masjid`

### Premium unavailable

For both `erasmia-abu-bakr` and `fawkner-rahman`, the Premium page returned the application-level error page:

```text
MasjidBoard live - 500
This masjid does not exist
```

and exposed neither `boardId` nor `theInfo`.

Both boards nevertheless remained usable through the public Core board page.

### Premium probe state

Premium capability should be represented with at least three states:

```text
available
unavailable
unknown
```

`unknown` is required for transient DNS, timeout, transport, malformed-response or unexpected-server failures. A temporary retrieval failure must not be interpreted as proof that Premium is unavailable.

HTTP status alone is not sufficient; the generated Premium page structure must also be validated.

## Capability Matrix

| Capability | FindMasjid | Core board | Premium board |
|---|---|---|---|
| Catalogue discovery | Yes | No | Via FindMasjid |
| Public `web_url` / `mid` | Yes | Used for retrieval | Yes |
| Masjid/city discovery metadata | Yes | Page metadata exists | Yes |
| Timezone offset | Yes | Not in observed `data` object | Yes |
| Last-updated value | Yes | Yes where supplied | Not yet defined as catalogue metadata |
| Fajr Adhan | No | Yes | Yes |
| Fajr Jamaah | Yes | Yes | Yes |
| Dhuhr Adhan | No | Yes | Yes |
| Dhuhr Jamaah | Yes | Yes | Yes |
| Asr Adhan | No | Yes | Yes |
| Asr Jamaah | Yes | Yes | Yes |
| Maghrib Adhan | Yes | Yes | Yes |
| Maghrib Jamaah | No reliable field | No dedicated field observed | Yes where supplied |
| Esha Adhan | No | Yes | Yes |
| Esha Jamaah | Yes | Yes | Yes |
| Sunset | Yes | Yes | Yes |
| Sehri/Fajr start | No | Yes | Yes |
| Sunrise/Ishraaq/Duha | No | Yes | Yes |
| Istiwa/Zawaal | No | Yes | Yes |
| Asr Shafi'i/Hanafi | No | Yes | Yes |
| Jumu'ah summary | Yes | Yes | Yes |
| Jumu'ah configured slots | No | Yes | Yes |
| Dedicated Premium Jumu'ah fields | No | No | Yes |
| Khateeb | No | Not observed | Yes where supplied |
| Live-stream metadata | Not established | Yes where configured | Yes where configured |
| Alternate-language values | No | Not observed | Yes |
| Announcements/community content | No | Not observed in `data` | Yes |
| Posters/programmes/notices | No | Not observed in `data` | Yes |
| Premium opaque `boardId` | No | No | Yes |
| Embedded `theInfo` | No | No | Yes |

## Same-board Core/Premium Validation

### Brits Jamia — `brits-jamia`

Core and Premium matched across all standard timetable fields checked:

| Field | Core | Premium |
|---|---:|---:|
| Fajr Adhan | 05:40 | 05:40 |
| Fajr Jamaah | 06:00 | 06:00 |
| Dhuhr Adhan | 13:00 | 13:00 |
| Dhuhr Jamaah | 13:20 | 13:20 |
| Asr Adhan | 16:40 | 16:40 |
| Asr Jamaah | 17:00 | 17:00 |
| Maghrib Adhan | 17:54 | 17:54 |
| Esha Adhan | 19:15 | 19:15 |
| Esha Jamaah | 19:30 | 19:30 |
| Sunset | 17:51 | 17:51 |

Core Jumu'ah data was:

```text
jumuahTime1    = 12:25
jumuahTime2    = 12:40
jumuahTime3    = 13:00
jumuahHeadings = 0,1,6
```

Using the verified heading mapping:

```text
0 = Adhan
1 = Lecture
6 = Khutbah
```

this produces:

```text
12:25 Adhan
12:40 Lecture
13:00 Khutbah
```

which matches the Premium Jumu'ah event sequence exactly.

### Masjid Aaisha — `erasmia-aaisha`

A second independent Premium-capable board produced the same result.

| Field | Core | Premium |
|---|---:|---:|
| Fajr start | 05:14 | 05:14 |
| Fajr Adhan | 05:45 | 05:45 |
| Fajr Jamaah | 06:00 | 06:00 |
| Sunrise | 06:32 | 06:32 |
| Ishraaq | 06:47 | 06:47 |
| Duha | 09:21 | 09:21 |
| Dhuhr Adhan | 12:25 | 12:25 |
| Dhuhr Jamaah | 12:40 | 12:40 |
| Asr Shafi'i | 15:25 | 15:25 |
| Asr Hanafi | 16:13 | 16:13 |
| Asr Adhan | 16:30 | 16:30 |
| Asr Jamaah | 16:45 | 16:45 |
| Sunset | 17:49 | 17:49 |
| Maghrib Adhan | 17:52 | 17:52 |
| Esha start | 19:06 | 19:06 |
| Esha Adhan | 19:10 | 19:10 |
| Esha Jamaah | 19:20 | 19:20 |

Core Jumu'ah data was:

```text
jumuahTime1    = 12:20
jumuahTime2    = 12:40
jumuahTime3    = 12:45
jumuahHeadings = 0,3,6
```

Using the verified heading mapping:

```text
0 = Adhan
3 = Sunan
6 = Khutbah
```

this produces:

```text
12:20 Adhan
12:40 Sunan
12:45 Khutbah
```

which again matches the Premium event sequence.

### Conclusion from same-board comparisons

Two independent Premium-capable boards now show Core and Premium agreeing field-for-field for the standard timetable and the astronomical/Jumu'ah values compared.

This is enough to adopt **Core as the working primary timetable source** for MasjidPi research.

It does not prove Core and Premium can never diverge, so production code should still be defensive and the decision should remain revisitable if future captures contradict it.

## Jumu'ah Semantics

The previously agreed Friday timetable behaviour remains unchanged:

```text
Friday standard five-prayer timetable:

Dhuhr Adhan  -> Jumu'ah Adhan
Dhuhr Salaah -> Jumu'ah Salaah/Jamaah

If Jumu'ah Salaah/Jamaah is unavailable:
    use Jumu'ah Khutbah as the Salaah-time fallback.
```

Core exposes three Jumu'ah slots and `jumuahHeadings`. Slot meaning must therefore be resolved from the heading configuration rather than inferred from position alone.

The Premium provider additionally exposes dedicated Jumu'ah fields and richer event data, which remains useful for the separate Friday detail element.

A Core or FindMasjid field labelled Khutbah must not automatically be treated as semantically identical to Jumu'ah Salaah/Jamaah merely because particular boards use the same time.

## Sunset / Iftar versus Maghrib Adhan

Core data confirms that sunset and Maghrib Adhan are separate concepts.

`fawkner-rahman`:

```text
sunset       = 17:50
maghribAthan = 17:55
```

`erasmia-abu-bakr`:

```text
sunset       = 17:49
maghribAthan = 17:52
```

MasjidPi must therefore retain separate semantic fields for sunset/Iftar and Maghrib Adhan.

During Ramadan, a masjid may display Iftar at sunset while delaying Maghrib Adhan by several minutes. Ramadan-specific data should be revalidated when boards are actively publishing those fields.

## Working Source Responsibilities

```text
MasjidBoard Live
        |
        +-- FindMasjid
        |       +--> discover boards
        |       +--> generate/update catalogue
        |       +--> location and summary metadata
        |
        +-- Core board page
        |       +--> PRIMARY standard timetable provider
        |       +--> Adhan/Jamaah values
        |       +--> astronomical values
        |       +--> Jumu'ah slots/configuration
        |       +--> optional live-stream metadata
        |
        +-- Premium board
                +--> OPTIONAL enrichment provider
                +--> richer Jumu'ah content
                +--> alternate-language values
                +--> announcements
                +--> posters
                +--> programmes
                +--> Nikah/funeral/community content
                +--> richer display/configuration data
```

Discovery, Core retrieval and Premium enrichment should remain separate responsibilities.

## Working Source-precedence Decision

The current Stage 2 decision is:

```text
FindMasjid
    -> discover board and obtain web_url

Core board page
    -> primary standard timetable source

Premium available
    -> enrich the board with Premium-only content
```

MasjidPi should not switch its standard timetable provider merely because Premium is available.

This has several advantages:

- Core-only and Premium-capable boards follow the same standard timetable path;
- Premium availability is no longer a prerequisite for normal prayer-time display;
- a Premium outage does not automatically remove standard timetable capability if Core remains available;
- the timetable provider can have one normalised schema and one defensive parsing path; and
- Premium can remain focused on genuinely richer content instead of duplicating timetable responsibility.

The production implementation should still define how to behave if future validation finds a Core/Premium disagreement.

## `MBL_ID` Capability Caveat

Known Premium boards have appeared with both `PRM` and `PRP` suffixes, while other entries expose suffixes including `CRM`, `CRP`, `CRS`, `EXT` and `PTA...EXT` forms.

MasjidPi must not derive Premium capability solely from these suffixes.

`MBL_ID` should remain opaque upstream metadata unless authoritative semantics are established later.

Premium capability should continue to be determined by the Premium capability probe.

## Remaining Stage 2 Questions

Stage 2 is not complete. Remaining research includes:

1. broader Core schema sampling across more countries and `MBL_ID` variants;
2. operational meaning and reliability of Core `last_updated`;
3. why some boards expose an empty `last_updated` value;
4. whether Core and Premium can temporarily diverge for the same board;
5. defined production behaviour if such a divergence is observed;
6. exact Core Jumu'ah behaviour across more `jumuahHeadings` combinations;
7. Ramadan-specific values, especially Iftar versus Maghrib behaviour;
8. how Core live-stream metadata relates to MasjidPi's existing audio-stream subsystem; and
9. which FindMasjid fields belong in the final user-facing catalogue versus selected-board data.

## Current Stage 2 Direction

The standard timetable source question is now sufficiently validated for research purposes:

**Core is the working primary timetable source; Premium is an optional enrichment source.**

The next Stage 2 work should focus on freshness, defensive parsing, broader schema sampling and remaining Ramadan/Jumu'ah semantics before production implementation is frozen.

No MasjidBoard work from this research branch is intended to be merged into the next MasjidPi release until the module is substantially further developed and validated.
