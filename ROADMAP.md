# MasjidPi Roadmap

## ✅ Completed

- [x] Project structure
- [x] AGPL-3.0 licensing
- [x] Configuration management
- [x] Structured logging
- [x] HTTP server
- [x] MPV process management
- [x] MPV IPC communication
- [x] MPV controller abstraction

## 🚧 Current Milestone

### Play the first stream

- [ ] Implement `Play(url)`
- [ ] Play a hard-coded test stream
- [ ] Verify audio output

## 📋 Upcoming Milestones

### Player Control

- [ ] Stop playback
- [ ] Pause / Resume
- [ ] Volume control
- [ ] Mute

### Stream Library

- [ ] Stream model
- [ ] Local stream database
- [ ] Import Livemasjid streams
- [ ] Background refresh

### REST API

- [ ] GET /api/streams
- [ ] POST /api/player/play
- [ ] POST /api/player/stop
- [ ] POST /api/player/volume

### Web UI

- [ ] Stream list
- [ ] Search
- [ ] Current playing indicator
- [ ] Volume slider

### Raspberry Pi

- [ ] systemd service
- [ ] Read-only friendly configuration
- [ ] Auto-start on boot
- [ ] Pi Zero 2 optimisation
