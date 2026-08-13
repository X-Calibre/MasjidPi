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

The important discovery is that **a structured JSON endpoint exists**. The response from `mblapi` was `application/json` and contained 29 top-level arrays. The same data was also embedded directly in the HTML page as the JavaScript variable `theInfo`.

The page also embeds:

```javascript
let boardId = "1asEQ0Ju83TPqBFHw7NbBAihAxMt5JQ2bJkbaWnwKf7k";
let mblVersion = "14";
```

The `boardId` value is the identifier used by `mblapi` in this capture. This is distinct from the public masjid URL identifier (`mid=erasmia-aaisha`).

The HAR confirmed that the `mblapi` response and the embedded `theInfo` value are identical for this board. This means MasjidPi can potentially consume the structured endpoint without scraping the rendered board page.

### `mblapi` response structure

The response is currently **not a self-describing JSON object**. It is a positional array-of-arrays. The first-level array contains 29 rows, with each row containing positional fields whose meaning is defined by the MasjidBoard Live frontend.

Examples observed in the capture include:

- Row 0: default/empty values.
- Row 1: Jumu'ah-related values including `Adhān`, `Sunan`, `Khutbah`, and associated times.
- Row 2: date, astronomical/moon-related values, masjid ID, city, timezone, language, and Islamic-time display settings.
- Row 3: daily Salah Adhan/Jamaah values for Fajr, Zuhr, Asr, Maghrib and Esha, plus additional settings.
- Row 6: English/Arabic masjid naming and Islamic calendar/settings values.
- Rows 7–9: Ayah, Hadith and Sunnah configuration/content indicators.
- Row 12: banking information and multiple configurable announcement/programme headings and HTML content.
- Row 14: Nikah-related information.
- Row 15: funeral-related information.
- Row 16: other programme/activity information.
- Row 17: poster/image identifiers and visibility flags.
- Row 18: an additional event/notice block.
- Row 20: contribution information.
- Row 21: configurable name/heading fields.
- Rows 22–23: additional daily time data and Arabic masjid names.
- Row 24: poster/image identifiers and visibility flags.
- Row 25: another image identifier/visibility block.
- Row 27: carousel/translation/display configuration.
- Row 28: translation setting names.

These row mappings are **observed from one live board and are not yet a stable MasjidPi schema**. The frontend's `handleResults(theInfo)` function and related code need to be studied to map every field reliably.

### Images/media

The board requests images through:

```text
https://api.masjidboardlive.com/imageproxy?id=<image-id>
```

The HAR contains successful PNG and JPEG responses from this endpoint. At least some board content therefore references media by opaque image IDs rather than direct public filenames.

### Static/client code

The page loads:

```text
https://api.masjidboardlive.com/mblfileapi
```

as a JavaScript `<script>` resource. The captured response is a large JavaScript resource containing client-side translations/configuration. The main page then calls:

```javascript
handleResults(theInfo);
startMBL();
```

and loads `functions_uo_latest.js?109`.

The captured `functions_uo_latest.js` response was served from cache in the HAR, so its body was not available in the capture. Its source still needs to be obtained to map the positional `theInfo` array completely.

### Important architectural conclusion

The previous assumption that MasjidBoard Live might require HTML scraping is no longer correct for this board. A structured `mblapi` endpoint is demonstrably used by the live application.

However, the endpoint's response schema is an opaque positional array rather than a clean public API model. We should therefore **not expose this structure directly to the MasjidPi application**. The eventual provider should translate it into a normalised MasjidPi data model.

## Initial Findings

- The public board uses a masjid-specific `mid` value in its URL.
- Live boards contain loading placeholders and subsequently populate board content, indicating that data is loaded dynamically.
- Live boards expose separate areas for prayer information, announcements, daily Ayah, Hadith, community broadcasts, Du'a, images, and New Moon information.
- Initial inspection has also shown Jumu'ah/Khateeb information, Nikah notices, weekly programmes, masjid information, and contribution information.

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
| Consume structured `mblapi` data rather than scrape rendered HTML | Proposed | HAR evidence shows a structured endpoint is used by the live application. |
| Normalise MasjidBoard Live's positional schema inside the provider | Proposed | Keeps the opaque upstream schema out of the rest of MasjidPi. |

## Open Questions

- How is the public `mid` mapped to the internal `boardId`?
- Is `boardId` stable for a masjid?
- Can `mblapi?id=<boardId>` be requested directly without first loading the Premium page?
- What authentication, rate limits, or access restrictions apply to `mblapi`?
- What exact versioning guarantees exist for `mblVersion` and the positional schema?
- What exact field mapping is implemented by `handleResults()` / `functions_uo_latest.js`?
- How are Ayah, Hadith, Sunnah and Du'a content populated?
- How are announcement schedules and expiry represented?
- How are Jumu'ah services represented when there are multiple services?
- How are posters and images associated with individual content items?
- What is the exact meaning of every row and field in `theInfo`?
- How frequently does the live application refresh `mblapi`?
- Does the board make additional requests after the initial load that were not captured in this HAR?
- Which data is generated/calculated by the client rather than supplied by `mblapi`?
- Can the current board be reproduced completely from `mblapi` plus the referenced media assets?

## Next Investigation Steps

1. Obtain and inspect the actual `functions_uo_latest.js` source used by the board.
2. Map every `theInfo` row/field to its semantic meaning using the frontend code.
3. Test `mblapi` against several different masjids/board IDs to establish which fields are optional and how the schema varies.
4. Determine how a public `mid` is resolved to a `boardId`.
5. Test direct access to `mblapi` and `imageproxy`, including behaviour without browser session state.
6. Determine refresh/polling behaviour.
7. Build a draft normalised MasjidPi data model from the verified field mapping.
8. Only then begin production MasjidBoard module implementation.

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
- Established that the next priority is mapping the positional schema through the frontend JavaScript.
