# MasjidBoard Initial Display Layout

**Status:** Initial display frontend implemented  
**Branch:** `research/masjidboard-live`

## Purpose

Define the first read-only MasjidBoard screen layout and its implemented frontend behaviour.

The initial layout is designed around the most constrained supported case: **three selected masjids displayed simultaneously**. One-board and two-board views expand from the same information hierarchy rather than introduce different semantics.

## Primary Use Case

The screen should help a household answer a practical question quickly:

> If I cannot make the Jamaah at one nearby masjid, what is the next practical option?

This makes side-by-side comparison of Jamaah times the primary design goal.

## Initial Content

The default screen includes:

```text
current local time / date
selected masjid names
Fajr
Dhuhr
Asr
Maghrib
Esha
Jumu'ah information
per-board stale/unavailable indication
```

The default screen does not show the extended astronomical/calculation set. Suhur, Fajr Start, Sunrise, Ishraaq, Duha, Istiwa/Zawaal, Shafi'i/Hanafi Asr calculation times, Sunset and Esha Start remain available for later layouts/preferences.

## Implemented Frontend

The first read-only display page is:

```text
/masjidboard.html
```

It uses only:

```text
GET /api/masjidboard/display
```

and does not expose configuration or audio controls.

The frontend files are:

```text
frontend/masjidboard.html
frontend/masjidboard-display.css
frontend/masjidboard-display.js
```

The display refreshes presentation data periodically while maintaining an independent local clock. If the API connection is interrupted after usable data has already been rendered, the existing timetable remains on screen and a small connection warning is shown rather than blanking the display.

## Three-Board Layout

The structure is a comparison grid with one column per selected masjid and the same prayer rows aligned horizontally.

Conceptually:

```text
                         19:25
                   Wednesday 19 August

              Darul Uloom      Brits Jamia       Masjid Taqwa

Fajr
              Adhan  05:30     Adhan  05:40      Adhan  05:40
              Jamaah 05:45     Jamaah 06:00      Jamaah 06:00

Dhuhr
              Adhan  12:30     Adhan  13:00      Adhan  13:00
              Jamaah 12:45     Jamaah 13:20      Jamaah 13:20

Asr
              Adhan  16:30     Adhan  16:40      Adhan  16:30
              Jamaah 16:45     Jamaah 17:00      Jamaah 16:50

Maghrib
              Adhan  17:54     Adhan  17:54      Adhan  17:54

Esha
              Adhan  19:15     Adhan  19:15      Adhan  19:15
              Jamaah 19:30     Jamaah 19:30      Jamaah 19:30

Jumu'ah
              available data   12:25 / 13:00     12:25 / 13:00
```

## Reading Order and Emphasis

Within every prayer cell, reading order is:

```text
Adhan
Jamaah
```

When both values exist:

- Adhan appears first with secondary emphasis;
- Jamaah appears second and is visually dominant, larger and bolder.

When only one value exists, the remaining value adopts the dominant styling.

The display does not render placeholder times for absent data. For example, Maghrib does not show an empty `Jamaah --:--` row when no Jamaah value is supplied. Its available Adhan time becomes the dominant value in that cell.

## Alignment

Prayer names form a stable vertical rhythm shared across all selected boards. Equivalent prayer values from different masjids remain horizontally comparable without requiring the viewer to scan separate cards vertically.

The display preserves the user's configured board order. It does not reorder columns automatically by the next Jamaah time.

## Board Names

Each selected board gets one stable column header. Long names wrap within their column rather than changing identity.

The API-provided selected order is authoritative.

## Jumu'ah

Jumu'ah is shown as a separate section rather than replacing the normal Dhuhr row in the general weekly layout.

`effective_salaah` is preferred as the compact congregational-time value where available. Adhan remains first and Salaah receives dominant emphasis when both are present. Additional timed Jumu'ah events, such as a lecture or Sunan time, may be shown secondarily.

A board may legitimately expose Jumu'ah headings with no associated times. That is not a stale/update failure and no fabricated time is shown.

## Stale and Unavailable State

Status is per board, not global.

A stale board continues displaying its last-known-good timetable and gets a subtle warning attached to that board/column.

An unavailable board with no live or cached timetable keeps its configured column/slot so the other boards do not move unexpectedly. No prayer times are invented.

## Current Time and Date

The screen contains a clearly visible browser-local clock and date as common screen-level elements rather than repeating them in each board column.

The board timezone remains available in the presentation API for later multi-timezone layout refinement.

## One- and Two-Board Adaptation

One-board and two-board modes preserve the same semantics:

- same prayer order;
- same Adhan-before-Jamaah reading order;
- same Jamaah emphasis;
- same missing-value omission rules;
- same stale/unavailable semantics; and
- same selected-board ordering.

The grid automatically adapts its column count to the number of selected boards, allowing fewer boards to use more width without adding extra information by default.

## Non-Goals for Initial Layout

The first layout does not include:

- configuration controls;
- board/location selection;
- automatic ranking of masjids;
- extended astronomical/calculation rows;
- announcements/media/banking content;
- audio controls; or
- provider diagnostic messages.

These are either administrative WebUI responsibilities or future display-layout features.

## Validation Boundary

The frontend should now be live-tested with the existing three-board Brits selection by opening:

```text
http://localhost:8080/masjidboard.html
```

The next refinement should be based on the actual rendered screen rather than adding more display data speculatively.
