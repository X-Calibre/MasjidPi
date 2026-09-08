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

### Appliance 720

The `appliance-720` profile is the native portrait presentation for the official
7-inch Raspberry Pi Touch Display 2. MasjidPi selects it when a connected DSI
connector advertises the panel's native `720x1280` mode.

This profile:

- launches Cog with `renderer=gles` and no output rotation;
- renders the appliance interface at the native 720x1280 portrait viewport;
- enables the same slideshow and touch controls as the existing appliance profile;
- relies on the Raspberry Pi OS DSI display and touch drivers, without the
  Waveshare USB calibration rule.

First-run setup and the reopened Change Wi-Fi workflow retain the detected
profile. The Touch Display 2 receives a dedicated 720x1280 setup presentation,
including larger network rows, form controls, picker sheets, on-screen keyboard,
masjid choices and success screen; the 600x1024 setup presentation is unchanged.

The DSI mode is checked before the Waveshare HDMI profile so the Touch Display 2
wins if both displays happen to be attached during startup.

## Remote browser preview

Profile detection controls only the local Cog display runtime. The appliance presentation remains directly accessible for development and troubleshooting from another computer:

```text
/masjidboard.html?profile=appliance
```

The Touch Display 2 presentation can be previewed at:

```text
/masjidboard.html?profile=appliance-720
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
standard      -> responsive landscape presentation
appliance     -> 600x1024 Waveshare portrait presentation
appliance-720 -> 720x1280 Touch Display 2 portrait presentation
```

A future conventional monitor/TV portrait presentation can therefore be added without redefining the appliance profile or restoring a saved layout selector.
