# MasjidPi v1.6.0 Release Acceptance Record

This record covers the v1.6.0 release-candidate cycle. `v1.6.0-rc.1` introduced first-run touchscreen onboarding for the MasjidFrame appliance. `v1.6.0-rc.2` adds post-setup network management, touch-control refinements and four additional light Board themes.

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

## Stable-release hardware follow-up

The first-run and RC2 enhancement flows have been functionally accepted on Raspberry Pi 4 source installations. Before stable v1.6.0 promotion:

- confirm the first-run flow on the intended Pi 3B appliance hardware or explicitly record its deferral;
- review Pi 3B memory headroom during setup and normal Board operation; and
- complete any fixes discovered during the RC soak period.

Release-candidate tags are immutable and must not be moved or reused. Any code change after RC2 requires a new release-candidate tag.
