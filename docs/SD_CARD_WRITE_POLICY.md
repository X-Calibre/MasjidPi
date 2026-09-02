# SD Card Write Policy

MasjidPi is intended to run continuously on Raspberry Pi hardware where the primary persistent storage may be an SD card. Avoiding unnecessary flash writes is therefore an appliance reliability requirement.

## Principles

- Runtime state should remain in memory whenever persistence is not required.
- User settings and required application state should only be persisted when they change.
- Routine operational logging must not generate high-frequency persistent writes.
- Temporary/runtime files should live outside persistent application state where practical.
- Persistent data updates should be atomic so power loss cannot leave corrupt state.
- Optimisations must not sacrifice the ability to diagnose real failures.

## Implemented

### Volume state

`volume.json` is no longer rewritten when the requested device volume already matches the persisted value. Intermediate slider movement changes live hardware volume without persistence, while the final slider value is persisted once when the control change completes.

Persisted volume state is loaded once and cached in memory, eliminating repeated filesystem reads during normal runtime.

### Other persistent state

Preferences, selected audio device, last playback stream, and favourites are loaded once and cached in memory. Unchanged state updates are discarded in memory without rereading or replacing the state file.

Changed persistent state is written using temporary files and atomic replacement. File contents are flushed before replacement, and the containing directory is flushed after rename or deletion so the metadata change survives sudden power loss.

The playback clear operation treats an already-empty state as a no-op.

### Catalogue updates

Automatic catalogue refresh runs weekly, with manual refresh available from the Web UI.

The LiveMasjid source HTML is downloaded and parsed in memory rather than being persisted as `page.html`. Generated catalogue content is compared with the existing `catalogue.json`; unchanged content produces no persistent write. Changed catalogues are written atomically.

After refresh, the in-memory stream store is updated directly instead of rereading `catalogue.json` from disk.

Unused legacy catalogue JavaScript download code has been removed.

### LiveMasjid MQTT logging

Individual MQTT events are logged at `DEBUG` rather than `INFO`. Connection and subscription failures remain visible at `WARN`/`ERROR` levels.

### systemd journal protection

The MasjidPi service applies systemd journal rate limiting so a persistent failure cannot generate an uncontrolled burst of journal writes.

### Read-only boot firmware

On Raspberry Pi installations, `/boot/firmware` is remounted read-only before normal MasjidPi operation. Board installations explicitly open and close a controlled writable window for splash and initramfs changes. APT/DPKG hooks do the same around package operations that may update the kernel or initramfs, then flush and restore the read-only mount.

This protects the FAT boot filesystem, which cannot provide the same power-loss guarantees as the journalled EXT4 root filesystem.

### Volatile runtime files

The mpv IPC socket now lives in `/run/masjidpi`, backed by volatile runtime storage. systemd creates the directory for the service on every boot. Existing installations using MasjidPi's previous `/tmp/masjidpi.sock` default are migrated automatically, while custom socket paths are preserved.

## Application-level audit status

The application-level filesystem audit is complete. Production runtime code no longer contains direct `os.WriteFile()` state writes or unnecessary persistent download files. Remaining temporary-file and rename operations correspond to intentional atomic persistence of changed state.

Normal steady-state operation should not generate meaningful application data writes while streams are playing, MQTT status traffic is received, the Web UI is polled, or the catalogue refresh timer is idle.

## Remaining OS-level work

OS-level behaviour will be assessed separately after application-level optimisation. This includes:

- measurement-led systemd/journald retention policy for a dedicated appliance installation;
- swap behaviour;
- OS services and timers that may write frequently to the SD card;
- measurement of real filesystem writes during long-running Raspberry Pi tests.

OS-level changes should be based on measurement and should be applied only where appropriate for a dedicated MasjidPi appliance rather than changing normal Linux behaviour unnecessarily.

## Validation goal

A normal MasjidPi installation should be able to run continuously without generating meaningful background application writes apart from intentional configuration/state changes, changed catalogue refreshes, package/update activity, and other explicitly expected persistence.
