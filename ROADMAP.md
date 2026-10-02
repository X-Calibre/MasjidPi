# MasjidFrame Roadmap

MasjidFrame is a lightweight appliance for live masjid audio and prayer-time information.

## Current release

**v1.6.3 is the current stable release.**

MasjidFrame currently provides:

- independent Listen, Board and combined appliance profiles;
- touchscreen first-run Wi-Fi, location and primary-masjid setup for the portrait appliance;
- priority masjid audio with optional secondary Islamic Radio;
- saved favourites, scheduling, source volumes and audio-output recovery;
- discovery and selection of up to three MasjidBoard Live masjids;
- responsive TV/Monitor presentation and a dedicated 720 × 1280 Raspberry Pi Touch Display 2 profile;
- prayer, Jumu'ah, Daily Times and next-event information;
- supported community notices and optional shared Islamic content;
- last-known-good data during temporary upstream failures; and
- release packages for Linux ARM64 and AMD64; and
- a signed Pi 3 A/B appliance image with automatic health confirmation and rollback.

Completed release details belong in [GitHub Releases](https://github.com/X-Calibre/MasjidPi/releases) and the relevant acceptance records under `docs/`.

## Current priorities

### Reliability and validation

- Continue longer-duration Raspberry Pi 3B and Pi 4 soak monitoring.
- Test HDMI disconnect and reconnect behavior.
- Validate suitable 512 MB ARM64 devices, particularly Pi Zero 2 W and Pi 3A+.
- Measure OS-level SD-card writes before changing journald, swap or other system services.
- Keep provider parsing defensive as MasjidBoard Live payloads evolve.

### Appliance product work

- Observe a changed shared daily item during a scheduled refresh on the appliance; automated tests cover repeated same-day fetches and changed-content persistence.
- Retain the accepted v1.6.0 image, update, rollback, soak and interrupted-update recovery evidence.
- Finalise the portrait enclosure, display, audio and power design.
- Validate the selected USB audio path, speakers and physical enclosure.
- Investigate HDMI-CEC behavior on intended displays.

### MasjidBoard

- Add poster/media support only when retrieval, caching and presentation are safe.
- Improve Arabic/RTL coverage using additional real-world content.
- Revisit Maghrib/Iftar semantics when representative Ramadan source data is available.
- Extend structured upstream content only when field meaning and privacy are established.
- Add display layouts or preferences only where they materially improve appliance use.

### Listen

- Continue validating Radio endpoints and fallback behavior.
- Consider richer current-programme metadata where stations expose reliable information.
- Explore advanced audio controls only if they remain simple on appliance hardware.

### SmartBilal integration (deferred)

Investigate SmartBilal as an additional Listen upstream in a future release; this is outside the v1.6.8 UI work.

Findings from the 2 October 2026 investigation:

- Public masjid profile pages expose structured masjid IDs, names, regions, stream URLs and a live flag. Example: [Nur ul Islam Masjid Lenasia](https://media.smartbilal.com/masjid/nurulislam).
- A sampled active stream decoded as Opus audio in Ogg, compatible with mpv.
- The observed [Icecast status endpoint](http://41.185.71.90:8000/status-json.xsl) lists active mounts. Offline streams are absent, and active entries have generic names, so it cannot supply a complete, named masjid catalogue.
- Initial website inspection found no complete catalogue API or sitemap. Subsequent static inspection of the user-supplied SmartBilal 2.0.8 Android APK identified working catalogue APIs; the website-only finding is superseded by the results below.

Android APK follow-up on 2 October 2026:

- The production API base is `https://media.smartbilal.com/api/`. Catalogue and detail requests tested without credentials returned data; the app also supports bearer authentication for other requests.
- `GET /api/v3/dashboard` exposes top-level categories. `GET /api/v2/dashboardList?page=1&type=collectionGroup&category_id=9` returned South Africa's regional groups, including offline masjids.
- The South Africa response contained 350 entries across 17 groups, representing 332 unique provider IDs. These are catalogue entries, not 332 independently verified playable streams; duplicate category membership must be reconciled.
- The grouped South Africa response was identical for `page=1` and `page=2`. Do not assume the page parameter advances grouped results; validate endpoint-specific pagination and traversal before claiming completeness.
- `GET /api/users/1471` returned Nur ul Islam's stream metadata, including mount `m15`, its channel URL and `is_live: false`.
- Tests of the list/filter endpoints used the app's `x-api-version: 36` header. An unqualified list request returned an empty result, so reproduce the appropriate request type and category rather than treating that as an empty catalogue.
- The APK references MQTT and stream-start/stream-stop events. This is static evidence only; broker access, subscriptions, credentials and event behavior have not been validated.
- These internal app APIs provide a promising enumeration path, but worldwide completeness, long-term stability and supported third-party use remain unconfirmed.

Before implementation:

- Establish a supported catalogue API or maintained export with stable IDs, names, locations and stream URLs, including offline masjids. Confirm pagination, completeness and refresh behavior.
- Evaluate curated public profile links as an initial alternative, clearly distinguishing a maintained subset from automatic complete discovery.
- Add provider-specific availability monitoring; SmartBilal cannot rely on LiveMasjid's MQTT status feed. Validate polling cadence, status freshness and broadcast start/stop detection.
- Validate playback, reconnects, network failures and Radio interruption/resumption on Raspberry Pi hardware.
- Keep provider identities distinct and reconcile masjids available through multiple upstreams without duplicate favourites or ambiguous status.
- Cache last-known-good catalogue data and retain it during upstream failures.

## Platform targets

### Linux x86-64 appliance

Official AMD64 packages already support compatible 64-bit Linux systems. Further work may make repurposed laptops and small PCs easier to use as dedicated appliances through:

- automated display and audio setup;
- boot-to-appliance behavior;
- broader built-in, USB and Bluetooth audio validation; and
- a simplified installation profile for old computers.

### Windows x64

Windows desktop support remains exploratory. It would require Windows-compatible mpv management, IPC, persistent paths, service/startup behavior and packaging. Linux remains the preferred appliance platform.

### Hardware controls

Possible future hardware integrations include:

- OLED status displays;
- physical playback and volume controls;
- amplifier/EQ integration; and
- enclosure-specific status indicators.

### Home Assistant

Potential integration could expose:

- media-player controls and playback state;
- selected source and masjid status;
- volume and audio-output information;
- failure/recovery events;
- prayer-time sensors and triggers; and
- MQTT or a native Home Assistant integration.

## Architecture guardrails

MasjidFrame remains one repository with shared Core functionality and two independently operable capabilities:

- **Core** — configuration, persistent state, APIs and platform integration
- **Listen** — stream discovery, priority playback, Radio, mpv and audio devices
- **Board** — MasjidBoard retrieval, caching, configuration and HDMI presentation

Listen must continue operating when Board or its upstream providers are unavailable, and Board must not depend on Listen.

## Project principles

MasjidFrame should remain:

- simple to install and operate;
- reliable through network and device interruptions;
- lightweight enough for supported Raspberry Pi hardware;
- usable without a separate server;
- conservative with persistent writes;
- resilient when optional services fail; and
- focused on a dependable home-appliance experience.
