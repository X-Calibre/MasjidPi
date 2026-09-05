# MasjidBoard Display Profiles

## Purpose

MasjidPi selects the presentation used on its local HDMI output from attached hardware. Display profile is a runtime property, not a saved user preference and not a synonym for orientation.

## Profiles

### Standard

The `standard` profile is the fallback for every display that is not recognised as MasjidPi appliance hardware.

- uses the responsive landscape MasjidBoard presentation
- launches Cog with `renderer=gles`
- does not rotate the DRM output
- does not require appliance touchscreen hardware

Opening `/masjidboard.html` in a normal browser also uses the standard profile.

### Appliance

The `appliance` profile is the dedicated presentation for the validated Waveshare 7-inch appliance display.

At display-service startup, MasjidPi selects this profile only when both conditions are true:

1. a Waveshare WS170120 USB touchscreen is present (`0eef:0005`); and
2. a connected HDMI connector advertises the native `1024x600` mode.

The appliance profile:

- launches Cog with `renderer=gles,rotation=1`
- renders the current 7-inch interface at the effective 600x1024 portrait viewport
- enables the appliance touch-control UI
- uses the Waveshare libinput calibration matrix `0 -1 1 1 0 0`

If either hardware condition is absent, MasjidPi falls back to the standard profile.

## Remote browser preview

Profile detection controls only the local Cog display runtime. The appliance presentation remains directly accessible for development and troubleshooting from another computer:

```text
/masjidboard.html?profile=appliance
```

This URL forces the appliance frontend presentation but does not rotate the remote computer's display or alter any saved MasjidPi setting.

Normal browser access remains:

```text
/masjidboard.html
```

which renders the standard presentation.

## Persistence

MasjidPi does not persist `standard` or `appliance` as a user preference. Existing selection files written by releases that stored `layout: landscape` or `layout: portrait` are accepted during upgrade because unknown JSON fields are ignored. The obsolete layout field disappears the next time the selection is saved.

Theme, slide duration, Daily Ayah/Hadith/Sunnah visibility, Dua-after-Adhan visibility and Islamic Economic Indicators visibility are persisted preferences shared by both profiles.

## Configuration UI

The Board Display page no longer offers a layout/profile selector. It reports that the local display profile is automatic and continues to configure the shared theme, slide duration and optional content.

A direct Appliance Preview link is provided for browser-based development and troubleshooting.

## Hardware ownership

The display launcher owns local hardware detection and Cog rotation. The backend does not attempt to determine the attached display and the frontend does not infer appliance hardware from viewport dimensions.

Touch calibration is installed as a narrow udev rule matching the validated Waveshare touchscreen. It is independent of the display profile decision so input coordinates are already correct when Cog opens the appliance frontend.

## Future orientation support

Profile and orientation are intentionally separate concepts. The current mapping is:

```text
standard  -> responsive landscape presentation
appliance -> dedicated 7-inch portrait presentation
```

A future conventional monitor/TV portrait presentation can therefore be added without redefining the appliance profile or restoring a saved layout selector.
