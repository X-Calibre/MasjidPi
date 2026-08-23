# MasjidBoard Live Data Inventory

**Status:** Initial discovery document  
**Date:** 23 August 2026  
**Purpose:** Record information supplied by MasjidBoard Live that MasjidPi Board does not yet display, before deciding how it should appear in the portrait and detailed layouts.

## Important distinction

MasjidBoard Live advertises both a free **Core** masjid listing and a separate **Premium** on-site noticeboard product. A feature being present on a MasjidBoard Live webpage or Premium board does **not** automatically mean it is available through the Core data response currently consumed by MasjidPi.

This inventory therefore uses three verification states:

- **Confirmed in current MasjidPi data:** already observed in the response returned by MasjidPi's `/api/masjidboard/display` endpoint.
- **Confirmed MasjidBoard Live content:** publicly documented or visible on MasjidBoard Live, but its raw API field still needs to be captured and mapped.
- **Premium/availability uncertain:** visible in the Premium product; we must verify that Core-listed masaajid expose it to MasjidPi.

## Current MasjidPi baseline

The Board currently concentrates on:

- masjid identity;
- daily adhan and iqamah times;
- Jumu'ah times where supplied;
- astronomical/perpetual times such as suhur, sunrise, ishraaq, duha, istiwa, zawaal and sunset;
- the next prayer/event and countdown;
- upcoming Fajr, Asr and Esha time changes where available.

These fields are the baseline and are not the main subject of this document.

## High-value information not yet displayed

| Content type | Information potentially available | Availability status | Value for a home appliance | Notes for later verification |
| --- | --- | --- | --- | --- |
| Funeral notices | Funeral announcement, deceased person, funeral/janazah time and venue, burial information, related message or poster | Confirmed MasjidBoard Live content; API shape unknown | **Very high** | MasjidBoard Live supports synchronised suburb-based funeral notices, so a notice may originate from another masjid in the same community. We need source masjid, locality, activation/expiry and media/text fields. |
| Nikaah notices | Nikaah announcement and event details | Confirmed MasjidBoard Live content; Premium/Core exposure uncertain | **High** | MasjidBoard Live says the notice remains visible on the day of the nikaah. We need event date, names/title, venue, time and expiry behaviour. |
| General masjid announcements | Free-form notices, announcements and urgent changes | Confirmed as part of the Core personal webpage feature set; API shape unknown | **Very high** | Likely the broadest and most useful category. Determine whether entries are text, poster images, or both. |
| Community broadcasts | Notices shared between synchronised boards in the same suburb | Confirmed MasjidBoard Live content; API shape unknown | **High** | Must retain the originating masjid and avoid presenting a community notice as if it came directly from the selected masjid. |
| Da'wah activities | Programmes, lectures, classes and other activities | Confirmed MasjidBoard Live content; Premium/Core exposure uncertain | **High** | Determine whether these are structured events or ordinary notice/poster entries. |
| Custom scheduled notices | Two custom announcement slides with scheduling | Confirmed Premium feature; Core exposure uncertain | **Medium–high** | Important fields would be start/end dates, display schedule, title/body and attached poster. |
| Uploaded posters | Full-colour HD/A4 notice images uploaded by the masjid | Confirmed Premium feature; Core exposure uncertain | **High** | MasjidPi will need image URL, dimensions/type, scheduling and a safe cached fallback. Portrait screens may need contain/crop rules. |
| Eid salaah notice | Eid-ul-Fitr or Eid-ul-Adha prayer times and related details | Confirmed Premium feature; Core exposure uncertain | **High, seasonal** | May contain multiple jamaats, venue, date and special instructions. It should probably override ordinary low-priority content near Eid. |
| Temporary salaah changes | Short-term overrides for unexpected changes or Ramadan schedules | Confirmed MasjidBoard Live feature; data exposure uncertain | **Very high** | This is operational rather than promotional content. MasjidPi must distinguish an override from the perpetual timetable and show its validity period. |
| Jumu'ah khateeb | Khateeb name associated with Jumu'ah | Confirmed on Core personal webpages and visible in Premium layout | **Medium–high** | We already use Jumu'ah time data, but should verify whether khateeb names are present in our response and whether up to three jamaats each have separate details. |
| Islamic/Hijri date | Locally adjusted Hijri date, including Arabic and English/transliterated forms | Confirmed MasjidBoard Live content; current MasjidPi exposure still to verify | **High** | MasjidBoard Live allows manual adjustment according to local moon sighting, making its date preferable to a purely calculated date. |

## Additional information and behaviours worth investigating

These are available in MasjidBoard Live but are lower priority for the first notices implementation.

| Content or behaviour | What MasjidBoard Live provides | Suggested MasjidPi priority |
| --- | --- | --- |
| Daily ayah, hadith and sunnah | Rotating educational content | Later; useful but can compete with urgent community information |
| Masnoon du'a after adhan | Du'a displayed for five minutes from adhan | Later; event-triggered display rather than persistent data |
| Full-colour Islamic reminders | Date/time-sensitive reminders tied to important Islamic occasions | Later |
| New moon information | Birth, age at sunset, moonset, best viewing date, angle, elevation and direction | Later or specialist view |
| Current moon phase | Moon-phase graphic | Low priority |
| Cellphone reminder | Reminder shown before iqamah | Low priority for a device located in the home |
| Zawaal alert | Visually prominent warning during zawaal | Medium; the underlying time is already present, so MasjidPi can derive the state locally |
| Multiple Jumu'ah jamaats | Up to three Jumu'ah services | Medium–high where supplied |
| Navigation/location | Link or coordinates for directions to the masjid | Medium; useful for notices and events away from the normal venue |
| Live audio link | Masjid audio stream information | Already belongs to MasjidPi Listen; verify whether the Board source contains a useful cross-reference |

## Proposed priority for display planning

### Priority 1 — community-critical

1. Funeral notices
2. Urgent/general masjid announcements
3. Temporary salaah-time changes
4. Eid salaah notices

### Priority 2 — community events

1. Nikaah notices
2. Da'wah programmes and activities
3. Community broadcasts
4. Uploaded event posters

### Priority 3 — useful enrichment

1. Locally adjusted Islamic date
2. Jumu'ah khateeb and multiple-jamaat detail
3. Daily educational material
4. New-moon and moon-phase information

## Data questions to answer before UI work

For each notice-like item, capture and document:

1. **Type:** funeral, nikaah, Eid, activity, general announcement, poster or another category.
2. **Identity:** stable notice ID for deduplication and caching.
3. **Source:** selected masjid or another synchronised community masjid.
4. **Content:** title, body, names, venue, contact information and any structured times.
5. **Media:** image/poster URL, MIME type, dimensions and alternative text.
6. **Timing:** event date/time, publication time, display-from time and expiry time.
7. **Priority:** whether MasjidBoard Live supplies an urgency or ordering value.
8. **State:** active, upcoming, expired, cancelled or replaced.
9. **Links:** directions, registration, livestream or other action URL.
10. **Localisation:** language, Arabic content and any HTML/formatting in the body.

We also need to test:

- whether empty notices are omitted, `null`, blank strings or empty arrays;
- whether the API returns HTML rather than plain text;
- whether poster URLs require authentication or expire;
- whether community notices are duplicated across selected masaajid;
- how quickly additions, changes and removals propagate;
- whether old/expired notices remain in the payload;
- whether notices are part of the Core endpoint or require a separate Premium endpoint;
- how the current MasjidPi cache behaves when a notice is withdrawn while the device is offline.

## Recommended API sampling exercise

Before selecting a layout, record raw responses from at least four cases:

1. a masjid with an active funeral notice;
2. a masjid with an active nikaah or community-event notice;
3. a masjid using an uploaded poster;
4. a normal masjid with no active notices.

Compare the complete JSON trees, not only known prayer-time fields. This will reveal the actual field names, nesting, nullability, notice lifecycle and whether a second endpoint is needed. The resulting field map should be added to this document before implementation begins.

## Preliminary product conclusion

The strongest expansion for MasjidPi Board is not more prayer-time detail; it is **time-sensitive community information**. Funeral notices, urgent announcements, temporary salaah changes and Eid notices are the best fit for a home appliance whose purpose is to keep a household connected to its masjid.

The interface should eventually treat these as prioritised, expiring content rather than permanently adding more small text to the existing prayer board. The exact display mechanism—banner, card, ticker, modal takeover or carousel—should be decided only after the live payload structure and typical content length are known.

## Sources

- [MasjidBoard Live overview and package details](https://masjidboardlive.com/)
- [MasjidBoard Live features, page 1](https://masjidboardlive.com/faq-items/)
- [MasjidBoard Live features, page 2](https://masjidboardlive.com/faq-items/page/2/)
- [MasjidBoard Live features, page 3](https://masjidboardlive.com/faq-items/page/3/)
- [MasjidBoard Live Premium example](https://premium.masjidboardlive.com/v2/?mid=actonville-jaame)
