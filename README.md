# MasjidPi

MasjidPi is a lightweight internet radio for Raspberry Pi, designed for listening to live masjid streams from around the world.

It provides a simple web interface for finding masājid, selecting streams, managing favourites and controlling playback.

![MasjidPi Web Interface](docs/images/masjidpi-ui.png)

## Features

- Browse and search live masjid streams
- Save favourite masājid
- Play and stop streams
- Volume control up to 125%
- Select an audio output device
- Automatically recover when an audio device is disconnected
- Automatically reconnect when a stream or network connection is interrupted
- Remember the last selected stream
- Resume playback automatically after a reboot
- Persistent settings shared across phones, tablets and computers
- Runs as a system service and starts automatically with the Raspberry Pi

## How it works

MasjidPi runs entirely on the Raspberry Pi.

The Raspberry Pi stores the application's settings, favourites and playback state. The web browser is simply a remote control.

This means you can open MasjidPi from your laptop, phone or tablet and see the same settings and favourites.

The Web UI does not need to remain open for playback or automatic recovery to work.

## Requirements

MasjidPi is designed for Linux and Raspberry Pi.

Tested on:

- Raspberry Pi Zero
- Raspberry Pi Zero 2 W
- Raspberry Pi 3
- Raspberry Pi 4
- Raspberry Pi 5

A Linux system with an ALSA-compatible audio device can also be used for development.

## Installation

Clone the repository:

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
```

Run the installer:

```bash
sudo ./scripts/install.sh
```

Once installation is complete, open:

```text
http://<raspberry-pi-ip>:8080
```

For example:

```text
http://192.168.1.50:8080
```

MasjidPi will run automatically as a system service and start after reboot.

## Updating

From an existing installation:

```bash
cd ~/MasjidPi
git pull
sudo ./scripts/install.sh
```

The installer handles stopping the existing service, rebuilding MasjidPi and starting the updated version.

## Useful Commands

Check the service:

```bash
sudo systemctl status masjidpi
```

View the logs:

```bash
sudo journalctl -u masjidpi -f
```

Restart MasjidPi:

```bash
sudo systemctl restart masjidpi
```

## Development

Run MasjidPi directly from the source tree:

```bash
make run
```

Run the tests:

```bash
make test
```

## Project Status

MasjidPi is currently in active development.

**Current version: v0.5.0**

The next development phase focuses on production hardening and making MasjidPi a reliable, unattended Raspberry Pi appliance.

See [ROADMAP.md](ROADMAP.md) for the development roadmap.

## License

MasjidPi is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

See [LICENSE](LICENSE) for details.
