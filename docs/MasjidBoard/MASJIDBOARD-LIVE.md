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

**Important source-code observation:** `handleResults()` assigns `spreadsheetArray[3][6]` to both `iftar` and `maghribAthan`. This does not establish that Iftar and Maghrib are semantically identical. A masjid may publish Iftar at a time that differs from the Maghrib Adhan by a few minutes, depending on its timetable/practice. MasjidPi must therefore retain this as an open research question rather than collapsing the two concepts in the normalised model. During Ramadhaan, compare real MasjidBoard Live boards that explicitly display both Iftar and Maghrib times to determine whether MasjidBoard Live supports distinct values or only exposes a shared source value. Row 23 separately contains `maghribAthanI` and other alternate-language Salah values.

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

## Existing Payload Content Audit

An August 2026 review confirmed that a new MasjidBoard Live discovery exercise is **not** required before work begins on notices and other community content. The repository already contains captured 29-row Premium responses in `docs/MasjidBoard/*.json` and `docs/MasjidBoard/fixtures/*.json`, and the row mapping above already documents the relevant content.

The existing captures contain the following information:

| Content | Existing upstream location | Current MasjidPi state |
|---|---|---|
| Announcements | Rows 11–12; heading, HTML content and display flag for slots 2–10 | Present in captured payloads; not parsed or exposed by the display API |
| Nikah notice | Row 13; names/relations, bride, date, time, popup value and visibility flags | Present in captured payloads; `model.NoticeTypeNikah` exists, but the provider does not populate it |
| Funeral notice | Row 14; deceased name/relation, address, pickup, cemetery, salaah venue/time and visibility flag | Present in captured payloads; `model.NoticeTypeFuneral` exists, but the provider does not populate it |
| Taleem/Dawah/Gasht programmes | Rows 10 and 15; programme details, dates/times and visibility flags | Present in captured payloads; `model.Programme` exists, but the provider does not populate it |
| Community posters | Row 16; ten image identifiers with visibility flags | Present in captured payloads; `model.Media` exists, but the provider does not populate it |
| Eid notice | Row 17; date, venue, address, lecture, salaah and visibility flag | Present in captured payloads; `model.NoticeTypeEid` exists, but the provider does not populate it |
| Standard posters | Row 18; ten image identifiers, visibility flags and layout settings | Present in captured payloads; not parsed or exposed |
| Large posters | Row 19; ten image identifiers, visibility flags and durations | Present in captured payloads; not parsed or exposed |
| Banking/contributions | Row 20; heading, bank/account details, branch code and display flag | Present in captured payloads; `model.Banking` exists, but the provider does not populate it |
| Sickness/well-wishes | Row 21; ten configurable messages and visibility state | Present in captured payloads; `model.NoticeTypeWellWish` exists, but the provider does not populate it |
| New moon information | Row 2 and related settings; birth, set, age, azimuth, altitude and visibility dates | Present in captured payloads; `model.NewMoon` exists only as a placeholder and is not populated |

### What is already implemented

The normalised `model.Board` already reserves collections or objects for `Announcements`, `Programmes`, `Notices`, `Media`, `Banking` and `NewMoon`. Notice types already include general, Nikah, funeral, well-wishes and Eid. This means the high-level domain boundary was deliberately prepared for this content and should be extended rather than redesigned.

The current provider parser intentionally normalises only the verified core rows:

- row 1: Jumu'ah;
- row 2: clock/identity context;
- row 3: daily prayer times;
- row 5: astronomical times; and
- row 6: masjid identity/configuration.

The current display view and `/api/masjidboard/display` response expose identity, dates, prayers, Jumu'ah and astronomical information only. They do not yet include announcement, programme, notice, media, banking or new-moon fields.

### Core versus Premium payloads

MasjidPi currently retrieves the public Core board data for normal operation. The Core validation found no announcements, community content, posters, programmes or notices in its embedded `data` object. The richer content listed above is confirmed in the captured Premium 29-row payloads.

### 2026-08-23 — Premium access path revalidated

The public Premium page and API path were rechecked for five previously researched boards:

- `brits-jamia`;
- `brits-taqwa`;
- `erasmia-aaisha`;
- `zakariyya-park-duzak`; and
- `fawkner-masjid`.

For every board, `https://premium.masjidboardlive.com/v2/?mid=<mid>` still returned a generated page containing both `let boardId = "<opaque-id>"` and `let theInfo = [...]`. Calling `https://api.masjidboardlive.com/mblapi?id=<opaque-id>` returned a valid 29-row JSON array for every resolved ID. Existing repository fixtures also demonstrate that MasjidBoard Live may append a 30th row, so the provider accepts 29 or more rows while continuing to parse only verified positions.

This confirms the existing technical access path remains operational and unauthenticated. It does not by itself establish a contractual entitlement or guarantee that every Core-listed masjid has a Premium board. MasjidPi must therefore treat Premium enrichment as optional and must retain the current Core provider as the reliable timetable fallback.

The provider includes a `PremiumClient` that resolves the opaque ID and embedded payload from the stable public `mid` on each fetch. The opaque ID is not persisted, so a server-side board rebuild cannot leave MasjidPi tied to a stale implementation identifier.

Runtime integration uses an `EnrichedClient` with explicit fallback semantics:

1. Fetch the public Core board first.
2. If Core fails, preserve the existing last-known-good runtime/cache behaviour; Premium is not used as a timetable substitute.
3. If Core succeeds, attempt Premium enrichment.
4. If Premium succeeds, merge only announcements, programmes, notices, media, banking and new-moon content into the Core board.
5. If Premium is absent or fails, return the successful Core board unchanged and keep the board status `current`.

Core therefore remains authoritative for identity, dates, prayer times, Jumu'ah and astronomical times. Premium cannot replace or contradict the operational timetable.

The read-only `/api/masjidboard/display` contract exposes active normalised `announcements` and `notices` when enrichment is available. The Detailed layout renders this content in an adaptive three-slot, theme-aware panel occupying the right quarter of the screen, with the timetable using the remaining three quarters. Three compact cards can appear together; dense content can span two slots; one item uses the full panel; and two remaining items use equal halves. Additional pages rotate automatically, duplicate active items are suppressed, and upstream HTML is converted to plain text rather than inserted into the DOM.

## Notice and Announcement Display Fixtures

The Premium payloads contain useful historical content even when its upstream visibility flag is `Hide`. This content should be retained as **development fixture material only**: it is valuable for layout design and automated tests, but it must never be presented to users as a current live notice.

The following inventory was revalidated against the live Premium payloads on 23 August 2026.

### Content-size range

| Content type | Observed plain-text range | Typical structure |
|---|---:|---|
| General announcements | approximately 43–272 characters | Heading plus free-form body, often containing dates, times and several sentences |
| Nikah notices | approximately 114–137 characters | Names and family relationships, event date, time and sometimes venue |
| Funeral notices | approximately 86–164 characters | Deceased person, relationship, address/pickup, cemetery, Janazah venue and time |
| Eid notices | approximately 75–93 characters | Date, venue, address, lecture/translation time and Salaah time |

These ranges describe the currently inspected samples, not hard maximums. Layouts must remain defensive when content exceeds them.

### Representative fixtures by board

| Board | Content worth retaining as a fixture | Display condition exercised |
|---|---|---|
| `zakariyya-park-duzak` | Active short masjid announcement; a hidden 272-character Salaah-time-change notice; Arabic announcement content | Live end-to-end content, duplicate content, long body, mixed scripts and RTL |
| `fawkner-masjid` | Long safety/access notices of approximately 194–255 characters; detailed funeral notice; Eid and Jumu'ah notices | Multi-sentence notices, long venue text, operational instructions and Australian terminology |
| `azaadville-darul-uloom` | Detailed funeral notice of approximately 164 characters; Taraweeh and weekly-programme announcements | Dense structured funeral information and programme content |
| `brits-jamia` | Weekly programmes, Eid-style announcement, Nikah, funeral and Eid fields | Medium-length general and structured notices |
| `brits-taqwa` | Short urgent announcement, weekly programmes, Nikah, funeral and Eid fields | Very short urgent content and structured notices |
| `erasmia-aaisha` | Banking/contribution information stored in announcement slots; Nikah, funeral and Eid fields | Content-category ambiguity and financial information that should not be mistaken for an urgent notice |

### Required fixture states

The display design and frontend tests should cover at least these states:

1. **No content** — the normal state for many Core-only boards and Premium boards with everything hidden.
2. **Single short announcement** — one heading and a one-line or short body.
3. **Single long announcement** — approximately 250–300 characters with several times or sentences.
4. **Multiple announcements** — ordering and rotation without overcrowding the prayer board.
5. **Duplicate active announcements** — Zakariyya currently publishes the same active announcement in two slots; the UI must not visibly repeat it without purpose.
6. **Funeral notice** — structured identity, location and Janazah details with visually high priority.
7. **Nikah notice** — structured names, relationships, date, time and venue where supplied.
8. **Eid notice** — event date, venue and multiple event times.
9. **Arabic/RTL content** — correct direction, shaping, alignment and font behaviour.
10. **Mixed-language content** — Arabic and English within the same notice set.
11. **Legacy HTML content** — line breaks, bold/underline tags and malformed markup converted safely to display text.
12. **Generic or blank heading** — headings such as `Masjid Announcement`, `Change Heading 2`, or an empty heading must not dominate the design.
13. **Sensitive personal information** — funeral and Nikah fixtures may contain names and addresses; screenshots, demos and public test data should use anonymised equivalents.
14. **Ambiguous category** — banking, programme or Salaah-change content placed in a generic announcement slot needs presentation based on supplied data, not an invented urgency level.

### Fixture-handling rules

- Preserve the source payload separately from its active/hidden state.
- Never convert a historical `Hide` entry into a live notice in production.
- Use anonymised copies for screenshots, demos and committed frontend fixtures when the original contains personal information.
- Preserve representative length, structure, HTML irregularities and multilingual behaviour when anonymising.
- Treat upstream HTML as untrusted input. Do not render it with `innerHTML` without an explicit sanitisation policy.
- Do not assume that announcement slot order implies urgency.
- Do not assume that repeated content represents two distinct events.
- Keep structured notice fields available so the final layout does not have to reverse-engineer a concatenated text body.

### Current live test source

At the time of this review, `zakariyya-park-duzak` supplied two active announcement slots with the same heading and body. It is the best existing end-to-end test board for retrieval, normalisation, cache and display-API validation. The other reviewed boards had no active announcement, Nikah, funeral or Eid entries, but their hidden historical data remains useful for anonymised layout fixtures.

Therefore the next investigation is narrowly defined:

1. Confirm whether the existing 29-row Premium endpoint can be used reliably and legitimately for selected Core-listed boards, or whether a separate provider/capability boundary is required.
2. Revalidate the already-mapped rows against a small number of current payloads, focusing on visibility flags, optional fields and lifecycle behaviour rather than rediscovering field names.
3. Add defensive parsers that populate the existing normalised model.
4. Extend caching and the display API only after parsing tests prove the content can be normalised safely.

The existing fixtures should be the starting point for parser tests. New captures are needed only to validate freshness and edge cases, not to repeat the completed schema discovery.

### Detailed-layout fixture mode

The Detailed display supports an explicit development-only fixture mode:

```text
/masjidboard.html?layout=detailed&notice-fixtures=1
```

This injects anonymised funeral, Nikah, Eid, long-announcement, Arabic and compact community samples into the frontend card renderer. The set exercises two-thirds/one-third, three-card and full-panel pages. The fixtures are derived from the historical content shapes documented above, are labelled as non-live layout fixtures, and never enter the provider, normalised model, runtime cache or display API.

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

### Iftar and Maghrib remain distinct concepts pending Ramadhaan validation

The current `functions_uo_latest.js` assigns the same upstream row-3 column to both `iftar` and `maghribAthan`. This is a **source mapping observation, not a semantic decision**.

MasjidPi should treat Iftar and Maghrib Athaan as distinct concepts until the upstream behaviour is verified. A masjid may publish Iftar a few minutes before or otherwise separately from the Maghrib Athaan according to its own timetable/practice.

The planned validation is to observe multiple MasjidBoard Live boards during Ramadhaan where both Iftar and Maghrib are explicitly displayed. The investigation should determine whether:

- MasjidBoard Live supports distinct Iftar and Maghrib values;
- both are always derived from one shared source value; or
- the frontend displays a single value under different labels even when the underlying masjid timetable distinguishes them.

Until that evidence exists, the normalised model must not silently equate Iftar with Maghrib Athaan.

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

## Research Plan

MasjidBoard remains a **research/development effort** and is not ready to be merged into `main` or included in the next MasjidPi release. The work should proceed through the following four research stages before production integration is considered.

### Stage 1 — Complete the upstream schema

Produce a definitive mapping of the MasjidBoard Live `theInfo` response and the corresponding frontend usage.

For every relevant row and column, record:

- row and column position;
- source variable name;
- observed meaning;
- example values;
- confidence level;
- whether it is data, configuration, presentation state, or derived information;
- known dependencies on other fields; and
- unresolved questions.

The output of this stage is the **MasjidBoard Live data contract** as currently understood by MasjidPi.

This stage must also trace the remaining Jumu'ah heading/code semantics and investigate fields that `handleResults()` does not fully populate or explain.

### Stage 2 — Cross-board validation

Validate the Stage 1 mapping against multiple real MasjidBoard Live boards rather than relying on a single capture.

Compare the same fields across the currently verified boards and add further representative boards where necessary. Particular attention should be given to:

- Jumu'ah heading/code combinations;
- optional and absent values;
- alternate-language values;
- announcements and programmes;
- posters and media;
- Islamic-date adjustments;
- board configuration differences;
- fields that appear only on some boards; and
- Iftar versus Maghrib behaviour during Ramadhaan.

The objective is to distinguish **actual upstream semantics** from values or behaviours that merely occur on one board.

### Stage 3 — Define the normalised MasjidPi model

Only after Stages 1 and 2 should the generic MasjidBoard domain model be finalised.

The model should:

- represent semantic concepts rather than the 29-row transport structure;
- keep MasjidBoard Live-specific parsing inside the provider;
- make fundamental identity and daily prayer times reliable;
- make all other content optional;
- preserve detailed Jumu'ah information without confusing it with the five daily prayers;
- preserve Iftar and Maghrib as separate concepts unless upstream evidence proves they are equivalent; and
- leave insufficiently understood upstream fields outside the model until their semantics are established.

The result should be a provider-independent model that could accept another MasjidBoard data source later.

### Stage 4 — Provider implementation and validation

Implement the MasjidBoard Live provider against the normalised model.

The provider should:

```text
public mid
    |
    v
GET /v2/?mid=<mid>
    |
    +--> resolve boardId
    |
    +--> obtain theInfo
    |
    v
parse upstream data
    |
    v
normalised Board
    |
    +--> tests using captured fixtures
```

Implementation should then proceed through:

1. Complete provider parsing for the verified schema.
2. Fixture-based tests using representative real boards.
3. Last-known-good local caching.
4. Refresh and recovery behaviour.
5. Media/image acquisition and caching.
6. Validation of the provider independently from the display layer.
7. Only after the provider is stable, implementation of the HDMI display scheduler/rendering layer.

Completion of Stage 4 does **not** automatically mean MasjidBoard is ready for release. Release integration should only be considered after the provider, cache, display behaviour, and real-board validation are sufficiently mature for production use.

## Current Status

The upstream investigation has established the main 29-row schema and the public-`mid`/opaque-`boardId` relationship. The domain model now explicitly distinguishes the standard five-prayer timetable from the detailed Friday Jumu'ah information.

On Friday, Jumu'ah supplies the Adhan and Salaah/Jamaah values displayed in the Dhuhr slot of the standard timetable, with Khutbah as the Salaah fallback when no Salaah time is supplied. A separate Friday element will display the richer Jumu'ah information.

Iftar versus Maghrib remains an open upstream-semantics question. The current source maps both names to one value, but MasjidPi will not assume they are semantically identical until Ramadhaan observations establish how MasjidBoard Live represents boards that distinguish the two times.

The MasjidBoard implementation remains isolated on `research/masjidboard-live` and is **not part of the current MasjidPi release path**.

## Immediate Next Step

Begin **Stage 1 — Complete the upstream schema** by tracing the remaining `functions_uo_latest.js` assignments and DOM/display usage, with particular focus on:

1. Jumu'ah heading-code semantics.
2. Rows 22–28 and fields whose assignments are incomplete or commented out.
3. Media/image identifiers and their retrieval contract.
4. Fields whose values are derived rather than directly supplied.
5. Any upstream fields required to reproduce the full range of MasjidBoard Live content.

No further production-facing MasjidPi integration should be undertaken until the Stage 1 and Stage 2 research has established a sufficiently reliable upstream data contract.
