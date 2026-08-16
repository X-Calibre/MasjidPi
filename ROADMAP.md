# MasjidPi Roadmap

MasjidPi is a lightweight Raspberry Pi home appliance that keeps people connected to their masjid through live audio and masjid information.

The roadmap focuses on the current product, remaining production work, and planned future features.

---

## Current Status

**v1.0.8 is released and verified on Raspberry Pi hardware.**

The current release provides:

- LiveMasjid stream catalogue and updates
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

The current release is considered functionally complete and reliable enough for normal use.

---

## Remaining Work

### Updates & Recovery

- Safe update workflow — Complete
- Safer update and recovery behaviour — Complete
- Validate a newly installed version before considering an update successful — Complete
- Provide a recovery/rollback path if an update fails — Complete
- Improve the first-run and production installation workflow where required — Complete

The remaining work should prioritise reliability and simplicity rather than unnecessary configuration options.

---

## Planned Architecture

### MasjidPi Core / Listen / Board

MasjidPi v2 will evolve from the current audio-focused application into a modular home appliance built around a shared MasjidPi core and independent product capabilities.

The repository will remain a single MasjidPi repository initially, with the code organised conceptually as:

```text
MasjidPi/
├── core/
├── listen/
├── board/
├── cmd/
├── frontend/
├── configs/
├── hardware/
├── scripts/
└── docs/
```

#### MasjidPi Core

Core will contain functionality shared by Listen and Board, such as:

- Configuration
- Persistent state
- Masjid selection and shared masjid data
- Catalogue and shared data access
- Networking and common service behaviour
- Common API/client functionality
- Shared Raspberry Pi/platform integration where appropriate

Core should provide reusable foundations rather than becoming a general-purpose dumping ground for unrelated code.

#### MasjidPi Listen

Listen will contain the audio-specific functionality:

- Stream playback
- MPV integration
- Volume control
- Audio-device handling
- Playback recovery and reconnection
- Audio-specific UI and controls

#### MasjidPi Board

Board will contain the MasjidBoard-specific functionality:

- Prayer times
- Jumu'ah times
- Next-prayer and countdown information
- MasjidBoard data retrieval and caching
- Display-oriented UI
- External HDMI display support
- Future display/hardware integration

Listen and Board must remain independent capabilities. Neither subsystem should require the other to operate. In particular, audio playback must continue operating if MasjidBoard is unavailable.

The intended high-level architecture is:

**MasjidPi Core → Listen**

and independently:

**MasjidPi Core → Board**

The application/UI layer may combine information from both capabilities, but Board must not directly control Listen or playback.

The initial implementation should remain a single repository. Core should only be extracted into a separate repository if a future technical need justifies doing so.

This restructuring is a v2 architectural target and should not unnecessarily disrupt the stable v1.x release.

---

## Planned Features

### Raspberry Pi Appliance Image

- Deployable Raspberry Pi appliance image
- Preconfigured operating environment
- Repeatable image-based installation workflow

### Hardware & Advanced Audio

- OLED display
- Push-button controls
- Hardware playback controls
- Display current station and playback state
- Audio equaliser / EQ processing
- Advanced audio controls

### MasjidBoard

MasjidBoard will be a separate subsystem from audio streaming and playback.

Potential capabilities:

- Search for and select a masjid
- Persist the selected masjid
- Display daily prayer times
- Display Jumu'ah times
- Display the next prayer and countdown
- Cache prayer information locally
- Continue displaying cached information during temporary outages
- Periodically refresh MasjidBoard data
- Expose masjid and prayer information through the MasjidPi API
- Display prayer information in the Web UI
- Support future OLED/hardware displays
- Support an external HDMI display connected to the Raspberry Pi

MasjidBoard must remain independent of the playback subsystem. Audio playback must continue operating if MasjidBoard is unavailable.

The intended architecture is:

**MasjidBoard → Application/UI**

and independently:

**Stream Catalogue → Playback → MPV → Audio Device**

The application/UI layer may combine information from both systems, but MasjidBoard must not directly control playback.

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
- Lightweight enough for Raspberry Pi hardware
- Simple to install and operate
- Reliable during network interruptions
- Usable without a separate server
- Focused on listening to mosque streams while expanding to masjid information
- Easy to maintain and update
- Independent of browser-local configuration
- Modular, with Listen and Board capabilities separated from shared Core functionality
- Resilient when optional external services are unavailable

The project should prioritise a reliable home appliance experience over unnecessary complexity.
