# MasjidBoard Display Layouts

**Status:** Responsive TV / Monitor and dedicated 7-inch Appliance Display layouts implemented
**Branch:** `docs/masjidboard-live-data-inventory` / PR #38

## Purpose

Define the read-only MasjidBoard layouts and their implemented frontend behaviour. Both layouts consume the same normalized display API and share timetable, date, countdown, theme and stale-state semantics.

The default layout is designed around the most constrained supported case: **three selected masjids displayed simultaneously**. One-board and two-board views expand from the same information hierarchy rather than introduce different semantics.

## Primary Use Case

The screen should help a household answer practical questions quickly, including:

> If I cannot make the Jamaah at one nearby masjid, what is the next practical option?

Side-by-side comparison of prayer and Jamaah times is therefore the primary design goal.

## Default Content

Both layouts include:

```text
current local time / date
selected masjid names
Fajr
Dhuhr (Saturday–Thursday)
Jumu'ah (Friday, replacing Dhuhr)
Asr
Maghrib
Esha
next-timetable-event countdown for each selected masjid
per-board stale/unavailable indication
```

Landscape additionally includes a full-width Daily Times footer sourced consistently from the first selected masjid. It shows Sehri end, Fajr start, Sunrise, Ishraaq, Duha/Chaasht, Zawaal/Istiwa, Shafi'i and Hanafi Asr start, Sunset and Esha start.

When optional Premium enrichment is available, Landscape also allocates an adaptive right-hand panel to rotating, source-labelled community cards. Supported categories include announcements, Nikah, funerals, Eid, upcoming Salaah changes, well-wishes, Taleem, Dawah/Gasht, three-day Jamaat, contributions and calculated new-moon information. Missing Premium content does not affect the Core timetable.

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
frontend/masjidboard-detailed.css
frontend/masjidboard-detailed.js
frontend/masjidboard-portrait.css
frontend/masjidboard-portrait.js
```

The display refreshes presentation data periodically while maintaining an independent local clock. If the API connection is interrupted after usable data has already been rendered, the existing timetable remains on screen and a small connection warning is shown rather than blanking the display.

## Three-Board Layout

The structure is a comparison grid with one column per selected masjid and the same prayer rows aligned horizontally.

On ordinary days:

```text
Fajr
Dhuhr
Asr
Maghrib
Esha
```

On Friday:

```text
Fajr
Jumu'ah
Asr
Maghrib
Esha
```

Jumu'ah is therefore not displayed as an additional sixth row. It is hidden on non-Friday days and replaces Dhuhr on Friday.

## Reading Order and Emphasis

Within every normal prayer cell, reading order is:

```text
Adhan
Jamaah
```

When both values exist:

- Adhan appears first with secondary emphasis;
- Jamaah appears second and is visually dominant, larger and bolder.

When only one value exists, the remaining value adopts the dominant styling.

The display does not render placeholder times for absent data. For example, Maghrib does not show an empty `Jamaah --:--` row when no Jamaah value is supplied. Its available Adhan time becomes the dominant value in that cell.

## Next Timetable Event Countdown

Each selected masjid independently identifies its next upcoming **visible timed event** from the timetable currently being displayed.

The countdown is not limited to Adhan or Jamaah. On ordinary days the candidate events are the visible Adhan and Jamaah times. On Friday, the visible Jumu'ah events participate instead of Dhuhr, including supplied events such as Adhan, Lecture, Sunan, Khutbah and explicit Salaah/Jamaah.

Only the next event for that masjid gets a countdown. The countdown is rendered directly below that event's time inside the relevant prayer card so there is no ambiguity about which masjid or event it belongs to.

Examples:

```text
Adhan       12:16
            in 15 min
Jamaah      13:15
```

After that Adhan has passed, the countdown moves to the next event:

```text
Adhan       12:16
Jamaah      13:15
            in 59 min
```

For longer intervals the display uses natural compact wording:

```text
in 1 hr
in 1 hr 15 min
```

When the event is due, the countdown may display:

```text
now
```

The countdown is visually subordinate to the event time but remains clearly readable. It updates from the browser-local clock without requiring an API request every second.

The countdown follows the same visibility rules as the timetable itself: Dhuhr is excluded on Friday, Jumu'ah takes its place, and no hidden or fabricated event is used as a countdown target.

When a masjid has no remaining visible event for the current day, its countdown rolls forward to the **following day's first Fajr event** rather than disappearing overnight. This rollover is calculated independently for each selected masjid, so one column can already be counting down to tomorrow's Fajr while another still has a later Esha/Jamaah event remaining tonight.

## Jumu'ah Behaviour

Jumu'ah is shown only on Friday and occupies the normal Dhuhr row position.

Within a Jumu'ah cell, timed items are displayed in chronological order regardless of the order returned by the provider. Source labels are preserved rather than inferred or replaced.

The visual hierarchy is:

```text
explicit Salaah/Jamaah      strongest
Adhan                       secondary
other timed events          slightly smaller than Adhan
```

If no explicit Salaah/Jamaah time is supplied but a Khutbah time is supplied, **Khutbah remains labelled Khutbah** and takes the dominant visual styling that would otherwise be used for Salaah. It is not reinterpreted as a Salaah time.

Additional timed events such as Lecture and Sunan remain clearly readable and are not treated as tiny metadata. Duplicate events that resolve to the same displayed Adhan or explicit Salaah time are suppressed.

A board may legitimately expose Jumu'ah headings with no associated times. That is not a stale/update failure and no fabricated time or placeholder is shown.

The Friday row is allowed slightly more vertical space than an ordinary Dhuhr row so that multiple timed Jumu'ah events can remain legible without shrinking the overall information hierarchy.

The next-event countdown participates in the Friday sequence. Manual validation has confirmed that it advances correctly through the supplied Jumu'ah events and then continues to later daily events.

## Alignment and Responsiveness

Prayer names form a stable vertical rhythm shared across all selected boards. Equivalent prayer values from different masjids remain horizontally comparable without requiring the viewer to scan separate cards vertically.

The display preserves the user's configured board order. It does not reorder columns automatically by prayer time.

The page is responsive. Card dimensions, typography and spacing adapt with the available viewport while preserving the same information hierarchy. This responsive behaviour is an intentional requirement and should not be lost during future refinements.

## Board Names

Each selected board gets one stable column header. Long names wrap within their column rather than changing identity.

The API-provided selected order is authoritative.

## Stale and Unavailable State

Status is per board, not global.

A stale board continues displaying its last-known-good timetable and gets a subtle warning attached to that board/column.

An unavailable board with no live or cached timetable keeps its configured column/slot so the other boards do not move unexpectedly. No prayer times are invented.

## Current Time and Date

The screen contains a clearly visible browser-local clock and date as common screen-level elements rather than repeating them in each board column.

The weekday shown by that local clock determines whether the Friday Jumu'ah row or the ordinary Dhuhr row is rendered. The display refresh loop re-evaluates this automatically, so the row changes without restarting MasjidPi when the appliance crosses into or out of Friday.

The board timezone remains available in the presentation API for later multi-timezone layout refinement.

## Development Date/Time Overrides

The display supports query-string overrides for repeatable browser testing without changing the host clock.

A Friday can be simulated with:

```text
/masjidboard.html?date=2026-08-21
```

A specific clock time can also be supplied:

```text
/masjidboard.html?date=2026-08-21&time=12:10
```

The `time=` override drives both the large displayed clock and the next-event calculation. This allows Jumu'ah transitions and countdown movement to be tested deterministically.

Invalid date/time override values are ignored and normal browser-local time is used instead.

These parameters are development/test aids rather than user-facing configuration preferences.

## One- and Two-Board Adaptation

One-board and two-board modes preserve the same semantics:

- same prayer order;
- same Friday Jumu'ah replacement rule;
- same Adhan-before-Jamaah reading order;
- same Jamaah emphasis;
- same next-timetable-event countdown behaviour;
- same overnight rollover behaviour;
- same missing-value omission rules;
- same stale/unavailable semantics; and
- same selected-board ordering.

The grid automatically adapts its column count to the number of selected boards, allowing fewer boards to use more width without adding extra information by default.

## Layout Selection

The configuration UI presents TV / Monitor and 7-inch Appliance Display modes while persisting the existing internal values `landscape` and `portrait`. The display observes the saved preference and switches without restarting the display service.

TV / Monitor mode responsively supports landscape displays from 1366 × 768 through 4K. The 7-inch Appliance Display targets the integrated 600 × 1024 screen in the physical MasjidPi appliance. Further hardware-specific layouts may be added where they provide clear value.

The 7-inch mode also provides an appliance-specific touch sheet opened by swiping up, with a small persistent text hint identifying the gesture. Its Masjid and Radio tabs reuse the existing Listen APIs for live status, favourite/station selection, three volume controls and playback modes. A third Theme tab updates the same saved Board preference used by the configuration Web UI. The slideshow pauses while the sheet is open and resumes when it closes. These controls are deliberately absent from TV / Monitor mode.

Potential alternative layouts may make different use of the same normalized display data, for example:

- more compact single-masjid layouts;
- layouts that expose selected astronomical/calculation times;
- different emphasis for wall-mounted/large-screen use;
- alternate visual density or typography; and
- layouts intended for different aspect ratios or display sizes.

Layout selection should remain a frontend/presentation concern. Alternative layouts should consume the same stable MasjidBoard display API rather than duplicating provider or timetable logic.

## Non-Goals

The TV / Monitor display does not include appliance controls. The 7-inch mode exposes only the simplified touch controls described above and does not include:

- full Masjid catalogue search or favourite management;
- board/location selection;
- automatic ranking of masjids;
- poster/image media;
- audio-device, schedule-time or resume-delay configuration; or
- provider diagnostic messages.

These are either administrative WebUI responsibilities or future display-layout features.

## Validation Boundary

The frontend should be live-tested with an existing one-, two- or three-board selection by opening:

```text
http://localhost:8080/masjidboard.html
```

Non-Friday validation should show Dhuhr and no Jumu'ah row. Friday validation should show Jumu'ah in the Dhuhr position and no separate Dhuhr row.

For each board, the countdown should appear only beneath that board's next upcoming visible timetable event and should move automatically to the following event after the current one has passed. After the final event of the day, it should roll forward to the following day's first Fajr event.
