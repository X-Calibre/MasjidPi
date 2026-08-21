# MasjidPi Roadmap

MasjidPi is a lightweight home appliance that keeps people connected to their masjid through live audio and masjid information.

The roadmap focuses on the current product, remaining production work, and planned future features.

---

## Current Status

**v1.1.0 is the current release.**

The current release provides:

- LiveMasjid stream catalogue and updates
- Weekly automatic catalogue refresh plus manual refresh from the Web UI
- Masjid search and discovery
- Favourites
- Web UI
- Stream playback via MPV
- Volume and audio-device control
- Persistent settings
- Playback, network and MPV recovery
- Audio-device recovery
- Raspberry Pi runtime and systemd integration
- Pre-built release packages
- Bundled installer
- Upgrade installation preserving configuration and runtime data
- Safe update workflow
- Safer update and recovery behaviour
- Validation of a newly installed version before considering an update successful
- Recovery/rollback path if an update fails
- First-run and production installation workflow improvements
- Application-level SD-card write optimisations, including in-memory state caching, reduced logging, atomic state writes, and in-memory catalogue processing
- Security dependency update for the Go network/HTML parsing stack

The v1.0.8 release was tested successfully on a Raspberry Pi 3 B and established the current stable appliance baseline. v1.0.10 added application-level SD-card write and filesystem-activity optimisations, and v1.0.11 adds the dependency security fix while retaining that baseline.

The Updates & Recovery work and application-level SD-card write optimisation are now considered complete. Remaining reliability work should prioritise measured OS-level appliance behaviour and simplicity rather than unnecessary configuration options.

MasjidBoard is being developed on `research/masjidboard-live`. The research implementation has progressed beyond the original browser-demo stage into a working Raspberry Pi appliance implementation and is now in final pre-integration validation.

---

## Remaining Work

There are no outstanding application-level Updates & Recovery or SD-card write optimisation items at this stage.

OS-level SD-card behaviour should be measured on Raspberry Pi hardware before deciding whether appliance-specific changes to journald, swap, or other system services are justified.

Current MasjidBoard work is focused on final production validation and integration rather than additional core functionality.

---

## Architecture

### MasjidPi Core / Listen / Board

MasjidPi is evolving from the original audio-focused application into a modular home appliance built around shared MasjidPi functionality and two independent capabilities: Listen and Board.

The repository remains a single MasjidPi repository. The architecture is conceptually:

**MasjidPi Core → Listen**

and independently:

**MasjidPi Core → Board**

#### MasjidPi Core

Core contains functionality shared by Listen and Board, including configuration, persistent state, networking, common API behaviour and Raspberry Pi/platform integration where appropriate.

#### MasjidPi Listen

Listen contains the audio-specific functionality:

- Stream playback
- MPV integration
- Volume control
- Audio-device handling
- Playback recovery and reconnection
- Audio-specific UI and controls

#### MasjidPi Board

Board contains the MasjidBoard-specific functionality:

- Prayer and Jumu'ah times
- Next-event and countdown information
- MasjidBoard data retrieval and caching
- MasjidBoard configuration UI
- HDMI display UI
- Appliance display runtime

Listen and Board remain independent capabilities. Neither subsystem requires the other to operate, and audio playback must continue operating if MasjidBoard or its upstream provider is unavailable.

The installer supports three appliance profiles:

- Listen
- Board
- Listen + Board

Only the dependencies, backend subsystems, APIs, configuration pages and appliance services required by the selected profile should be active.

---

## Planned Features

### Raspberry Pi Appliance Image

- Deployable Raspberry Pi appliance image
- Preconfigured operating environment
- Repeatable image-based installation workflow

### Linux x86-64 Appliance / Old Laptop Support

MasjidPi should eventually support repurposing an older x86-64 laptop or PC as a dedicated MasjidPi appliance.

The preferred approach is to use a lightweight Linux operating system rather than requiring Windows. This would allow the machine to boot directly into a simple, appliance-like MasjidPi environment without requiring users to manage a normal desktop operating system.

Potential capabilities:

- Support for standard Linux x86-64 hardware
- Dedicated MasjidPi appliance installation/image
- Automatic MasjidPi startup at boot
- Audio output through built-in, USB, Bluetooth, or other supported Linux audio hardware
- Optional HDMI display output for MasjidBoard
- Ability to run MasjidPi Listen and MasjidPi Board together on the same machine
- Reuse of the same MasjidPi Core, Listen, and Board architecture as Raspberry Pi deployments
- Simple installation suitable for repurposing older laptops and PCs

This should be treated as a platform/appliance target rather than a separate MasjidPi product. The application should remain shared across Raspberry Pi and x86-64 Linux, with platform-specific behaviour isolated where necessary.

### Windows x64 Desktop Support

Potential future support for running MasjidPi directly on Windows x64 as a desktop application.

This is primarily intended as a convenience, development, testing, and desktop-listening platform rather than the preferred dedicated appliance platform. A Linux-based x86-64 appliance remains the preferred approach for repurposing an old laptop as a dedicated MasjidPi device.

Potential capabilities:

- Windows x64 build
- Bundled or managed MPV runtime
- Windows-compatible playback IPC and filesystem handling
- Simple desktop installation
- Local Web UI access

### Hardware & Advanced Audio

- OLED display
- Push-button controls
- Hardware playback controls
- Display current station and playback state
- Audio equaliser / EQ processing
- Advanced audio controls

### MasjidBoard

MasjidBoard is a separate subsystem from audio streaming and playback. The active `research/masjidboard-live` implementation now provides the main appliance functionality.

Implemented and validated capabilities include:

- Location hierarchy and scoped masjid discovery
- Selection and ordering of up to three masjids
- Persisted location scope and selection
- Live prayer and Jumu'ah timetable retrieval from MasjidBoard Live
- Normalized provider data and resilient Jumu'ah handling
- Dedicated MasjidBoard configuration Web UI
- Read-only HDMI display page
- Responsive one-, two- and three-board comparison layout
- Friday Jumu'ah replacement of Dhuhr
- Chronological Jumu'ah event presentation
- Per-board next-event countdown, including overnight rollover to the next day's Fajr
- 30-minute automatic timetable refresh and manual refresh
- Per-board last-known-good cache fallback during provider outages
- Stale/current recovery after provider outages
- Reduced cache writes when timetable data is unchanged
- Independence from Listen when MasjidBoard Live is unavailable
- Raspberry Pi OS Lite appliance display using Cog/WPE directly on DRM/KMS
- Automatic systemd-managed HDMI display startup and recovery
- Component-aware installation of Listen, Board or Listen + Board
- Component-aware runtime dependencies, APIs and configuration UI
- Transactional component-profile changes with update validation and rollback protection
- Successful Raspberry Pi appliance validation of Listen-only, Board-only and combined Listen + Board profiles

The current display is the **default layout**, not the only long-term layout. Alternative user-selectable layouts should consume the same normalized MasjidBoard display API so presentation remains separate from provider, caching and timetable logic.

Remaining pre-integration work is intentionally narrow:

- Verify production ownership, permissions and preservation of MasjidBoard state across fresh installs and upgrades
- Complete clean first-run validation with no existing MasjidBoard configuration or cache
- Complete final user-facing error-message cleanup while retaining useful diagnostics
- Run a final Listen/audio regression pass
- Perform longer-duration appliance and HDMI disconnect/reconnect testing where practical
- Complete final documentation/integration review and prepare the research branch for merge/release planning

Future MasjidBoard enhancements may include:

- Additional user-selectable display layouts
- Layouts exposing selected astronomical/calculation times
- OLED or other hardware displays
- Richer announcement/media content where useful
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
