# SD Card Write Policy

MasjidPi is intended to run continuously on Raspberry Pi hardware where the primary persistent storage may be an SD card. Avoiding unnecessary flash writes is therefore an appliance reliability requirement.

## Principles

- Runtime state should remain in memory whenever persistence is not required.
- User settings and required application state should only be persisted when they change.
- Routine operational logging must not generate high-frequency persistent writes.
- Temporary/runtime files should live outside persistent application state where practical.
- Persistent data updates should be atomic so power loss cannot leave corrupt state.
- Optimisations must not sacrifice the ability to diagnose real failures.

## Implemented on `research/sd-card-write-optimisation`

### Volume state

`volume.json` is no longer rewritten when the requested device volume already matches the persisted value. This prevents redundant flash writes from repeated identical volume updates.

Web UI volume changes are now split into two operations: intermediate slider movement changes the live hardware volume without persistence, while the final slider value is persisted once when the control change completes. A short client-side debounce also prevents excessive API requests while the slider is being dragged.

### LiveMasjid MQTT logging

Individual MQTT events are now logged at `DEBUG` rather than `INFO`. Connection and subscription failures remain visible at `WARN`/`ERROR` levels.

### systemd journal protection

The MasjidPi service now applies systemd journal rate limiting so a persistent failure cannot generate an uncontrolled burst of journal writes.

## Still to audit

- Other persistent state writers for unchanged-value writes.
- Catalogue update write patterns and atomicity.
- systemd/journald persistent-vs-volatile storage policy for the appliance image.
- OS-level services that may write frequently to the SD card.
- Measurement tooling for filesystem writes over long-running tests.

## Validation goal

A normal MasjidPi installation should be able to run continuously without generating meaningful background SD-card writes apart from intentional configuration/state changes, catalogue refreshes, package/update activity, and other explicitly expected persistence.
