# MasjidBoard Live Research

**Status:** Research / discovery
**Branch:** `research/masjidboard-live`

## Purpose

Record findings, assumptions, decisions, and open questions for integrating MasjidBoard Live as the primary data source for the MasjidBoard module in MasjidPi V2.

## Confirmed Requirements

- Display the full range of information and content presented by MasjidBoard Live.
- Display it on an external HDMI-connected screen attached to the Raspberry Pi.
- Use MasjidBoard Live as the primary data source.
- Cache downloaded data locally so the board can continue operating during temporary connectivity loss.
- Keep MasjidBoard independent from the MasjidPi audio playback subsystem.

## Initial Architecture

```text
MasjidBoard Live
       |
       v
Data Provider -> Normalised Board Data -> Local Cache -> Display Scheduler -> HDMI

Audio playback remains a separate subsystem.
```

The display layer should not depend directly on MasjidBoard Live's response format. A provider boundary should allow another data source to be added later.

## Information Scope

The investigation must account for all content currently presented by MasjidBoard Live, including prayer and astronomical times, Adhan/Iqamah, Jumu'ah and Khateeb information, daily Ayah, Hadith, Sunnah, Du'a, community broadcasts, announcements, request-for-Du'a information, Nikah and funeral notices, weekly programmes, masjid information, contribution information, posters/images, New Moon information, Eid information, and other board-specific content where supplied.

This list is provisional and must be expanded if the underlying data reveals additional content.

## Verified Findings

### 2026-08-13 — HAR/network capture

A Firefox HAR capture of a live Premium board was analysed. The captured board was:

```text
https://premium.masjidboardlive.com/v2/?mid=erasmia-aaisha
```

The browser made the following MasjidBoard-specific requests:

```text
GET https://api.masjidboardlive.com/mblfileapi
GET https://api.masjidboardlive.com/mblapi?id=1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k
GET https://api.masjidboardlive.com/imageproxy?id=<image-id>
```

The `mblapi` response is structured JSON and contains 29 top-level arrays. The same data is embedded in the page as the JavaScript variable `theInfo`.

The page embeds:

```javascript
let boardId = "1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k";
let mblVersion = "14";
```

The `boardId` used by `mblapi` is distinct from the public URL identifier (`mid=erasmia-aaisha`).

### 2026-08-13 — `functions_uo_latest.js` schema mapping

The actual `functions_uo_latest.js` source was captured and inspected. Its `handleResults(spreadsheetArray)` function maps the positional API response into named JavaScript variables. This gives us a substantially more reliable description of the upstream schema than the earlier row-level observations.

The following mapping is **verified from the supplied JavaScript source**. Row and column numbers below are zero-based and refer to the top-level `mblapi` array.

#### Row 0 — upcoming Salah changes

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `fajrNextDate` | Next Fajr date when next-change display is disabled |
| 1 | `fajrNextTime` | Next Fajr time when next-change display is disabled |
| 2 | `fajrNextMilli` | Next Fajr millisecond/time value |
| 3 | `asrNextDate` | Next Asr date when next-change display is disabled |
| 4 | `asrNextTime` | Next Asr time when next-change display is disabled |
| 5 | `asrNextMilli` | Next Asr millisecond/time value |
| 6 | `eshaNextDate` | Next Esha date when next-change display is disabled |
| 7 | `eshaNextTime` | Next Esha time when next-change display is disabled |
| 8 | `eshaNextMilli` | Next Esha millisecond/time value |
| 9 | `nextChangeDisplay` | Whether next-change values should use the alternate date/time columns |
| 10–12 | `fajrNextDate/Time/Milli` | Alternate next Fajr values when enabled |
| 13–15 | `asrNextDate/Time/Milli` | Alternate next Asr values when enabled |
| 16–18 | `eshaNextDate/Time/Milli` | Alternate next Esha values when enabled |

#### Row 1 — Jumu'ah

| Column | Variable | Meaning |
|---:|---|---|
| 0–5 | `jumuahHead1/Time1` through `jumuahHead3/Time3` | Up to three Jumu'ah heading/time pairs |
| 6 | `jumuahKhateeb` | Khateeb |
| 7 | `jumuahAdhan` | Jumu'ah Adhan |
| 8 | `jumuahJamaah` | Jumu'ah Jamaah/Salah |
| 9 | `jumuahAdhanI` | Alternate-language Jumu'ah Adhan |
| 10 | `jumuahJamaahI` | Alternate-language Jumu'ah Jamaah/Salah |
| 11 | `jumuahHeadingsArray` | Jumu'ah headings configuration/array |

This confirms that Jumu'ah is not limited to one service or one simple time value.

#### Row 2 — moon/location/masjid media data

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `moonBirth` | Moon birth/new-moon value |
| 1 | `moonSet0` | Moon set value, first calculation |
| 2 | `moonAge0` | Moon age, first calculation |
| 3 | `moonAzimuth0` | Moon azimuth, first calculation |
| 4 | `moonAltitude0` | Moon altitude, first calculation |
| 5 | `bestVisibil` | Best visibility value |
| 6 | `moonSet1` | Moon set value, second calculation |
| 7 | `moonAge1` | Moon age, second calculation |
| 8 | `moonAzimuth1` | Moon azimuth, second calculation |
| 9 | `moonAltitude1` | Moon altitude, second calculation |
| 10 | `moonBirthScript` | Moon-birth script/content |
| 11 | `moonVisibilScript` | Moon-visibility script/content |
| 12 | `masjidImage` | Masjid image identifier |

#### Row 3 — daily Salah

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `fajrAthan` | Fajr Adhan |
| 1 | `fajrJamaah` | Fajr Jamaah |
| 2 | `dhuhrAthan` | Zuhr/Dhuhr Adhan |
| 3 | `dhuhrJamaah` | Zuhr/Dhuhr Jamaah |
| 4 | `asrAthan` | Asr Adhan |
| 5 | `asrJamaah` | Asr Jamaah |
| 6 | `iftar` / `maghribAthan` | Iftar value; the source code also assigns this same column to `maghribAthan` |
| 7 | `maghribJamaah` | Maghrib Jamaah |
| 8 | `eshaAthan` | Esha Adhan |
| 9 | `eshaJamaah` | Esha Jamaah |
| 10 | `sundayDhuhr` | Sunday Zuhr value |
| 11 | `salaahLanguage` | Salah-language setting |
| 12 | `sundayzuhr` | Additional Sunday Zuhr value |
| 13 | `highlightSalaah` | Highlight Salah setting |
| 14 | `highlightOwwal` | Highlight Owwal/early-time setting |
| 15 | `liveStreamingServer` | Live-stream server |
| 16 | `liveStreamingURL` | Live-stream URL |

**Important source-code observation:** `handleResults()` assigns `spreadsheetArray[3][6]` to both `iftar` and `maghribAthan`. This should not be treated as proof that Iftar and Maghrib are semantically identical. Row 23 separately contains `maghribAthanI` and the other alternate-language Salah values. MasjidPi should preserve the upstream distinction where possible and verify this against actual board data.

#### Row 4 — display/theme/time configuration

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `istiwaCaution` | Istiwa caution value |
| 1 | `zawaalEnd` | Zawaal end |
| 2 | `theme` | Board theme |
| 3 | `displayLanguage` | Display language |
| 4 | `istiwa` | Istiwa |
| 6 | `islamicTime` | Display Islamic time |
| 7 | `salaah12HourFormat` | Salah 12-hour formatting |
| 8 | `pageHardReload` | Hard-reload setting |
| 9 | `islamicTimeInterval` | Islamic-time rotation interval, multiplied by 1000 by the client |
| 10 | `westernTimeInterval` | Western-time rotation interval, multiplied by 1000 by the client |
| 11 | `puAlternatingTime` | Perpetual/upcoming time alternation setting |
| 12 | `islamicOwwalTimes` | Islamic Owwal-times setting |
| 13 | `rightToLeft` | RTL display setting |
| 14 | `islamicTimesArabic` | Arabic Islamic-time display setting |
| 15 | `islamicTimeMaghribFirst` | Maghrib-first Islamic-time setting |
| 5 | — | `boardLayout` is commented out in the source and is not currently read by `handleResults()` |

#### Row 5 — astronomical/prayer-start times and poster controls

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `sehriEnds` | Sehri/Suhur end |
| 1 | `fajrStarts` | Fajr start |
| 2 | `sunrise` | Sunrise |
| 3 | `ishraaq` | Ishraaq |
| 4 | `duha` | Duha |
| 5 | `asrShafi` | Shafi'i Asr |
| 6 | `asrHanafi` | Hanafi Asr |
| 7 | `sunset` | Sunset |
| 8 | `eshaStarts` | Esha start |
| 9 | `secondLanguage` | Second language |
| 10 | `mblPosterHide` | Poster-hide list |
| 11 | `mblPosterTime` | Poster timing |
| 12 | `disableBoard` | Board disable flag when supplied |

#### Row 6 — masjid identity and clock configuration

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `masjidName1` | Primary masjid name |
| 1 | `masjidName2` | Secondary/alternate masjid name |
| 2 | `masjidUrl` | Masjid URL |
| 3 | `timeZone` | Time zone |
| 4 | `timeZoneMilli` | Time-zone offset/value in milliseconds |
| 8 | `islamicDateAdjust` | Islamic date adjustment |
| 9 | `forceDate30` | Force Islamic month/day length setting |
| 10 | `masjidname1TopPad` | Masjid-name top padding |
| 11 | `puIslamicTimes` | Perpetual/upcoming Islamic-times setting |
| 12 | `bigClockTopPad` | Big-clock top padding |
| 13 | `clockNameFormat` | Clock/name formatting |
| 14 | `islamicTime12Hour` | Islamic time 12-hour formatting |
| 15 | `defaultRightToLeft` | Default RTL setting |
| 5–7 | — | Islamic year/month/date values are commented out in the source |

#### Rows 7–9 — religious content configuration

| Row | Fields | Meaning |
|---:|---|---|
| 7 | `slideDuration`, `ayahTrueFalse` | Slide duration and Ayah enable flag |
| 8 | `hadithTrueFalse` | Hadith enable flag |
| 9 | `sunnahTrueFalse`, `sunnahRefCheck` | Sunnah enable/reference settings |

The source shown here does not assign the actual Ayah/Hadith/Sunnah text in `handleResults()`. This indicates that some content may be obtained or generated elsewhere in the frontend and requires further investigation.

#### Row 10 — Taleem

Two configurable Taleem entries are mapped:

```text
Taleem 1: time, date, resident, address, enabled
Taleem 2: time, date, resident, address, enabled
```

#### Rows 11–12 — announcements

The source maps announcement slots 2–10. Each slot has a heading, content, and enabled flag. The first announcement slot is not mapped in this section and requires investigation elsewhere in the frontend.

```text
Row 11: announcements 2–5
Row 12: announcements 6–10
```

#### Row 13 — Nikah

Fields include:

```text
nikahNameOne
nikahGroomRelation
nikahRelationOne
nikahRelationTwo
nikahNameTwo
nikahDate
nikahTime
nikahPopUp
nikahTrueFalse
nikahPopUpTrueFalse
nikahBride
```

#### Row 14 — funeral

Fields include:

```text
funeralName
funeralRelation
funeralAddy
funeralPickup
funeralCemetery
funeralSalaahVenue
funeralSalaahTime
funeralTrueFalse
```

#### Row 15 — Taleem/Dawah/Gasht/three-day programme

Fields include Masjid Taleem, Gasht in/out days and times, an enable flag, a three-day header, two three-day programme entries with dates, and a Dawah enable flag.

#### Row 16 — community posters

Ten community poster image identifiers are mapped, each paired with an enable/visibility flag:

```text
commPosterImage1..10
commPosterTrueFalse1..10
```

The older text-based community-broadcast fields in this section are commented out in the source.

#### Row 17 — Eidgah

Fields include:

```text
eidDate
eidVenue
eidAddress
eidLecture
eidSalaah
eidgahTrueFalse
```

#### Row 18 — standard posters

Ten poster image identifiers are mapped with visibility flags. Additional layout settings are supplied for the first poster:

```text
posterImage1..10
posterTrueFalse1..10
posterTopPadding
posterHeight
posterWidth
```

#### Row 19 — large posters

Ten large-poster image identifiers are mapped with visibility flags and individual durations. Additional settings include:

```text
oldTotalPosterDuration
bigPosterSecondScreen
```

This confirms that poster presentation is configurable rather than simply a single image.

#### Row 20 — refresh and banking

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `sheetRefreshRate` | Client API refresh interval; the source converts this value to its timer behaviour |
| 1 | `bankHeader` | Banking heading |
| 2 | `bankName` | Bank name |
| 3 | `accountName` | Account name |
| 4 | `bankCode` | Bank code |
| 5 | `accountNum` | Account number |
| 6 | `bankTrueFalse` | Banking display flag |
| 7 | `bankAUSBSB` | Australian BSB value |

The API polling function calls `getAPI(boardId)` again after `sheetRefreshRate`, confirming that the refresh interval is supplied by the board data rather than being hard-coded by the client.

#### Row 21 — sickness/well-wishes messages

Ten configurable message fields are mapped:

```text
sick1..sick10
sickTrueFalse
```

#### Row 22 — alternate/secondary astronomical times

This row contains a second set of astronomical/prayer-start values:

```text
sunsetI
eshaStartsI
suhurEndsI
fajrStartsI
sunriseI
ishraaqI
duhaI
istiwaCI
istiwaI
zuhrStartsI
asrShafiI
asrHanafiI
```

It also contains ticker timing and translated/secondary values:

```text
tickerSpeed
tsehriEnds
tfajrStarts
tsunrise
tishraaq
tduha
tisiwaCaution
tistiwa
tzuhrStarts
tasrShafi
tasrHanafi
tsunset
teshaStarts
```

#### Row 23 — alternate-language daily Salah and masjid names

| Column | Variable | Meaning |
|---:|---|---|
| 0 | `maghribAthanI` | Alternate-language Maghrib Adhan |
| 1 | `maghribJamaahI` | Alternate-language Maghrib Jamaah |
| 2 | `eshaAthanI` | Alternate-language Esha Adhan |
| 3 | `eshaJamaahI` | Alternate-language Esha Jamaah |
| 4 | `fajrAthanI` | Alternate-language Fajr Adhan |
| 5 | `fajrJamaahI` | Alternate-language Fajr Jamaah |
| 6 | `dhuhrAthanI` | Alternate-language Zuhr/Dhuhr Adhan |
| 7 | `dhuhrJamaahI` | Alternate-language Zuhr/Dhuhr Jamaah |
| 8 | `asrAthanI` | Alternate-language Asr Adhan |
| 9 | `asrJamaahI` | Alternate-language Asr Jamaah |
| 10 | `masjidName1I` | Alternate-language masjid name 1 |
| 11 | `masjidName2I` | Alternate-language masjid name 2 |
| 12 | `masjidNameLanguage` | Masjid-name language setting |
| 13 | `masjidNameOnTop` | Masjid-name-on-top setting |

#### Row 24 — not read by `handleResults()`

The supplied `functions_uo_latest.js` contains no `spreadsheetArray[24]` access in `handleResults()`. The row exists in the 29-row API response but its semantics remain unresolved.

#### Row 25 — configurable secondary/continuous large posters

Ten additional poster image identifiers, visibility flags, and durations are mapped:

```text
cBigPosterImage0..9
cBigPosterTrueFalse0..9
cBigPosterDuration0..9
```

#### Rows 26–28 — not read by `handleResults()`

The supplied `functions_uo_latest.js` contains no `spreadsheetArray[26]`, `[27]`, or `[28]` accesses in `handleResults()`. Their exact purpose therefore remains unresolved from this source alone. They may be reserved, consumed by another code path, or represent data that is not currently used by the main client mapping.

## API polling behaviour

The captured source shows the client requesting:

```text
https://api.masjidboardlive.com/mblapi?id=${boardId}
```

The JSON response is passed to `handleResults()`, then `checkForChange()`, and the client schedules another `getAPI(boardId)` using `sheetRefreshRate`.

This confirms that MasjidPi should support periodic synchronisation and should not assume that one daily download is sufficient.

The exact unit/value used by `sheetRefreshRate` should be verified against an actual API response before choosing the MasjidPi default polling interval.

## Images/media

The board requests referenced images through:

```text
https://api.masjidboardlive.com/imageproxy?id=<image-id>
```

The JavaScript source also constructs image-proxy URLs from the opaque image identifiers stored in the board data. MasjidPi should therefore treat image IDs as upstream media references and cache the resulting assets locally.

## Client behaviour / display model

The source confirms that MasjidBoard Live is a dynamic multi-content presentation with configurable slides, durations, poster visibility, large-poster durations, language/time settings, and special content sections.

MasjidPi should therefore model the board as **structured content plus display configuration**, followed by its own display scheduler. It should not reproduce the upstream positional array throughout the application and should not simply render a static prayer-time dashboard.

## Important architectural conclusions

1. **A structured upstream data source is available.** MasjidPi does not need to scrape the rendered HTML for the core board data.
2. **The upstream schema is positional and opaque.** The provider must translate it into a normalised MasjidPi model.
3. **The provider should preserve optional/configurable content.** A masjid may use only a subset of the available fields.
4. **Display configuration is data.** Slide duration, poster visibility, language, RTL, time-formatting, and related settings are part of the upstream board state.
5. **The MasjidPi display should remain independent from the upstream presentation implementation.** We should reproduce the information and functionality, but not copy the JavaScript UI architecture.
6. **Local caching remains required.** The board depends on remote data and media, while the Raspberry Pi should continue displaying the last known valid state when connectivity is lost.

## Decisions

| Decision | Status | Reason |
|---|---|---|
| MasjidBoard is independent from audio playback | Confirmed | Audio must continue if MasjidBoard is unavailable. |
| MasjidBoard Live is the primary initial data source | Confirmed | User requirement. |
| HDMI external display is the target | Confirmed | User requirement. |
| Full MasjidBoard content is in scope | Confirmed | User requirement. |
| Local caching is required | Confirmed | Reliable appliance/offline behaviour. |
| Data-provider abstraction | Proposed | Allows another source later without redesigning the display layer. |
| Native MasjidPi display rather than a browser wrapper | Proposed | Better control of offline behaviour, resources, layout, and integration. |
| Consume structured `mblapi` data rather than scrape rendered HTML | Confirmed by investigation | HAR and frontend source show the live application itself consumes this endpoint. |
| Normalise MasjidBoard Live's positional schema inside the provider | Proposed | Keeps the opaque upstream schema out of the rest of MasjidPi. |
| Treat the HDMI board as a scheduled presentation | Proposed | The live board is a dynamic multi-content presentation with configurable slide behaviour. |

## Open Questions

- How is the public `mid` mapped to the internal `boardId`?
- Is `boardId` stable for a masjid?
- Can `mblapi?id=<boardId>` be requested directly without first loading the Premium page?
- What authentication, rate limits, or access restrictions apply to `mblapi`?
- What exact versioning guarantees exist for `mblVersion` and the positional schema?
- What is the exact semantic meaning of rows 24 and 26–28?
- Where are the actual Ayah/Hadith/Sunnah text values populated?
- Where is announcement slot 1 populated?
- How are announcement schedules and expiry represented?
- How are Jumu'ah services represented when there are multiple services?
- How are posters and images associated with individual content items?
- What is the exact unit of `sheetRefreshRate` and what values are normally supplied?
- Which data is generated/calculated by the client rather than supplied by `mblapi`?
- What are the exact carousel/slide rules used by the current MasjidBoard Live frontend?
- How does the board handle date changes, Ramadan/Eid changes, and other special-day transitions?
- What data should be cached indefinitely versus refreshed frequently?

## Next Investigation Steps

1. Obtain the API response for several different masjids/board IDs and compare the 29-row schema.
2. Determine how a public `mid` is resolved to a `boardId`.
3. Test direct access to `mblapi` and `imageproxy`, including behaviour without browser session state.
4. Trace the remaining frontend functions that populate Ayah, Hadith, Sunnah, announcements, and other content not assigned directly by `handleResults()`.
5. Investigate rows 24 and 26–28.
6. Determine the exact refresh interval/value semantics from multiple boards.
7. Determine the frontend carousel/slide scheduling rules.
8. Build a draft normalised MasjidPi data model from the verified field mapping.
9. Only then begin production MasjidBoard module implementation.

## Implementation Guardrail

Do not implement the production MasjidBoard module until the core API/data investigation is complete enough to define a stable internal data model.

The eventual implementation should separate:

```text
MasjidBoard Live provider
    -> Normalised board data
        -> Persistent cache
            -> Board state / scheduler
                -> HDMI display
```

## Research Log

### 2026-08-13

- Created `research/masjidboard-live` from `main`.
- Confirmed full MasjidBoard content is required.
- Confirmed external HDMI display is the target.
- Confirmed MasjidBoard Live is the primary data source.
- Confirmed MasjidBoard must remain independent from audio playback.
- Established that API/data investigation should precede production implementation.
- Analysed a Firefox HAR capture from a live Premium board.
- Confirmed the structured `https://api.masjidboardlive.com/mblapi?id=<boardId>` endpoint.
- Confirmed the endpoint returns JSON containing the same 29-row `theInfo` structure embedded in the page.
- Confirmed the page exposes `boardId` and `mblVersion` values.
- Confirmed `imageproxy?id=<image-id>` is used for referenced media.
- Confirmed the live board is a dynamic, multi-content presentation rather than a static HTML page.
- Reviewed an existing Home Assistant integration as supporting evidence; it polls every 600 seconds and exposes only the five daily prayer Adhan/Jamaah pairs.
- Captured and inspected `functions_uo_latest.js`.
- Verified the `handleResults()` mapping for rows 0–23 and 25.
- Confirmed the client schedules subsequent API requests using the board-supplied `sheetRefreshRate`.
- Identified unresolved rows 24 and 26–28 and content populated outside the main `handleResults()` mapping.
- Established that the next priority is comparing the schema across multiple boards and tracing the remaining frontend content paths.
