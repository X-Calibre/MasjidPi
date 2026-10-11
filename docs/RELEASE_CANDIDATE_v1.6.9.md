# MasjidFrame v1.6.9 Release Acceptance Record

Released on 11 October 2026. This record distinguishes completed publication and preview checks from outstanding hardware acceptance.

## Changes

- Portrait contributions use stacked labels and larger banking values.
- Portrait and landscape notices use clearer groups, fitted text and full-height detailed cards.
- Landscape daily Ayah, Hadith and Sunnah share a slide, with content-dependent heights and independent text fitting; failed minimum-size checks fall back to individual slides.
- Landscape Economic Indicators follow the portrait table style.
- Landscape Dua after Adhan independently fits Arabic and English text, with the explanation at the bottom.
- Landscape text-fitting limits and relative spacing scale through 4K.
- Funeral Notice and Nikaah Notice category labels appear above the names in both layouts.
- Appliance Audio offers only detected MasjidFrame Built-In Audio (Waveshare USB ID `0c76:1203`) and External Stereo Output (Pi 3.5 mm jack). Other installations retain general audio choices.

## Completed evidence

- [x] Pi 4 portrait contributions and revised notice previews reviewed and accepted.
- [x] Pi 4 landscape daily content, Economic Indicators, Dua, contributions and revised notice previews reviewed and accepted.
- [x] Full HD and 4K automated text-sizing checks passed; field grouping preserves supplied details.
- [x] All backend Go tests and standalone JavaScript tests passed locally.
- [x] GitHub tests and code analysis passed for the merged release.
- [x] PR #122 merged to main, including the retired-profile test correction.
- [x] Immutable tag v1.6.9 pushed at `22bc292561ed42b85614108ef15c57f0710509a0`.
- [x] Release workflow successfully published ARM64 and AMD64 application packages and SHA256SUMS.
- [x] Pi 3 A/B image built from the tagged source with runtime `v1.6.9-image`.
- [x] Bundle builder verified embedded source commit, build version, runtime version and image metadata against the checkout.
- [x] Fresh-install image passed Zstandard integrity and SHA256 checks on the build machine.
- [x] Update bundle passed Minisign signature and trusted-comment verification using the tracked public key.
- [x] All four Pi 3 image/update assets uploaded and independently confirmed on the stable release page.

## Published artifacts

Release: https://github.com/X-Calibre/MasjidPi/releases/tag/v1.6.9

Source commit: `22bc292561ed42b85614108ef15c57f0710509a0`

| Artifact | Published bytes | SHA256 recorded from build output |
|---|---:|---|
| `masjidframe-v1.6.9-linux-amd64.tar.gz` | 8,340,929 | See application-package `SHA256SUMS`. |
| `masjidframe-v1.6.9-linux-arm64.tar.gz` | 7,726,571 | See application-package `SHA256SUMS`. |
| `SHA256SUMS` | 208 | Checksums for application archives. |
| `masjidframe-v1.6.9-pi3-ab.img.zst` | 940,160,752 | See the accompanying image `.sha256` file; its check passed on the build machine. |
| `masjidframe-v1.6.9-pi3-ab.img.zst.sha256` | 100 | Image checksum file. |
| `masjidframe-update-v1.6.9-pi3.tar.zst` | 373,129,683 | `58bf329c7b106a77bbd3ff95a07f8939f59bc28d7a42c0686edb89691ff60ef3` |
| `masjidframe-update-v1.6.9-pi3.tar.zst.minisig` | 332 | `76115cfa66ebe981e881100932c4f73452ddcf91ead9bbb924f59ef1be6487ff` |

The compressed image expands to 14,235,467,776 bytes. No v1.6.9 on-device update or complete payload-verifier result is asserted by the publication checks above.

## Build-version exception

The immutable v1.6.9 source still contains `version.json` with v1.6.8. The release workflow embeds the tag into application packages. For the appliance build, the operator explicitly used:

```bash
MASJIDFRAME_IMAGE_VERSION=v1.6.9-image ./appliance-image/build-pi3-ab-prototype.sh
```

The embedded release/build version was v1.6.9 and the runtime version was v1.6.9-image. Use the same clean tagged checkout and override to reproduce the image; do not move the release tag. Update version metadata as part of the next code release.

## Remaining acceptance

- [ ] Visually confirm native 3840 × 2160 rendering at 100% browser zoom.
- [ ] Confirm the two appliance audio choices and actual playback through both physical outputs.
- [ ] Flash a spare card and verify v1.6.9 fresh-install setup and persistence.
- [ ] Verify the published signed bundle through the appliance payload verifier and install through the Pi 3 updater.
- [ ] Confirm ten-minute probation, the running slot, preserved settings and a subsequent normal reboot.
- [ ] Complete v1.6.9 hardware soak validation.

## Documentation and screenshots

Current usage and technical documentation describe v1.6.9. Existing user-guide screenshots remain explicitly labelled as v1.6.7 examples; they are not evidence of the v1.6.9 audio-output list or revised cards. Earlier release acceptance records and dated provider/hardware research remain historical.
