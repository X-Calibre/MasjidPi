# MasjidBoard documentation

The maintained description of the implemented subsystem is split across:

- `MASJIDBOARD-ARCHITECTURE.md` — subsystem boundaries and dependencies
- `MASJIDBOARD-IMPLEMENTATION-STATUS.md` — implemented capability status
- `MASJIDBOARD-DISPLAY-BOUNDARY.md` — display/API ownership boundary
- `MASJIDBOARD-DISPLAY-RUNTIME.md` — HDMI appliance runtime
- `MASJIDBOARD-DISPLAY-LAYOUT.md` — current display layouts
- `MASJIDBOARD-LAST-KNOWN-GOOD-CACHE.md` — outage and cache behaviour
- `MASJIDBOARD-LIVE.md` — Core timetable and optional Premium enrichment contract

The remaining design and discovery documents record implementation research and
decisions. Captured `.html` and `.json` responses are provider research fixtures;
they are not shipped in release packages or used as application runtime data.

When behaviour changes, update the maintained documents above first. Retain a
research document only when it explains a provider-specific constraint or a
decision that is not evident from the current implementation.
