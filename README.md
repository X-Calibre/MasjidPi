# MasjidPi

MasjidPi is an open-source internet radio designed specifically for listening to live audio streams from masājid around the world.

The project is intentionally lightweight and designed to run on low-powered Raspberry Pi hardware while providing a simple web interface for selecting and controlling mosque audio streams.

---

## Screenshot

![MasjidPi Web Interface](docs/images/masjidpi-ui.png)

---

# About this project

This project has a unique story.

MasjidPi is being developed by someone with **no prior software development experience**. Every line of code has been written through an iterative collaboration with AI, while the project vision, feature decisions, testing, and overall direction remain human-driven.

The goal is to demonstrate that modern AI tools can empower people without formal programming backgrounds to build high-quality open-source software by combining domain knowledge, curiosity, careful testing, and incremental development.

Rather than using AI to generate an application in one step, every feature is designed, implemented, tested and refined through many small iterations.

If the project inspires someone else to learn software development, then it has already been a success.

---

# Vision

A Raspberry Pi connected to speakers that allows users to:

- Browse live masjid streams
- Listen to live broadcasts
- Save favourite masājid
- Automatically reconnect after network outages
- Resume playback after power loss
- Be managed entirely through a web browser

MasjidPi is designed to run comfortably on hardware as small as the Raspberry Pi Zero 2 W while remaining extensible for future features.

---

# Project Philosophy

MasjidPi follows a few simple principles.

- Keep the software lightweight.
- Prioritise reliability over complexity.
- Design for Raspberry Pi first.
- Build features only when they solve a real problem.
- Prefer simple, maintainable code over clever code.
- Keep the application fully self-contained whenever practical.
- Minimise external dependencies.

---

# Current Features

- Web-based player interface
- Live playback status
- Volume control (0–125%)
- MPV audio playback
- Local stream catalogue
- Automatic catalogue generation from LiveMasjid
- Play streams by catalogue selection
- Remember the last selected stream
- Optional automatic playback on startup
- Responsive mobile-friendly interface

---

# Planned Features

- One-click catalogue update
- Reload catalogue without restarting
- Searchable catalogue
- Favourite masājid
- Graceful handling of offline streams
- OLED display support
- Push-button controls
- Audio equaliser
- Multi-language interface

---

# Installation

## Requirements

Currently MasjidPi requires:

- Linux
- Go 1.26 or newer
- MPV
- Git

On Debian or Raspberry Pi OS:

```bash
sudo apt update
sudo apt install git mpv
```

Install Go by following the official Go installation instructions.

---

## Clone the repository

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
```

---

## Run MasjidPi

From the project root:

```bash
make run
```

or manually:

```bash
cd backend
go run ./cmd/masjidpi
```

If everything starts correctly you should see output similar to:

```text
Starting application
Configuration loaded
Loaded stream catalogue
Connected to MPV
Starting HTTP server :8080
```

---

## Open the Web Interface

Open your browser and navigate to:

```
http://localhost:8080
```

If running on another machine, replace `localhost` with the Raspberry Pi's IP address.

---

## Updating the Stream Catalogue

MasjidPi generates its catalogue from LiveMasjid.

At the moment this is performed during development.

Future versions will include a one-click **Update Catalogue** button in the web interface.

---

# Development Status

MasjidPi is currently in active development.

The project is functional, but many planned features are still under development.

Breaking changes may occur until version 1.0.

---

# Contributing

Bug reports, feature requests and pull requests are welcome.

Please see:

- CONTRIBUTING.md

---

# License

MasjidPi is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

See the LICENSE file for details.