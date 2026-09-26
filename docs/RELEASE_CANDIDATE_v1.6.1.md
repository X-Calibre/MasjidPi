# MasjidFrame v1.6.1 Release Acceptance Record

v1.6.1 is a reliability and product-identity release built on the accepted v1.6.0 appliance baseline. It completes the MasjidPi-to-MasjidFrame rename, hardens touchscreen Wi-Fi setup, and corrects the on-display update installation lifecycle found during v1.6.0 validation.

## Release scope

### MasjidFrame identity

- apply the approved MasjidFrame name, logos and splash artwork throughout the application and appliance image;
- rename runtime executables, services, A/B tooling, updater assets and persistent paths; and
- preserve compatibility needed to migrate an existing v1.6.0 appliance safely.

### Wi-Fi robustness

- validate WPA-PSK credentials at the API and NetworkManager boundaries;
- accept 8–63 character passphrases and 64-character hexadecimal keys; and
- reject malformed credentials before attempting to change the active connection.

### Update installation lifecycle

- detach an approved installation from the browser request that initiated it;
- persist the installing state before returning HTTP 202;
- keep the appliance display polling through installation, reboot and A/B probation;
- bound technical failure details presented in the UI;
- retain the downloaded bundle and signature through the rollback window; and
- remove those assets only after the new slot is confirmed healthy.

## Automated validation

- [x] Feature branches are fully contained in the v1.6.1 integration branch.
- [x] Go formatting and static analysis pass.
- [x] Race-enabled Go tests pass, including updater lifecycle, cleanup and Wi-Fi validation coverage.
- [x] Frontend and shell regression tests pass.
- [x] GitHub Actions passes on the integrated feature head.
- [ ] GitHub Actions passes on the release-preparation pull request.
- [ ] GitHub Actions passes on the resulting `main` commit.

## Appliance validation

- [ ] Install the published signed v1.6.1 update from the appliance display.
- [ ] Navigate away from the update page after approval and confirm installation continues.
- [ ] Confirm the appliance boots the trial slot, remains healthy for ten minutes and commits automatically.
- [ ] Confirm a subsequent reboot remains on the committed slot.
- [ ] Confirm machine identity, SSH host keys, component selection, Board configuration and audio settings are preserved.
- [ ] Confirm display and audio operation after the update.
- [ ] Confirm downloaded update assets are removed after confirmation and not before it.
- [ ] Exercise rollback or a controlled failed trial and confirm retry assets are retained.
- [ ] Re-test visible and hidden Wi-Fi replacement with valid and invalid WPA credentials.
- [ ] Confirm saved Wi-Fi reconnects after reboot.
- [ ] Confirm no failed services, power throttling or unexpected storage growth.

## Publication checklist

- [x] Version metadata is set to `v1.6.1`.
- [x] Release scope and validation gates are documented.
- [ ] Merge the release-preparation pull request after CI passes.
- [ ] Create immutable tag `v1.6.1` from the accepted `main` commit.
- [ ] Verify ARM64 and AMD64 archives plus `SHA256SUMS`.
- [ ] Build, sign and verify the Pi 3 A/B image and update bundle.
- [ ] Publish the release only after the required appliance validation passes.

The `v1.6.1` tag must be created only from the accepted release-preparation commit and must remain immutable.
