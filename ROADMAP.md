# MasjidPi Roadmap

MasjidPi is a lightweight Raspberry Pi home appliance that keeps people connected to their masjid through live audio and masjid information.

The roadmap focuses on the current product, remaining production work, and planned future features.

---

## Current Status

**v1.0.11 is the current release.**

The current release provides:

- LiveMasjid stream catalogue and updates
- Weekly automatic catalogue refresh plus manual refresh from the Web UI
- Masjid search and discovery
- Favourites
- Web UI
- Stream playback via MPV
- Volume and audio-device control
- Persistent settings
- Playback, network and MPV recovery
- Audio-device recovery
- Raspberry Pi runtime and systemd integration
- Pre-built release packages
- Bundled installer
- Upgrade installation preserving configuration and runtime data
- Safe update workflow
- Safer update and recovery behaviour
- Validation of a newly installed version before considering an update successful
- Recovery/rollback path if an update fails
- First-run and production installation workflow improvements
- Application-level SD-card write optimisations, including in-memory state caching, reduced logging, atomic state writes, and in-memory catalogue processing
- Security dependency update for the Go network/HTML parsing stack

The v1.0.8 release was tested successfully on a Raspberry Pi 3 B and established the current stable appliance baseline. v1.0.10 added application-level SD-card write and filesystem-activity optimisations, and v1.0.11 adds the dependency security fix while retaining that baseline.

The Updates & Recovery work and application-level SD-card write optimisation are now considered complete. Remaining reliability work should prioritise measured OS-level appliance behaviour and simplicity rather than unnecessary configuration options.

---

## Remaining Work

There are no outstanding application-level Updates & Recovery or SD-card write optimisation items at this stage.

OS-level SD-card behaviour should be measured on Raspberry Pi hardware before deciding whether appliance-specific changes to journald, swap, or other system services are justified.

Future work should otherwise focus on planned capabilities and further reliability improvements where they provide clear user value.

---

## Planned Architecture

### MasjidPi Core / Listen / Board

MasjidPi v2 will evolve from the current audio-focused application into a modular home appliance built around a shared MasjidPi core and independent product capabilities.

The repository will remain a single MasjidPi repository initially, with the code organised conceptually as:

```text
MasjidPi/
├── core/
├── listen/
├── board/
├── cmd/
├── frontend/
├── configs/
├── hardware/
├── scripts/
└── docs/
```

#### MasjidPi Core

Core will contain functionality shared by Listen and Board, such as:

- Configuration
- Persistent state
- Masjid selection and shared masjid data
- Catalogue and shared data access
- Networking and common service behaviour
- Common API/client functionality
- Shared Raspberry Pi/platform integration where appropriate

Core should provide reusable foundations rather than becoming a general-purpose dumping ground for unrelated code.

#### MasjidPi Listen

Listen will contain the audio-specific functionality:

- Stream playback
- MPV integration
- Volume control
- Audio-device handling
- Playback recovery and reconnection
- Audio-specific UI and controls

#### MasjidPi Board

Board will contain the MasjidBoard-specific functionality:

- Prayer times
- Jumu'ah times
- Next-event and countdown information
- MasjidBoard data retrieval and caching
- Display-oriented UI
- External HDMI display support
- Future display/hardware integration

Listen and Board must remain independent capabilities. Neither subsystem should require the other to operate. In particular, audio playback must continue operating if MasjidBoard is unavailable.

The intended high-level architecture is:

**MasjidPi Core → Listen**

and independently:

**MasjidPi Core → Board**

The application/UI layer may combine information from both capabilities, but Board must not directly control Listen or playback.

The initial implementation should remain a single repository. Core should only be extracted into a separate repository if a future technical need justifies doing so.

This restructuring is a v2 architectural target and should not unnecessarily disrupt the stable v1.x release.

---

## Planned Features

### Raspberry Pi Appliance Image

- Deployable Raspberry Pi appliance image
- Preconfigured operating environment
- Repeatable image-based installation workflow

### Linux x86-64 Appliance / Old Laptop Support

MasjidPi should eventually support repurposing an older x86-64 laptop or PC as a dedicated MasjidPi appliance.

The preferred approach is to use a lightweight Linux operating system rather than requiring Windows. This would allow the machine to boot directly into a simple, appliance-like MasjidPi environment without requiring users to manage a normal desktop operating system.

Potential capabilities:

- Support for standard Linux x86-64 hardware
- Dedicated MasjidPi appliance installation/image
- Automatic MasjidPi startup at boot
- Audio output through built-in, USB, Bluetooth, or other supported Linux audio hardware
- Optional HDMI display output for MasjidBoard
- Ability to run MasjidPi Listen and MasjidPi Board together on the same machine
- Reuse of the same MasjidPi Core, Listen, and Board architecture as Raspberry Pi deployments
- Simple installation suitable for repurposing older laptops and PCs

This should be treated as a platform/appliance target rather than a separate MasjidPi product. The application should remain shared across Raspberry Pi and x86-64 Linux, with platform-specific behaviour isolated where necessary.

### Windows x64 Desktop Support

Potential future support for running MasjidPi directly on Windows x64 as a desktop application.

This is primarily intended as a convenience, development, testing, and desktop-listening platform rather than the preferred dedicated appliance platform. A Linux-based x86-64 appliance remains the preferred approach for repurposing an old laptop as a dedicated MasjidPi device.

Potential capabilities:

- Windows x64 build
- Bundled or managed MPV runtime
- Windows-compatible playback IPC and filesystem handling
- Simple desktop installation
- Local Web UI access

### Hardware & Advanced Audio

- OLED display
- Push-button controls
- Hardware playback controls
- Display current station and playback state
- Audio equaliser / EQ processing
- Advanced audio controls

### MasjidBoard

MasjidBoard will be a separate subsystem from audio streaming and playback.

The active `research/masjidboard-live` implementation has now established the main functional shape of the feature, including:

- location hierarchy and scoped masjid discovery;
- selection and ordering of up to three masjids;
- persisted MasjidBoard location scope and selection;
- live prayer and Jumu'ah timetable retrieval from the MasjidBoard Live Core board source;
- normalized provider data and resilient Jumu'ah label handling;
- a dedicated MasjidBoard configuration WebUI;
- a separate read-only HDMI/browser display page;
- a responsive default one-, two- or three-board comparison layout;
- Friday Jumu'ah replacing Dhuhr;
- chronological Jumu'ah event presentation with source labels preserved;
- a per-board countdown to the next visible timetable event, including overnight rollover to the next day's Fajr;
- 30-minute automatic timetable refresh plus manual timetable refresh;
- per-board last-known-good cache fallback during provider outages;
- stale/current recovery behaviour validated in both automated tests and a real runtime outage test;
- reduced cache writes when timetable data is unchanged; and
- continued independence from audio playback when MasjidBoard Live is unavailable.

The current display should be treated as the **default layout**, not the only long-term layout. Alternative user-selectable layouts are intentionally deferred until the overall MasjidBoard appliance path is more production-ready. They should consume the same normalized MasjidBoard display API so presentation choices remain separate from provider, caching and timetable logic.

Before this work is considered production-ready, the research implementation still needs a focused pre-integration pass covering:

- installer/upgrade handling for all MasjidBoard runtime directories and persisted files;
- production ownership and permissions for hierarchy, scope, catalogue, selection and per-board cache data;
- first-run behaviour with no MasjidBoard configuration and no existing cache;
- combined Listen + Board validation on Raspberry Pi hardware, including reboot/restart and resource usage;
- final user-facing error-message cleanup while retaining diagnostic detail for troubleshooting;
- a full Listen/audio regression pass to confirm the stable playback appliance path remains unaffected; and
- final merge/release planning for the research branch.

Future MasjidBoard enhancements may include:

- additional user-selectable display layouts;
- layouts exposing selected astronomical/calculation times;
- future OLED/hardware displays;
- richer announcement/media content where useful;
- additional display preferences; and
- broader appliance/platform validation.

MasjidBoard must remain independent of the playback subsystem. Audio playback must continue operating if MasjidBoard is unavailable.

The intended architecture is:

**MasjidBoard → Application/UI**

and independently:

**Stream Catalogue → Playback → MPV → Audio Device**

The application/UI layer may combine information from both systems, but MasjidBoard must not directly control playback.

### Home Assistant Integration

Potential capabilities:

- Home Assistant media-player entity
- Play / stop control
- Volume control
- Current stream and masjid information
- Playback state
- Audio output device
- Stream and playback status
- Playback failure and recovery events
- Audio-device loss and recovery events
- MQTT integration
- Optional native Home Assistant integration
- Optional prayer-time and masjid sensors
- Prayer-related events and triggers

---

## Project Principles

MasjidPi should remain:

- A simple home appliance for people who want to stay connected to their masjid
- Lightweight enough for Raspberry Pi hardware
- Simple to install and operate
- Reliable during network interruptions
- Usable without a separate server
- Focused on listening to mosque streams while expanding to masjid information
- Easy to maintain and update
- Independent of browser-local configuration
- Modular, with Listen and Board capabilities separated from shared Core functionality
- Resilient when optional external services are unavailable
- Portable across appropriate hardware platforms without unnecessarily duplicating the application
- Focused on appliance simplicity rather than exposing unnecessary configuration

The project should prioritise a reliable home appliance experience over unnecessary complexity.
