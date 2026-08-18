# MasjidBoard Live Capability Matrix Research

**Status:** Research / Stage 2  
**Branch:** `research/masjidboard-live`

## Purpose

Record verified comparisons between MasjidBoard Live discovery, Core board and Premium board data, and define the capabilities that MasjidPi can expect from each access path.

This document is intentionally separate from board discovery and from the positional Premium parser research.

## Reference Boards

The research has used both Core-only candidates and known Premium boards:

- Core-only candidate: `erasmia-abu-bakr`
- Core-only candidate: `fawkner-rahman`
- Known Premium board: `erasmia-aaisha`
- Known Premium board: `brits-jamia`
- Known Premium board: `fawkner-masjid`

The same-board FindMasjid/Premium comparisons were performed for `erasmia-aaisha` and `brits-jamia`.

The public Core board data interface was independently verified for `erasmia-abu-bakr` in South Africa and `fawkner-rahman` in Australia.

## Three Distinct Public Data Paths

Research now supports three separate upstream responsibilities:

```text
FindMasjid endpoint
    -> board discovery and catalogue metadata

https://masjidboardlive.com/boards/?<web_url>
    -> Core board data
    -> embedded `let data = {...}` object
    -> full standard timetable and astronomical values

https://premium.masjidboardlive.com/v2/?mid=<web_url>
    -> Premium board data when available
    -> embedded `boardId` and `theInfo`
    -> richer board/community/display content
```

The Core board page must therefore not be confused with the smaller FindMasjid catalogue record. A board can be useful to MasjidPi without having Premium capability.

## Premium Capability Probe

A Core catalogue `web_url` can be probed against:

```text
https://premium.masjidboardlive.com/v2/?mid=<web_url>
```

### Verified Premium success

A successful Premium page exposes both a server-supplied `boardId` and an embedded `theInfo` payload.

For example, `erasmia-aaisha` exposes:

```javascript
let boardId = "1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k";
let theInfo = [...];
```

`brits-jamia` and `fawkner-masjid` also resolve as valid Premium boards.

### Verified Premium absence

For `erasmia-abu-bakr`, the Premium page renders:

```text
MasjidBoard live - 500
This masjid does not exist
```

and does not expose `boardId` or `theInfo`.

The same negative pattern was independently observed for `fawkner-rahman` in Australia.

Both boards nevertheless expose usable Core board data through the public `/boards/?<web_url>` page.

### Probe-state decision

Premium capability must not be represented as a simple boolean derived from HTTP status.

The research model should distinguish at least:

```text
available
unavailable
unknown
```

Suggested semantics:

```text
available
    Generated Premium page contains a valid boardId and theInfo payload.

unavailable
    Premium page explicitly reports that the masjid does not exist and does not expose boardId/theInfo.

unknown
    Timeout, DNS failure, transport failure, unexpected server response,
    malformed page, or other condition where Premium capability cannot be
    determined reliably.
```

A transient failure must never be interpreted as proof that a board is Core-only.

## Verified Core Board Interface

The public Core board page:

```text
https://masjidboardlive.com/boards/?<web_url>
```

embeds a JavaScript object of the form:

```javascript
let data = {
    ...
}
```

This object is present in the generated HTML and is consumed by `/boards/script.js`. No additional timetable API call is required to obtain the observed Core data.

The same field set was captured from both `erasmia-abu-bakr` and `fawkner-rahman`:

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

There were no schema-key differences between these two captures despite the boards being in different countries.

This is strong evidence for treating the public Core board page as a distinct individual-board data source. It does not yet prove that every MasjidBoard Live Core board always exposes exactly this schema, so the eventual parser must remain defensive.

### Core placeholder values

Core values are not guaranteed to be usable clock times. For example, `fawkner-rahman` exposed:

```text
jumuahTime1 = "~~~~"
jumuahTime2 = "~~~~"
jumuahTime3 = "~~~~"
```

Values such as `~~~~`, `Hide`, empty strings and other established upstream placeholders must be normalised as absent/unavailable where appropriate rather than treated as malformed prayer times.

## Revised Capability Matrix

| Capability | FindMasjid discovery | Core board page | Premium board |
|---|---|---|---|
| Public catalogue discovery | Yes | No; selected board retrieval | Yes, through FindMasjid catalogue entry |
| Public slug (`web_url` / `mid`) | Yes | Used to retrieve board | Yes |
| Masjid display name | Yes | Page-level metadata exists; model use not yet defined | Yes |
| City | Yes | Page-level metadata exists; model use not yet defined | Yes |
| Timezone offset | Yes | Not yet established as an embedded `data` field | Yes |
| Last-updated timestamp | Yes | Yes, where supplied | Not yet defined as a Premium catalogue field |
| Fajr Jamaah | Yes | Yes | Yes |
| Dhuhr Jamaah | Yes | Yes | Yes |
| Asr Jamaah | Yes | Yes | Yes |
| Maghrib Adhan | Yes | Yes | Yes |
| Esha Jamaah | Yes | Yes | Yes |
| Fajr Adhan | No | Yes | Yes |
| Dhuhr Adhan | No | Yes | Yes |
| Asr Adhan | No | Yes | Yes |
| Esha Adhan | No | Yes | Yes |
| Maghrib Jamaah | No reliable field observed | No dedicated field observed | Yes where supplied |
| Jumu'ah summary/Khutbah | Yes | Up to three Jumu'ah slots plus heading configuration | Yes |
| Detailed Jumu'ah heading/time events | No | Limited/configured slots | Yes |
| Jumu'ah Adhan/Salaah semantics | No | Requires mapping from Core slots/configuration | Explicit dedicated values where supplied |
| Khateeb | No | Not observed in embedded Core object | Yes where supplied |
| Sunset | Yes | Yes | Yes |
| Broader astronomical times | No | Yes | Yes |
| Sehri/Fajr start | No | Yes | Yes |
| Ishraaq/Duha | No | Yes | Yes |
| Asr Shafi'i/Hanafi | No | Yes | Yes |
| Istiwa/Zawaal values | No | Yes | Yes |
| Alternate-language Salah values | No | No observed equivalent | Yes |
| Board identity/configuration | Limited | Limited | Yes |
| Announcements/community content | No | No observed equivalent in `data` | Yes |
| Live-stream metadata | Not established | Fields present where configured | Yes where supplied |
| Premium opaque `boardId` | No | No | Yes |
| Embedded `theInfo` payload | No | No | Yes |

## FindMasjid Record Fields Observed

The discovery endpoint exposes records containing fields such as:

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

This remains valuable for catalogue generation and lightweight summaries, but it is no longer the only verified non-Premium timetable source.

## Premium Data Observed

Premium `theInfo` payloads contain substantially richer information, including:

- full five-prayer Adhan/Jamaah data;
- detailed Jumu'ah events;
- Jumu'ah Adhan and Jamaah;
- astronomical times;
- alternate-language data;
- masjid identity and display configuration;
- announcements;
- Nikah and funeral data;
- Taleem/programmes;
- posters and media identifiers;
- contribution/banking information;
- moon/Eid/display configuration; and
- live-stream information where configured.

Premium should therefore be treated as the richer content source, but it is not required for a useful standard timetable.

## Same-board FindMasjid/Premium Validation

### Masjid Aaisha — `erasmia-aaisha`

The FindMasjid catalogue and Premium `theInfo` agreed on every overlapping timetable value checked:

| Field | FindMasjid | Premium |
|---|---:|---:|
| Fajr Jamaah | 06:00 | 06:00 |
| Dhuhr Jamaah | 12:40 | 12:40 |
| Asr Jamaah | 16:45 | 16:45 |
| Maghrib Adhan | 17:52 | 17:52 |
| Esha Jamaah | 19:20 | 19:20 |
| Jumu'ah Khutbah | 12:45 | 12:45 |
| Sunset | 17:49 | 17:49 |
| Timezone offset | +02:00 | +02:00 |

### Brits Jamia — `brits-jamia`

A second independent same-board comparison produced the same result:

| Field | FindMasjid | Premium |
|---|---:|---:|
| Fajr Jamaah | 06:00 | 06:00 |
| Dhuhr Jamaah | 13:20 | 13:20 |
| Asr Jamaah | 17:00 | 17:00 |
| Maghrib Adhan | 17:54 | 17:54 |
| Esha Jamaah | 19:30 | 19:30 |
| Jumu'ah Khutbah | 13:00 | 13:00 |
| Sunset | 17:51 | 17:51 |
| Timezone offset | +02:00 | +02:00 |

This provides cross-board evidence that the FindMasjid timetable summary is consistent with corresponding Premium board data for the overlapping fields tested.

This does **not** prove that the interfaces can never diverge, nor that they have identical semantics or refresh behaviour.

## Sunset / Iftar versus Maghrib Adhan

The Core board research provides direct evidence that sunset and Maghrib Adhan must remain separate timetable concepts.

For `fawkner-rahman`:

```text
sunset       = 17:50
maghribAthan = 17:55
```

For `erasmia-abu-bakr`:

```text
sunset       = 17:49
maghribAthan = 17:52
```

Operationally, Iftar normally occurs at sunset, while a masjid may schedule the Maghrib Adhan at the same time or a few minutes later. MasjidPi must therefore not collapse sunset/Iftar and Maghrib Adhan into one semantic field merely because some upstream JavaScript or some boards use the same value.

Ramadan-specific behaviour still warrants later validation when boards actively expose Ramadan/Iftar data.

## Jumu'ah Comparison

Brits Jamia provided a useful Premium Jumu'ah comparison.

The Premium board exposed the detailed sequence:

```text
Adhan      12:25
Lecture    12:40
Khutbah    13:00
```

and the dedicated Jumu'ah fields contained:

```text
jumuahAdhan   = 12:25
jumuahJamaah  = 13:00
```

The FindMasjid catalogue exposed:

```text
jumuah_khutbah = 13:00
```

In this sample, FindMasjid `jumuah_khutbah`, the Premium Khutbah event and Premium `jumuahJamaah` happen to contain the same time.

This must **not** be generalised into a rule that Khutbah and Jumu'ah Salaah/Jamaah are always identical. They remain separate semantic values.

Core board pages additionally expose `jumuahTime1`, `jumuahTime2`, `jumuahTime3` and `jumuahHeadings`. Their precise semantic mapping must be handled using the verified JavaScript/display research rather than inferred solely from slot position.

The previously agreed Friday timetable rule remains:

```text
Friday standard five-prayer timetable:

Dhuhr Adhan  -> Jumu'ah Adhan
Dhuhr Salaah -> Jumu'ah Salaah/Jamaah

If Jumu'ah Salaah/Jamaah is unavailable:
    use Jumu'ah Khutbah as the Salaah-time fallback.
```

The dedicated Friday Jumu'ah element can preserve and display the richer event sequence separately.

## Source Responsibilities

The current research supports the following separation:

```text
MasjidBoard Live
        |
        +-- FindMasjid
        |       +--> board discovery
        |       +--> MasjidPi catalogue
        |       +--> lightweight summary metadata
        |
        +-- Core board page
        |       +--> selected Core board retrieval
        |       +--> full standard timetable
        |       +--> astronomical times
        |       +--> Jumu'ah slots/configuration
        |       +--> optional stream/configuration fields
        |
        +-- Premium board
                +--> selected Premium board retrieval
                +--> full timetable
                +--> richer Jumu'ah
                +--> alternate-language values
                +--> announcements/community content
                +--> richer board/display configuration
```

Discovery, Core retrieval and Premium retrieval should remain separate provider responsibilities.

## Revised Source Precedence Decision

The earlier working assumption that Premium should always be the only full board-data source is no longer accurate.

For an individual selected board, the research now supports:

```text
FindMasjid
    -> discover board and obtain web_url

Core board page
    -> standard timetable source available independently of Premium

Premium board available
    -> optional richer capability/source for Premium-only content
```

MasjidPi should not merge Core and Premium timetable values field-by-field without a defined precedence policy. Same-board comparisons should continue before production behaviour is frozen.

**Working decision:**

- FindMasjid is primarily the **discovery/catalogue source**.
- The public Core board page is a verified **individual timetable source**.
- Premium is an optional **richer board/content source** where available.
- Premium is **not required** for a useful full standard timetable.

## `MBL_ID` Capability Caveat

Known Premium boards have been observed with both `PRM` and `PRP` suffixes, while other entries expose suffixes including `CRM`, `CRP`, `CRS` and `EXT`.

Premium capability must therefore **not** be derived solely from the `MBL_ID` suffix. The suffix should remain opaque upstream metadata unless authoritative semantics are established later.

Premium availability should be determined by the Premium capability probe described above.

## Important Open Questions

Stage 2 is not complete. The following still require validation:

1. How stable the embedded Core `data` schema is across a larger sample of boards.
2. How often Core board values are updated and what `last_updated` means operationally when populated.
3. Why some Core boards expose an empty `last_updated` value and whether that affects freshness guarantees.
4. Whether Core and Premium can temporarily diverge for the same selected board.
5. Which source should take precedence for overlapping timetable values on a Premium-capable board.
6. Whether a Premium board can temporarily stop resolving while its Core board remains valid.
7. Whether Premium access can ever require authentication or another non-public access path.
8. Exact Core Jumu'ah slot semantics across different `jumuahHeadings` configurations.
9. Ramadan-specific Core values, especially Iftar versus Maghrib behaviour, when Ramadan boards are actively populated.
10. Whether ladies-facility, moon-seen, date-adjust and similar FindMasjid fields belong in the final user-facing catalogue or only in board data.

## Current Stage 2 Direction

Stage 2 has now established that the public Core board page exposes a substantially richer timetable than the FindMasjid catalogue record, and that the same Core schema was observed on two Core-only candidates in different countries.

The next research should broaden Core schema validation, compare Core and Premium data for the same Premium-capable boards, and refine Jumu'ah and Ramadan semantics before production provider precedence is frozen.

No production implementation is frozen yet; this remains research on `research/masjidboard-live`.
