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

Examples: the Gateway Musallah board exposes the full perpetual-time section and a Nikah announcement; the Musjidus Salaam board exposes announcements, multiple Nikah announcements, a weekly programme and contribution information; and the Husami Masjid board exposes Jumu'ah, several change/notice headings, weekly programmes, contribution information, and a Nikah announcement. citeturn0search1turn0search2turn0search5

### Masjid identifier

The Premium board URL uses a masjid-specific `mid` query parameter, for example `?mid=umhlanga-gateway` or `?mid=cravenby-estate-husami`. This confirms that MasjidPi will need to retain a MasjidBoard Live masjid identifier as part of its configuration. citeturn0search1turn0search5

### No documented public API established yet

The investigation has **not yet established a documented/public JSON API** that can be relied upon for the complete board. Search results for `api.masjidboardlive.com` did not provide a verifiable public API specification, and the current evidence is insufficient to claim that such an API exists.

This is important: we should not design the MasjidPi provider around an assumed API. The next technical investigation should inspect the live application's actual network requests and/or its JavaScript assets to determine whether a structured backend endpoint exists.

### Existing third-party integration

A third-party Home Assistant integration for MasjidBoard Live exists and is useful as corroborating evidence. It retrieves prayer information from the online board, but its scope is limited to the five daily prayers and therefore does **not** satisfy MasjidPi's full-board requirement.

This means it is useful as a reference for programmatic access, but not as the implementation basis for MasjidPi's complete MasjidBoard integration.

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
- What refresh intervals are expected?
- What data should be cached permanently versus temporarily?
- Can the current board be reproduced from structured data without rendering the website?
- What terms/usage restrictions apply to programmatic consumption of MasjidBoard Live data?

## Next Investigation Step

Inspect the live MasjidBoard application's JavaScript/network behaviour directly. Use several representative boards and identify the actual requests made for:

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
- Established that API/data investigation should precede production implementation.
- Investigated multiple current Premium boards and confirmed substantial variation and a much broader content set than the daily prayer times.
- Confirmed the masjid-specific `mid` identifier used by the Premium board URLs.
- Did not find sufficient evidence to claim a documented public JSON API; direct network/JavaScript inspection is the next step.
