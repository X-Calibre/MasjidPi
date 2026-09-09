# MasjidBoard Display Profiles

## Purpose

MasjidPi selects the presentation used on its local display from attached
hardware. The display profile is a runtime property, not a saved preference and
not a synonym for orientation.

## Supported profiles

### Standard

The `standard` profile is the fallback for every display that is not the
supported Raspberry Pi Touch Display 2.

- uses the responsive landscape MasjidBoard presentation;
- launches Cog with `renderer=gles`;
- does not rotate the DRM output; and
- does not enable appliance touch controls.

Opening `/masjidboard.html` in a normal browser also uses this profile.

### Appliance 720

The `appliance-720` profile is the native portrait presentation for the
official 7-inch Raspberry Pi Touch Display 2. MasjidPi selects it when a
connected DSI connector advertises the panel's native `720x1280` mode.

This profile:

- launches Cog with `renderer=gles` and no output rotation;
- renders at the native 720x1280 portrait viewport;
- enables the slideshow and appliance touch controls;
- provides persistent backlight brightness and colour-temperature correction;
  and
- relies on the Raspberry Pi OS DSI display and touch drivers.

First-run setup and Change Wi-Fi always use the dedicated 720x1280 setup
presentation, including its larger controls, picker sheets and on-screen
keyboard.

## Retired 600x1024 profile

The former `appliance` profile for the rotated 1024x600 Waveshare HDMI panel is
no longer supported. Its USB/HDMI detector, Cog rotation path, touch
calibration contract, frontend route and configuration preview have been
removed. A 1024x600 HDMI display now receives the responsive `standard`
profile.

A direct browser request using the retired
`/masjidboard.html?profile=appliance` value also falls back to the standard
presentation.

## Remote browser preview

The supported portrait presentation can be previewed from another computer:

```text
/masjidboard.html?profile=appliance-720
```

The query parameter changes only that browser view. It does not alter saved
settings or rotate the remote display.

### Visual-review gallery

Open the dedicated review launcher to inspect every supported Touch Display 2
state from one page:

```text
/appliance-720-demo.html
```

The launcher embeds an exact 720x1280 viewport and links to:

- the live slideshow, all synthetic community-card types, Dua-after-Adhan and a
  detailed Jumu'ah fixture;
- each touch-control tab;
- Masjid playback, Radio playback and Radio-resume notifications; and
- network selection, password, hidden-network, success, location and Masjid
  first-run screens.

Card-fixture mode pauses automatic rotation; use the slide dots or swipe to
inspect each card. Setup, control and notification fixtures are read-only and
do not change Wi-Fi, playback, display settings or saved preferences. The live
slideshow and non-community cards still use the configured device data.

Normal browser access remains:

```text
/masjidboard.html
```

## Persistence

MasjidPi does not persist `standard` or `appliance-720` as a user preference.
Old selection files containing `layout: landscape` or `layout: portrait`
remain upgrade-safe because unknown JSON fields are ignored. The obsolete field
disappears the next time the selection is saved.

Theme, slide duration, Daily Ayah/Hadith/Sunnah visibility, Dua-after-Adhan
visibility and Islamic Economic Indicators visibility are shared preferences.

## Configuration UI

The Board Display page does not offer a layout selector. It reports automatic
local profile selection and provides one direct Touch Display 2 preview link.

## Hardware ownership

The display launcher owns local hardware detection. The backend does not infer
attached display hardware, and the frontend does not infer a profile from
viewport dimensions.

The current mapping is:

```text
standard      -> responsive landscape presentation
appliance-720 -> 720x1280 Touch Display 2 portrait presentation
```
