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

- Download LiveMasjid catalogue
- Parse LiveMasjid HTML
- Generate local catalogue
- Preserve LiveMasjid stream order
- Generate relay URLs automatically
- Normalise stream names and locations
- Reload catalogue without restarting the application
- Catalogue update API
- Update Catalogue button in the web UI

---

## 🚧 Milestone 0.5

Playback improvements

### Playback reliability

- Detect streams that are offline
- Gracefully handle unavailable streams
- Display playback errors in the UI
- Prevent endless reconnect attempts

### Stream information

- Display current stream name and relay URL
- Show whether a stream is online or offline
- Indicate buffering and reconnecting states

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
- Push buttons

### Audio

- Equaliser
- Presets

### Raspberry Pi

- Kiosk mode
- Read-only filesystem
- Automatic startup
- Systemd service

### Settings

- Persist application settings
- Remember volume
- Remember autoplay
- Remember last selected stream