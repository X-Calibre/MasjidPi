# MasjidFrame v1.6.5 Release Acceptance Record

v1.6.5 corrects the first-run audio selection in fresh Pi 3 appliance images. The v1.6.4 WaveShare USB label change made the old `USB Audio` first-run match miss the built-in adapter.

## Validation

- [x] PR #114 CI passed.
- [x] On a Pi 3 spare card, a temporary ARM64 build with no saved audio choice selected `alsa/plughw:CARD=Device,DEV=0`, saved it, and played audio.
- [x] The original v1.6.4 binary and saved audio state were restored on the spare card; services remained active.
- [ ] Build v1.6.5 Pi 3 A/B image from the immutable tag and inspect provenance and checksum.
- [ ] Flash a spare card, complete first-run setup, confirm built-in audio is selected automatically, and verify Wi-Fi and Board state after reboot.
- [ ] Publish the verified fresh-install image and checksum.
- [ ] Build, sign, verify, and publish the Pi 3 update bundle only after hardware validation.
