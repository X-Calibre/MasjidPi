# MasjidPi Validation Checklist

This document tracks physical and deferred validation work that should be revisited before future releases. Update the checkboxes as tests are completed.

## Completed v1.5.2 boot validation

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

- [x] **Power-interruption resilience**
  - `/boot/firmware` remains mounted read-only during normal operation and after reboot.
  - A source update temporarily enabled boot-firmware writes, regenerated Plymouth/initramfs assets, and restored the read-only mount.
  - A controlled `raspi-firmware` reinstall copied firmware, kernel and initramfs files successfully; the APT/DPKG hooks restored `/boot/firmware` read-only afterward.
  - The mpv IPC socket is created by systemd under `/run/masjidpi/mpv.sock`, and playback works normally.
  - Board theme and slide duration, Radio volume and Radio operating mode were changed, restarted, restored and restarted again; every persisted JSON file remained valid.
  - A controlled abrupt power loss recovered through normal ext4 orphan cleanup with no FAT repair, filesystem errors or I/O errors.
  - Both portrait splash stages, Emerald Board startup and Radio playback operated normally after abrupt-power recovery; all services were healthy with zero restarts and no failed units.

- [x] **USB audio hot-plug restoration**
  - A Logitech H390 USB headset connected after boot was discovered as `alsa/plughw:CARD=Headset,DEV=0` without restarting MasjidPi.
  - The new USB output was selected and saved, exposed hardware-volume support, and played Radio audio normally.
  - Disconnecting it during playback left the saved preference intact, marked the device unavailable, and safely fell back to automatic output while playback continued.
  - Reconnecting it caused MasjidPi to rediscover and automatically restore the saved device and working audio.
  - The service remained healthy with zero restarts throughout; the original `alsa/plughw:CARD=Headphones,DEV=0` output and scheduled Radio mode were restored after testing.

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

- [ ] **Published v1.5.2 release-candidate package**
  - Install the published ARM64 archive on the Raspberry Pi 4 test appliance rather than relying only on a source installation.
  - Confirm the archive checksum and reported version are `v1.5.2-rc.2`.
  - Confirm the existing Listen + Board profile, settings, selected audio output and cached data survive the upgrade.
  - Confirm both splash stages, application services, Board display and Listen playback operate normally after reboot.
  - Confirm `/boot/firmware` returns to read-only after installation.

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
