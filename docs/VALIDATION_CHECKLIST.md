# MasjidPi Validation Checklist

This living checklist records completed v1.5.2 hardware evidence and validation that remains useful for later releases. Release-specific sign-off belongs in the matching release acceptance record.

## Completed for v1.6.0-rc.2

### Appliance network and controls

- [x] Network setup can be reopened from the configured Appliance Display.
- [x] The current Wi-Fi connection remains active until a replacement connection succeeds.
- [x] Visible and hidden replacement Wi-Fi networks connect successfully.
- [x] The Network tab shows the current IPv4 address and genuine network-issued FQDN.
- [x] The installer summary uses the genuine network-issued FQDN rather than an assumed `.local` address.
- [x] A downward drag from the control-sheet handle/header closes the sheet under Cog/WPE.
- [x] The control sheet opens without an initial focus highlight on the close button.
- [x] Highlighted controls use readable text in light themes.

### Board themes

- [x] Light Gold, Ivory, Sage, Sky and Rose render correctly.
- [x] All ten themes can be selected and persist across display refreshes.
- [x] The compact theme selector remains usable on the 7-inch Appliance Display.

## Completed for v1.6.0-rc.1

### First-run MasjidFrame setup

- [x] A clean Pi 4 with no saved Wi-Fi profile opens setup automatically.
- [x] Nearby 2.4 GHz networks scan and render correctly in Cog/WPE on DRM.
- [x] The embedded keyboard accepts letters, numbers and symbols and supports password visibility.
- [x] A visible password-protected network connects successfully.
- [x] A hidden password-protected network connects successfully after manual SSID entry.
- [x] Touch-native location and masjid pickers work without unsupported browser dropdown popups.
- [x] Location scoping, primary-masjid selection and initial timetable retrieval complete successfully.
- [x] The completion screen shows the current IPv4 address and DHCP/network-issued FQDN without inventing `.local`.
- [x] The normal Board opens after setup and opens directly on subsequent appliance startup.
- [x] Setup controls fit the 1024 × 600 portrait touchscreen, including the masjid save action.

## Completed for v1.5.2

### Boot and display

- [x] Portrait 7-inch appliance boots through correctly oriented Plymouth and Cog stages.
- [x] Landscape HDMI displays use the standard responsive profile.
- [x] Cog/WPE warm-up and Plymouth-to-Cog DRM handoff complete without persistent mode-setting failures.
- [x] Touch input is correctly calibrated in the portrait appliance profile.
- [x] Multiple cold boots and an abrupt-power recovery completed normally.
- [x] /boot/firmware remains read-only during normal operation and returns to read-only after installer/package write windows.

### Persistence and recovery

- [x] Atomic JSON state remains valid across service restarts and abrupt power loss.
- [x] Source update, staged activation, self-test and rollback behavior were exercised.
- [x] The mpv IPC socket is created under /run/masjidpi/mpv.sock.
- [x] Board preferences, Listen settings and component profile survive upgrades and restarts.
- [x] Last-known-good Board and daily-content caches survive upstream failures without being overwritten.

### Audio

- [x] A USB audio device connected after boot is discovered without restarting MasjidPi.
- [x] Playback falls back safely when the selected USB device is removed.
- [x] The saved device is restored automatically after reconnection.
- [x] Playback and services remain healthy throughout device loss and recovery.

### MasjidBoard

- [x] Detailed Jumu'ah schedules render in standard and appliance profiles.
- [x] Supported structured community cards render in both profiles.
- [x] Daily Ayah, Hadith and Sunnah preferences default correctly, persist and fall back to cache.
- [x] Dua after Adhan takes priority for five minutes and then releases the display.
- [x] Zawaal/Istiwaa warning behavior is correct at both interval boundaries.
- [x] Special-day Dhuhr is shown in Daily Times only when distinct from normal Dhuhr.
- [x] Eleven Daily Times entries fit on one landscape row.
- [x] Community content follows selected-masjid and priority ordering.

## Release artifact verification

For every stable release:

- [ ] Verify the release page contains ARM64 and AMD64 archives plus SHA256SUMS.
- [ ] Verify archive checksums before installation.
- [ ] Install the published ARM64 archive on a Raspberry Pi test appliance.
- [ ] Confirm /api/version reports the expected stable version.
- [ ] Confirm the installed component profile and persistent settings survive.
- [ ] Reboot and verify applicable boot, Listen and Board services.
- [ ] Confirm /boot/firmware returns to read-only where protection is enabled.

## Future platform validation

- [ ] Run a fresh long-duration soak of the final feature set on Raspberry Pi 3B.
- [ ] Validate Listen and Board memory headroom on Raspberry Pi Zero 2 W.
- [ ] Validate Listen and Board memory headroom on Raspberry Pi 3A+.
- [ ] Test HDMI disconnect/reconnect behavior.
- [ ] Test HDMI-CEC behavior on intended displays.
- [ ] Validate the Waveshare display audio path and physical volume controls.

The original ARMv6 Raspberry Pi Zero W and Raspberry Pi 1 are not current release targets and should not be listed as pending stable-release validation.

## Notes

- Do not remove `console=tty1`; physical testing showed that configuration caused a failed boot/display state.
- Do not mask `getty@tty1.service`; disabling it is sufficient for the dedicated Board boot experience.
- Cog must start only after Plymouth releases DRM.
- Temporary splash artwork remains scheduled for replacement with final product branding.
