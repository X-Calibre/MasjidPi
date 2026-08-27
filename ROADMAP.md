# MasjidPi Roadmap

MasjidPi is a lightweight home appliance that keeps people connected to their masjid through live audio and masjid information.

The roadmap focuses on the current product, remaining reliability work, and planned future features.

---

## Current Status

**v1.4.0 is the current stable release. v1.4.1 is being prepared.**

The current product provides:

- LiveMasjid stream catalogue, search, favourites and playback
- Automatic catalogue refresh every 28 days plus manual refresh from the Web UI
- MPV playback with volume and audio-device control
- Playback, network, MPV and audio-device recovery
- Persistent settings and playback state
- MasjidBoard Live masjid discovery and selection
- Prayer and Jumu'ah timetable display for up to three selected masjids
- Next-event countdowns and timetable refresh
- Last-known-good MasjidBoard cache fallback
- Dedicated HDMI Board appliance runtime using Cog/WPE on DRM/KMS
- Selectable Listen, Board, and Listen + Board appliance profiles
- Raspberry Pi runtime and systemd integration
- Pre-built ARM64 and AMD64 release packages
- Bundled installer and public one-line installation
- Upgrade installation preserving configuration and runtime data
- Safe update workflow with validation and rollback
- Application-level SD-card write optimisations

### v1.2.0 additions

v1.2.0 extends the MasjidBoard appliance experience with:

- User-selectable **Standard** and **Detailed** HDMI layouts
- Detailed layout with shared Adhan/Jamaah headings and a Daily Times panel
- Full Gregorian date display
- MasjidBoard-derived Islamic date calculation, including sunset rollover and upstream date-adjustment/forced-month-end behaviour
- Islamic weekday transliteration on the Detailed display
- 1080p layout refinements validated on the Raspberry Pi 3B reference appliance
- Six curated Board colour themes: **Emerald, Midnight, Slate, Ruby, Light, and Black & White**
- Persisted HDMI layout and theme preferences
- Live theme updates and automatic Standard/Detailed switching from the Web UI without restarting the display service
- MPV IPC response-ordering reliability improvements with concurrent-response tests

The published v1.2.0 ARM64 release was validated on the Raspberry Pi 3B reference appliance with a native 1080p TV. The public `install-latest.sh` upgrade path preserved the Listen + Board profile, selected audio device, selected masjids, saved Detailed layout and saved colour theme. Listen playback and the HDMI MasjidBoard display were verified after installation.

### v1.2.1 maintenance release

v1.2.1 improves maintainability and upgrade persistence without changing the Core / Listen / Board architecture:

- Web UI preferences now use persistent runtime storage under `/var/lib/masjidpi`
- Existing preferences are migrated automatically from the legacy installation path
- API construction and application component startup are simpler and more explicit
- Playback runtime state and retry handling are easier to maintain
- Atomic JSON file replacement is shared across Listen and MasjidBoard stores
- MasjidBoard configuration actions and messages are owned directly by the main page controller
- CI now checks Go formatting, vet, race-enabled tests, ShellCheck and JavaScript syntax

The v1.2.1 source upgrade was validated on the Raspberry Pi 3B reference appliance. Preference migration, Listen and Board service startup, saved-state recovery after reboot, the MasjidBoard configuration page and the HDMI display all passed appliance testing.

### v1.3.0 additions

v1.3.0 adds:

- Landscape (1920 × 1080) and Portrait (600 × 1024) Board orientations
- Landscape timetable and full-width Daily Times presentation refinements
- Optional Premium community-content enrichment while retaining Core as the authoritative timetable source
- Adaptive, rotating and source-labelled cards for announcements, Nikah, funerals, Eid, Salaah changes, well-wishes, Taleem, Dawah/Gasht, three-day Jamaat, contributions and calculated new-moon information
- Duplicate suppression, safe plain-text conversion and anonymised development fixtures for community cards
- Improved first-run guidance when Board has not yet been configured
- Cog GLES rendering for reliable Raspberry Pi 4 DRM/KMS output
- Raspberry Pi 4 validation covering installation, reinstallation, reboot recovery, saved preferences, simultaneous Listen + Board operation and 1080p HDMI output

Automated CI and the Raspberry Pi 4 acceptance pass completed before release. The published artifacts are built automatically from the signed release tag and include checksum verification.

### v1.4.1 maintenance release

v1.4.1 adds:

- LiveMasjid relay-to-Icecast playback failover using matching cached mount URLs
- A two-second readiness delay after a new MQTT mount-start event
- Active-endpoint reporting in Listen status and structured playback logs
- Temporary preference for an endpoint that has played successfully
- A 28-day automatic Listen catalogue refresh interval with the manual update option retained
- Upgrade migration from the previous seven-day stock refresh interval while preserving custom values
- Adaptive Landscape Board proportions for one, two and three selected masjids
- Bundled Fira Sans display typography and lighter large-time weights for improved television legibility
- Dedicated, centred Salaah-time-change cards with complete formatted dates
- Clearer standalone Maghrib Adhan presentation
- Complete Islamic Economic Indicator fields with standard South African Rand and accounting presentation
- Separate effective and retrieval dates for economic data
- Effective-date-driven economic refreshes from 09:00 South African time with bounded weekend attempts

The updated Listen catalogue, manual refresh, normal relay playback and automatic Icecast failover were validated on the Raspberry Pi 4 test appliance. Automated Go tests, race-enabled playback tests, vet checks, shell syntax checks and frontend JavaScript validation also passed.

---

## Remaining Work

There are no known release-blocking items for v1.4.1.

Ongoing reliability work:

- Longer-duration appliance soak testing
- HDMI disconnect/reconnect testing
- Broader Raspberry Pi hardware validation
- Measure OS-level SD-card behaviour before deciding whether changes to journald, swap or other system services are justified

---

## Architecture

### MasjidPi Core / Listen / Board

MasjidPi remains a single repository with shared Core functionality and two independent capabilities:

**MasjidPi Core → Listen**

and independently:

**MasjidPi Core → Board**

#### MasjidPi Core

Core contains shared configuration, persistent state, networking, API behaviour and platform integration.

#### MasjidPi Listen

Listen contains:

- Stream playback
- MPV integration
- Volume control
- Audio-device handling
- Playback recovery and reconnection
- Audio-specific UI and controls

#### MasjidPi Board

Board contains:

- Prayer and Jumu'ah times
- Islamic/Gregorian date presentation
- Next-event and countdown information
- MasjidBoard data retrieval and caching
- MasjidBoard configuration UI
- HDMI display layouts and themes
- Appliance display runtime

Listen and Board remain independent capabilities. Neither subsystem requires the other to operate, and audio playback continues operating if MasjidBoard or its upstream provider is unavailable.

The installer supports three appliance profiles:

- Listen
- Board
- Listen + Board

Only the dependencies, backend subsystems, APIs, configuration pages and appliance services required by the selected profile are active.

---

## Planned Features

### Raspberry Pi Appliance Image

- Deployable Raspberry Pi appliance image
- Preconfigured operating environment
- Repeatable image-based installation workflow

### Linux x86-64 Appliance / Old Laptop Support

MasjidPi should eventually support repurposing an older x86-64 laptop or PC as a dedicated appliance using a lightweight Linux operating system.

Potential capabilities:

- Standard Linux x86-64 hardware support
- Dedicated MasjidPi appliance installation/image
- Automatic MasjidPi startup at boot
- Audio through built-in, USB, Bluetooth or other supported Linux hardware
- Optional HDMI MasjidBoard output
- Listen and Board together on the same machine
- Simple installation suitable for repurposing older laptops and PCs

This should remain a platform target for the same MasjidPi application rather than a separate product.

### Windows x64 Desktop Support

Potential future support for running MasjidPi directly on Windows x64 as a desktop application.

Potential capabilities:

- Windows x64 build
- Bundled or managed MPV runtime
- Windows-compatible playback IPC and filesystem handling
- Simple desktop installation
- Local Web UI access

A Linux-based x86-64 appliance remains the preferred approach for repurposing older hardware as a dedicated MasjidPi device.

### Hardware & Advanced Audio

- OLED display
- Push-button controls
- Hardware playback controls
- Display current station and playback state
- Audio equaliser / EQ processing
- Advanced audio controls

### MasjidBoard

Implemented and validated MasjidBoard capabilities now include:

- Location hierarchy and scoped masjid discovery
- Selection and ordering of up to three masjids
- Persisted location scope and selection
- Live prayer and Jumu'ah timetable retrieval from MasjidBoard Live
- Normalized provider data and resilient Jumu'ah handling
- Dedicated MasjidBoard configuration Web UI
- Landscape and Portrait HDMI layouts
- Responsive one-, two- and three-board comparison layouts
- Friday Jumu'ah replacement of Dhuhr
- Chronological Jumu'ah event presentation
- Per-board next-event countdown, including overnight rollover to the next day's Fajr
- Detailed Daily Times sourced consistently from Masjid 1
- Gregorian and masjid-adjusted Islamic date display
- Six curated Board colour themes
- Live HDMI preference updates from the Web UI
- Optional Premium community-content enrichment with Core timetable fallback
- Adaptive rotating cards for supported announcements, programmes, notices, contributions and new-moon information
- Automatic timetable refresh and manual refresh
- Per-board last-known-good cache fallback during provider outages
- Stale/current recovery after provider outages
- Reduced cache writes when timetable data is unchanged
- Independence from Listen when MasjidBoard Live is unavailable
- Raspberry Pi OS Lite appliance display using Cog/WPE directly on DRM/KMS
- Automatic systemd-managed HDMI display startup and recovery
- Component-aware installation of Listen, Board or Listen + Board
- Transactional component-profile changes with update validation and rollback protection

Future MasjidBoard enhancements may include:

- Additional display layouts or resolution-specific choices where they provide clear value
- OLED or other hardware displays
- Poster/media support where it can be retrieved and displayed safely
- Additional display preferences
- Broader appliance/platform validation

### Home Assistant Integration

Potential capabilities:

- Home Assistant media-player entity
- Play / stop control
- Volume control
- Current stream and masjid information
- Playback state
- Audio output device
- Stream and playback status
- Playback failure and recovery events
- Audio-device loss and recovery events
- MQTT integration
- Optional native Home Assistant integration
- Optional prayer-time and masjid sensors
- Prayer-related events and triggers

---

## Project Principles

MasjidPi should remain:

- A simple home appliance for people who want to stay connected to their masjid
- Lightweight enough for supported Raspberry Pi hardware
- Simple to install and operate
- Reliable during network interruptions
- Usable without a separate server
- Focused on bringing live masjid audio and useful masjid information into the home
- Easy to maintain and update
- Independent of browser-local configuration
- Modular, with Listen and Board capabilities separated from shared functionality
- Resilient when optional external services are unavailable
- Portable across appropriate hardware platforms without unnecessarily duplicating the application
- Focused on appliance simplicity rather than exposing unnecessary configuration

The project should prioritise a reliable home appliance experience over unnecessary complexity.
