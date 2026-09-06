# MasjidPi Roadmap

MasjidPi is a lightweight appliance for live masjid audio and prayer-time information.

## Current release

**v1.5.2 is the current stable release. v1.6.0-rc.1 is the current release candidate.**

MasjidPi currently provides:

- independent Listen, Board and combined appliance profiles;
- touchscreen first-run Wi-Fi, location and primary-masjid setup for the portrait appliance;
- priority masjid audio with optional secondary Islamic Radio;
- saved favourites, scheduling, source volumes and audio-output recovery;
- discovery and selection of up to three MasjidBoard Live masjids;
- responsive TV/Monitor and portrait Appliance Display presentations;
- prayer, Jumu'ah, Daily Times and next-event information;
- supported community notices and optional shared Islamic content;
- last-known-good data during temporary upstream failures; and
- release packages for Linux ARM64 and AMD64.

Completed release details belong in [GitHub Releases](https://github.com/X-Calibre/MasjidPi/releases) and the relevant acceptance records under `docs/`.

## Current priorities

### Reliability and validation

- Continue longer-duration Raspberry Pi 3B and Pi 4 soak monitoring.
- Test HDMI disconnect and reconnect behavior.
- Validate suitable 512 MB ARM64 devices, particularly Pi Zero 2 W and Pi 3A+.
- Measure OS-level SD-card writes before changing journald, swap or other system services.
- Keep provider parsing defensive as MasjidBoard Live payloads evolve.

### Appliance product work

- Complete published-package validation of v1.6.0-rc.1 on Raspberry Pi 4 and confirm acceptable Pi 3B resource use.
- Replace temporary splash artwork with final MasjidFrame branding.
- Finalise the portrait enclosure, display, audio and power design.
- Validate the selected Waveshare display's audio path and physical controls.
- Investigate HDMI-CEC behavior on intended displays.
- Build a repeatable Raspberry Pi appliance image.

### MasjidBoard

- Add poster/media support only when retrieval, caching and presentation are safe.
- Improve Arabic/RTL coverage using additional real-world content.
- Revisit Maghrib/Iftar semantics when representative Ramadan source data is available.
- Extend structured upstream content only when field meaning and privacy are established.
- Add display layouts or preferences only where they materially improve appliance use.

### Listen

- Continue validating Radio endpoints and fallback behavior.
- Consider richer current-programme metadata where stations expose reliable information.
- Explore advanced audio controls only if they remain simple on appliance hardware.

## Platform targets

### Linux x86-64 appliance

Official AMD64 packages already support compatible 64-bit Linux systems. Further work may make repurposed laptops and small PCs easier to use as dedicated appliances through:

- automated display and audio setup;
- boot-to-appliance behavior;
- broader built-in, USB and Bluetooth audio validation; and
- a simplified installation profile for old computers.

### Windows x64

Windows desktop support remains exploratory. It would require Windows-compatible mpv management, IPC, persistent paths, service/startup behavior and packaging. Linux remains the preferred appliance platform.

### Hardware controls

Possible future hardware integrations include:

- OLED status displays;
- physical playback and volume controls;
- amplifier/EQ integration; and
- enclosure-specific status indicators.

### Home Assistant

Potential integration could expose:

- media-player controls and playback state;
- selected source and masjid status;
- volume and audio-output information;
- failure/recovery events;
- prayer-time sensors and triggers; and
- MQTT or a native Home Assistant integration.

## Architecture guardrails

MasjidPi remains one repository with shared Core functionality and two independently operable capabilities:

- **Core** — configuration, persistent state, APIs and platform integration
- **Listen** — stream discovery, priority playback, Radio, mpv and audio devices
- **Board** — MasjidBoard retrieval, caching, configuration and HDMI presentation

Listen must continue operating when Board or its upstream providers are unavailable, and Board must not depend on Listen.

## Project principles

MasjidPi should remain:

- simple to install and operate;
- reliable through network and device interruptions;
- lightweight enough for supported Raspberry Pi hardware;
- usable without a separate server;
- conservative with persistent writes;
- resilient when optional services fail; and
- focused on a dependable home-appliance experience.
