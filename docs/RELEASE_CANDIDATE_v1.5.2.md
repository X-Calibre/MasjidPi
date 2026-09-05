# MasjidPi v1.5.2 Release Acceptance Record

This record covers promotion of the v1.5.2 release after three release candidates and the final MasjidBoard feature-validation cycle.

## Release scope

### Appliance resilience

- quiet branded Plymouth and Cog startup stages for portrait appliance and landscape HDMI profiles;
- WebKit warm-up and ordered Plymouth-to-Cog DRM handoff;
- a normally read-only Raspberry Pi boot firmware filesystem with controlled write windows;
- an mpv IPC socket under the systemd-managed `/run/masjidpi` directory;
- durable atomic JSON replacement and fewer unnecessary persistent-state writes;
- safer source-update migration, incomplete-update detection, rollback and self-test behaviour;
- hardened Jamiat Islamic Economic Indicator date parsing;
- automatic audio-output fallback and restoration; and
- reduced WebKit prayer-grid replacement and display-response churn.

### MasjidBoard content and presentation

- detailed, per-masjid Jumu'ah schedule cards during the Islamic-Friday interval;
- structured announcements, class changes, weekly and Ramadan programmes, Arabic notices, funeral, Nikah, Eid, Taleem/Jamaat and contribution content;
- optional Dua-after-Adhan display as an exclusive five-minute card;
- daily Ayah, Hadith and Sunnah content;
- deterministic notice ordering by selected masjid and priority;
- flashing red clock/date warning during the published Zawaal/Istiwaa interval;
- special-day Dhuhr data in Daily Times when it differs from the normal timetable; and
- responsive single-row Daily Times presentation with up to 11 landscape items.

## Automated validation

- [x] Go formatting passes.
- [x] Go vet passes.
- [x] Go tests pass on the development machine and Raspberry Pi 4.
- [x] GitHub Actions passes race-enabled Go tests, ShellCheck, installer tests and frontend JavaScript tests.
- [x] Release-package regression checks are included in CI.

## Raspberry Pi 4 validation

- [x] Source installation completes and the installer self-test passes.
- [x] MasjidPi and display services remain active with no unexpected restarts or logged errors.
- [x] Existing display preferences persist across service restart.
- [x] Detailed Jumu'ah schedules render in portrait and landscape layouts.
- [x] Structured Section 10 community cards render in both layouts.
- [x] Dua after Adhan remains exclusively visible for five minutes, spans the landscape notice column and has no unnecessary source attribution.
- [x] Zawaal/Istiwaa warning activates at the published interval boundaries in both layouts.
- [x] Non-duplicate special Dhuhr appears in Daily Times; a time matching normal Dhuhr Jamaah is suppressed.
- [x] Eleven Daily Times items fit on one landscape row and the appliance layout remains correct.

## Release decisions

- Raspberry Pi 3 soak work completed during the RC cycle identified and drove the WebKit churn reduction. A new soak of the final Board feature set is deferred as post-release support evidence rather than a v1.5.2 blocker.
- Raspberry Pi Zero Listen-only validation is deferred; this release does not claim new Pi Zero-specific validation.
- The temporary splash artwork remains accepted for this maintenance release and is not the final MasjidFrame branding.

## Publication checklist

- [x] The v1.5.2 version and release documentation are prepared on a dedicated branch.
- [ ] The release-preparation pull request is merged after CI passes.
- [ ] Annotated tag `v1.5.2` is created from the accepted `main` commit.
- [ ] The release workflow publishes ARM64 and AMD64 archives plus `SHA256SUMS`.
- [ ] Published checksums and release metadata are verified.

## Publication procedure

After the release-preparation pull request is merged and CI passes:

```bash
git switch main
git pull --ff-only origin main

git tag -a v1.5.2 -m "MasjidPi v1.5.2"
git push origin v1.5.2
```

The tag starts the release workflow. Do not reuse or move any release-candidate tag.
