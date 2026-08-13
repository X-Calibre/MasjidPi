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

## Initial Findings

- The public board uses a masjid-specific `mid` value in its URL.
- Live boards contain loading placeholders and subsequently populate board content, indicating that data is loaded dynamically.
- Live boards expose separate areas for prayer information, announcements, daily Ayah, Hadith, community broadcasts, Du'a, images, and New Moon information.
- Initial inspection has also shown Jumu'ah/Khateeb information, Nikah notices, weekly programmes, masjid information, and contribution information.

These observations must be verified against the underlying data before becoming implementation requirements.

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

## Open Questions

- What exact public/client API endpoints does MasjidBoard Live use?
- What authentication or tokens are required?
- What is the complete data schema?
- How are content schedules and expiry represented?
- How are images and posters delivered?
- How are prayer-time changes represented?
- How is time zone information represented?
- How are multiple Jumu'ah services represented?
- What refresh intervals are expected?
- What data should be cached permanently versus temporarily?
- Which display resolutions should be officially supported?
- How should content rotate on the HDMI display?
- Can the current board be reproduced from structured data without rendering the website?

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
