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
- MPV stream/playback recovery
- Preserve the selected stream during playback failures
- Resume the last stream after a system reboot
- Resume playback automatically without requiring the Web UI to be opened

MPV **process lifecycle recovery** is being hardened separately in Phase 6. A Phase 6 failure test confirmed that killing the MPV process leaves MasjidPi connected to a dead IPC socket unless the MPV process is recreated.

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

## ✅ Phase 5 — Web UI & Stream Discovery

Phase 5 transformed MasjidPi from a functional streaming service into a more practical radio-style appliance.

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
- Better volume control
- Select the audio output device from the Web UI
- Show the currently selected audio device
- Persist the selected audio device
- Persist volume across restarts and reboots
- Share persistent audio settings across devices
- Recover gracefully when an audio device disappears
- Automatically restore a selected audio device when it becomes available again
- Keep playback/resume independent of whether the Web UI is open

### Phase 5.3 — Audio Controls — Complete

The audio-control work is complete for the v0.x product scope.

- Audio output device discovery
- Audio output device selection
- Restore the selected device after reboot
- Keep device selection available alongside playback controls
- Improve audio-device error handling and feedback
- Persist volume
- Persist audio-device selection
- Automatic audio-device disappearance/reappearance recovery

The following are intentionally **not** part of v0.x:

- Mute/unmute — Stop and Play provide the required control
- EQ/audio processing — deferred to v2.0.0

---

## 🚧 Phase 6 — Production Hardening

Phase 6 is focused on making MasjidPi a dependable, unattended Raspberry Pi appliance rather than adding a large configuration system.

The core principle is to keep configuration simple and keep persistent state owned by the Raspberry Pi rather than by individual browsers.

### Reliability & Recovery

- Gracefully handle missing or corrupt persistent state
- Validate configuration and recover from invalid values
- Improve service recovery after unexpected failures
- Improve diagnostics and useful error reporting
- Continue strengthening unattended operation after reboot and network outages
- **MPV process lifecycle recovery — Phase 6 failure confirmed:** killing MPV leaves MasjidPi running but retrying against the dead `/tmp/masjidpi.sock`; MasjidPi must recreate the MPV process, recreate the IPC connection, and resume playback without restarting the MasjidPi service

### Audio Hardware & Volume Reliability

MasjidPi should use the selected audio device's hardware volume as the user-facing volume control rather than maintaining two independent volume stages.

The intended audio path is:

**Stream → MPV at 100% software gain → ALSA hardware mixer → Audio device**

- Keep MPV software gain at 100%
- Control the selected audio device's ALSA hardware volume from the MasjidPi UI where supported
- Maintain a separate persistent volume level for each audio device
- Use a newly discovered device's current hardware volume as its initial MasjidPi volume
- Restore the saved volume when switching back to a previously configured device
- Handle devices without a controllable hardware mixer gracefully
- Keep hardware-volume state stored on the Raspberry Pi rather than in the browser
- Verify reliable volume and device behaviour across USB, HDMI and other supported outputs

### Updates & Installation

- Safe automatic updates
- Safer update/recovery workflow
- Improve first-run installation experience
- Production installation workflow
- Raspberry Pi appliance image/workflow

### Appliance Operation

- Kiosk/appliance mode
- Read-only filesystem support
- Minimise maintenance required from the end user
- Ensure the system can operate reliably without the Web UI being open

Phase 6 should prioritise reliability and simplicity over adding configuration options that are not required for normal radio operation.

---

## 🔮 v2.0.0 — Hardware, Advanced Audio, MasjidBoard & Home Automation

V2 will add capabilities that go beyond the core radio appliance while keeping the v0.x product simple and reliable.

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

Example automations could include turning on an amplifier when MasjidPi starts playing, turning equipment off when playback stops, reacting to a particular stream being selected, or using prayer-time information to control other home-automation devices.

---

## Current Release

**v0.5.0** represents the completion of the core player, LiveMasjid integration, playback reliability, Raspberry Pi runtime foundations, stream discovery, favourites, and the initial audio-control experience.

The next target is **v0.6.0**, focused on production hardening and appliance reliability.

MasjidBoard integration is planned for **v2.0.0** and is intentionally not part of the v0.x scope.

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
