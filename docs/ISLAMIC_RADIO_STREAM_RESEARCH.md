# South African Islamic Radio Stream Research

**Status:** Active research / Raspberry Pi 4 validation complete for eight streams  
**Last validated:** 2026-08-27  
**Purpose:** Evaluate South African Islamic radio stations for inclusion in MasjidPi Listen as secondary audio sources.

## Objective

MasjidPi Listen currently focuses on live audio streams from individual masjids. A secondary catalogue of continuous Islamic radio stations would provide users with an alternative source of Islamic programming.

This investigation identifies suitable South African Islamic radio stations, resolves direct audio endpoints where possible, records measured stream properties, and validates compatibility with MasjidPi's `mpv`-based playback architecture.

## Current Validation Summary

Eight stations now have direct endpoints that have been successfully tested with `mpv` on both an x86_64 MasjidPi development environment and the aarch64 Raspberry Pi 4 test platform.

| Station | Region | Direct endpoint | Codec | Sample rate | Channels | Bitrate | Container | Pi 4 5-minute soak |
|---|---|---|---|---:|---|---:|---|---|
| Channel Islam International | Gauteng | `https://edge.iono.fm/xice/109_medium.aac` | AAC | 44.1 kHz | Stereo | ~69.5 kbps | AAC | PASS |
| Radio Islam International | Gauteng | `https://cast1.my-control-panel.com/proxy/netmoham/radioislam.mp3` | MP3 | 48 kHz | Mono | 64 kbps | MP3 | PASS |
| Radio 786 | Western Cape | `https://stream.krypton.co.za/proxy/radio786?mp=/stream` | MP3 | 44.1 kHz | Mono | 56 kbps | MP3 | PASS |
| Voice of the Cape (VOC FM) | Western Cape | `https://streaming.fabrik.fm/voc/echocast/audio/low/index.m3u8` | AAC | 44.1 kHz | Stereo | Not reported | HLS | PASS |
| Radio Al Ansaar | KwaZulu-Natal | `https://edge.iono.fm/xice/467_medium.aac` | AAC | 44.1 kHz | Stereo | ~80.1 kbps | AAC | PASS |
| Markaz Sahaba Online Radio | KwaZulu-Natal | `http://zas4.ndx.co.za:9088/stream` | MP3 | 44.1 kHz | Mono | 160 kbps | MP3 | PASS |
| Sirius FM 105.7 | Gauteng | `http://s8.voscast.com:7112/;` | MP3 | 44.1 kHz | Mono | 64 kbps | MP3 | PASS |
| Salaamedia | Gauteng | `http://capeant.antfarm.co.za:1935/salaam/salaam.stream/playlist.m3u8` | AAC | 44.1 kHz | Stereo | ~129.4 kbps | HLS | PASS |

The remaining candidates are not yet ready for production catalogue inclusion:

| Station | Status |
|---|---|
| Radio 786 fallback | Primary stream resolved and validated; genuine independent Stream 2/fallback remains unresolved |
| IFM 88.3 | Current online player confirmed; raw stream unresolved |
| Islam Alive Radio | Lower priority; not proposed for initial catalogue |

## Validation Environments

### MasjidPi-Dev

The resolved streams were tested using `mpv` with audio output disabled (`--ao=null`) to isolate network, transport, demuxing and decoding behaviour.

The original seven streams passed a 30-second screening test. Salaamedia subsequently passed a five-minute `mpv` soak test on MasjidPi-Dev.

`ffprobe` was used to measure codec, sample rate, channel layout, bitrate and container information. The measured properties are recorded in the validation table above.

### Raspberry Pi 4 — MasjidPi-Test

The resolved streams were tested on the Raspberry Pi test platform:

```text
Linux MasjidPi-Test 6.18.39+rpt-rpi-v8
Debian / Raspberry Pi kernel
Architecture: aarch64
```

Each stream was played through `mpv` with `--ao=null` for a five-minute soak period.

Results:

| Station | Duration | Result |
|---|---:|---|
| Markaz Sahaba | 300 s | PASS |
| Channel Islam International | 301 s | PASS |
| Voice of the Cape | 300 s | PASS |
| Radio Islam International | 300 s | PASS |
| Radio 786 | 300 s | PASS |
| Radio Al Ansaar | 301 s | PASS |
| Sirius FM | 300 s | PASS |
| Salaamedia | 300 s | PASS |

This confirms successful operation through the Raspberry Pi network, TLS/HTTP/HLS transport, `mpv`/FFmpeg demuxing and audio decoding stack.

Radio 786 emitted the following message when the timed test ended:

```text
[ffmpeg] tls: Error decoding the received TLS packet.
```

The stream nevertheless remained operational for the complete 300-second test. Similar TLS teardown output was observed during forced termination of another persistent stream during development testing. At present this is classified as non-fatal connection teardown noise rather than a playback failure. It should be monitored during longer uninterrupted testing.

## Channel Islam International

Preferred direct endpoint:

```text
https://edge.iono.fm/xice/109_medium.aac
```

Measured properties:

```text
Codec:       AAC
Sample rate: 44100 Hz
Channels:    Stereo
Bitrate:     ~69.5 kbps
Container:   AAC
```

**Assessment:** Validated candidate.

## Radio Islam International

Direct endpoint:

```text
https://cast1.my-control-panel.com/proxy/netmoham/radioislam.mp3
```

Measured properties:

```text
Codec:       MP3
Sample rate: 48000 Hz
Channels:    Mono
Bitrate:     64 kbps
Container:   MP3
```

**Assessment:** Validated candidate.

## Radio 786

Validated primary endpoint:

```text
https://stream.krypton.co.za/proxy/radio786?mp=/stream
```

The underlying HTTP stream is also represented as:

```text
http://stream.krypton.co.za:8040/stream/
```

These appear to be two access paths to the same Krypton source and should not be treated as independent fallback services.

Measured properties:

```text
Codec:       MP3
Sample rate: 44100 Hz
Channels:    Mono
Bitrate:     56 kbps
Container:   MP3
```

The genuine independent Stream 2/fallback endpoint exposed by Radio 786's web player remains to be identified.

**Assessment:** Primary stream validated; independent fallback unresolved.

## Voice of the Cape (VOC FM)

Current endpoint:

```text
https://streaming.fabrik.fm/voc/echocast/audio/low/index.m3u8
```

Measured properties:

```text
Codec:       AAC
Sample rate: 44100 Hz
Channels:    Stereo
Bitrate:     Not reported by ffprobe
Container:   HLS
```

VOC confirms that the proposed radio catalogue does not need to be restricted to Icecast/SHOUTcast-style streams.

**Assessment:** Validated candidate.

## Radio Al Ansaar

Preferred direct endpoint:

```text
https://edge.iono.fm/xice/467_medium.aac
```

Measured properties:

```text
Codec:       AAC
Sample rate: 44100 Hz
Channels:    Stereo
Bitrate:     ~80.1 kbps
Container:   AAC
```

The older SimpleStreaming endpoint:

```text
https://al-ansaar.simplestreaming.co.za/listen/al-ansaar_radio/radio.mp3
```

should no longer be treated as the preferred long-term endpoint because the station's current first-party web presence uses iono.fm. It may still be investigated as a possible fallback if it remains operational and independent.

**Assessment:** Preferred iono.fm endpoint validated.

## Markaz Sahaba Online Radio

Direct stream:

```text
http://zas4.ndx.co.za:9088/stream
```

Measured properties:

```text
Codec:       MP3
Sample rate: 44100 Hz
Channels:    Mono
Bitrate:     160 kbps
Container:   MP3
```

At 160 kbps, Markaz Sahaba currently has the highest measured bitrate among the validated stations.

**Assessment:** Validated candidate.

## Sirius FM 105.7

Direct endpoint:

```text
http://s8.voscast.com:7112/;
```

Measured properties:

```text
Codec:       MP3
Sample rate: 44100 Hz
Channels:    Mono
Bitrate:     64 kbps
Container:   MP3
```

**Assessment:** Validated candidate.

## Salaamedia

Salaamedia's current direct HLS endpoint has now been verified as:

```text
http://capeant.antfarm.co.za:1935/salaam/salaam.stream/playlist.m3u8
```

The endpoint returns HTTP 200 from Wowza Streaming Engine 4.7.6 with MIME type:

```text
application/vnd.apple.mpegurl
```

The HLS master playlist returned during validation was:

```text
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=128285,CODECS="mp4a.40.2"
chunklist_w1784185512.m3u8
```

The child playlist name is dynamically generated and should not be stored in the catalogue. The stable master playlist URL above should be used.

Measured properties:

```text
Codec:       AAC
Sample rate: 44100 Hz
Channels:    Stereo
Bitrate:     ~129.4 kbps
Container:   HLS
Protocol:    HTTP
```

The advertised HLS bandwidth of 128,285 bps closely matches the approximately 129,370 bps measured by `ffprobe`.

The endpoint passed a five-minute `mpv` soak test on both MasjidPi-Dev and the Raspberry Pi 4 test platform.

The stream currently uses plain HTTP rather than HTTPS. This is not a blocker for `mpv` playback but should be recorded as part of the station's transport characteristics.

**Assessment:** Validated candidate.

## IFM 88.3

IFM 88.3 continues to provide a first-party Listen Live facility. External listings report an Internet stream around 128 kbps, but a sufficiently trustworthy raw endpoint has not yet been identified.

Current state:

```text
Raw URL:     Unresolved
Codec:       Unknown
Container:   Unknown
Bitrate:     128 kbps reported externally; not measured
```

**Assessment:** Raw endpoint unresolved.

## Islam Alive Radio

Islam Alive Radio appears in South African Internet-radio listings but currently has less first-party evidence and infrastructure available than the other candidates.

**Assessment:** Do not include in the initial catalogue. Reassess later.

## Proposed Initial Catalogue

### Validated candidates

1. Channel Islam International
2. Radio Islam International
3. Salaamedia
4. Voice of the Cape
5. Radio 786 — primary stream
6. Radio Al Ansaar
7. Markaz Sahaba Online Radio
8. Sirius FM 105.7

### Pending endpoint discovery

9. IFM 88.3

Radio 786's genuine independent fallback stream also remains to be resolved.

Islam Alive Radio should remain outside the initial catalogue pending further investigation.

## Proposed MasjidPi Architecture

Radio stations should not be represented as ordinary masjid catalogue entries.

MasjidPi Listen should distinguish between:

- **Masjids** — individual masjid live streams
- **Radio Stations** — continuous Islamic radio broadcasters

The validation work shows that radio streams may use different transports and codecs, including MP3, raw AAC and HLS/AAC. MasjidPi should continue to delegate transport and codec handling to `mpv` rather than restricting catalogue entries to one stream type.

## Fallback Stream Support

The radio catalogue should support more than one endpoint per station.

Radio 786 provides the clearest first-party justification because its own player exposes multiple streams specifically so listeners can switch when one is unavailable or buffering. CII also has multiple provider variants, and Radio Al Ansaar has changed preferred streaming infrastructure over time.

A simple ordered stream list is preferred unless later requirements justify per-stream metadata:

```yaml
type: radio
name: Channel Islam International
country: ZA
region: Gauteng
streams:
  - https://edge.iono.fm/xice/109_medium.aac
  - https://edge.iono.fm/xice/109_medium.mp3
```

Array order should define priority.

Fallback should be triggered by a genuine connection or playback failure, not by a brief pause or transient buffering event.

Where two URLs are merely different proxy/access paths to the same upstream service, they should not be represented as independent resilience fallbacks unless testing demonstrates a real availability benefit.

## Potential Secondary-Stream Behaviour

A future enhancement could allow a selected radio station to operate as a secondary audio source when a selected masjid is not broadcasting:

```text
Selected Masjid
      |
      +-- Broadcasting ----> Play masjid
      |
      +-- Not broadcasting -> Selected radio station
```

This behaviour should not be coupled to the initial radio catalogue implementation.

Recommended implementation sequence:

1. Complete endpoint discovery and hardware validation.
2. Establish a radio catalogue with ordered stream support.
3. Add Radio Stations as a separate selectable category.
4. Reuse the existing MasjidPi `mpv` playback architecture.
5. Validate normal manual radio playback.
6. Consider automatic fallback/secondary-stream behaviour as a separate enhancement.

## Remaining Technical Validation

The following work remains:

- Extract Radio 786's genuine independent Stream 2/fallback endpoint.
- Identify IFM 88.3's current raw stream.
- Test the eight validated streams on Raspberry Pi 3 B hardware.
- Perform longer soak testing, prioritising VOC, Radio 786, CII, Salaamedia and Markaz Sahaba.
- Measure Raspberry Pi CPU and memory utilisation during playback.
- Validate playback through the production ALSA output path rather than `--ao=null`.
- Test reconnection behaviour after deliberate network interruption.
- Determine whether Radio Al Ansaar's older SimpleStreaming endpoint remains useful as an independent fallback.
- Determine whether ICY/current-programme metadata should be retained for a future Now Playing feature.

## Validation Criteria

For production catalogue entries, validation should ultimately record:

- DNS resolution
- HTTP/HTTPS/TLS compatibility
- redirect behaviour
- MIME type
- direct versus player-only URL
- stream transport/container
- audio codec
- bitrate
- sample rate
- channel layout
- startup latency
- short playback stability
- extended playback stability
- reconnection behaviour
- Raspberry Pi CPU and memory impact
- Raspberry Pi 4 compatibility
- Raspberry Pi 3 B compatibility
- production ALSA output compatibility

## Conclusion

The investigation has progressed from preliminary stream discovery to successful hardware validation.

Eight South African Islamic radio stations now have direct endpoints that have successfully played through MasjidPi's `mpv`/FFmpeg stack on both the x86_64 development environment and an aarch64 Raspberry Pi 4. The test set includes MP3, raw AAC and HLS/AAC streams.

Salaamedia's previously uncertain Antfarm endpoint has now been confirmed as a valid HLS/AAC stream and has passed five-minute playback testing on both environments.

This provides strong evidence that a dedicated Radio Stations catalogue is technically viable in MasjidPi Listen.

Production implementation should wait until the remaining endpoint discovery is completed and the validated streams have been tested on Raspberry Pi 3 B hardware, through the production ALSA output path, and under longer/reconnection test conditions.
