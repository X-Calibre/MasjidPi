# MasjidPi

MasjidPi is an open-source internet radio designed specifically for listening to live audio streams from masājid around the world.

The project is intentionally lightweight and designed to run on low-powered Raspberry Pi hardware while providing a simple web interface for selecting and controlling mosque audio streams.

---

## About this project

This project has a unique story.

MasjidPi is being developed by someone with **no prior software development experience**. Every line of code has been written through an iterative collaboration with AI, while the project vision, feature decisions, testing, and overall direction remain human-driven.

The goal is to demonstrate that modern AI tools can empower people without formal programming backgrounds to build high-quality open-source software by combining domain knowledge, curiosity, careful testing, and incremental development.

Every feature is developed in small, testable steps, with an emphasis on clean architecture, maintainability, and learning along the way.

## Vision

A Raspberry Pi connected to speakers that allows users to:

- Browse masjid streams
- Listen to live broadcasts
- Save favourites
- Automatically reconnect after network outages
- Resume playback after power loss
- Be managed entirely through a web browser

MasjidPi is designed to be lightweight enough to run on a Raspberry Pi Zero 2 W while remaining extensible through future modules.

## Project Philosophy

MasjidPi follows a few simple principles:

- Keep the software lightweight.
- Prioritise reliability over complexity.
- Design for Raspberry Pi first.
- Build features only when they solve a real problem.
- Prefer simple, maintainable code over clever code.
- Keep the application fully self-contained whenever practical.

## Features

Current functionality includes:

- Web-based player interface
- Live player status
- Volume control (0–125%)
- Automatic playback using MPV
- Stream catalogue
- Play streams by catalogue ID
- Automatic restoration of the last selected masjid
- Optional auto-play on startup
- Responsive mobile-friendly interface

## Planned Features

- One-click catalogue updater
- LiveMasjid catalogue scraper
- Searchable stream catalogue
- Favourite masājid
- Physical hardware controls
- OLED display support
- Rotary encoder support
- Audio equaliser
- Multi-language interface

## Project Goals

- Lightweight
- Reliable
- Offline-friendly
- Headless-first
- Privacy-respecting
- Community-driven
- 100% Open Source

## License

MasjidPi is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).
