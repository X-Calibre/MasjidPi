# MasjidFrame v1.6.3 Release Acceptance Record

v1.6.3 refreshes MasjidBoard Live shared Ayah, Hadith and Sunnah content at the same time as selected masjid notices and timetables. The normal upstream polling interval is 30 minutes; startup, board selection and local-date refresh paths can also trigger a check.

## Scope

- Remove the once-per-Johannesburg-day fetch limit for enabled shared daily content.
- Compare the newly fetched items with cached items and persist only when they change.
- Continue displaying the last good content if the source is unavailable or returns invalid data.
- Skip requests when all three shared daily items are disabled.

## Automated validation

- [x] Service test confirms multiple checks on the same day and persistence of a changed Hadith.
- [x] Service test confirms disabled-content behavior and last-known-good fallback.
- [x] Focused MasjidBoard service and app Go tests pass locally.
- [x] Go formatting, vet, race-enabled tests and shell/frontend CI pass on the release-preparation pull request.
- [x] Release workflow succeeds for the tagged main commit.

## Appliance validation

- [x] On the Pi 3, the temporary v1.6.3 binary started with all daily items enabled, the API and services remained healthy, and the startup fetch had the Hadith subsequently observed by the independent monitor at 12:30 SAST on 27 September 2026.
- [x] The 12:47 SAST notice refresh ran for all three selected masjids with no daily-content error.
- [ ] Confirm a changed Hadith is picked up by a periodic fetch without a reboot. The release owner accepted proceeding without waiting for this live observation; repeated same-day fetches and changed-content persistence are covered by the passing automated tests.
- [ ] Verify unchanged items remain stable and a temporary source failure keeps the last good content.
- [x] Verify normal notice and timetable refresh, display and audio remain healthy. The installed v1.6.3 image runs both services, and the owner confirmed board and audio operation.
- [x] Install the published signed update through the WebUI on the Pi 3. The appliance reported `v1.6.3-image` in `system_a`, completed ten-minute probation, and confirmed the slot (`active_slot=a`, `rollback_slot=a`, `upgrade_available=0`, `bootcount=0`). After a normal reboot, both services remained active and no units had failed.
- [x] Explain the cleared installation status after reboot: an automatic update check recognizes v1.6.3 as installed and clears the prior release offer and installation record; this does not indicate a rollback.

## Publication

- [x] Version metadata is prepared as `v1.6.3` on the enhancement branch.
- [x] Merge pull request #111 into main at `f47d8aeace519509d04fa6c4827b8d3746b5da3b`.
- [x] Create immutable tag `v1.6.3` from the accepted main commit.
- [x] Verify the published ARM64 and AMD64 archives, SHA256SUMS and release metadata.
- [x] Build, sign and verify the Pi 3 A/B image and update bundle; upload the image, signed bundle and Pi 3 checksums to the release.

The published tag is immutable; do not move or reuse it.
