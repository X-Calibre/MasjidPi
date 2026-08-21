# MasjidBoard Live Core Provider Validation

**Status:** Research milestone validated  
**Branch:** `research/masjidboard-live`

## Purpose

Record the current experimental validation status of the MasjidBoard Live Core provider path and its handoff from FindMasjid discovery.

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

## Live Core End-to-End Validation

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

## Live Discovery-to-Provider Validation

Also on 2026-08-18, the structured FindMasjid discovery path was tested live for Brits, South Africa.

The endpoint returned three boards:

```text
Jamiah Yusuf Darul Uloom Brits  -> brits-darul-uloom
Brits Jamia Masjid              -> brits-jamia
Masjid Taqwa                    -> brits-taqwa
```

`brits-jamia` was selected from the returned `CatalogueEntry` and passed through the formal provider handoff:

```text
FindMasjid
    -> CatalogueEntry
    -> NewCoreClientFromCatalogue(...)
    -> CoreClient
    -> live Core board
    -> model.Board
```

The handoff derived the provider identity without caller-side timezone construction:

```text
ID        brits-jamia
Name      Brits Jamia Masjid
Timezone  GMT+02:00
```

The subsequently retrieved live board preserved the same identity and returned the expected timetable:

```text
Fajr Jamaah      06:00
Dhuhr Jamaah     13:20
Asr Jamaah       17:00
Maghrib Adhan    17:54
Esha Jamaah      19:30
```

The discovery record and the independently retrieved Core board also agreed on the opaque upstream identifier:

```text
FindMasjid MBL_ID   MBL11517PRP
Core mbl_number     MBL11517PRP
```

The handoff code preserves the full millisecond timezone offset and formats fixed-offset zones as `GMT±HH:MM`, including fractional offsets such as `+05:30` and `-03:30`.

This validates the discovery -> catalogue entry -> Core provider path experimentally against the live service. It does not make MasjidBoard release-ready; catalogue persistence/refresh, selected-board persistence, caching, broader failure handling, Premium enrichment and display integration remain separate work.

## Next Work Item

The next boundary is **Stage 3 catalogue design and persistence strategy**.

Before adding user-facing selection or scheduled catalogue refresh, define:

1. the stable MasjidPi catalogue record and which FindMasjid fields belong in it;
2. which upstream fields remain provider-only metadata;
3. how catalogue entries are keyed, renamed and reconciled;
4. how discovery timestamps and validation state are recorded;
5. how a selected board is persisted independently from catalogue refresh;
6. how last-known-good catalogue data is retained if discovery fails; and
7. how catalogue refresh remains independent from retrieval of an already-selected board.

The intended boundary remains:

```text
FindMasjid discovery
    -> normalised local catalogue
    -> user selects board
    -> persist selected board identity
    -> Core provider
    -> cached normalised Board
```

Discovery and individual-board retrieval must remain independent so a catalogue refresh failure does not prevent an already-selected board from continuing to operate from cached configuration/data.
