# MasjidPi v1.6.0 Release Acceptance Record

This record covers the v1.6.0 release-candidate cycle. `v1.6.0-rc.1` introduces first-run touchscreen onboarding for the MasjidFrame appliance and includes the accumulated Board fixes completed after v1.5.2.

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

## Automated validation

- [x] Go formatting passes.
- [x] Go vet passes.
- [x] Go tests pass, including NetworkManager, API and credential-safety coverage.
- [x] Frontend JavaScript tests cover first-run routing, keyboard layouts, touch pickers, hidden networks and access URLs.
- [x] Installer, runtime, boot, display-profile and release-package shell tests pass.
- [ ] GitHub Actions passes on the integrated `main` commit.

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
- [ ] Integrated changes are merged to `main` after local validation.
- [ ] Tag `v1.6.0-rc.1` is created from the accepted `main` commit.
- [ ] The release workflow publishes ARM64 and AMD64 archives plus `SHA256SUMS` as a prerelease.
- [ ] The published ARM64 artifact is installed and validated on the Pi 4 test appliance.

## RC1 hardware follow-up

The first-run flow has been functionally accepted on Raspberry Pi 4 source installation. Before stable v1.6.0 promotion:

- install and validate the published ARM64 RC1 artifact on the Pi 4;
- confirm the first-run flow on the intended Pi 3B appliance hardware or explicitly record its deferral;
- review Pi 3B memory headroom during setup and normal Board operation; and
- complete any fixes discovered during the RC soak period.

The RC tag is immutable and must not be moved or reused. Any code change after RC1 requires a new release-candidate tag.
