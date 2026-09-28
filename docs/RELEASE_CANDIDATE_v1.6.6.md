# MasjidFrame v1.6.6 Release Acceptance Record

v1.6.6 includes the v1.6.5 first-run built-in audio fix and repairs deferred Board setup on the Pi 3 portrait appliance.

## Evidence

- [x] Clean v1.6.5 image selected and saved `alsa/plughw:CARD=Device,DEV=0` on first boot without manual audio selection.
- [x] MasjidBoard Live outage exposed a blank Board screen after choosing `Set up Board later`.
- [x] PR #116 frontend changes on the spare Pi displayed the unconfigured Board panel and reopened location setup with Step 2 alone, including a repeat defer/resume cycle.
- [x] PR #116 CI passed on its final commit.

## Publication gates

- [ ] Build the image from the immutable v1.6.6 tag, verify embedded release provenance and compressed-image checksum.
- [ ] Flash a spare card and validate built-in audio, Wi-Fi, offline Board deferral, repeat resume, and persistence after reboot.
- [ ] Publish the tested fresh-install Pi 3 A/B image with its checksum.
- [ ] Build, sign, verify, and publish the Pi 3 update bundle after image acceptance.
