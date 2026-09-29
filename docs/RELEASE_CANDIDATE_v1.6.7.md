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

- [ ] Build a clean image from the immutable v1.6.7 tag and verify embedded version, source commit, and compressed image checksum.
- [ ] Flash a spare card and validate first-run Board discovery and selection from the clean image.
- [ ] Publish the tested Pi 3 A/B fresh-install image with its checksum.
- [ ] Build, sign, verify, and publish the Pi 3 update bundle after image acceptance.
