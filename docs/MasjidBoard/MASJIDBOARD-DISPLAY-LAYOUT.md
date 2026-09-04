# MasjidBoard Display Presentation

**Status:** Responsive standard display and dedicated 7-inch Appliance profile implemented

## Purpose

This document describes the presentation behaviour shared by MasjidBoard display profiles. Hardware/profile selection is documented separately in `MASJIDBOARD-DISPLAY-PROFILES.md`.

Both presentations consume the same normalized display API:

```text
GET /api/masjidboard/display
```

Display profile is not persisted in the API. The local display runtime selects the profile from attached hardware, while a remote browser may explicitly request the Appliance presentation with:

```text
/masjidboard.html?profile=appliance
```

Normal browser access uses the standard presentation:

```text
/masjidboard.html
```

## Frontend files

```text
frontend/masjidboard.html
frontend/masjidboard-display.css
frontend/masjidboard-display.js
frontend/masjidboard-detailed.css
frontend/masjidboard-detailed.js
frontend/masjidboard-appliance.css
frontend/masjidboard-appliance.js
frontend/masjidboard-touch-controls.js
```

The old generic `portrait` presentation name is no longer used. The current 600x1024 interface is specifically the Appliance profile.

## Shared timetable semantics

Both profiles present up to three selected masjids in the configured order and use the same normalized prayer data. Ordinary days show:

```text
Fajr
Dhuhr
Asr
Maghrib
Esha
```

During the Islamic Friday interval, Jumu'ah replaces Dhuhr. Jumu'ah is not displayed as an additional sixth prayer row.

Within normal prayer data, Adhan is presented before Jamaah. When both values exist, Jamaah receives the stronger visual emphasis. When only one value exists, the available value becomes dominant rather than rendering an empty placeholder.

For Maghrib specifically, an available Adhan remains the displayed dominant time when no separate Jamaah time is supplied.

## Next-event countdown

Each selected masjid independently determines its next visible timed event from the timetable currently being displayed. The countdown follows the visible sequence rather than a fixed prayer position.

On ordinary days, visible Adhan and Jamaah times participate. During Jumu'ah, supplied visible events such as Adhan, Lecture, Sunan, Khutbah and explicit Salaah/Jamaah participate in chronological order.

After the final event of the day, the countdown rolls forward to the following day's first Fajr event.

## Jumu'ah behaviour

Jumu'ah replaces Dhuhr during the Friday display interval. Supplied timed items are shown chronologically and source labels are preserved.

An explicit Salaah/Jamaah value receives strongest emphasis. If the provider supplies Khutbah but no explicit Salaah/Jamaah time, Khutbah remains labelled Khutbah and receives the dominant styling; it is not reinterpreted as Salaah.

Duplicate events resolving to the same visible Adhan or Salaah time are suppressed.

Each selected masjid also contributes a detailed Jumu'ah notice card during this interval when its default-enabled per-masjid option is active. The card preserves the provider's configured event headings and times, includes optional Khateeb or explanatory text, and uses dedicated Jumu'ah Adhan/Salaah values only when they add information not already represented by the detailed events. Outside the Islamic Friday interval, the card is omitted from both display profiles.

For visual validation only, `jumuah-fixture=khateeb` adds representative Khateeb text to the generated card in a browser preview. It does not modify the provider response, cache or saved selection.

## Standard profile presentation

The standard profile is the normal responsive TV/monitor presentation. It supports one to three selected masjids and scales from smaller landscape displays through Full HD and 4K while preserving the same information hierarchy.

The main timetable is a comparison grid with equivalent prayer rows aligned horizontally across selected masjids. The first selected masjid also supplies the Daily Times information, including:

```text
Sehri end
Fajr start
Sunrise
Ishraaq
Duha / Chaasht
Zawaal / Istiwa
Asr Shafi'i
Asr Hanafi
Sunset
Esha start
```

When Premium enrichment is available, the standard presentation also shows rotating community information such as announcements, Nikah, funerals, Salaah changes, new-moon information and Islamic Economic Indicators.

The standard profile contains no appliance Listen controls.

## Appliance profile presentation

The Appliance profile targets the validated 7-inch Waveshare display at an effective 600x1024 portrait viewport. It uses the same timetable and notice data but adapts it to a slideshow-oriented compact interface.

The top area includes the clock, Gregorian date, Islamic date, primary masjid name and next-event information. Slides include one salaah-times slide per selected masjid, Daily Times, community notices and optional Islamic Economic Indicators.

The default slide duration is 15 seconds and remains configurable between 5 and 60 seconds.

### Appliance touch controls

The Appliance profile adds a touch sheet opened by swiping up and closed by swiping down or after inactivity. It contains:

- Masjid controls for favourite selection, playback and Masjid volume;
- Radio controls for station selection, Scheduled Play, Play Now, Stop and Radio volume;
- Theme selection; and
- Master volume where the active audio device supports it.

The slideshow pauses while the touch sheet is open and resumes when it closes.

These controls deliberately do not expose full catalogue search, board/location configuration, schedule-time editing or audio-device administration. Those remain WebUI responsibilities.

## Themes and optional information

Theme, slide duration and Islamic Economic Indicators visibility are shared saved display preferences. They are independent of the automatically selected display profile.

Changing a theme therefore affects both standard and Appliance presentations. Selecting an Appliance profile is never a saved user preference.

## Stale and unavailable data

Status is per board rather than global. A stale board continues displaying its last-known-good timetable with a subtle warning. An unavailable selected board retains its configured slot so other boards do not move unexpectedly, and no prayer times are invented.

If the display API temporarily becomes unreachable after usable data has been rendered, the existing timetable remains on screen while the frontend continues retrying.

## Development date/time overrides

The display supports deterministic browser test overrides:

```text
/masjidboard.html?date=2026-08-21
/masjidboard.html?date=2026-08-21&time=12:10
/masjidboard.html?profile=appliance&date=2026-08-21&time=12:10
```

Invalid date/time overrides are ignored and browser-local time is used.

## Profile and orientation boundary

Profile and orientation are separate concepts. The current implementation maps:

```text
standard  -> responsive landscape presentation
appliance -> dedicated 7-inch portrait presentation
```

A future conventional portrait monitor/TV presentation may be added without redefining the Appliance profile. Hardware detection, Cog rotation and touchscreen calibration remain responsibilities of the local display runtime rather than the presentation API.
