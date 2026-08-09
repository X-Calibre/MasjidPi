# MasjidPi

MasjidPi is an open-source internet radio designed specifically for listening to live audio streams from masājid around the world.

The project is intentionally lightweight and designed to run on low-powered Raspberry Pi hardware while providing a simple web interface for selecting and controlling mosque audio streams.

---

## Screenshot

![MasjidPi Web Interface](docs/images/masjidpi-ui.png)

---

## About this Project

MasjidPi has a unique story.

It is being developed by someone with **no prior software development experience**. Every line of code has been developed through an iterative collaboration with AI, while the project vision, feature decisions, testing, and overall direction remain human-driven.

Rather than asking AI to generate an application in one step, every feature is designed, implemented, tested, refined and documented through many small iterations.

The goal is not only to build a useful piece of software, but also to demonstrate that modern AI tools can empower people without formal programming backgrounds to create high-quality open-source software through curiosity, persistence and careful testing.

If this project inspires someone else to begin learning software development, then it has already been a success.

---

## Vision

A Raspberry Pi connected to speakers that allows users to:

- Browse live masjid streams
- Listen to live broadcasts
- Save favourite masājid
- Automatically reconnect after network outages
- Resume playback after power loss
- Be managed entirely through a web browser

MasjidPi is designed to run on low-powered Raspberry Pi hardware while remaining extensible for future features.

---

## Project Philosophy

MasjidPi follows a few simple principles:

- Keep the software lightweight.
- Prioritise reliability over complexity.
- Design for Raspberry Pi first.
- Build features only when they solve a real problem.
- Prefer simple, maintainable code over clever code.
- Keep the application self-contained whenever practical.
- Minimise external dependencies.
- Build incrementally and test every feature.

---

## Current Features

- Responsive web-based player interface
- Live playback status
- Volume control (0–125%)
- MPV-based audio playback
- ALSA audio output
- Local stream catalogue
- LiveMasjid catalogue updates
- Preserve LiveMasjid stream ordering
- Generate relay URLs automatically
- Play streams by catalogue selection
- Remember the last selected stream
- Automatically resume the last playing stream after restart or reboot
- Automatically retry streams that stop unexpectedly
- Automatically reconnect when a selected masjid broadcast becomes available again
- MQTT LiveMasjid status monitoring with automatic reconnect/resubscribe
- Runtime catalogue updates without restarting MasjidPi
- Installable as a Linux systemd service
- Automatic startup on boot
- Installation self-test
- Automated Go tests

---

## Recent Improvements

- Automatic LiveMasjid catalogue updates
- Runtime path abstraction
- Persistent last-playback state across restarts
- Automatic recovery from offline streams
- MQTT status-feed reconnect handling
- "Waiting for Masjid" playback state
- Persistent frontend settings
- Responsive web interface
- Installation self-test
- systemd service support
- Automated API and playback tests

---

## Planned Features

- Automatic application updates
- Searchable stream catalogue
- Favourite masājid
- OLED display support
- Push-button controls
- Audio equaliser
- Multi-language interface
- Read-only Raspberry Pi mode

---

## Requirements

MasjidPi currently supports Linux systems with:

- Go 1.26 or newer
- MPV
- Git
- ALSA-compatible audio hardware for playback
- systemd for the packaged service installation

The installer can install missing packages and Go where supported.

Supported hardware currently includes:

- Raspberry Pi Zero
- Raspberry Pi Zero 2 W
- Raspberry Pi 3
- Raspberry Pi 4
- Raspberry Pi 5
- Standard Linux PCs (development)

The project is primarily tested on Raspberry Pi and is designed with low-powered hardware in mind.

---

# Installation

The recommended installation method is the included installer.

### 1. Clone the repository

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
```

### 2. Run the installer

```bash
sudo ./scripts/install.sh
```

The installer will:

1. Verify that it is running as root.
2. Detect the Linux distribution and CPU architecture.
3. Install required system packages.
4. Install Go if it is not already available.
5. Build MasjidPi from the current source tree.
6. Stop an existing MasjidPi service when updating an installation.
7. Install the runtime under `/opt/masjidpi`.
8. Install the configuration under `/etc/masjidpi`.
9. Install persistent application data under `/var/lib/masjidpi`.
10. Install and enable the systemd service.
11. Start MasjidPi.
12. Run the installation self-test.

A successful installation ends with the service active and the HTTP interface responding.

### 3. Open the Web UI

After installation, open:

```text
http://<raspberry-pi-ip>:8080
```

For a local installation, this can also be:

```text
http://localhost:8080
```

---

## Installation Verification

Check the service:

```bash
sudo systemctl status masjidpi --no-pager
```

Check the player API:

```bash
curl -s http://127.0.0.1:8080/api/player/status
```

Run the automated Go tests:

```bash
cd backend
go test ./...
```

Follow the service log:

```bash
sudo journalctl -u masjidpi -f
```

The installer also performs its own service, HTTP and audio-device checks.

---

## Updating an Existing Installation

If you are working from a cloned development repository, update the source first:

```bash
cd ~/MasjidPi
git pull
```

Then run the installer again:

```bash
sudo ./scripts/install.sh
```

The installer stops the running service before replacing the runtime binary, installs the updated systemd unit, starts the service again and runs the self-test.

After an update, verify the service with:

```bash
sudo systemctl status masjidpi --no-pager
```

If the service fails to start, inspect the recent logs:

```bash
sudo journalctl -u masjidpi -n 50 --no-pager
```

---

## Runtime Layout

The installed runtime is laid out as follows:

```text
/opt/masjidpi
├── bin
│   └── masjidpi
└── frontend
    ├── index.html
    ├── app.js
    └── style.css
```

Configuration is installed separately:

```text
/etc/masjidpi/config.yaml
```

Persistent application data is stored under:

```text
/var/lib/masjidpi
├── catalogue.json
└── playback.json
```

`playback.json` contains the selected stream ID used to resume playback after a normal restart or reboot. Explicitly stopping playback clears the saved playback state.

The systemd service runs the application from `/opt/masjidpi` and starts automatically at boot.

---

## Configuration

The default configuration is installed to:

```text
/etc/masjidpi/config.yaml
```

The default settings include:

```yaml
http:
  address: ":8080"

player:
  socket: "/tmp/masjidpi.sock"
  volume: 100

streams:
  refresh_interval: "6h"

playback:
  retry_interval: "5s"
  reconnect_delay: "5s"
```

Changes to the configuration generally require restarting the service:

```bash
sudo systemctl restart masjidpi
```

---

## Catalogue

MasjidPi uses the LiveMasjid catalogue to populate its stream list and automatically generates the relay URLs used for playback.

The catalogue refresh interval is currently **6 hours**. Catalogue updates are performed by MasjidPi at runtime and do not require restarting the service.

The installed catalogue is stored at:

```text
/var/lib/masjidpi/catalogue.json
```

---

## Playback Recovery

MasjidPi is designed to keep a selected stream playing when possible.

If a stream stops unexpectedly, the playback manager enters a retrying state and attempts playback again.

If a selected masjid is temporarily offline, MasjidPi waits for the stream to become available and automatically retries it.

The selected stream is persisted so that a normal restart or Raspberry Pi reboot can resume the last stream. The stream is not resumed when the user explicitly presses Stop, because that action clears the saved playback state.

---

## Network and MQTT Recovery

MasjidPi monitors LiveMasjid stream status through MQTT.

If the MQTT status connection is lost, MasjidPi automatically reconnects and resubscribes. A temporary network outage therefore does not require manually restarting the service.

The audio player and MQTT status feed recover independently, so a network interruption can temporarily place playback into a retrying state while the service itself remains running.

---

## Development

To run directly from the source tree:

```bash
make run
```

or:

```bash
cd backend
go run ./cmd/masjidpi
```

Format the Go source:

```bash
make fmt
```

Update Go dependencies:

```bash
make tidy
```

Run the automated tests:

```bash
make test
```

The test suite does not require MPV or an audio device for the unit-tested packages.

---

## Web Interface

After installation, open:

```text
http://localhost:8080
```

When accessing MasjidPi from another device, replace `localhost` with the Raspberry Pi's IP address, for example:

```text
http://192.168.1.50:8080
```

The Web UI provides stream selection, playback status and volume control.

---

## Troubleshooting

### Service is not running

```bash
sudo systemctl status masjidpi --no-pager
sudo journalctl -u masjidpi -n 50 --no-pager
```

### HTTP interface is not responding

Check that the service is running and that port 8080 is listening:

```bash
sudo systemctl status masjidpi --no-pager
curl -i http://127.0.0.1:8080
```

### Playback is retrying

Check the player status:

```bash
curl -s http://127.0.0.1:8080/api/player/status
```

Then inspect recent playback messages:

```bash
sudo journalctl -u masjidpi --since "10 minutes ago" --no-pager \
  | grep -Ei "retry|reconnect|playback|error"
```

A `retrying` state is expected when the selected stream is temporarily unavailable.

### MQTT status feed disconnected

Check for disconnect/reconnect events without the high-volume per-mount notifications:

```bash
sudo journalctl -u masjidpi --since "10 minutes ago" --no-pager \
  | grep -Ei "LiveMasjid status feed|MQTT.*(disconnect|connect|reconnect)" \
  | grep -v "mounts/"
```

A disconnected status feed should be followed by a successful connection message once network connectivity is restored.

### No audio output

Check that Linux detects an audio device:

```bash
aplay -l
```

MasjidPi uses MPV with ALSA audio output in the packaged systemd service.

---

## Development Status

MasjidPi is currently in active development and is versioned below 1.0.

The following components are complete:

- Core player
- Stream catalogue
- Catalogue updater
- Installer
- Runtime environment
- Systemd integration
- Persistent playback state
- MQTT reconnect handling
- Automated tests

Current development is focused on:

- Automatic application updates
- Search
- Favourites
- Further Raspberry Pi optimisation

Breaking changes may still occur before version 1.0.

---

## Contributing

Bug reports, feature requests and pull requests are welcome.

Please see:

- `CONTRIBUTING.md`

---

## License

MasjidPi is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

See the `LICENSE` file for details.
