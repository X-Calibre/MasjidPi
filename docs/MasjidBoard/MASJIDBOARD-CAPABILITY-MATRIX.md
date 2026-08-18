# MasjidBoard Live Capability Matrix Research

**Status:** Research / Stage 2  
**Branch:** `research/masjidboard-live`

## Purpose

Record verified comparisons between MasjidBoard Live Core/FindMasjid data and Premium board data, and define the capabilities that MasjidPi can expect from each access path.

This document is intentionally separate from board discovery and from the positional Premium parser research.

## Reference Boards

The research has used both Core-only candidates and known Premium boards:

- Core-only candidate: `erasmia-abu-bakr`
- Core-only candidate: `fawkner-rahman`
- Known Premium board: `erasmia-aaisha`
- Known Premium board: `brits-jamia`
- Known Premium board: `fawkner-masjid`

The same-board Core/Premium comparisons were performed for `erasmia-aaisha` and `brits-jamia`.

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

## Stage 2 Capability Matrix

| Capability | Core directory / Core-only | Premium board |
|---|---|---|
| Public catalogue discovery | Yes | Yes, through Core catalogue entry |
| Public slug (`web_url` / `mid`) | Yes | Yes |
| Masjid display name | Yes | Yes |
| City | Yes | Yes |
| Timezone offset | Yes | Yes |
| Last-updated timestamp | Yes | Not yet defined as a Premium catalogue field |
| Fajr Jamaah | Yes | Yes |
| Dhuhr Jamaah | Yes | Yes |
| Asr Jamaah | Yes | Yes |
| Maghrib Adhan | Yes | Yes |
| Esha Jamaah | Yes | Yes |
| Fajr Adhan | Not observed in Core record | Yes |
| Dhuhr Adhan | Not observed in Core record | Yes |
| Asr Adhan | Not observed in Core record | Yes |
| Esha Adhan | Not observed in Core record | Yes |
| Maghrib Jamaah | Not observed as a reliable Core field | Yes where supplied |
| Jumu'ah Khutbah | Yes | Yes |
| Detailed Jumu'ah heading/time events | No | Yes |
| Jumu'ah Adhan | No | Yes |
| Jumu'ah Salaah/Jamaah | No | Yes |
| Khateeb | No | Yes where supplied |
| Sunset | Yes | Yes |
| Broader astronomical times | No | Yes |
| Alternate-language Salah values | No | Yes |
| Board identity/configuration | Limited | Yes |
| Islamic/date-display configuration | Limited | Yes |
| Announcements | No | Yes |
| Nikah notices | No | Yes |
| Funeral notices | No | Yes |
| Taleem/programme data | No | Yes |
| Posters / images / community posters | No | Yes |
| Banking / contribution information | No | Yes |
| Eid / moon / display settings | Limited or absent | Yes |
| Live-stream metadata | Not established | Yes where supplied |
| Premium opaque `boardId` | No | Yes |
| Embedded `theInfo` payload | No | Yes |

## Core Record Fields Observed

The Core discovery endpoint exposes records containing fields such as:

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

This means a Core-only board is not merely a directory entry. It already contains enough data for a basic timetable-oriented MasjidBoard experience, subject to later validation of freshness and field semantics.

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

The Premium provider should therefore be treated as the richer board-data source rather than as the only usable MasjidBoard source.

## Same-board Core/Premium Validation

### Masjid Aaisha — `erasmia-aaisha`

The Core catalogue and Premium `theInfo` agreed on every overlapping timetable value checked:

| Field | Core | Premium |
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

| Field | Core | Premium |
|---|---:|---:|
| Fajr Jamaah | 06:00 | 06:00 |
| Dhuhr Jamaah | 13:20 | 13:20 |
| Asr Jamaah | 17:00 | 17:00 |
| Maghrib Adhan | 17:54 | 17:54 |
| Esha Jamaah | 19:30 | 19:30 |
| Jumu'ah Khutbah | 13:00 | 13:00 |
| Sunset | 17:51 | 17:51 |
| Timezone offset | +02:00 | +02:00 |

This provides cross-board evidence that the Core timetable summary is consistent with the corresponding Premium board data for the overlapping fields tested.

This does **not** prove that Core and Premium can never diverge, nor that the interfaces have identical semantics or refresh behaviour.

## Premium Contains Richer Timetable Data

The verified Brits Jamia Premium Salah row contained:

```text
Fajr     Adhan 05:40   Jamaah 06:00
Dhuhr    Adhan 13:00   Jamaah 13:20
Asr      Adhan 16:40   Jamaah 17:00
Maghrib  Adhan 17:54   Jamaah unavailable
Esha     Adhan 19:15   Jamaah 19:30
```

The corresponding Core catalogue entry exposes only:

```text
Fajr Jamaah       06:00
Dhuhr Jamaah      13:20
Asr Jamaah        17:00
Maghrib Adhan     17:54
Esha Jamaah       19:30
```

Core therefore remains useful for discovery, summary and potentially a limited fallback representation, but it is not a replacement for the full Premium board provider.

## Jumu'ah Comparison

Brits Jamia provided a useful Jumu'ah comparison.

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

The Core catalogue exposed:

```text
jumuah_khutbah = 13:00
```

In this sample, Core `jumuah_khutbah`, the Premium Khutbah event and Premium `jumuahJamaah` happen to contain the same time.

This must **not** be generalised into a rule that Khutbah and Jumu'ah Salaah/Jamaah are always identical. They remain separate semantic values.

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
        +-- FindMasjid / Core
        |       |
        |       +--> board discovery
        |       +--> MasjidPi catalogue
        |       +--> lightweight timetable summary
        |       +--> discovery metadata
        |
        +-- Premium board
                |
                +--> full board retrieval
                +--> Adhan and Jamaah times
                +--> detailed Jumu'ah
                +--> astronomical times
                +--> Islamic date/context
                +--> announcements and other board content
```

This reinforces the architectural distinction between discovering which boards exist and retrieving/parsing an individual board.

These should remain separate provider responsibilities.

## Source Precedence Decision

For an individual selected board, the working source precedence is:

```text
Premium board available
        |
        +--> use Premium as the full board-data source
        |
        +--> retain Core as catalogue/discovery metadata

Premium board unavailable
        |
        +--> Core may provide a limited board representation
             for fields whose semantics have been verified
```

MasjidPi should not merge Core and Premium timetable values field-by-field by default. Doing so would introduce unnecessary ambiguity over which source is authoritative.

**Working decision:**

- Core is primarily the **discovery/catalogue source**.
- Premium is the preferred **full board-data source** when available.
- Core remains potentially useful as a limited fallback source when Premium is unavailable.

## `MBL_ID` Capability Caveat

Known Premium boards have been observed with both `PRM` and `PRP` suffixes, while other Core entries expose suffixes including `CRM`, `CRP`, `CRS` and `EXT`.

Premium capability must therefore **not** be derived solely from the `MBL_ID` suffix. The suffix should remain opaque upstream metadata unless authoritative semantics are established later.

Premium availability should be determined by the Premium capability probe described above.

## Important Open Questions

Stage 2 is not complete. The following still require validation:

1. Whether Core-only boards are sufficiently fresh/reliable for continuous timetable display.
2. Whether Core exposes additional per-board data through another endpoint or board page.
3. How often Core records are updated and what `last_updated` means operationally.
4. Whether Core and Premium can temporarily diverge despite the two successful same-board comparisons.
5. Whether a Premium board can temporarily stop resolving while its Core entry remains valid.
6. Whether Premium access can ever require authentication or another non-public access path.
7. Whether Core has a distinct Basic/Free board page whose structure should also be parsed.
8. Whether Ramadan-specific values expose additional Core capabilities.
9. Whether ladies-facility, moon-seen, date-adjust and similar fields belong in the final user-facing catalogue or only in board data.

## Current Stage 2 Direction

The same-board comparison question is now sufficiently validated to adopt a working source-precedence rule: Core for catalogue/discovery and Premium for full board data when available.

The next Stage 2 work should investigate the operational quality of Core-only data, especially freshness, update behaviour, and whether the public Core board page exposes additional data beyond the FindMasjid catalogue record.

No production implementation is frozen yet; this remains research on `research/masjidboard-live`.
