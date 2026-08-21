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
- MasjidBoard Live masjid discovery and selection
- Prayer and Jumu'ah timetable display
- Up to three selected MasjidBoard displays
- Next-event countdowns and timetable refresh
- Last-known-good MasjidBoard cache fallback
- Dedicated HDMI Board appliance runtime using Cog/WPE on DRM/KMS
- Selectable Listen, Board, and Listen + Board appliance profiles
- Raspberry Pi runtime and systemd integration
- Pre-built release packages
- Bundled installer
- Upgrade installation preserving configuration and runtime data
- Safe update workflow with validation and rollback
- Application-level SD-card write optimisations
- Security dependency updates

v1.1.0 has been production-validated on a Raspberry Pi 3B running 64-bit Raspberry Pi OS Lite. The public one-line installer was tested from a clean OS installation using the Listen + Board profile. Installation, release checksum verification, both systemd services, Listen audio playback through the selected output device, and the HDMI MasjidBoard display were all verified successfully. Listen-only, Board-only, and profile-transition behaviour were also validated during pre-release appliance testing.

The Updates & Recovery work, application-level SD-card write optimisation, and initial MasjidBoard appliance integration are considered complete for v1.1.0. Further reliability work should be driven by measured appliance behaviour rather than unnecessary configuration.

---

## Remaining Work

There are no outstanding release-blocking items for v1.1.0.

OS-level SD-card behaviour should be measured on Raspberry Pi hardware before deciding whether appliance-specific changes to journald, swap, or other system services are justified.

Longer-duration appliance testing, HDMI disconnect/reconnect behaviour, and broader hardware validation can continue as ongoing reliability work without blocking the current release.

---

## Architecture

### MasjidPi Core / Listen / Board

MasjidPi is a modular home appliance built around shared MasjidPi functionality and two independent capabilities: Listen and Board.

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

MasjidBoard is a separate subsystem from audio streaming and playback and is included in v1.1.0.

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
- Successful clean production installation of v1.1.0 using the public release installer
- Verified HDMI display output and Listen audio playback on the production installation

The current display is the **default layout**, not the only long-term layout. Alternative user-selectable layouts should consume the same normalized MasjidBoard display API so presentation remains separate from provider, caching and timetable logic.

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
