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

### 2026-08-14 — Public `mid` and API `boardId` relationship resolved

The relationship between the public MasjidBoard URL identifier and the identifier used by `mblapi` has now been established from the generated board HTML.

A Premium board URL has the form:

```text
https://premium.masjidboardlive.com/v2/?mid=<public-mid>
```

For example:

```text
https://premium.masjidboardlive.com/v2/?mid=fawkner-masjid
```

The generated HTML contains a server-supplied JavaScript variable:

```javascript
let boardId = "170sRYVcxfOC-l3IGK0FeqWX8D8rM35afyl7nyL2SWHI";
let mblVersion = "15";
```

The same HTML contains the complete board data as:

```javascript
let theInfo = [...];
```

and then invokes:

```javascript
handleResults(theInfo);
startMBL();
```

The captured page therefore establishes the following flow:

```text
public mid
    |
    | GET /v2/?mid=<mid>
    v
server-generated board HTML
    |
    +--> boardId
    |
    +--> theInfo[29 rows]
    |
    v
functions_uo_latest.js
    |
    +--> handleResults(theInfo)
    |
    +--> getAPI(boardId)
              |
              v
        /mblapi?id=<boardId>
```

The important architectural conclusion is that **`boardId` is an opaque, server-supplied identifier**. It is not derived by the frontend JavaScript from the public `mid`.

Therefore MasjidPi must **not attempt to calculate or transform a public `mid` into an API ID**.

The public `mid` should remain the stable board identifier used by our MasjidBoard catalogue/provider configuration. The opaque `boardId` should be treated as a MasjidBoard Live implementation detail discovered from the generated board page.

### Confirmed board mappings

The following mappings were verified directly from the corresponding Premium board pages/network requests:

| Public `mid` | MasjidBoard Live `boardId` |
|---|---|
| `fawkner-masjid` | `170sRYVcxfOC-l3IGK0FeqWX8D8rM35afyl7nyL2SWHI` |
| `zakariyya-park-duzak` | `1GcUzmzuO3XM-xXblab-Mq0lE2ddp_qGZ0w-eeMYwlK8` |
| `erasmia-aaisha` | `1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k` |
| `brits-taqwa` | `1ZK8NtqROdU3Ww4THcHkHyDJN2gu98HC1ovBbGO7iooY` |
| `brits-jamia` | `1lTEEzl7sefO4W72c9iKxcUkB1ZReHhtZt9DFmCFfT_0` |
| `azaadville-darul-uloom` | `1Zpg5LKfd_ZoEQsA0rsyWNBrUgY6QVaHnGdPfuKHF24A` |

Earlier captured JSON files were temporarily associated with the wrong board names because the opaque IDs had been assigned to filenames without verifying the corresponding generated page. Those filename associations must not be treated as authoritative. The identity contained in the generated board data and the directly observed page/request relationship are authoritative.

### 2026-08-14 — `theInfo` is available directly in generated HTML

The generated Premium page contains the same 29-row board data that is subsequently consumed by `handleResults()`.

This is significant for the provider design because it gives us a reliable acquisition path based on the public `mid`:

```text
GET /v2/?mid=<public-mid>
        |
        +--> boardId
        |
        +--> theInfo
```

The provider does not need to know the opaque API ID in advance. It can obtain it from the page if it subsequently needs to use `mblapi` for refreshes.

For initial board discovery and association, the generated HTML is preferable to maintaining a separate hard-coded `mid -> boardId` catalogue.

### Provider acquisition decision

**Decision:** The MasjidBoard Live provider should use the public `mid` as its external board identifier and treat the generated board page as the authoritative mechanism for resolving the current `boardId` and obtaining the initial `theInfo` payload.

The provider should conceptually operate as:

```text
MasjidBoard public mid
        |
        v
GET /v2/?mid=<mid>
        |
        +--> extract boardId
        |
        +--> extract theInfo
        |
        v
parse theInfo
        |
        v
Normalised MasjidBoard data
```

If the provider subsequently uses `/mblapi`, it should use the `boardId` obtained from the generated page rather than a hard-coded value.

### Source-code caveat

The captured `functions_uo_latest.js` confirms that `getAPI(url)` calls:

```text
https://api.masjidboardlive.com/mblapi?id=<url>
```

and the application passes `boardId` to `getAPI()`.

The source therefore confirms the API relationship, but the provider should not depend on frontend implementation details beyond the externally observed page/API contract unless necessary.

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

It also contains ticker timing and translated/secondary values. The exact complete mapping of the remaining columns is still subject to verification.

#### Row 23 — alternate-language Salah values

This row contains alternate-language versions of the daily Salah/related values, including Maghrib Adhan and Jamaah and other time fields. The exact complete mapping should be verified against the full `handleResults()` source before it is exposed in the normalised model.

#### Row 24 onward — remaining configuration/data

Rows 24–28 contain additional board configuration, translation, ticker, poster, and display settings. These require no assumption-driven modelling at this stage. Only fields whose semantics are verified from the source and/or multiple live boards should be promoted into the normalised MasjidBoard model.

## Architectural Decisions

### MasjidBoard is a separate application capability

MasjidBoard must be capable of operating as a separate application from MasjidPi audio playback. An end user should ultimately be able to run:

- audio only;
- MasjidBoard display only; or
- both together.

The MasjidBoard display must not depend on the audio subsystem being active.

### MasjidBoard Live supplies data; MasjidPi renders the board

MasjidBoard Live is the primary data source. MasjidPi will **render the board itself** rather than embedding the MasjidBoard Live webpage.

This gives us control over the display, allows MasjidBoard to run as a standalone application, and avoids making the HDMI display dependent on a full web browser rendering the upstream site.

### Fundamental data is mandatory; other content is optional

The only mandatory board content is:

- fundamental board identity; and
- daily prayer times.

Everything else must be capable of being absent without making the board invalid. Optional content includes Jumu'ah, astronomical information, announcements, programmes, Eid, Nikah, funeral notices, posters, banking/contributions, moon information, Ayah/Hadith/Sunnah, and other board-specific content.

### Jumu'ah and the standard five-prayer timetable

**Decision:** The normalised domain retains the five daily prayer slots: Fajr, Dhuhr, Asr, Maghrib and Esha. Jumu'ah is not a sixth daily prayer and does not replace the Dhuhr field in the underlying domain model.

On Friday, Jumu'ah replaces **Dhuhr only in the standard five-prayer timetable presentation**, specifically for the Adhan and Salaah/Jamaah times shown in that Dhuhr slot.

The Friday resolution rules are:

1. Use the Jumu'ah Adhan when it is available.
2. Use the Jumu'ah Salaah/Jamaah time when it is available.
3. If a Jumu'ah Salaah/Jamaah time is not supplied, use the Jumu'ah Khutbah time as the Salaah-time fallback.
4. The underlying Dhuhr field remains Dhuhr and may retain the supplied Dhuhr/astronomical calculation data.
5. The Friday timetable therefore presents the Dhuhr slot using the Jumu'ah congregational times rather than Dhuhr Jamaah.

Conceptually:

```text
Normal day:
    Dhuhr slot -> Dhuhr Adhan + Dhuhr Salaah

Friday:
    Dhuhr slot -> Jumu'ah Adhan + Jumu'ah Salaah
                           |
                           +-> Khutbah fallback if Salaah unavailable
```

### Detailed Friday Jumu'ah information

A separate Friday Jumu'ah element will display the richer Jumu'ah information supplied by MasjidBoard Live. This detailed element is independent of the five-prayer timetable representation and may show the detailed sequence of Jumu'ah headings/times, Khutbah information, Khateeb information, and other supplied Friday-specific data.

The conceptual separation is:

```text
PrayerTimes
    |
    +-- Fajr
    +-- Dhuhr
    +-- Asr
    +-- Maghrib
    +-- Esha
    |
    +-- Friday presentation of Dhuhr slot
            -> Jumu'ah Adhan
            -> Jumu'ah Salaah
            -> Khutbah fallback

Jumu'ah details
    +-- Adhan
    +-- detailed heading/time events
    +-- Khutbah
    +-- Salaah/Jamaah
    +-- Khateeb
    +-- other supplied Friday information
```

This distinction is intentional: the standard timetable needs only the Jumu'ah Adhan and Salaah semantics needed to replace the Friday Dhuhr presentation, while the dedicated Friday element can preserve and display richer upstream Jumu'ah information.

### Jumu'ah may contain multiple detailed entries/events

The MasjidBoard Live upstream mapping supports up to three Jumu'ah heading/time pairs and separate Khateeb, Adhan and Jamaah values. The detailed normalised representation should therefore preserve these as optional detailed Jumu'ah information rather than collapsing the complete upstream row into a single prayer time.

The exact Go representation of the detailed events should be frozen only after the remaining JavaScript assignment/display path has been traced. The parser must not invent semantics for positional values that have not been verified.

### Do not mirror the 29-row structure in the domain model

The 29-row response is an upstream transport/configuration format. It should not become the MasjidPi domain model.

The provider should parse the upstream response into a normalised model containing semantically verified fields. Unknown or insufficiently understood upstream fields should remain outside the domain model until their meaning is established.

### Ayah/Hadith/Sunnah are deferred from the initial implementation

Ayah, Hadith and Sunnah are not critical to the initial implementation. Their provider mappings can remain under investigation while the core board is implemented.

### Local prayer-time representation

Prayer times in the normalised domain are represented as local clock values using hour and minute rather than `time.Time`. The board timezone is stored separately on the board identity.

```go
type PrayerTime struct {
    Hour   int
    Minute int
}
```

This is appropriate because a prayer time is fundamentally a local mosque clock time; the calendar date and timezone are contextual.

## Current Status

The upstream investigation has established the main 29-row schema and the public-`mid`/opaque-`boardId` relationship. The domain model now explicitly distinguishes the standard five-prayer timetable from the detailed Friday Jumu'ah information.

On Friday, Jumu'ah supplies the Adhan and Salaah/Jamaah values displayed in the Dhuhr slot of the standard timetable, with Khutbah as the Salaah fallback when no Salaah time is supplied. A separate Friday element will display the richer Jumu'ah information.

## Next Step

Before adding further parser fields, use the generated HTML (`View Source`) and the captured `functions_uo_latest.js` together as the authoritative upstream fixtures. This avoids the earlier problem of associating an opaque API response with the wrong public board.

Then:

1. Trace the remaining Jumu'ah assignments and DOM/display usage so each detailed heading/time field is understood.
2. Freeze the normalised detailed Jumu'ah representation around verified semantics.
3. Add parser fixtures/tests using complete generated-page data for representative boards.
4. Implement the MasjidBoard Live provider using the public `mid`, resolving `boardId` from the generated page rather than hard-coding it.
5. Implement last-known-good local caching and refresh/recovery behaviour.
6. Build the standalone HDMI display application against the normalised model.
