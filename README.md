# MasjidPi

MasjidPi is a lightweight home appliance for live masjid audio and prayer-time information.

It can run either capability independently or install both together:

- **Listen** prioritises a selected LiveMasjid stream and can play Islamic Radio while the masjid is offline.
- **Board** displays prayer times, Jumu'ah schedules and supported community content from MasjidBoard Live on an attached HDMI display.

> New to MasjidPi? Start with the [User Guide](docs/USER_GUIDE.md).

## Features

### Listen

- Search LiveMasjid streams and save favourites
- Give the selected masjid immediate priority over Radio
- Resume Radio after a configurable delay or within a daily schedule
- Use scheduled, immediate or stopped Radio modes
- Set independent Masjid and Radio volumes
- Select an ALSA audio output and recover from device interruptions
- Restore saved settings and unattended playback after reboot

### Board

- Select and order up to three MasjidBoard Live masjids
- Use a responsive TV/Monitor layout or the portrait 7-inch Appliance Display
- Show prayer times, next-event countdowns, Daily Times and detailed Friday Jumu'ah schedules
- Show supported announcements, programmes, funeral, Nikah, Eid, contribution and other community cards
- Show optional Daily Ayah, Hadith, Sunnah, Islamic Economic Indicators and Dua after Adhan
- Highlight the published Zawaal/Istiwaa interval
- Continue with last-known-good timetable data during temporary upstream outages
- Use six colour themes and integrated touch controls on the Appliance Display

Content availability depends on what each upstream masjid publishes.

## Screenshots

| Listen | Board |
|---|---|
| ![MasjidPi Listen interface](docs/images/masjidpi-listen-v1.3.0.png) | ![MasjidBoard HDMI display](docs/images/masjidboard-display-v1.3.0.png) |

![MasjidBoard configuration interface](docs/images/masjidboard-configuration-v1.3.0.png)

## Install

MasjidPi supports 64-bit ARM Linux and 64-bit x86 Linux. Raspberry Pi 3B and Raspberry Pi 4 are the production-validated appliance platforms.

Install the latest stable release:

```bash
curl -fsSL https://raw.githubusercontent.com/X-Calibre/MasjidPi/main/scripts/install-latest.sh | sudo bash
```

The installer prompts for one of three profiles:

1. Listen
2. Board
3. Listen + Board

After installation, open the configuration interface from another device on the same network:

```text
http://<masjidpi-ip-address>:8080
```

See the [Installation Guide](docs/INSTALL.md) for supported systems, installer behavior, updates and troubleshooting. Hardware compatibility and measured Raspberry Pi performance are documented in the [Hardware Guide](docs/HARDWARE.md).

## Documentation

- [User Guide](docs/USER_GUIDE.md) — Listen, Radio, Board and troubleshooting
- [Installation Guide](docs/INSTALL.md) — installation, updates and component profiles
- [Hardware Guide](docs/HARDWARE.md) — supported platforms and measured performance
- [Roadmap](ROADMAP.md) — current priorities and future work
- [MasjidBoard technical documentation](docs/MasjidBoard/README.md)
- [SD-card write policy](docs/SD_CARD_WRITE_POLICY.md)

## Development

MasjidPi uses a Go backend and a browser-based frontend.

Run the automated tests:

```bash
make test
```

For source-based development installation:

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
sudo ./scripts/install.sh --source
```

Contributions are described in [CONTRIBUTING.md](CONTRIBUTING.md).

## Release

**Current stable release: v1.5.2**

v1.5.2 improves appliance startup and resilience and substantially expands MasjidBoard content and presentation. See the [v1.5.2 acceptance record](docs/RELEASE_CANDIDATE_v1.5.2.md) and [GitHub Releases](https://github.com/X-Calibre/MasjidPi/releases) for details and downloads.

## Data sources and acknowledgements

MasjidPi was inspired by the [eBilal project](https://github.com/Muslims-in-IT/ebilal).

MasjidPi uses [LiveMasjid](https://www.livemasjid.com/) for live masjid streams, [MasjidBoard Live](https://masjidboardlive.com/) for timetable and community data, and [Jamiatul Ulama South Africa](https://www.jamiatsa.org/category/islamic-economic-indicators/) for optional Islamic Economic Indicators.

MasjidPi is independent and is not affiliated with or endorsed by those projects or services.

## Licence

MasjidPi is licensed under the [GNU Affero General Public License v3.0](LICENSE).
