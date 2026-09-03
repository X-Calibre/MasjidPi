# MasjidPi v1.5.2 Release-Candidate Acceptance Record

This document tracks the checks required to promote `v1.5.2-rc.1` to the v1.5.2 maintenance release.

## Scope

v1.5.2 improves existing appliance behaviour rather than introducing a new user workflow:

- quiet branded Plymouth and Cog startup stages for portrait appliance and landscape HDMI profiles;
- WebKit warm-up and ordered Plymouth-to-Cog DRM handoff;
- a normally read-only Raspberry Pi boot firmware filesystem with controlled package and installer write windows;
- an mpv IPC socket under the systemd-managed `/run/masjidpi` runtime directory;
- more durable atomic JSON replacement and reduced unnecessary persistent-state writes;
- safer source-update migration, incomplete-update detection, rollback and self-test behaviour;
- corrected and hardened Jamiat Islamic Economic Indicator date parsing; and
- automatic fallback from an unavailable selected audio output and restoration when it reconnects.

## Completed evidence

- Automated Go, vet, race, shell, installer and frontend validation passed on `main` before candidate preparation.
- Portrait boot behaviour was validated on the Raspberry Pi 4 with the 7-inch display.
- Landscape boot behaviour was validated on the Raspberry Pi 4 with a 1440p HDMI display.
- Controlled shutdown, cold start and abrupt power-loss recovery completed with normal splash and application behaviour.
- The boot filesystem remained read-only during normal use and returned to read-only after source installation and a real `raspi-firmware` reinstall.
- Persisted Board, Listen and audio settings survived service restarts; all stored JSON remained valid after the abrupt-power test.
- A Logitech H390 connected after boot was discovered and selected without restarting MasjidPi. Playback fell back safely when it was removed and returned automatically when it reconnected, with zero service restarts.
- Jamiat economic dates accept standard and observed month spellings, spacing and dash variants while rejecting impossible dates.

Detailed physical checks are recorded in [VALIDATION_CHECKLIST.md](VALIDATION_CHECKLIST.md).

## Required before stable promotion

- [ ] The release workflow publishes ARM64 and AMD64 `v1.5.2-rc.1` archives and `SHA256SUMS` as a GitHub prerelease.
- [ ] The published ARM64 artifact upgrades the Raspberry Pi 4 test appliance successfully and reports `v1.5.2-rc.1`.
- [ ] Existing component selection, user settings, audio output and cached state remain intact after the packaged upgrade.
- [ ] Both boot splash stages, Board display, Listen playback and read-only boot mount remain correct after reboot.
- [ ] The post-change Raspberry Pi 3 soak data is reviewed with no release-blocking memory growth, restarts, throttling, storage or network failures.
- [ ] Raspberry Pi Zero Listen-only validation is completed, or its absence is explicitly documented as deferred support evidence for this release.
- [ ] No release-blocking defect remains open.

The temporary splash artwork may be accepted for this maintenance release only by an explicit product decision; it is not the final MasjidFrame branding.

## Automated validation

```bash
cd ~/MasjidPi
git switch release/v1.5.2-rc.1
git pull --ff-only

make test

cd backend
go vet ./...
test -z "$(gofmt -l .)"
```

CI must also pass its race-enabled Go tests, ShellCheck, installer regression tests and frontend JavaScript checks.

## Candidate publication

After the release-preparation commit is merged into `main` and CI passes:

```bash
git switch main
git pull --ff-only origin main

git tag -a v1.5.2-rc.1 -m "MasjidPi v1.5.2-rc.1"
git push origin v1.5.2-rc.1
```

The release workflow must publish ARM64 and AMD64 archives plus `SHA256SUMS`, and mark the GitHub release as a prerelease.

## Stable promotion

Once every required gate above is resolved, set `version.json` to `v1.5.2`, update this record with the acceptance outcome, merge that change to `main`, and publish the annotated `v1.5.2` tag. Do not reuse or move the RC tag.
