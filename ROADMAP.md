# MasjidPi Roadmap

MasjidPi is a lightweight Raspberry Pi internet radio appliance for LiveMasjid streams.

The roadmap focuses on the current product, remaining production work, and planned future features.

---

## Current Status

**v1.0.7 is released and verified on Raspberry Pi hardware.**

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

- Safe update workflow
- Safer update and recovery behaviour
- Validate a newly installed version before considering an update successful
- Provide a recovery/rollback path if an update fails
- Improve the first-run and production installation workflow where required

The remaining work should prioritise reliability and simplicity rather than unnecessary configuration options.

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

- Lightweight enough for Raspberry Pi hardware
- Simple to install and operate
- Reliable during network interruptions
- Usable without a separate server
- Focused on listening to mosque streams
- Easy to maintain and update
- Independent of browser-local configuration
- Modular, with external services isolated from core playback
- Resilient when optional external services are unavailable

The project should prioritise a reliable radio experience over unnecessary complexity.
