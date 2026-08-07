# MasjidPi Roadmap

## ✅ Milestone 0.1

Project foundation

- Repository structure
- Backend skeleton
- Frontend skeleton
- MPV integration

---

## ✅ Milestone 0.2

Player API

- Play
- Stop
- Status
- Volume
- Responsive web UI

---

## ✅ Milestone 0.3

Stream catalogue

- Stream model
- Local catalogue
- Play by stream ID
- Stream API
- Remember last stream
- Auto play

---

## ✅ Milestone 0.4

Catalogue updater

### Completed

- Download LiveMasjid catalogue
- Parse LiveMasjid HTML
- Generate local catalogue
- Preserve LiveMasjid stream order
- Generate relay URLs automatically
- Normalise stream names and locations
- Reload catalogue without restarting
- Catalogue update API
- Update Catalogue button in the web interface

---

## ✅ Milestone 0.5

Installer & Runtime

### Completed

- Runtime directory layout
- Runtime path abstraction
- Automatic dependency installation
- Go installation
- Build automation
- Runtime installation
- systemd service installation
- Automatic startup on boot
- Installer self-test

---

## 🚧 Milestone 0.6

Update mechanism

### Planned

- Runtime layout abstraction
- Separate application from user data
- Detect newer releases
- Download release package
- Install update
- Restart service
- Preserve local configuration
- Preserve catalogue and user data
- End-user installer (future)

---

## ✅ Milestone 0.7

Playback reliability

### Completed

- Detect offline streams
- Gracefully handle unavailable streams
- Display playback errors
- Prevent endless reconnect attempts
- Automatic recovery after temporary failures

---

## Future

### Search

- Instant filtering
- Search by mosque
- Search by location

### Favourites

- Favourite stations
- Quick access

### Hardware

- OLED display
- Push-button controls

### Audio

- Equaliser
- Presets

### Raspberry Pi

- Kiosk mode
- Read-only mode
- Automatic updates

### User Experience

- Multi-language interface
- First-run setup wizard
- Web-based configuration