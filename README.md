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

Rather than using AI to generate an application in one step, every feature is designed, implemented, tested and refined through many small iterations.

The goal is not only to build a reliable Raspberry Pi streaming appliance, but also to demonstrate how AI can help people learn software engineering through practical, real-world projects.

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
- Current stream name and relay URL
- Volume control (0–125%)
- MPV audio playback
- Local JSON stream catalogue
- Automatic catalogue generation from LiveMasjid
- One-click catalogue update
- Hot reload of the catalogue without restarting
- Automatically generated relay URLs
- Play streams by catalogue selection
- Remember the last selected stream
- Optional automatic playback on startup
- Responsive mobile-friendly interface

---

# Planned Features

- Detect offline streams
- Gracefully handle unavailable streams
- Searchable catalogue
- Favourite masājid
- OLED display support
- Push-button controls
- Audio equaliser
- Multi-language interface

---

# Installation

MasjidPi provides an automated installer that downloads all required dependencies, installs Go (if necessary), builds the application, and performs a self-test.

## Supported Platforms

The installer currently supports:

- Debian 12+
- Ubuntu 24.04+
- Raspberry Pi OS Bookworm

Additional Linux distributions may work but have not yet been officially tested.

---

## Clone the repository

```bash
git clone https://github.com/X-Calibre/MasjidPi.git
cd MasjidPi
```

---

## Run the installer

The installer requires sudo privileges to install system packages.

```bash
chmod +x install.sh
./install.sh
```

The installer will automatically:

- Detect your operating system
- Install all required system packages
- Install Go (if it is not already installed)
- Build MasjidPi
- Verify the installation by running a self-test
- Confirm that audio playback is available

A successful installation will end with output similar to:

```text
[ OK ] Dependencies installed.
[ OK ] Go 1.26.5 installed.
[ OK ] Build complete.
[ OK ] Application started successfully.
[ OK ] Audio device detected.
[ OK ] Installation completed.
```

---

## Running MasjidPi

After installation, start the application from the project directory:

```bash
make run
```

or manually:

```bash
cd backend
./masjidpi
```

When MasjidPi starts successfully you should see:

```text
Starting application
Configuration loaded
Loaded stream catalogue
Connected to MPV
Starting HTTP server
```

---

## Open the Web Interface

Open your web browser and navigate to:

```
http://localhost:8080
```

If MasjidPi is running on another computer or Raspberry Pi, replace `localhost` with its IP address.

For example:

```
http://192.168.1.100:8080
```
---

# Updating the Stream Catalogue

MasjidPi builds its local catalogue directly from the LiveMasjid website.

To refresh the catalogue:

1. Open the web interface.
2. Click **Update Catalogue**.
3. MasjidPi downloads the latest catalogue from LiveMasjid.
4. The local catalogue is regenerated and reloaded automatically.
5. No application restart is required.

---

# Current Limitations

MasjidPi currently assumes that every stream in the catalogue is available.

Support for detecting offline streams and providing user-friendly playback errors is currently under development.

---

# Development Status

MasjidPi is currently in active development.

The application is already suitable for everyday listening, while additional reliability improvements and Raspberry Pi-specific features continue to be added.

Breaking changes may still occur before the first stable release.

---

# Contributing

Bug reports, feature requests and pull requests are welcome.

Please see:

- CONTRIBUTING.md

---

# License

MasjidPi is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

See the LICENSE file for details.