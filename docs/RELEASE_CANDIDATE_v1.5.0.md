# MasjidPi v1.5.0 Release Candidate Checklist

This document defines the release-candidate acceptance checks for MasjidPi v1.5.0.

The current planned release candidate is `v1.5.0-rc.4`. RC1 completed initial functional testing, RC2 carried the runtime and Web UI optimisation pass, RC3 addressed issues identified during RC2 hardware soak testing, and RC4 adds the responsive and touch-oriented Board appliance experience.

## Scope

v1.5.0 introduces secondary Islamic radio playback to MasjidPi Listen while preserving the selected masjid as the priority audio source.

The release also includes the outstanding enhancements and bug fixes already incorporated into `feature/radio-secondary-stream`.

## RC2 optimisation delta

In addition to the original v1.5.0 Radio scope, RC2 includes:

- one shared Listen status poll instead of independent polling by each control module;
- event-driven playback retry and startup-delay timers instead of a 50 ms idle loop;
- adaptive audio-device checks with fast missing-device recovery and slower healthy-state checks;
- one shared MasjidBoard display view supplying layout, theme, slide duration and renderer data;
- transactional preference mutations that preserve unrelated concurrent setting changes;
- concurrent retrieval of up to three selected boards while preserving configured display order;
- Radio timing controls grouped under Radio and a dedicated Audio tab;
- keyboard-accessible Listen tabs with session-only active-tab memory; and
- version footer retrieval through the component-independent version API.

## RC3 hardware-soak delta

RC3 adds:

- synchronization of the Audio output control with the persisted active backend device;
- preservation of the selected audio output across release upgrades;
- idempotent disabling of NetworkManager Wi-Fi power saving on Raspberry Pi appliances without restarting NetworkManager;
- reachable IPv4 and local hostname Web UI addresses in the installer completion summary;
- Jumu'ah presentation from the Thursday Islamic-date rollover through the Friday rollover;
- temporary IBM Plex Sans and Source Sans 3 display-preview options for typography evaluation;
- targeted one-second clock and countdown updates instead of rebuilding the complete prayer grid;
- structural Board rendering only when API data or minute/event state changes; and
- automatic saving for HDMI display settings and selected-masjid changes, while retaining explicit saving for location scope changes.

## RC4 Board and appliance delta

RC4 adds:

- responsive TV / Monitor rendering validated at 1366 × 768, 1920 × 1080, 2048 × 1152 and 3840 × 2160;
- dedicated naming for the 7-inch Appliance Display while retaining the existing internal layout value;
- clearer Listen power toggles and placement of the live Radio resume countdown in Now Playing;
- named source-transition notifications on both Board layouts, with a persistent Radio-resume countdown;
- a theme-aware touch control panel for the 7-inch Appliance Display;
- favourite Masjid and Radio selection, playback actions and all three Radio operating modes from the touch panel;
- Master, Masjid and Radio volume controls from the touch panel;
- all six saved display themes from the touch panel;
- consistent selected, playing and unavailable-control states; and
- automatic touch-panel closure after 60 seconds without user activity.

## Radio catalogue

The bundled South African Islamic radio catalogue contains the streams validated during development:

- Channel Islam International
- Radio Islam International
- Radio 786
- Voice of the Cape
- Radio Al Ansaar
- Markaz Sahaba Online Radio
- Sirius FM 105.7
- Salaamedia

The underlying stream research and validation notes are recorded in `docs/ISLAMIC_RADIO_STREAM_RESEARCH.md`.

## Listen behaviour to validate

### Masjid priority

- Select a masjid and a radio station.
- With the masjid offline, Radio should play when Radio operation permits it.
- When the masjid comes online, the transition from Radio to Masjid must be immediate.
- When the masjid goes offline, Radio must follow the configured post-masjid resume delay in scheduled mode.
- If the masjid comes back online during the delay, the pending Radio resume must be cancelled immediately.

### Source volumes

- Masjid and Radio have independent software-volume controls from 0% to 150%.
- Values above 100% are visibly identified as boosted in the Web UI.
- Switching between Masjid and Radio must restore the correct source-specific software volume.
- Master Volume remains a separate hardware-volume control where the selected ALSA device exposes a controllable mixer.

### Radio operation modes

Validate all three Radio modes:

**Play on Schedule**

- Follows the optional Radio daily schedule.
- Uses the configured 1–30 minute post-masjid resume delay.
- Masjid playback always interrupts immediately.

**Play Now**

- Starts Radio immediately when the masjid is offline.
- Overrides current Radio quiet time.
- Overrides a pending post-masjid resume delay.
- Ends at the next masjid-online event or the next Radio schedule boundary.
- Returns to scheduled operation after the override ends.

**Stop Radio**

- Radio remains stopped indefinitely.
- Masjid playback continues to operate normally.
- Radio remains stopped across masjid online/offline transitions.
- Radio remains stopped across service/appliance restart.
- Only Play on Schedule or Play Now re-enables Radio operation.

### Radio schedule

- Schedule can be disabled completely.
- Start and stop times are optional behavior controls when scheduling is enabled.
- Standard daytime windows such as 06:00–22:00 work.
- Overnight windows such as 22:00–02:00 work.
- Outside the configured window Radio remains silent in scheduled mode.
- Masjid playback is never blocked by Radio quiet time.

### Module power

Validate the persistent Masjid and Radio module power switches:

- Masjid ON / Radio ON works.
- Masjid ON / Radio OFF works.
- Turning Masjid OFF forces Radio OFF and stops the Listen controller completely.
- With Masjid OFF, Listen status must report `listening=false`, `radio_enabled=false` and `active_source=none`.
- Turning Radio ON while Masjid is OFF automatically powers Masjid ON first and informs the user.
- Re-enabling Masjid after a manual Masjid power-off does not independently power Radio back on unless the user also enables Radio.
- Module power state survives service restart and appliance reboot.

## Web UI acceptance

- Listen sub-tabs are Masjid, Radio and Audio.
- Masjid and Radio power controls are rendered as on/off toggle switches.
- The Masjid tab always shows a clear Selected Masjid summary even when nothing is playing.
- The Radio tab contains Radio operation mode, Radio Volume, Radio Station selection, Radio Resume Delay and Radio scheduling.
- The Audio tab contains audio output and Master Volume.
- Left and right arrow keys move between the Listen tabs, and the active tab is retained across a page refresh.
- Powered-off Radio controls remain visually stable with no polling-induced flicker.
- The Masjid catalogue can be freely scrolled without the one-second status poll snapping the list back to the selected masjid.
- The same non-snapping behavior applies to the Radio station list.
- Notifications correctly describe module power dependencies and mode changes.

## Persistence and restart

With representative settings configured, restart `masjidpi.service` and reboot the appliance. Confirm persistence of:

- selected masjid
- selected radio station
- Masjid volume
- Radio volume
- Radio resume delay
- Radio schedule and times
- persistent Radio mode (`schedule` or `stopped`)
- Masjid module power
- Radio module power
- Listen stopped/running state where applicable

`play_now` is intentionally temporary and must not be restored after restart; scheduled operation is restored instead.

## Stream validation

On Raspberry Pi hardware, allow each bundled radio station to play long enough to establish stable playback. Confirm:

- successful connection
- expected codec/container handling by mpv
- no immediate reconnect loop
- no unexpected CPU or memory regression

The detailed codec/bitrate measurements from development are recorded separately in `docs/ISLAMIC_RADIO_STREAM_RESEARCH.md`.

## Automated validation

Before merging to `main`:

```bash
cd ~/MasjidPi
git pull --ff-only origin enhancements/v1.5.0-rc4

make test

cd backend
go vet ./...

cd ..
git status -sb
```

The working tree must be clean and all tests/vet checks must pass.

The pull request CI run must also pass. CI runs Go formatting, vet and race-enabled tests, ShellCheck, installer tests, frontend JavaScript syntax validation and every JavaScript regression test in `tests/`.

## RC publication

After the feature branch is merged into `main` and `main` is validated:

```bash
git switch main
git pull --ff-only origin main

git tag -a v1.5.0-rc.4 -m "MasjidPi v1.5.0-rc.4"
git push origin v1.5.0-rc.4
```

The release workflow must publish ARM64 and AMD64 archives plus `SHA256SUMS`, and the GitHub release must be marked as a prerelease.

## Raspberry Pi acceptance

Install the actual `v1.5.0-rc.4` release artifact on the Raspberry Pi 4 test appliance rather than performing the final acceptance only from a source build.

Recommended minimum RC acceptance:

1. installation/upgrade succeeds and self-test passes;
2. existing Listen + Board configuration is preserved;
3. all three Radio operation modes behave correctly;
4. Radio-to-Masjid interruption is immediate;
5. Masjid-to-Radio resume delay works;
6. quiet-time and Play Now behavior work;
7. module power dependency and persistence work;
8. source volumes, including >100% boost, work;
9. service restart and full reboot restore the intended persistent state;
10. Board remains operational alongside Listen.
11. Listen controls remain synchronized without flicker and generate only one status request per polling cycle.
12. TV / Monitor and 7-inch Appliance Display layout, theme and slide-duration changes reach the HDMI display through the shared Board view.
13. Audio-device fallback and reconnection still work with adaptive checking.
14. Multiple selected boards refresh successfully and remain in configured order.
15. The persisted HDMI audio output remains selected after upgrade, service restart and reboot.
16. Wi-Fi power saving is disabled after installation and remains disabled after reboot.
17. The installer summary reports a reachable IPv4 Web UI URL and local hostname URL where available.
18. Jumu'ah replaces Dhuhr from the Thursday Islamic-date rollover until the Friday rollover.
19. Board display settings and selected-masjid actions save automatically; location scope changes still require explicit saving.
20. Cog/WPE RSS remains bounded during the soak test and does not repeat RC2's rapid prayer-grid allocation growth.
21. On a 600 × 1024 touch display, the swipe-up hint is visible, swipe-up opens the appliance control sheet, slide rotation pauses, and swipe-down/close resumes it.
22. The touch panel selects only favourited Masjids, selects available Radio stations, and changes Master, Masjid and Radio volumes without desynchronising the full Web UI.
23. Touch-panel Masjid start/Listen stop and Radio Scheduled Play, Play Now and Stop Radio actions preserve the established source-priority semantics.
24. All six themes can be applied and persisted from the Theme tab, and the complete control sheet remains legible in every theme.
25. Board source notifications name the active Masjid or Radio station, remain visible throughout a pending Radio-resume countdown, and do not appear merely because the Board loaded or reconnected.
26. The appliance control sheet closes after 60 seconds without activity and remains open while the user taps, scrolls, swipes, types or adjusts a control.
27. The TV / Monitor layout remains unclipped and proportionate at 1366 × 768, 1920 × 1080, 2048 × 1152 and 3840 × 2160.

If release-blocking defects are found, fix them on the development branch and publish a subsequent release candidate after repeating the automated and hardware acceptance gates.

## Soak monitoring

Run the release candidate on Raspberry Pi hardware for at least 24 hours, preferably 48 hours, while exercising normal Listen and Board behaviour.

Use the temporary systemd-based monitor documented in [`RC_SOAK_MONITORING.md`](RC_SOAK_MONITORING.md). It continues running across SSH disconnects and reboots and records service restart counts, resource usage, Raspberry Pi temperature/throttling state, Listen status, Board status summaries, and recent warnings/errors.

During the soak period confirm:

- no unexpected `masjidpi.service` or `masjidpi-display.service` restarts;
- no sustained process RSS growth;
- no unexpected memory or swap pressure;
- Raspberry Pi throttling remains `0x0`;
- no repeated audio reconnect loops or source oscillation;
- Listen and Board HTTP endpoints remain available;
- Board provider failures remain isolated and last-known-good fallback continues to operate.

Stop and disable the temporary monitor after RC acceptance as described in the monitoring guide.

## Stable promotion

Promote to `v1.5.0` only after the RC has passed Raspberry Pi hardware validation and no release-blocking regressions remain.
