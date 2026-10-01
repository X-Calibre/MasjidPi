# MasjidFrame v1.6.7 Release Acceptance Record

v1.6.7 includes the v1.6.6 first-run and deferred Board setup fixes, plus the MasjidBoard Live blank-city discovery fix in PR #118. The v1.6.6 image was tested on a spare card but its Pi 3 image and update bundle were not published.

## Evidence

- [x] Clean v1.6.6 Pi 3 image booted through both splash screens, Wi-Fi setup, and Step 2.
- [x] Without a valid location hierarchy, Set up Board later displayed the unconfigured Board panel; Configure Board reopened Step 2 alone.
- [x] Built-in USB audio was selected and saved automatically on first boot.
- [x] Deferred Board state, audio selection, and both services persisted after reboot.
- [x] PR #118 CI passed and its temporary ARM64 backend on the spare Pi refreshed 24 countries; blank Gauteng cities were counted as unresolved.
- [x] Step 2 displayed the location choices and allowed a masjid to be selected on the Pi.

## Publication gates

- [x] Build a clean image from the immutable v1.6.7 tag and verify embedded version, source commit, and compressed image checksum.
- [x] Flash a spare card and validate first-run Board discovery and selection from the clean image.
- [x] Publish the tested Pi 3 A/B fresh-install image with its checksum.
- [x] Build, sign, verify, and publish the Pi 3 update bundle after image acceptance.

## Final acceptance — 1 October 2026

- [x] Clean v1.6.7 image completed first-run Wi-Fi, loaded the location hierarchy, and allowed a masjid to be selected. Board was configured with one board; built-in audio was selected automatically.
- [x] Fresh-install configuration and audio selection persisted after reboot; both services remained active.
- [x] The original v1.6.4 installation updated to v1.6.7 through the built-in updater during the overnight window. All three configured boards and the built-in audio selection were preserved; audio playback worked.
- [x] The upgraded original card passed a further manual reboot. Running and active slot: system_a / a; rollback_slot=a; upgrade_available=0; bootcount=0.
- [x] Update status reported current_version=v1.6.7-image, installation=null, and last_check_error=null. Both services were active and systemctl reported no failed units.
- [x] Development-machine cleanup removed the uncompressed image and build filesystem after retained release artifacts passed checksum verification; 19 GB remained free.

## Published artifacts

Release: https://github.com/X-Calibre/MasjidPi/releases/tag/v1.6.7

Source commit: `4533ac86854104d1725924e0384c0c7c93071a1c`

| Artifact | SHA256 |
| --- | --- |
| `masjidframe-v1.6.7-pi3-ab.img.zst` | `71506d9acadd6900f001ac93a7ed43a6db0ba90580368a57750fbb849e6e27ca` |
| `masjidframe-update-v1.6.7-pi3.tar.zst` | `7d798837844e4d96178df2c56d256c703aa4714e1484bfd8e2b86e0b36d2486b` |
| `masjidframe-update-v1.6.7-pi3.tar.zst.minisig` | `a98d381b41df41f94f7a6474b6f604394398adbba961a7e5c8fe06e183f76b27` |

The signature, manifest, payload checksums, and decompressed root filesystem passed the update verifier before publication. `SHA256SUMS-PI3` is published alongside these artifacts.
