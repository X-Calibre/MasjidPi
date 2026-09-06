# MasjidPi Installation Guide

## Recommended: official release

For a supported 64-bit Linux system, install the latest official release with:

```bash
curl -fsSL https://raw.githubusercontent.com/X-Calibre/MasjidPi/main/scripts/install-latest.sh | sudo bash
```

On an interactive terminal, the production installer prompts for the appliance profile:

1. Listen
2. Board
3. Listen + Board

The bootstrap installer:

1. Detects the CPU architecture.
2. Retrieves the latest GitHub release.
3. Downloads the matching release archive and `SHA256SUMS`.
4. Verifies the release checksum before extraction.
5. Validates required release contents.
6. Reconnects the bundled installer to the controlling terminal when available so component selection remains interactive.
7. Runs the bundled production installer.

If no controlling terminal is available, the installer uses its documented non-interactive default profile.

No Git checkout or Go installation is required for an official release.

## Supported production systems

MasjidPi production releases currently provide binaries for:

- Linux `x86_64` / AMD64
- Linux `aarch64` / ARM64

The production installer expects:

- Debian, Ubuntu, Linux Mint or Raspberry Pi OS
- `systemd` running as PID 1
- `apt-get`

Component-specific requirements are installed only when needed:

- **Listen** installs MPV, FFmpeg and ALSA tools and requires an ALSA-compatible audio device for playback.
- **Board** installs Cog and WPE WebKit and uses a DRM/KMS display runtime for the dedicated HDMI board.
- **Listen + Board** installs both dependency sets.

Raspberry Pi 3B and Raspberry Pi 4 have been validated with 64-bit Raspberry Pi OS. Board has also been validated as a dedicated Raspberry Pi OS Lite HDMI appliance using Cog/WPE directly on DRM/KMS. See [HARDWARE.md](HARDWARE.md) for platform status and measured performance.

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

The installed component profile is stored under:

```text
/etc/masjidpi/components.env
```

Persistent runtime data is stored under:

```text
/var/lib/masjidpi
```

The main application service is:

```text
masjidpi.service
```

When Board is installed, the installer also installs and enables:

```text
masjidpi-display-warmup.service
masjidpi-display.service
```

It also installs the boot-splash and read-only boot-firmware protection assets used on Raspberry Pi Board appliances.

When Board is removed from the selected profile, the display service and launcher are removed and stale systemd failure state is cleared.

The Web UI listens on port `8080`.

Existing configuration and runtime data are preserved when upgrading an existing installation or changing component profile.

## First-run appliance setup

On a fresh Raspberry Pi Board appliance with no saved Wi-Fi profile, the attached display opens the MasjidFrame setup wizard automatically. The wizard is designed for the constrained Cog/WPE renderer and portrait touchscreen; it does not require a desktop environment or a separate operating-system keyboard.

The wizard can:

1. Scan for nearby 2.4 GHz Wi-Fi networks.
2. Connect to a visible network using the on-screen keyboard.
3. Add a hidden password-protected or open network by entering its SSID manually.
4. Select the initial country, province/region, town/city and MasjidBoard.
5. Start the appliance Board after downloading its first timetable.

After Wi-Fi connects, the wizard displays the appliance IPv4 URL and, when supplied by DHCP or network DNS, its FQDN URL. It does not invent a `.local` hostname. From a phone, tablet or computer on the same network, use either displayed address to configure Listen, audio, additional masjids and display preferences.

The normal appliance startup bypasses the wizard once NetworkManager has a saved Wi-Fi profile. Removing a profile is therefore a factory/reset operation rather than a normal way to reopen setup.

## Installation validation

The installer does not consider installation successful merely because the systemd service starts.

Every profile validates:

- the main systemd service is running
- the HTTP interface responds
- the `/api/version` endpoint responds
- the running version matches the expected release version

Listen profiles additionally validate:

- `/api/player/status` responds
- an audio device is reported when one is exposed by the system

Board profiles additionally validate:

- `/api/masjidboard/status` responds
- `masjidpi-display.service` reaches the running state

The self-test is component-aware, so APIs and hardware belonging to an uninstalled component are not required.

If a fresh installation fails, the installer removes the incomplete application runtime rather than leaving a partially installed `/opt/masjidpi` tree behind.

For an existing installation, the safe update workflow stages the new runtime, applies the selected component profile, validates the result and automatically restores the previous runtime/profile if validation fails.

## Changing component profile

Re-run the installer on an existing installation to change between:

- Listen
- Board
- Listen + Board

The current profile is shown as the default selection. Profile changes are transactional and preserve persistent configuration/data.

For example, changing from Board to Listen removes the Cog display service and starts the MPV playback subsystem. Changing from Listen to Board removes the MPV runtime subsystem and installs the display service.

## Source installation

Source installation is intended for development and testing rather than normal production deployment:

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
sudo ./scripts/install.sh --source
```

It builds MasjidPi locally and then uses the same component selection, service installation and validation workflow as a release installation.

## Troubleshooting

Check the main service:

```bash
sudo systemctl status masjidpi --no-pager
```

Check the Board display service when Board is installed:

```bash
sudo systemctl status masjidpi-display --no-pager
```

View recent logs:

```bash
sudo journalctl -u masjidpi --no-pager -n 100
```

Follow Board display logs:

```bash
sudo journalctl -u masjidpi-display -f
```

Check installed components:

```bash
curl -s http://127.0.0.1:8080/api/components
```

For Listen:

```bash
curl -s http://127.0.0.1:8080/api/player/status
```

For Board:

```bash
curl -s http://127.0.0.1:8080/api/masjidboard/status
```

If installation stops before MasjidPi is running, correct the reported prerequisite or systemd problem and run the installer again.
