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

The investigation must account for all content currently presented by MasjidBoard Live, including prayer and astronomical times, Adhan/Iqamah, Jumu'ah and Khateeb information, daily Ayah, Hadith, Sunnah, Du'a, community broadcasts, announcements, request-for-Du'a information, Nikah and funeral notices, weekly programmes, masjid information, contribution information, posters/images, New Moon information, and Ramadan/Eid-related information where supplied.

This list is provisional and must be expanded if the underlying data reveals additional content.

## Investigation Areas

### Masjid identification

Determine how masjids are identified, what the `mid` value represents, how search works, what metadata is returned, and how MasjidPi should persist the selected masjid.

### Data/API discovery

Determine the underlying client requests used by the MasjidBoard Live application rather than relying on rendered-page scraping. Record endpoints, methods, parameters, authentication, response formats, relevant fields, errors, and update characteristics.

### Prayer data

Determine the exact representation and semantics of all prayer and astronomical times, including Adhan versus Iqamah, Asr calculations, time zones, future changes, missing values, and special-day behaviour.

### Jumu'ah

Determine how multiple Jumu'ah services and their Adhan, Sunan, lecture, Khutbah, Salah, Khateeb, and service-specific information are represented.

### Content and media

Determine how text, rich content, links, posters, images, ordering, visibility, scheduling, and expiry are represented and delivered.

### Refresh and offline behaviour

Determine refresh intervals and what should be cached. Define behaviour for network loss, service failure, malformed data, failed images, stale data, and reboot without connectivity.

## Findings

### Live board is data-rich and varies by masjid

Inspection of multiple current Premium boards shows that the content is substantially broader than the five daily prayers. Examples include:

- Daily Salah Adhan/Iqamah times.
- Suhur, Fajr start, Sunrise, Ishraq, Duha, Istiwa, Zuhr start, multiple Asr calculations, Sunset, and Isha start.
- Next Salah / upcoming change information.
- Jumu'ah with fields such as Adhan, Sunan, Lecture, Khutbah, Salah, and Khateeb; not every masjid supplies every field.
- Daily Ayah and Hadith.
- Sunnah content on some boards.
- Community broadcasts and general announcements.
- Weekly programmes.
- Masjid notices.
- Nikah announcements.
- Masjid contribution/bank information.
- Du'a after Adhan.
- Special notices such as Ramadan/Eid material and New Moon information.
- Images/posters.

Current live boards demonstrate that content availability differs between masjids. This means the MasjidPi data model must support optional sections rather than assuming every board has identical content. citeturn0search4turn0search5turn0search6

### Masjid identifier

The Premium board URL uses a masjid-specific `mid` query parameter, for example `?mid=cravenby-estate-husami`, `?mid=ridgeway-quba`, and `?mid=zeerust-jaamiah`. MasjidPi will therefore need to retain the MasjidBoard Live masjid identifier as part of its configuration. citeturn0search4turn0search5turn0search6

### MasjidBoard Live is intended to be remotely updated

MasjidBoard Live's own site describes remote editing/updating from a mobile device or computer and specifically mentions short-notice announcements and unexpected Salah-time changes. It also describes suburb-based synchronisation of community and funeral notices. This means our integration needs to treat board content as mutable rather than a static daily schedule. citeturn0search0turn0search11

### Full-HD is the native Premium target

MasjidBoard Live's published Premium hardware requirements specify a Full HD 1920×1080 monitor. This aligns directly with our planned HDMI target and gives us a sensible primary display resolution for MasjidPi. citeturn0search0

### Slide-based display model

MasjidBoard Live's FAQ describes the Premium display as a carousel and says slides can be hidden, except for the Ayah slide. It also documents dedicated behaviour for some slides, such as Du'a after Adhan displaying for five minutes and a cellphone reminder displaying for five minutes before Iqamah. Scheduled notices and custom slides are also supported. This is important for MasjidPi: the display should be modelled as a scheduler/carousel rather than one fixed dashboard. citeturn0search10

### Existing third-party programmatic integration uses HTML scraping

A public Home Assistant integration exists and provides useful evidence about the current board's HTML structure. It fetches the configured board URL and parses the returned HTML with BeautifulSoup rather than calling a documented JSON API. It extracts the masjid name from an element with ID `masjidName2`, and extracts the five daily prayers using IDs following the pattern `<prayer>Athan` and `<prayer>Jamaah`, where the prayer IDs are `fajr`, `zuhr`, `asr`, `maghrib`, and `esha`. fileciteturn23file0L2-L2

The same integration uses a configurable polling interval with a default of 600 seconds (10 minutes). This is useful evidence for a conservative initial refresh interval, but it is not evidence that 10 minutes is the correct interval for every MasjidBoard content type. fileciteturn25file0L2-L2

The integration is explicitly limited to Adhan/Jamaah sensors for the five daily prayers, so it is not sufficient for MasjidPi's full-board requirements. fileciteturn22file0L2-L2

### No documented public API established yet

The investigation has **not yet established a documented/public JSON API** that can be relied upon for the complete board. Search results for `api.masjidboardlive.com` did not provide a verifiable public API specification, and the existing third-party integration we inspected does not use one.

This is important: we should not design the MasjidPi provider around an assumed API. We need direct inspection of the live application's network requests and JavaScript assets to determine whether a structured backend endpoint exists.

### Dynamic content

The live board presents loading placeholders and then populated content, which indicates that the application has a dynamic data-loading mechanism. The exact underlying requests still need to be identified.

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
| Do not assume an undocumented API exists | Confirmed | No verifiable public API specification has been established. |
| Primary display target is 1920×1080 | Confirmed | Matches MasjidBoard Live's published Premium requirement. |
| Display architecture should support scheduled slides | Confirmed | MasjidBoard Live itself uses a configurable carousel/slide model. |

## Open Questions

- What exact network requests does the live board make after loading?
- Is there a structured JSON/API endpoint behind the board?
- Does the board use one endpoint or several endpoints for different content categories?
- What authentication, tokens, cookies, or headers are required?
- What is the complete response/data schema?
- How are content schedules and expiry represented?
- How are images and posters delivered?
- How are prayer-time changes represented?
- How is time zone information represented?
- How are multiple Jumu'ah services represented?
- What refresh intervals are expected for each content type?
- What data should be cached permanently versus temporarily?
- Can the current board be reproduced from structured data without rendering the website?
- What terms/usage restrictions apply to programmatic consumption of MasjidBoard Live data?

## Next Investigation Step

The remaining high-value investigation requires inspecting the live application's actual browser network traffic and JavaScript assets. Web search can confirm rendered content and external documentation, but it cannot reliably expose the browser's complete runtime request sequence.

The next capture should use several representative boards and identify the requests made for:

1. Masjid metadata.
2. Daily prayer times.
3. Jumu'ah.
4. Announcements and notices.
5. Religious content.
6. Images/media.
7. Special dates/events.

Record the exact endpoints, request parameters, response structures, and update behaviour here before designing the MasjidPi internal data model.

## Implementation Guardrail

Do not implement the production MasjidBoard module until the core API/data investigation is complete enough to define a stable internal data model.

The eventual implementation should separate:

```text
Data provider
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
- Investigated multiple current Premium boards and confirmed substantial variation and a much broader content set than the daily prayer times.
- Confirmed the masjid-specific `mid` identifier used by the Premium board URLs.
- Confirmed MasjidBoard Live's published Premium target is Full HD 1920×1080.
- Confirmed the Premium board is carousel/slide based and supports slide-specific timing and scheduling.
- Inspected an existing public Home Assistant integration and confirmed it currently obtains prayer data by parsing HTML rather than a documented JSON API.
- Confirmed that integration uses a default 10-minute polling interval and extracts the masjid name plus five daily Adhan/Jamaah pairs.
- Did not find sufficient evidence to claim a documented public JSON API; direct browser network/JavaScript inspection remains the next step.
