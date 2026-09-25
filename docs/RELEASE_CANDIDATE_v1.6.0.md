# MasjidPi v1.6.0 Release Acceptance Record

This record covers the v1.6.0 release-candidate cycle. `v1.6.0-rc.1` introduced first-run touchscreen onboarding for the MasjidFrame appliance. `v1.6.0-rc.2` added post-setup network management, touch-control refinements and four additional light Board themes. `v1.6.0-rc.3` added native Raspberry Pi Touch Display 2 support and on-device screen controls. `v1.6.0-rc.4` incorporated the physical-display legibility review, retired the older 600 × 1024 profile and refined the touch-control model. `v1.6.0-rc.5` completes the accepted post-RC4 interface refinements while keeping the automatic-updater and A/B appliance-image prototype outside the release. `v1.6.0-rc.6` integrates the signed automatic updater and Pi 3 A/B appliance image after end-to-end hardware validation, and adds resilient first-run handling for temporary MasjidBoard outages. `v1.6.0-rc.7` corrects post-RC6 clock-date refresh and stale updater-state issues found during published-image validation. `v1.6.0-rc.8` attempts to move verification extraction to persistent storage after the published RC7 bundle exposed the Pi 3 appliance’s insufficient RAM-backed `/tmp` capacity. Embedded-image inspection found malformed newline escapes before appliance assets were signed or published. `v1.6.0-rc.9` corrects the executable workspace statements and strengthens the regression test.

## Release scope

### First-run MasjidFrame onboarding

- automatic setup entry when no NetworkManager Wi-Fi profile is saved;
- nearby 2.4 GHz Wi-Fi scanning and connection;
- visible and hidden network support, including protected and open hidden networks;
- a lightweight on-screen keyboard with a permanent number row and separate symbol layout;
- password visibility and safe credential handling through standard input rather than process arguments;
- touch-native country, province/region, city and initial-masjid selection;
- initial timetable retrieval and transition to the appliance Board;
- DHCP/network-derived FQDN and IPv4 URLs for advanced configuration from another device; and
- no assumed or fabricated `.local` hostname.

### Included post-v1.5.2 fixes

- prevent duplicate Appliance slides after the Dua-after-Adhan interval;
- report the project version correctly from source builds;
- derive missing Zawaal warning boundaries from available Istiwaa data; and
- format Gregorian Board dates in day-first order, for example `Saturday, 5 September 2026`.

### RC2 appliance enhancements

- reopen visible or hidden Wi-Fi setup from the Appliance Display without deleting the active connection first;
- return from network setup to the Board without changing Wi-Fi;
- show the active network-issued FQDN and IPv4 URL in the Network tab for advanced configuration;
- use the same network-issued FQDN in the installer summary instead of assuming a `.local` hostname;
- reliably close the control sheet with a downward drag under Cog/WPE;
- avoid highlighting the close button when the control sheet first opens;
- use readable white text on highlighted controls in light themes;
- rename the displayed Light theme to Light Gold while retaining its compatible saved identifier; and
- add Ivory, Sage, Sky and Rose, expanding the curated Board theme set from six to ten.

## Automated validation

- [x] Go formatting passes.
- [x] Go vet passes.
- [x] Go tests pass, including NetworkManager, API and credential-safety coverage.
- [x] Frontend JavaScript tests cover first-run routing, keyboard layouts, touch pickers, hidden networks and access URLs.
- [x] Installer, runtime, boot, display-profile and release-package shell tests pass.
- [x] GitHub Actions passes on the integrated `main` commit.

## Raspberry Pi 4 source validation

- [x] Source installation completes and the component-aware installer self-test passes.
- [x] A factory-like reset with no saved Wi-Fi profile enters setup automatically after reboot.
- [x] Visible-network selection, alphanumeric password entry, password visibility and connection succeed.
- [x] Hidden-network manual SSID/password entry and connection succeed.
- [x] Touch location pickers and masjid selection succeed under Cog/WPE DRM.
- [x] The first timetable becomes current without cached-data fallback or update error.
- [x] The completion screen displays the DHCP-issued FQDN and current IPv4 address.
- [x] Advanced configuration is reachable from another device using the displayed address.
- [x] The setup override can be removed and normal appliance startup redirects to the configured Board.

## RC1 publication checklist

- [x] All outstanding feature and bug-fix work is integrated with the latest `main` documentation.
- [x] Version metadata is set to `v1.6.0-rc.1`.
- [x] Release documentation is prepared.
- [x] Integrated changes are merged to `main` after local validation.
- [x] Tag `v1.6.0-rc.1` is created from the accepted `main` commit.
- [x] The release workflow publishes ARM64 and AMD64 archives plus `SHA256SUMS` as a prerelease.
- [x] The published ARM64 artifact is installed and validated on the Pi 4 test appliance.

## RC2 source validation

- [x] Reopening network setup and returning to the Board work on the Pi 4 appliance.
- [x] Visible and hidden replacement networks connect successfully.
- [x] The Network tab displays the current DHCP/reverse-DNS FQDN and IPv4 URL.
- [x] The installer summary displays the real network-issued FQDN.
- [x] Downward control-sheet closing works reliably on the Cog/WPE touchscreen renderer.
- [x] The control sheet opens without an unwanted close-button focus highlight.
- [x] Highlighted controls remain readable in the light themes.
- [x] Light Gold, Ivory, Sage, Sky and Rose render, apply and persist correctly.
- [x] The ten-theme Appliance selector fits and works on the 7-inch display.

## RC2 publication checklist

- [x] Version metadata is set to `v1.6.0-rc.2`.
- [x] RC2 release scope and hardware results are documented.
- [x] Automated CI passes on the integrated `main` commit `efa580d02ffa6057d4bce62df564a2af552dc65b`.
- [x] The RC2 preparation branch is merged to `main`.
- [x] Tag `v1.6.0-rc.2` is created from the accepted `main` commit.
- [x] The release workflow publishes ARM64 and AMD64 archives plus `SHA256SUMS` as a prerelease.
- [x] The published ARM64 artifact is installed and validated on the Pi 4 test appliance.

## RC2 published-package validation

- [x] The published ARM64 archive matches its SHA-256 checksum `10b403f836f160055b82a396d0c6d1dd60bca27606d53191503f6e8f880544df`.
- [x] Release installation completes and the installer self-test passes on the Pi 4.
- [x] `/api/version` reports `v1.6.0-rc.2`.
- [x] The Listen + Board component profile and all saved Board settings survive the upgrade.
- [x] All three configured masjids return current data without cache fallback or update errors.
- [x] The Network tab data retains the expected IPv4 address and network-issued FQDN.
- [x] Cog starts the Appliance Display and `/boot/firmware` returns to read-only mode.
- [x] The warm-up oneshot completes with `Result=success` and `ExecMainStatus=0`.

## RC3 scope

- detect a connected 7-inch Raspberry Pi Touch Display 2 from its DRM DSI connector and exact 720 × 1280 mode;
- launch a dedicated `appliance-720` profile in the panel's native portrait orientation without Cog rotation;
- retire the rotated 600 × 1024 Waveshare appliance profile and its hardware-detection path;
- provide a purpose-built 720 × 1280 Board, touch-control sheet, first-run setup and Change Wi-Fi layout;
- select an upright boot splash for the native portrait DSI panel;
- expose persistent Touch Display 2 backlight brightness through the standard Linux kernel backlight interface;
- provide Off, Mild, Medium and Strong cool-white correction for the Board and Wi-Fi setup interfaces;
- correct keyboard width, economic-indicator vertical use and economic heading spacing at 720 × 1280; and
- refresh changed frontend asset keys so Cog/WPE cannot reuse controllers that predate `appliance-720`.

## RC3 automated and source validation

- [x] Go formatting passes.
- [x] The full Go test suite passes, including display-settings persistence, validation and backlight conversion coverage.
- [x] The ARM64 source build completes on Raspberry Pi 4.
- [x] Frontend JavaScript tests pass.
- [x] Installer, display-profile and other shell tests pass.
- [x] GitHub Actions passes on integrated `main` commit `eb7f998be51e54d3958aab65ba45ae68df34800c`.

## RC3 Raspberry Pi 4 hardware validation

- [x] DRM reports `card1-DSI-1` connected with the native `720x1280` mode.
- [x] The display launcher automatically selects `profile=appliance-720`.
- [x] Cog runs without `rotation=1` and loads the dedicated portrait presentation.
- [x] The Board fills the display and all portrait slides render legibly.
- [x] Touch directions, swipe-up controls and swipe-down closing work correctly.
- [x] First-run setup and Change Wi-Fi use the 720 × 1280 layout.
- [x] The on-screen keyboard fits without horizontal clipping.
- [x] Islamic Economic Indicators use the available height with corrected heading-to-date spacing.
- [x] The kernel exposes `panel_backlight@1` with a native brightness range of 0–31.
- [x] Backlight brightness and selectable cool-white correction work from the Display tab.

## RC3 publication checklist

- [x] Validated Touch Display 2 work is integrated with current `main`.
- [x] Version metadata is set to `v1.6.0-rc.3`.
- [x] RC3 scope and hardware results are documented.
- [x] GitHub Actions passes on release-preparation `main` commit `eb7f998be51e54d3958aab65ba45ae68df34800c`.
- [x] Tag `v1.6.0-rc.3` is created from the accepted `main` commit.
- [x] The release workflow publishes ARM64 and AMD64 archives plus `SHA256SUMS` as a prerelease.
- [x] The published ARM64 artifact is installed and validated on the Pi 4 test appliance.

## RC4 scope

- improve distance legibility throughout the native 720 × 1280 presentation;
- rebalance the header with a larger Masjid name, larger alternating Gregorian/Islamic date and clearer upcoming-event text;
- show source-switch notifications as top toast overlays without changing header geometry;
- enlarge community and Dua-after-Adhan content while preserving title hierarchy;
- make Detailed Jumu'ah schedules full-height cards with vertically stacked, enlarged event times;
- improve Salaah Time Change and Islamic Economic Indicator spacing;
- split touch controls into top-edge Quick Settings and a bottom-edge Appliance Controls sheet;
- provide brightness and Master, Masjid and Radio volume sliders in Quick Settings;
- duplicate Masjid and Radio playback actions in Quick Settings while retaining the full bottom controls;
- size Quick Settings dynamically to its content and label its Masjid and Radio action groups;
- remove the experimental software White Balance correction and retain only the kernel backlight control;
- retire the former 600 × 1024 Waveshare profile, rotated splash and calibration path;
- use the responsive standard presentation for unsupported displays; and
- remove the temporary visual-demo feature used during development.

## RC4 automated validation

- [x] Go formatting and vet pass.
- [x] The full race-enabled Go test suite passes.
- [x] Frontend JavaScript syntax and regression tests pass.
- [x] Installer, boot and display-profile shell tests pass.
- [x] GitHub Actions run 535 passes on the final feature head.
- [x] GitHub Actions passes on the RC4 release-preparation pull request.
- [x] GitHub Actions passes on the integrated RC4 `main` commit.

## RC4 Raspberry Pi 4 hardware validation

- [x] The Touch Display 2 is automatically detected as `appliance-720`.
- [x] The presentation fills the native portrait panel and touch directions are correct.
- [x] Header, Salaah, community, Dua, economic and Detailed Jumu'ah layouts were reviewed on the physical panel.
- [x] Kernel-backed brightness adjustment works and persists.
- [x] Top and bottom touch sheets open from their respective display edges.
- [x] Masjid and Radio actions operate from the touch controls.
- [ ] Recheck the final dynamically sized Quick Settings panel on the physical display.
- [ ] Recheck the final full-height Appliance Controls panel on all four tabs.
- [ ] Complete a reboot check with the final RC4 source before tagging.

## RC4 publication checklist

- [x] Completed Touch Display 2 refinements are merged to `main`.
- [x] Version metadata is set to `v1.6.0-rc.4`.
- [x] RC4 release scope and completed validation are documented.
- [x] Merge the RC4 release-preparation pull request after CI passes.
- [x] Confirm CI passes on the resulting `main` commit.
- [x] Create immutable tag `v1.6.0-rc.4` from the accepted `main` commit.
- [x] Verify ARM64 and AMD64 archives plus `SHA256SUMS`.
- [ ] Install and validate the published ARM64 archive on the Pi 4.

## RC5 scope

- make Daily Ayah, Daily Hadith and Daily Sunnah full-height cards on the native 720 × 1280 Appliance Display;
- dynamically fit daily-content text to the largest size that does not overflow the visible card;
- enlarge Detailed Jumu'ah labels to match the associated times;
- rebalance the upcoming-event block and top control-panel handle spacing;
- apply the selected Board theme to first-run and Change Wi-Fi setup;
- allow saved Listen favourites to be reordered and preserve that order in Appliance touch controls;
- sort the main Masjid catalogue alphabetically while pinning Qur'aan Recitation, Takbeer and Sautun Noor at the top;
- sort Radio choices alphabetically in both the main interface and Appliance controls; and
- exclude the automatic-updater and Pi 3 A/B appliance-image prototype from this release.

## RC5 automated validation

- [x] The full GitHub Actions test job passes on feature head `1adc3ec77dddf407cdafca904bbcc2f6f3819025`.
- [x] CodeQL reports no new alert in the pull request.
- [x] Go, JavaScript/TypeScript and Actions code analysis pass.
- [x] Pull request #83 is merged to `main` as `385371f6a6ce35a4459768abc481caa0a9e1ab67`.
- [ ] GitHub Actions passes on the RC5 release-preparation pull request.
- [ ] GitHub Actions passes on the integrated RC5 `main` commit.

## RC5 publication checklist

- [x] All completed non-updater branches are accounted for.
- [x] Automatic-updater and A/B image work remains isolated on `prototype/pi3-ab-image`.
- [x] Version metadata is set to `v1.6.0-rc.5`.
- [x] RC5 scope and validation status are documented.
- [ ] Merge the RC5 release-preparation pull request after CI passes.
- [ ] Confirm CI passes on the resulting `main` commit.
- [ ] Create immutable tag `v1.6.0-rc.5` from the accepted `main` commit.
- [ ] Verify ARM64 and AMD64 archives plus `SHA256SUMS`.
- [ ] Install and validate the published ARM64 archive on the Pi 4 daily-use appliance.

## RC6 scope

- integrate the signed stable-release discovery, approval, download, scheduling and installation workflow;
- check for stable releases weekly, with user-initiated checks available from the touchscreen and Web UI;
- expose a concise touchscreen update summary and a detailed Web UI Updates page;
- install verified updates into the inactive Pi 3 system slot and preserve appliance identity, SSH host keys, configuration and the `masjidframe` account password;
- require a ten-minute healthy trial before confirming a new system slot, with automatic rollback after two failed boots;
- use the documented factory account in new appliance images;
- allow first-run Board setup to be deferred when MasjidBoard is unavailable, with manual retry only and no automatic short-interval polling;
- restore deferred Board configuration once the provider is available; and
- package the complete Pi 3 A/B image and signed updater through the release workflow.

## RC6 automated validation

- [x] Go formatting and vet pass.
- [x] The complete race-enabled Go test suite passes.
- [x] Frontend JavaScript syntax and UI regression tests pass.
- [x] Installer, display, release-package, updater and password-preservation shell tests pass.
- [x] ShellCheck passes for the appliance scripts and regression tests.
- [x] GitHub Actions passes on the integrated updater pull request.
- [ ] GitHub Actions passes on the RC6 release-preparation pull request.
- [ ] GitHub Actions passes on the integrated RC6 `main` commit.

## RC6 Raspberry Pi 3 hardware validation

- [x] A signed update bundle built from `b79c1442c6e06caa55726701000b669ace92f332` verifies successfully on the appliance.
- [x] The updater plans the system-A to system-B transition without writes.
- [x] Installation writes and reads back the complete inactive root filesystem and arms system B with system A as rollback.
- [x] Machine identity, SSH host keys, release identity and the `masjidframe` password survive installation.
- [x] System B boots into probation and reports the signed release through the update API.
- [x] The ten-minute probation confirms system B automatically and clears the trial boot flags.
- [x] A subsequent reboot remains on confirmed system B with no failed units.
- [x] MasjidBoard selection and display configuration persist across the update and reboot.
- [x] First-run Board selection works when MasjidBoard is online.
- [x] The temporary-provider-outage path is covered by automated and embedded-image validation without aggressive automatic retries.

## RC6 publication checklist

- [x] Automatic-updater and Pi 3 A/B image work is merged to `main`.
- [x] Hardware validation is complete on the Pi 3 appliance candidate.
- [x] Version metadata is set to `v1.6.0-rc.6`.
- [x] RC6 scope and validation status are documented.
- [x] Merge the RC6 release-preparation pull request after CI passes.
- [x] Confirm CI passes on the resulting `main` commit.
- [x] Create immutable tag `v1.6.0-rc.6` from the accepted `main` commit.
- [x] Verify published release archives, A/B image, signed update bundle and checksums.
- [x] Confirm stable-only discovery intentionally excludes RC6 prerelease assets; retain the end-to-end GitHub discovery gate for stable v1.6.0.

## RC7 scope

- refresh MasjidBoard data after NTP or another system-clock correction changes the calendar date;
- reconcile persisted update candidates against the running version before automatic preparation;
- force release discovery when persisted availability is invalid, non-stable, equal to or older than the installed version; and
- stop safely without downloading or installing when reconciliation cannot reach GitHub.

## RC7 automated validation

- [x] Go formatting and vet pass.
- [x] The complete race-enabled Go test suite passes.
- [x] Frontend syntax, shell, installer, display and packaging checks pass.
- [x] Regression tests cover invalid laboratory versions, equal and older stable versions, matching-RC promotion and discovery failure.
- [x] GitHub Actions run 605 passes on updater reconciliation commit `f230a968986853e51beb36044ee8886f305fba2a`.
- [ ] GitHub Actions passes on the RC7 release-preparation pull request.
- [ ] GitHub Actions passes on the integrated RC7 `main` commit.

## RC7 Raspberry Pi 3 hardware validation

- [x] Published RC6 is confirmed in system A with `active_slot=a`, `rollback_slot=a`, `upgrade_available=0` and `bootcount=0`.
- [x] The persisted laboratory candidate reproduces the stale automatic-preparation attempt under the unmodified RC6 runtime.
- [x] A patched ARM64 binary built from merged commit `3cd46bdbd45d725f4337d7866c0217d927e981f5` refreshes discovery and clears the stale candidate.
- [x] The reconciled state contains no available release, approval, download or obsolete installation record.
- [x] Both MasjidPi services remain active with no failed systemd units.
- [x] A complete reboot retains confirmed system A and the reconciled update state without another preparation attempt.

## RC7 publication checklist

- [x] NTP date refresh and stale updater-state fixes are merged to `main`.
- [x] Pi 3 source-binary regression validation is complete.
- [x] Version metadata is set to `v1.6.0-rc.7`.
- [x] RC7 scope and validation status are documented.
- [ ] Merge the RC7 release-preparation pull request after CI passes.
- [ ] Confirm CI passes on the resulting `main` commit.
- [ ] Create immutable tag `v1.6.0-rc.7` from the accepted `main` commit.
- [ ] Verify published archives, A/B image, signed update bundle and checksums.
- [ ] Install and validate the published RC7 update bundle on the Pi 3 appliance.

## RC8 scope

- create verification extraction directories beside the downloaded update bundle on persistent storage;
- avoid the Pi 3 appliance’s 453 MB RAM-backed `/tmp`, which cannot hold both the compressed bundle and extracted payload;
- preserve automatic removal of verification and installation workspaces; and
- retain the inactive A/B slot as the rollback copy without retaining unnecessary extracted payloads.

## RC8 automated validation

- [x] Go formatting, vet and the complete race-enabled Go test suite pass on the updater fix pull request.
- [x] Shell syntax and ShellCheck pass.
- [x] A regression test requires verification and installation workspaces to use persistent storage.
- [x] GitHub Actions run 610 passes on pull-request head `9118e947d9eb09d25c758593d4b5a576584e6ae1`.
- [ ] GitHub Actions passes on the RC8 release-preparation pull request.
- [ ] GitHub Actions passes on the integrated RC8 `main` commit.

## RC8 Raspberry Pi 3 validation

- [x] The published RC7 bundle and signature match their release hashes.
- [x] Verification initially reproduces `No space left on device` with an empty 453 MB `/tmp`.
- [x] The identical bundle passes signature, manifest, payload and decompressed-rootfs checks when verification uses persistent storage.
- [x] The read-only plan targets inactive `SYSTEM_B` while confirmed `SYSTEM_A` and U-Boot state remain unchanged.
- [ ] Install the published RC8 bundle into inactive `SYSTEM_B`.
- [ ] Boot RC8 into probation and confirm it automatically after ten healthy minutes.
- [ ] Reboot again and confirm RC8 remains the stable slot with no failed units.

## RC8 publication checklist

- [x] Persistent verification-workspace fix is merged to `main`.
- [x] Version metadata is set to `v1.6.0-rc.8`.
- [x] RC8 scope and validation status are documented.
- [ ] Merge the RC8 release-preparation pull request after CI passes.
- [ ] Confirm CI passes on the resulting `main` commit.
- [ ] Create immutable tag `v1.6.0-rc.8` from the accepted `main` commit.
- [ ] Verify published archives, A/B image, signed update bundle and checksums.
- [ ] Install and validate the published RC8 update bundle on the Pi 3 appliance.

## RC9 scope

- replace malformed literal newline escapes in the updater workspace block with executable shell statements;
- require persistent-workspace assignments to occupy complete lines in regression coverage;
- reject literal newline escapes in the updater source; and
- supersede RC8 for appliance image and signed-update validation.

## RC9 automated validation

- [x] Go formatting, vet and the complete race-enabled Go test suite pass on the corrective pull request.
- [x] Shell syntax and ShellCheck pass.
- [x] The strengthened regression test requires executable persistent-workspace assignments.
- [x] GitHub Actions run 614 passes on corrective head `09079e8b69d4926bbe3d909c00444126e969efc9`.
- [x] GitHub Actions passes on the RC9 release-preparation pull request.
- [x] GitHub Actions passes on the integrated RC9 `main` commit `de96737af40eba22cb8e2efc63ef6573904ccb34`.

## RC9 Raspberry Pi 3 validation

- [x] RC8 embedded-image inspection detects the malformed updater before bundle signing or appliance publication.
- [x] No RC8 appliance update bundle or A/B image is published or installed.
- [x] The published updater contains separate executable persistent-workspace assignments and rejects literal newline escapes.
- [x] The published bundle and signature match their recorded SHA-256 hashes and pass Minisign, manifest, payload and decompressed-rootfs verification without overriding `TMPDIR`.
- [x] The read-only plan targets inactive `SYSTEM_B` while confirmed `SYSTEM_A` and U-Boot state remain unchanged.
- [x] RC9 installs into `SYSTEM_B`; the written root filesystem and staged boot payloads verify before the trial is armed.
- [x] Machine ID, SSH host keys, account password, component profile and Board selection survive the slot installation.
- [x] RC9 boots from `SYSTEM_B` and confirms automatically after the ten-minute probation with `active_slot=b`, `rollback_slot=b`, `upgrade_available=0` and `bootcount=0`.
- [x] A normal reboot remains on confirmed RC9 with both services active, no failed units, the Board visible and live Radio audio working through the selected USB ALSA device.
- [x] Two deliberately unconfirmed `SYSTEM_A` boots reach `bootcount=2`; the following boot automatically rolls back to confirmed RC9 on `SYSTEM_B` and clears the trial state.

## RC9 publication checklist

- [x] Corrective updater and regression-test changes are merged to `main`.
- [x] Version metadata is set to `v1.6.0-rc.9`.
- [x] RC9 scope and validation status are documented.
- [x] Merge the RC9 release-preparation pull request after CI passes.
- [x] Confirm CI passes on the resulting `main` commit `de96737af40eba22cb8e2efc63ef6573904ccb34`.
- [x] Create immutable tag `v1.6.0-rc.9` from the accepted `main` commit.
- [x] Verify the release contains the ARM64 and AMD64 archives, `SHA256SUMS`, Pi 3 A/B image and checksum, signed update bundle and Minisign signature.
- [x] Install and validate the published RC9 update bundle on the Pi 3 appliance.

## Post-RC9 power-loss validation

- [x] A physical power cut during the inactive `SYSTEM_B` root-filesystem write leaves confirmed RC6 on `SYSTEM_A` bootable with the trial state unarmed.
- [x] Both application services recover with no failed units, throttling or persistent-configuration loss.
- [x] Retrying the unchanged signed RC9 bundle rewrites and verifies `SYSTEM_B`, boots RC9 and confirms it after the normal ten-minute probation.
- [x] Temporary BOOT payloads are replaced and removed by the retry.
- [x] The interrupted extraction leaves a 372 MB `/persistent/updates/install.*` workspace, demonstrating that EXIT-trap cleanup alone is insufficient across power loss.
- [ ] Merge exclusive updater locking and pre-verification orphan cleanup.
- [ ] Validate the corrective updater on Pi 3 before creating the next immutable release candidate.

## Stable-release hardware follow-up

The first-run and RC2 enhancement flows have been functionally accepted on Raspberry Pi 4 source installations. Before stable v1.6.0 promotion:

- confirm the first-run flow on the intended Pi 3B appliance hardware or explicitly record its deferral;
- review Pi 3B memory headroom during setup and normal Board operation; and
- complete any fixes discovered during the RC soak period.

Release-candidate tags are immutable and must not be moved or reused. Any code change after RC9 requires a new release-candidate tag.
