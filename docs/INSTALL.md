# MasjidPi Installation Guide

## Recommended: official release

For a supported 64-bit Linux system, install the latest official release with:

```bash
curl -fsSL https://raw.githubusercontent.com/X-Calibre/MasjidPi/main/scripts/install-latest.sh | sudo bash
```

The bootstrap installer:

1. Detects the CPU architecture.
2. Retrieves the latest GitHub release.
3. Downloads the matching release archive and `SHA256SUMS`.
4. Verifies the release checksum before extraction.
5. Validates the release contents.
6. Runs the bundled production installer.

No Git checkout or Go installation is required for an official release.

## Supported production systems

MasjidPi production releases currently provide binaries for:

- Linux `x86_64` / AMD64
- Linux `aarch64` / ARM64

The production installer expects:

- Debian, Ubuntu, Linux Mint or Raspberry Pi OS
- `systemd` running as PID 1
- `apt-get`
- an ALSA-compatible audio device

Raspberry Pi 3B has been validated with 64-bit Raspberry Pi OS.

32-bit ARM (`armv6l` / `armv7l`) does not currently have an official pre-built release.

## What the installer changes

Application files are installed under:

```text
/opt/masjidpi
```

Persistent configuration is stored under:

```text
/etc/masjidpi/config.yaml
```

The active LiveMasjid catalogue is stored under:

```text
/var/lib/masjidpi/catalogue.json
```

The installer installs and enables:

```text
masjidpi.service
```

The Web UI listens on port `8080`.

Existing configuration and catalogue data are preserved when upgrading an existing installation.

## Installation validation

The installer does not consider installation successful merely because the systemd service starts.

After starting MasjidPi it verifies:

- the systemd service is running
- the HTTP interface responds
- the `/api/version` endpoint responds
- the running version matches the expected release version
- the player status endpoint responds
- an ALSA audio device is available when one is exposed by the system

If a fresh installation fails, the installer removes the incomplete application runtime rather than leaving a partially installed `/opt/masjidpi` tree behind.

For an existing installation, the safe update workflow stages the new runtime, validates it, and automatically rolls back to the previous runtime if validation fails.

## Source installation

Source installation is intended for development and testing rather than normal production deployment:

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
sudo ./scripts/install.sh --source
```

It builds MasjidPi locally and then uses the same service installation and validation workflow as a release installation.

## Troubleshooting

Check the service:

```bash
sudo systemctl status masjidpi --no-pager
```

View recent logs:

```bash
sudo journalctl -u masjidpi --no-pager -n 100
```

Follow logs while troubleshooting:

```bash
sudo journalctl -u masjidpi -f
```

Check the local API:

```bash
curl -s http://127.0.0.1:8080/api/version
curl -s http://127.0.0.1:8080/api/player/status
```

If installation stops before MasjidPi is running, correct the reported prerequisite or systemd problem and run the installer again.
