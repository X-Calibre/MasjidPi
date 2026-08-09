# MasjidPi Roadmap

MasjidPi is being developed as a lightweight Raspberry Pi internet radio appliance for LiveMasjid streams.

The roadmap is organised by product phase rather than historical version numbers. Completed phases document what is already implemented and verified; upcoming phases describe the next product priorities.

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

LiveMasjid is now integrated directly into MasjidPi without requiring a separate scraping service.

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

---

## ✅ Phase 3 — Playback Reliability

Playback and network recovery have been implemented and tested on Raspberry Pi.

- Detect playback failures
- Gracefully handle unavailable streams
- Display playback errors in the API/UI
- Controlled playback retry and reconnect behaviour
- Automatic recovery after temporary network failures
- Recovery after Ethernet disconnect/reconnect
- Recovery after extended network outages
- MQTT disconnect/reconnect handling
- MPV recovery
- Preserve the selected stream during playback failures
- Resume the last stream after a system reboot

---

## ✅ Phase 4 — Raspberry Pi Runtime

The Raspberry Pi installation and runtime workflow is now functional and has been tested on real hardware.

- Runtime directory layout
- Runtime path abstraction
- Automatic dependency installation
- Go installation/checking
- Build automation
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

---

## 🚧 Phase 5 — Web UI & Stream Discovery

**Next major development phase.** The goal is to make MasjidPi feel like a simple radio rather than a technical streaming tool.

### Stream Discovery

- Search streams instantly
- Search by mosque name
- Search by location
- Filter the stream catalogue
- Improve stream selection and browsing

### Favourites

- Favourite stations
- Persistent favourites
- Quick access to favourite stations

### Playback & Audio Controls

- Clear current-stream display
- Clear playing/reconnecting/error states
- Improved playback controls
- Better volume control
- **Select the audio output device from the Web UI**
- Show the currently selected audio device
- Persist the selected audio device
- More useful playback and audio feedback

---

## 🚧 Phase 6 — Configuration & Personalisation

Make common MasjidPi settings configurable without editing files manually.

- Web-based configuration
- First-run setup
- Playback settings
- Volume preferences
- Persist user configuration
- Additional user preferences

---

## 🔮 Phase 7 — Hardware Interface

Turn MasjidPi into a dedicated physical radio appliance.

- OLED display
- Push-button controls
- Hardware playback controls
- Display current station and playback state

---

## 🔮 Phase 8 — Appliance & Production Hardening

Prepare MasjidPi for simple, unattended deployment on Raspberry Pi hardware.

- Kiosk mode
- Read-only filesystem mode
- Safe automatic updates
- Update/recovery safety
- Production installation workflow
- Raspberry Pi appliance image/workflow

---

## Current Release

**v0.5.0** represents the completion of the core player, LiveMasjid integration, playback reliability, and Raspberry Pi runtime foundations.

The next target is **v0.6.0**, focused on Web UI improvements, stream discovery, favourites, and user-facing audio controls.

---

## Project Principles

MasjidPi should remain:

- Lightweight enough for Raspberry Pi hardware
- Simple to install and operate
- Reliable during network interruptions
- Usable without a separate server
- Focused on listening to mosque streams
- Easy to maintain and update

The project should prioritise a reliable radio experience over unnecessary complexity.
