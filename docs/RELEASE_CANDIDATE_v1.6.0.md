# MasjidPi v1.6.0 Release Acceptance Record

This record covers the v1.6.0 release-candidate cycle. `v1.6.0-rc.1` introduced first-run touchscreen onboarding for the MasjidFrame appliance. `v1.6.0-rc.2` added post-setup network management, touch-control refinements and four additional light Board themes. `v1.6.0-rc.3` added native Raspberry Pi Touch Display 2 support and on-device screen controls. `v1.6.0-rc.4` incorporates the physical-display legibility review, retires the older 600 × 1024 profile and refines the touch-control model.

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
- [ ] GitHub Actions passes on the RC4 release-preparation pull request.
- [ ] GitHub Actions passes on the integrated RC4 `main` commit.

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
- [ ] Merge the RC4 release-preparation pull request after CI passes.
- [ ] Confirm CI passes on the resulting `main` commit.
- [ ] Create immutable tag `v1.6.0-rc.4` from the accepted `main` commit.
- [ ] Verify ARM64 and AMD64 archives plus `SHA256SUMS`.
- [ ] Install and validate the published ARM64 archive on the Pi 4.

## Stable-release hardware follow-up

The first-run and RC2 enhancement flows have been functionally accepted on Raspberry Pi 4 source installations. Before stable v1.6.0 promotion:

- confirm the first-run flow on the intended Pi 3B appliance hardware or explicitly record its deferral;
- review Pi 3B memory headroom during setup and normal Board operation; and
- complete any fixes discovered during the RC soak period.

Release-candidate tags are immutable and must not be moved or reused. Any code change after RC2 requires a new release-candidate tag.
