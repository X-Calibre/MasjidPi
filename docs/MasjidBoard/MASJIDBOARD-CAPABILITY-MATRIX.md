# MasjidBoard Live Capability Matrix Research

**Status:** Research / Stage 2
**Branch:** `research/masjidboard-live`

## Purpose

Record the first verified comparison between a Core-only MasjidBoard Live entry and a Premium board, and define the capabilities that MasjidPi can expect from each access path.

This document is intentionally separate from board discovery and from the positional Premium parser research.

## Reference Boards

The first comparison uses two boards from the same city and timezone so that the upstream service category is the main variable:

- Core-only candidate: `erasmia-abu-bakr`
- Known Premium board: `erasmia-aaisha`

## Premium Capability Probe

A Core catalogue `web_url` can be probed against:

```text
https://premium.masjidboardlive.com/v2/?mid=<web_url>
```

### Verified Premium success

For `erasmia-aaisha`, the generated page contains both:

```javascript
let boardId = "1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k";
let theInfo = [...];
```

This is considered a successful Premium capability probe.

### Verified Premium absence

For `erasmia-abu-bakr`, the Premium page renders:

```text
MasjidBoard live - 500
This masjid does not exist
```

and does not expose `boardId` or `theInfo`.

The same negative pattern was also observed for `fawkner-rahman` in Australia, while `fawkner-masjid` is a known Premium board.

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
    Timeout, DNS failure, transport failure, unexpected server response, malformed page, or other condition where Premium capability cannot be determined reliably.
```

A transient failure must never be interpreted as proof that a board is Core-only.

## First Stage 2 Capability Matrix

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

The `erasmia-aaisha` Premium page contains a full `theInfo` payload with substantially richer information, including:

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

The Premium provider should therefore be treated as a richer capability layer rather than as the only usable MasjidBoard source.

## Architectural Consequence

The emerging provider architecture is:

```text
Core catalogue entry
        |
        +--> basic timetable and discovery metadata
        |
        +--> web_url
                |
                v
       Premium capability probe
                |
        +-------+-------+
        |               |
        v               v
   unavailable       available
        |               |
        v               v
   Core provider    Premium provider
        |               |
        +-------+-------+
                |
                v
       Normalised MasjidBoard model
```

MasjidPi should be able to use the Core path when Premium is unavailable, while preferring or augmenting with Premium data when the richer source is available.

## Important Open Questions

Stage 2 is not complete. The following still require validation:

1. Whether Core-only boards are sufficiently fresh/reliable for continuous timetable display.
2. Whether Core exposes additional per-board data through another endpoint or board page.
3. How often Core records are updated and what `last_updated` means operationally.
4. Whether Core and Premium times can diverge for the same Premium masjid and which source should win.
5. Whether a Premium board can temporarily stop resolving while its Core entry remains valid.
6. Whether Premium access can ever require authentication or another non-public access path.
7. Whether Core has a distinct Basic/Free board page whose structure should also be parsed.
8. Whether Ramadan-specific values expose additional Core capabilities.
9. Whether ladies-facility, moon-seen, date-adjust and similar fields belong in the final user-facing catalogue or only in board data.

## Current Stage 2 Direction

The next research step should compare the Core entry and Premium payload for the **same Premium masjid** rather than comparing two different masjids.

Recommended reference board:

```text
erasmia-aaisha
```

This will let us determine:

- whether the prayer times agree between Core and Premium;
- whether Core values are derived from Premium;
- which fields are duplicated;
- which source should be authoritative when both are available; and
- whether Core can serve as a fallback cache/source for a Premium board.

No production implementation decision should be frozen until that same-board comparison is complete.
