# MasjidPi Validation Checklist

This document tracks physical and deferred validation work that should be revisited before future releases. Update the checkboxes as tests are completed.

## Required before merging `enhancement/boot-splash`

- [x] **7-inch appliance boot splash — portrait profile**
  - Quiet boot enabled without removing `console=tty1`.
  - Raspberry Pi firmware splash suppressed.
  - Plymouth branded splash displays in the correct portrait orientation.
  - Headless Cog/WPE warm-up completes while Plymouth remains visible.
  - Plymouth releases DRM before the display Cog starts.
  - Brief blank DRM handoff is acceptable.
  - Cog startup splash displays correctly.
  - MasjidBoard loads in the appliance profile.
  - Touch remains correctly calibrated and functional.
  - Three consecutive cold boots completed successfully without DRM `Cannot set mode` / `Permission denied` failures.
  - Temporary JEDI TECH SOLUTIONS logo is fully visible and correctly sized in both splash stages.

- [x] **Standard/landscape Raspberry Pi boot splash**
  - Test on a Raspberry Pi connected to a normal HDMI landscape display.
  - Confirm Plymouth logo is upright, centred, fully visible, and correctly sized.
  - Confirm headless Cog/WPE warm-up completes successfully.
  - Confirm no DRM `Cannot set mode` / `Permission denied` failures.
  - Confirm the Cog startup splash is upright and correctly sized.
  - Confirm MasjidBoard remains in the standard responsive landscape profile.
  - Confirm SSH and normal Board operation remain available after boot.
  - Perform multiple cold boots to confirm repeatability.

## Required before a future release

- [ ] **Power-interruption resilience**
  - Confirm `/boot/firmware` is mounted read-only during normal operation.
  - Run a source update and confirm the installer temporarily remounts it read-write, updates Plymouth/initramfs, and restores read-only mode.
  - Run a package operation that updates boot firmware and confirm the APT/DPKG hooks restore read-only mode.
  - Confirm the mpv IPC socket is created under `/run/masjidpi` and playback works normally.
  - Change and clear persisted settings, then confirm all state files remain valid after restart.
  - Perform one controlled abrupt power-loss test, then review filesystem recovery, state integrity, service startup, playback and Board operation.

- [ ] **USB audio hot-plug restoration**
  - Baseline physical audio-device discovery has already passed.
  - Connect a supported USB audio device after boot and refresh devices.
  - Select and save the USB audio output.
  - Disconnect the device and verify MasjidPi handles its disappearance safely.
  - Reconnect it and verify the saved device can be rediscovered/restored as intended.
  - Verify playback after restoration.

- [ ] **Raspberry Pi Zero Listen-only validation**
  - Install the current release/branch in Listen-only mode on the intended Raspberry Pi Zero hardware.
  - Confirm architecture/package compatibility, including the intended armv6/armv7 target.
  - Confirm `masjidpi.service` starts normally.
  - Confirm catalogue refresh and stream discovery.
  - Confirm masjid stream playback and radio playback.
  - Confirm audio output discovery and volume behaviour.
  - Confirm WebUI responsiveness is acceptable for the hardware.
  - Run an appropriate stability/soak test.

- [ ] **Post-change Raspberry Pi soak test**
  - Run the current candidate on the long-soak Raspberry Pi.
  - Monitor service restarts, memory, CPU temperature, throttling, disk usage, network stability, and mpv behaviour.
  - Exercise scheduled radio start/stop and masjid interruption/resume behaviour.
  - Confirm the display service remains stable where Board is installed.
  - Review persistent journal logs before release sign-off.

## Deferred hardware / feature investigation

- [ ] **HDMI-CEC physical testing**
  - Verify CEC availability on the Raspberry Pi HDMI output.
  - Test whether the connected Waveshare 7-inch HDMI LCD responds to relevant CEC power/display commands.
  - If unsupported by the Waveshare panel, document that limitation rather than treating it as a MasjidPi fault.

- [ ] **Waveshare HDMI audio physical testing**
  - Confirm whether the 7-inch display exposes an HDMI audio sink to Raspberry Pi OS.
  - Test actual audio playback through the display's supported audio path/output.
  - Test any physical volume controls or display-side volume behaviour that are available.
  - Determine whether this is suitable for the planned appliance speaker design.

## Notes

- Do **not** remove `console=tty1` as part of quiet-boot work; physical testing showed that configuration caused a failed boot/display state.
- Do **not** mask `getty@tty1.service`; disabling it is sufficient for the dedicated Raspberry Pi Board boot experience.
- Cog's direct DRM backend must start only after Plymouth releases DRM. Starting it while Plymouth owns DRM produced `Cannot set mode (Permission denied)` and a persistent blank display.
- The current temporary splash artwork is not final appliance branding and is expected to be replaced later.
