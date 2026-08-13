# MasjidPi

MasjidPi is a lightweight internet radio for Raspberry Pi, designed for listening to live masjid streams from around the world.

It provides a simple web interface for finding masājid, selecting streams, managing favourites and controlling playback.

![MasjidPi Web Interface](docs/images/Updated-UI-v1.0.7.png)

## Features

- Browse and search live masjid streams
- Update the LiveMasjid stream catalogue from the Web UI
- Save favourite masājid
- Play and stop streams
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

### Supported operating systems

The installer currently supports:

- Debian
- Raspberry Pi OS / Raspbian
- Ubuntu
- Linux Mint

The system must have an ALSA-compatible audio device and `systemd`.

### Supported release architectures

Official pre-built releases are provided for:

- **Linux ARM64 (`aarch64`)** — recommended for 64-bit Raspberry Pi OS on Raspberry Pi 3B and newer
- **Linux AMD64 (`x86_64`)** — suitable for standard 64-bit Intel/AMD Linux systems

The official release packages are **ARM64 and AMD64 only**. ARMv7/ARMv6 release packages are not provided.

### Raspberry Pi testing

The following Raspberry Pi models have been tested:

- **Raspberry Pi 3B — works**
- **Raspberry Pi Zero W — not supported**

The Raspberry Pi Zero W has been tested and MasjidPi does not run reliably on that hardware. Do not use a Pi Zero W for a MasjidPi installation.

Other Raspberry Pi models and variants have not yet been fully tested and should not currently be considered officially supported.

## Installation

### Recommended: install the latest release

For normal users, the recommended installation method is a **single command**. You do not need to clone the Git repository, install Go, determine the CPU architecture manually, download an archive manually, or run a separate service verification step.

Run:

```bash
curl -fsSL https://raw.githubusercontent.com/X-Calibre/MasjidPi/main/scripts/install-latest.sh | sudo bash
```

The bootstrap installer automatically:

1. Detects whether the system is `aarch64` or `x86_64`.
2. Selects the latest official ARM64 or AMD64 release.
3. Downloads the release archive and `SHA256SUMS`.
4. Verifies the release checksum.
5. Extracts the release package.
6. Runs the complete bundled MasjidPi installer.

The bundled installer then installs the required runtime dependencies, MasjidPi, the Web UI and the systemd service, starts the service and runs the MasjidPi self-test.

Existing configuration and persistent runtime data are preserved during upgrades.

### Open the Web UI

From another computer, phone or tablet on the same network, open:

```text
http://<raspberry-pi-ip>:8080
```

For example:

```text
http://192.168.1.50:8080
```

The Web UI is a remote control. It does not need to remain open for MasjidPi to continue playing audio or recovering from stream/network interruptions.

## Configuration

The persistent configuration is stored at:

```text
/etc/masjidpi/config.yaml
```

The active stream catalogue is stored at:

```text
/var/lib/masjidpi/catalogue.json
```

Application files are installed under:

```text
/opt/masjidpi
```

The Web UI is available from the same MasjidPi service on port `8080`.

Configuration changes made through the Web UI are persistent and are shared across clients using the same MasjidPi installation.

### Updating the LiveMasjid catalogue

The Web UI includes an **Update Catalogue** button. The catalogue updater downloads the latest LiveMasjid data and writes the generated catalogue to the active runtime location:

```text
/var/lib/masjidpi/catalogue.json
```

The Web UI reloads the updated catalogue after a successful update. The updater does not use the legacy relative `data/` catalogue path for an installed deployment.

## Updating MasjidPi

The same single command can be used for both fresh installations and upgrades:

```bash
curl -fsSL https://raw.githubusercontent.com/X-Calibre/MasjidPi/main/scripts/install-latest.sh | sudo bash
```

The installer detects the existing installation and updates it in place. It replaces the application binary and Web UI while preserving `/etc/masjidpi` and `/var/lib/masjidpi`.

The installer also handles older release packages that do not contain a bundled catalogue by preserving the existing runtime catalogue rather than treating it as an installation failure.

Each official release provides a `SHA256SUMS` file, and the single-command installer verifies the downloaded archive automatically before installation.

### Installing from source

Source installation is intended for development and testing rather than normal users.

Clone the repository and run the installer with `--source`:

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
sudo ./scripts/install.sh --source
```

Source installation builds MasjidPi locally using Go. It is useful when testing unreleased changes or developing MasjidPi.

## Release packages

Each release includes pre-built archives for the supported architectures:

```text
masjidpi-vX.Y.Z-linux-arm64.tar.gz
masjidpi-vX.Y.Z-linux-amd64.tar.gz
```

A `SHA256SUMS` file is provided with each release.

Release archives contain:

- MasjidPi executable
- Default configuration
- Initial stream catalogue
- Web UI
- Version information
- Complete installer scripts
- Systemd service definition

## Useful Commands

Check the service:

```bash
sudo systemctl status masjidpi --no-pager
```

View the logs:

```bash
sudo journalctl -u masjidpi -f
```

Restart MasjidPi:

```bash
sudo systemctl restart masjidpi
```

Stop MasjidPi:

```bash
sudo systemctl stop masjidpi
```

Check playback status:

```bash
curl -s http://127.0.0.1:8080/api/player/status
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

**Current stable release: v1.0.6**

v1.0.6 has been validated on a Raspberry Pi 3B running a 64-bit ARM64 Linux environment. Validation includes installation from the official pre-built release, use of the bundled release installer, systemd service operation, MPV playback, the Web UI, persistent configuration, live masjid stream playback and successful LiveMasjid catalogue updates using the installed `/var/lib/masjidpi` runtime path.

The v1.0.6 release also establishes the release-first installation workflow for normal users. Official releases contain the pre-built application and complete installer, so normal users do not need to clone the source repository, install Go or build MasjidPi. The single-command bootstrap installer selects the correct official package for the host architecture, verifies its checksum and invokes the bundled release installer. Source installation remains available for developers with `--source`.

See [ROADMAP.md](ROADMAP.md) for the development roadmap.

## Acknowledgements

MasjidPi was inspired by the [eBilal project](https://github.com/Muslims-in-IT/ebilal). We gratefully acknowledge the eBilal contributors and the work they did to make a Raspberry Pi-based masjid audio receiver available as an open-source project.

MasjidPi relies on [LiveMasjid](https://www.livemasjid.com/) for the live masjid streams and the related stream and status data that make the service possible. We gratefully acknowledge the LiveMasjid team for providing and maintaining this service.

MasjidPi is an independent project and is not affiliated with or endorsed by eBilal or LiveMasjid.

## License

MasjidPi is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

See [LICENSE](LICENSE) for details.
