# MasjidPi Roadmap

MasjidPi is a lightweight Raspberry Pi internet radio appliance for LiveMasjid streams.

The roadmap is organised by product phase. Completed phases describe functionality that is implemented and verified. Current work focuses only on remaining production-hardening tasks; future features are intentionally deferred to v2.

---

## ✅ Phase 1 — Core Player

The core radio player and playback API are complete.

- MPV integration
- Stream catalogue and stream model
- Play by stream ID
- Play, stop, status and volume controls
- Responsive web interface
- Local catalogue
- Remember the last played stream
- Resume the last stream after reboot

---

## ✅ Phase 2 — LiveMasjid Integration

LiveMasjid is integrated directly into MasjidPi without requiring a separate scraping service.

- Download LiveMasjid catalogue
- Parse LiveMasjid HTML
- Generate local catalogue
- Preserve LiveMasjid stream order
- Generate relay URLs automatically
- Normalise stream names and locations
- Reload catalogue without restarting MasjidPi
- Catalogue update API
- Update Catalogue button in the web interface
- LiveMasjid MQTT status feed
- MQTT connection recovery
- Catalogue updater uses the same runtime paths as the installed application
- Web UI catalogue updates write to and reload the active `/var/lib/masjidpi` catalogue

The LiveMasjid catalogue runtime-path issue was fixed and verified in v1.0.5.

---

## ✅ Phase 3 — Playback Reliability

Playback and network recovery are implemented and tested on Raspberry Pi.

- Detect playback failures
- Gracefully handle unavailable streams
- Display playback errors in the API/UI
- Controlled playback retry and reconnect behaviour
- Automatic recovery after temporary network failures
- Recovery after Ethernet disconnect/reconnect
- Recovery after extended network outages
- MQTT disconnect/reconnect handling
- MPV stream/playback recovery
- MPV process lifecycle recovery and recreation after unexpected MPV exit
- Preserve the selected stream during playback failures
- Resume the last stream after a system reboot
- Resume playback automatically without requiring the Web UI to be opened

MPV process recovery is implemented in the player process layer and includes restart handling after an unexpected MPV exit.

---

## ✅ Phase 4 — Raspberry Pi Runtime

The Raspberry Pi installation and runtime workflow is functional and has been tested on real hardware.

- Runtime directory layout
- Runtime path abstraction
- Automatic dependency installation
- Go installation/checking for source builds
- Build automation for source builds
- Runtime installation
- systemd service installation
- Automatic startup on boot
- Service restart/stop handling during installation
- Installer self-test
- HTTP interface health check
- Audio device detection
- Clean repository/build-artifact handling
- MPV journal output cleanup
- ALSA audio selection for the appliance runtime
- Installation and reboot verification
- Pre-built release packaging for ARM64 and AMD64
- Release package checksum verification
- Bundled release installer
- Release installation without Git or Go
- Upgrade installation preserving persistent configuration and runtime data

The release-first installation workflow was introduced and verified in v1.0.6 on a Raspberry Pi 3B.

---

## ✅ Phase 5 — Web UI & Stream Discovery

Phase 5 transformed MasjidPi from a functional streaming service into a practical radio-style appliance.

### Stream Discovery — Complete

- Search streams instantly
- Search by mosque name
- Search by location
- Filter the stream catalogue
- Improve stream selection and browsing

### Favourites — Complete

- Favourite stations
- Persistent favourites
- Quick access to favourite stations
- Shared favourites across phones, tablets and computers
- Favourites stored on the Raspberry Pi rather than in browser-local storage

### Playback & Audio Controls — Complete

- Clear current-stream display
- Clear playing/reconnecting/error states
- Improved playback controls
- Select the audio output device from the Web UI
- Show the currently selected audio device
- Persist the selected audio device
- Persist volume across restarts and reboots
- Share persistent audio settings across devices
- Recover gracefully when an audio device disappears
- Automatically restore a selected audio device when it becomes available again
- Keep playback/resume independent of whether the Web UI is open

### Phase 5.3 — Audio Controls — Complete

The audio-control work is complete for the current product scope.

- Audio output device discovery
- Audio output device selection
- Restore the selected device after reboot
- Persist volume
- Persist audio-device selection
- Automatic audio-device disappearance/reappearance recovery
- ALSA hardware-volume control where supported
- Separate persistent volume level for each audio device
- MPV software gain kept at 100%
- Graceful handling of devices without a controllable hardware mixer
- Hardware verification across the supported USB/HDMI audio configurations is complete

The following are intentionally not part of the current product scope:

- Mute/unmute — Stop and Play provide the required control
- EQ/audio processing — deferred to v2

---

## 🚧 Phase 6 — Production Hardening

Phase 6 is now intentionally small. The core player, runtime, Web UI, recovery behaviour, audio hardware handling and release installation workflow are considered complete and reliable enough for the current product scope.

### Completed Hardening

The following are considered complete and require no further work at this stage:

- Persistent-state handling
- Configuration validation
- Service recovery
- Network recovery
- MPV process recovery
- Audio-device recovery
- ALSA hardware-volume handling
- Per-device persistent volume
- Hardware verification
- Reliable operation without the Web UI being open
- LiveMasjid catalogue runtime-path integration
- Web UI catalogue updates against the active `/var/lib/masjidpi` catalogue
- Pre-built ARM64/AMD64 release packages
- Release checksum verification
- Bundled release installer
- Release installation without Git or Go
- Upgrade installation preserving persistent configuration and runtime data
- v1.0.6 release installation and validation on Raspberry Pi 3B

### Remaining Work

#### Updates & Recovery

- Safe update workflow
- Safer update/recovery behaviour
- Validate a newly installed version before considering an update successful
- Provide a recovery/rollback path if an update fails
- Improve the first-run/production installation workflow where required

Phase 6 should prioritise reliability and simplicity over adding configuration options that are not required for normal radio operation.

---

## 🔮 v2 — Future Platform & Advanced Features

Features outside the current radio-appliance scope are deferred to v2.

### Raspberry Pi Appliance Image

- Build and maintain a deployable Raspberry Pi appliance image/workflow
- Preconfigure the operating environment for MasjidPi
- Provide a repeatable image-based installation path

### Hardware & Advanced Audio

- OLED display
- Push-button controls
- Hardware playback controls
- Display current station and playback state
- Audio equaliser / EQ processing
- Advanced audio controls

### MasjidBoard Integration

MasjidPi should be able to associate the selected mosque with its MasjidBoard Live listing and display the mosque's prayer information alongside the audio experience.

MasjidBoard must remain a **separate subsystem from audio streaming and playback**. The audio player must continue to operate normally if MasjidBoard Live is unavailable, changes its interface, or cannot be reached.

Potential capabilities:

- Search MasjidBoard Live for a mosque
- Select and persist the corresponding masjid
- Display daily prayer times
- Display Jumu'ah times
- Display the next prayer and relevant countdown information
- Display current masjid information alongside the selected audio stream
- Cache prayer information locally on the Raspberry Pi
- Continue displaying the last known valid prayer information during temporary network outages
- Periodically refresh MasjidBoard data without affecting audio playback
- Expose masjid and prayer information through the MasjidPi API
- Display prayer information in the Web UI
- Support a future OLED/hardware display for prayer information
- Investigate a reliable machine-readable MasjidBoard data/API interface rather than depending unnecessarily on HTML scraping

The intended architecture is:

**MasjidBoard subsystem → Masjid/application layer → Web UI/display**

and independently:

**Stream catalogue → Playback subsystem → MPV → Audio device**

The application layer may combine information from both subsystems, for example by showing the next prayer alongside the currently playing mosque, but MasjidBoard should not directly control MPV or the playback subsystem.

### Home Assistant Integration

Allow MasjidPi to participate in Home Assistant automations and expose its playback state and controls to the home-automation system.

Potential capabilities:

- Home Assistant media-player entity
- Play / stop control
- Volume control
- Current stream / mosque information
- Playback state
- Audio output device
- Stream and playback status sensors
- Events/triggers for stream started, stopped and changed
- Events/triggers for playback failure and reconnection
- Events/triggers for audio-device changes, loss and recovery
- Local MQTT integration, building on MasjidPi's existing MQTT architecture
- Option for a polished native Home Assistant integration in addition to MQTT
- Prayer-time and masjid information as optional Home Assistant sensors
- Events/triggers related to upcoming prayer times where appropriate

---

## Current Project Status

**v1.0.6 has been released and verified on Raspberry Pi hardware.** The release-first installation workflow has been validated using the official ARM64 release package and its bundled installer. The Web UI catalogue update workflow has also been verified against the active `/var/lib/masjidpi` catalogue.

The original v0.x roadmap has been substantially completed. The core radio player, LiveMasjid integration, playback reliability, Raspberry Pi runtime, Web UI, stream discovery, favourites, audio hardware controls and release installation workflow are implemented and verified.

The remaining pre-v2 work is now limited to:

1. Building a safe update/recovery workflow, including validation of newly installed versions and rollback handling for failed updates.

The following previously proposed Phase 6 items are explicitly considered complete and require no further work:

- Persistent-state hardening
- Configuration validation
- Service recovery
- Audio hardware and volume reliability
- LiveMasjid catalogue runtime-path integration
- Release packaging and bundled installation workflow

The following are no longer planned as standalone features:

- Kiosk/appliance mode
- Read-only filesystem support

The Raspberry Pi appliance image is moved to **v2**.

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
- Modular, with external services such as MasjidBoard isolated from core playback
- Resilient when optional external services are unavailable

The project should prioritise a reliable radio experience over unnecessary complexity.
