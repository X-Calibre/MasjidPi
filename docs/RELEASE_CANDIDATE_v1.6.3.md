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
- [ ] Go formatting, vet, race-enabled tests and shell/frontend CI pass on the release-preparation pull request.
- [ ] CI passes on the resulting main commit.

## Appliance validation

- [ ] With daily items enabled, verify a startup fetch and another fetch on the normal 30-minute notice refresh.
- [ ] Verify a changed Hadith appears on the Board after a successful fetch without a reboot.
- [ ] Verify unchanged items remain stable and a temporary source failure keeps the last good content.
- [ ] Verify normal notice and timetable refresh, display and audio remain healthy.

## Publication

- [x] Version metadata is prepared as `v1.6.3` on the enhancement branch.
- [ ] Merge the release-preparation pull request after CI and appliance validation.
- [ ] Create immutable tag `v1.6.3` from the accepted main commit.
- [ ] Verify the published ARM64 and AMD64 archives, SHA256SUMS and release metadata.
- [ ] Build, sign and verify the Pi 3 A/B image and update bundle where applicable.

Do not move or reuse the release tag. The release workflow publishes assets when the tag is created, so tagging follows acceptance.
