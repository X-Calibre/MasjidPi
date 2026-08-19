# MasjidBoard Initial Display Layout

**Status:** Initial display design decision  
**Branch:** `research/masjidboard-live`

## Purpose

Define the first read-only MasjidBoard screen layout before implementing the display frontend.

The initial layout is designed around the most constrained supported case: **three selected masjids displayed simultaneously**. One-board and two-board views should later expand from the same information hierarchy rather than introduce different semantics.

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

## Three-Board Layout

The preferred initial structure is a comparison grid with one column per selected masjid and the same prayer rows aligned horizontally.

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

The example is structural only; exact typography, spacing and visual styling will be determined during frontend implementation.

## Reading Order and Emphasis

Within every prayer cell, reading order is:

```text
Adhan
Jamaah
```

When both values exist:

- Adhan appears first with secondary emphasis;
- Jamaah appears second and is visually dominant, for example larger and/or bolder.

When only one value exists, the remaining value adopts the dominant styling.

The display must not render placeholders for absent data. For example, Maghrib should not show an empty `Jamaah --:--` row when no Jamaah value is supplied.

This rule prevents sparse upstream data from producing visually broken timetable rows.

## Alignment

Prayer names should form a stable vertical rhythm shared across all selected boards. Equivalent prayer values from different masjids should be horizontally comparable without requiring the viewer to scan separate cards vertically.

This is especially important for Jamaah comparison:

```text
Asr Jamaah
Darul Uloom    16:45
Masjid Taqwa   16:50
Brits Jamia    17:00
```

The display should preserve the user's configured board order. It must not reorder columns automatically by the next Jamaah time.

## Board Names

Each selected board gets one stable column header. Long names may need controlled wrapping or abbreviation at the presentation layer, but the frontend must not silently substitute a different mosque identity.

The API-provided selected order is authoritative.

## Jumu'ah

Jumu'ah should be shown as a separate section rather than replacing the normal Dhuhr row in the general weekly layout.

Available Jumu'ah information should be rendered without fabricating missing values. `effective_salaah` is the preferred compact congregational-time value when available, while richer event headings/times may be used where screen space permits.

A board may legitimately expose Jumu'ah headings with no associated times. That is not a stale/update failure.

## Stale and Unavailable State

Status is per board, not global.

A stale board should continue displaying its last-known-good timetable and show a subtle but visible warning attached to that board/column. The warning should not dominate the timetable.

Conceptually:

```text
Brits Jamia Masjid   [stale]
```

An unavailable board with no live or cached timetable should keep its configured column/slot so the other boards do not move unexpectedly. The slot should communicate that timetable data is unavailable without displaying invented times.

## Current Time and Date

The screen should contain a clearly visible current local clock and date. The exact placement can be refined during visual implementation, but these should be common screen-level elements rather than repeated in every board column when the selected boards share the same local context.

If future multi-location selections span different timezones, the board timezone remains available in the presentation API and the layout can adapt accordingly. The initial implementation should not assume provider-specific timezone behaviour.

## One- and Two-Board Adaptation

One-board and two-board modes should preserve the same semantics:

- same prayer order;
- same Adhan-before-Jamaah reading order;
- same Jamaah emphasis;
- same missing-value omission rules;
- same stale/unavailable semantics; and
- same selected-board ordering.

Fewer boards should use the additional width to improve readability rather than add more data by default.

The initial frontend should therefore be responsive to the number of selected boards, not switch to a different information model.

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

## Next Implementation Step

Implement the read-only display frontend against:

```text
GET /api/masjidboard/display
```

Start with the three-board comparison layout, then verify that the same component structure adapts cleanly to two and one selected board.
