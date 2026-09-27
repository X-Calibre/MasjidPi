# MasjidFrame v1.6.4 Release Acceptance Record

v1.6.4 refines the Listen and Board controls on the appliance. This record tracks the branch test separately from the signed A/B update and publication.

## Scope

- Remove redundant Listen Power switches; a Play action can restore a previously disabled source.
- Make Radio schedule time fields larger and save both times together on request.
- Show `MasjidFrame Built-In Audio` for the appliance USB audio adapter (`0c76:1203`).
- Simplify Board display settings, align Web UI theme descriptions with the touch panel, and hide the standard display preview when a 720×1280 DSI display is detected.
- Add spacing around the MasjidFrame display preview action.

## Automated validation

- [x] Draft pull request #113 runs Go formatting, vet, race-enabled tests, shell checks, installer/display tests and frontend tests successfully on the latest UI-fix commit.
- [x] The release-preparation commit `cb291b0` passed the complete pull-request CI workflow.
- [ ] Recheck CI on the merged main commit.

## Appliance branch validation

- [x] ARM64 test binary built from `fb2a72a` and temporarily installed on the Pi 3 with a full frontend copy; the subsequent preview spacing change from `1ebe65c` was copied to the same test appliance.
- [x] The API reported `v1.6.4-rc.1`, the Touch Display 2 profile `appliance-720`, and the audio device name `MasjidFrame Built-In Audio`; both application and display services were active.
- [x] The owner confirmed the larger Radio schedule controls after a browser refresh, the preview spacing, and the remaining Listen, schedule, Board preview and theme-description checks.
- [x] Restore the temporary test binary and original frontend to the confirmed v1.6.3 baseline. The API reports `v1.6.3-image`, the original frontend contains the v1.6.3 Radio schedule label, and both services are active. Restart the display browser before signed-update validation so it loads the restored assets.

## Release validation and publication

- [x] Prepare `version.json` as `v1.6.4` on the enhancement branch. Keep README's current-stable-release reference at v1.6.3 until publication.
- [ ] Merge pull request #113 after release-preparation CI passes and the release is accepted.
- [ ] Create immutable tag `v1.6.4` from the accepted main commit; never move or reuse the tag.
- [ ] Verify release archives and checksums, then build, sign and verify the Pi 3 A/B image and update bundle.
- [ ] Install the signed update from a confirmed v1.6.3 appliance, complete ten-minute probation and verify the slot remains confirmed after a normal reboot.
- [ ] Verify Listen, Board, saved settings, audio and network remain healthy after the signed update; confirm no failed services.
