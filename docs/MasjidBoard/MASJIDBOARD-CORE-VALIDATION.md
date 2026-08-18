# MasjidBoard Live Core Provider Validation

**Status:** Research milestone validated  
**Branch:** `research/masjidboard-live`

## Purpose

Record the current experimental validation status of the MasjidBoard Live Core provider path before discovery/catalogue integration continues.

## Working Source Architecture

The current research supports the following separation of responsibilities:

```text
FindMasjid structured directory
    -> discovery / catalogue
    -> public `web_url` slug
    -> board identity and geography

Core board page
https://masjidboardlive.com/boards/?<web_url>
    -> primary standard timetable source
    -> embedded `let data = {...}` object
    -> prayer times
    -> astronomical times
    -> Jumu'ah slots/configuration
    -> optional livestream metadata

Premium board
https://premium.masjidboardlive.com/v2/?mid=<web_url>
    -> optional enrichment source
    -> richer Jumu'ah and community/display content
    -> announcements, posters and other Premium-only material
```

Premium is therefore not required for a useful standard timetable.

## Core Schema Validation

The public Core `let data = {...}` object was sampled across 11 boards in South Africa and Australia.

Every usable sample exposed the same 36-field schema.

The sample included multiple observed `MBL_ID` forms and both Premium-capable and Core-only candidates.

The parser remains defensive because schema stability cannot be guaranteed solely from the sample, but the evidence is sufficient to treat the Core object as the working standard timetable interface.

## Same-board Core/Premium Validation

Two known Premium-capable boards were compared field-for-field between Core and Premium sources:

- `brits-jamia`
- `erasmia-aaisha`

For both boards, the overlapping standard prayer timetable and astronomical values matched. Jumu'ah slot values also matched the corresponding Premium events.

This is the basis for the working decision that Core is the primary standard timetable source and Premium is enrichment rather than a prerequisite for timetable retrieval.

## Freshness Findings

`last_updated` is shared upstream metadata between FindMasjid and Core when populated, but it may also be blank.

MasjidBoard Live supports perpetual/annual timetable behaviour as well as manual timetable maintenance and temporary overrides. MasjidPi therefore must not infer timetable validity from the age or presence of `last_updated` alone.

The provider should track its own retrieval/validation freshness independently.

Do not infer upstream timetable-maintenance mode from:

- `MBL_ID`;
- `last_updated`;
- `nextChangeDisplay`; or
- next-change fields.

MasjidPi consumes the resolved current timetable supplied by MasjidBoard Live.

## Jumu'ah and Placeholder Handling

Core Jumu'ah consists of up to three time slots plus a heading-code configuration. Slot meaning must be derived from the heading code rather than position alone.

Observed heading codes include:

```text
0 -> Adhan
1 -> Lecture
3 -> Sunan
6 -> Khutbah
```

Core placeholder values include empty strings and values such as `~~~~`. These are normalised as unavailable rather than parsed as clock times.

Khutbah remains semantically distinct from Jumu'ah Salaah/Jamaah. Where a dedicated Salaah/Jamaah value is unavailable, the existing domain-model fallback may use Khutbah as the effective Friday Salaah display time.

## Live End-to-End Validation

On 2026-08-18, the production-shaped Core path was tested live against `brits-jamia`.

The path successfully completed:

```text
MasjidBoard Live
    -> GET /boards/?brits-jamia
    -> HTML response
    -> extract embedded `let data = {...}`
    -> parse Core object
    -> normalise into `model.Board`
    -> preserve provider metadata
```

The live result matched the previously captured values, including:

```text
Fajr Adhan       05:40
Fajr Jamaah      06:00
Dhuhr Adhan      13:00
Dhuhr Jamaah     13:20
Asr Adhan        16:40
Asr Jamaah       17:00
Maghrib Adhan    17:54
Esha Adhan       19:15
Esha Jamaah      19:30
Sunrise          06:34
Ishraaq          06:49
Sunset           17:51
Esha Start       19:08
```

Jumu'ah resolved as:

```text
Adhan    12:25
Lecture  12:40
Khutbah  13:00
```

Provider metadata also matched the capture:

```text
MBL number       MBL11517PRP
Last updated     Sun, 22 Mar 2026, 12:47:25
Stream provider  SmartBilal
Stream URL       https://media.smartbilal.com/masjid/britsj
```

This validates the complete Core provider path experimentally. It does not make MasjidBoard release-ready; caching, discovery/catalogue integration, broader failure handling and display integration remain separate work.

## Next Work Item

The next implementation/research boundary is discovery/catalogue integration:

```text
FindMasjid directory
    -> discover board
    -> normalise catalogue entry
    -> obtain `web_url`
    -> construct selected-board identity
    -> CoreClient
    -> model.Board
```

Discovery must remain independent from individual-board retrieval so a catalogue refresh failure does not prevent an already-selected board from continuing to operate from cached configuration/data.
