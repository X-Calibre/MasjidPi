# South African Islamic Radio Stream Research

**Status:** Active research / Raspberry Pi 4 validation complete for seven streams  
**Last validated:** 2026-08-27  
**Purpose:** Evaluate South African Islamic radio stations for inclusion in MasjidPi Listen as secondary audio sources.

## Objective

MasjidPi Listen currently focuses on live audio streams from individual masjids. A secondary catalogue of continuous Islamic radio stations would provide users with an alternative source of Islamic programming.

This investigation identifies suitable South African Islamic radio stations, resolves direct audio endpoints where possible, records measured stream properties, and validates compatibility with MasjidPi's `mpv`-based playback architecture.

## Current Validation Summary

Seven stations now have direct endpoints that have been successfully tested with `mpv` on both an x86_64 MasjidPi development environment and the aarch64 Raspberry Pi 4 test platform.

| Station | Region | Direct endpoint | Codec | Sample rate | Channels | Bitrate | Container | Pi 4 5-minute soak |
|---|---|---|---|---:|---|---:|---|---|
| Channel Islam International | Gauteng | `https://edge.iono.fm/xice/109_medium.aac` | AAC | 44.1 kHz | Stereo | ~69.5 kbps | AAC | PASS |
| Radio Islam International | Gauteng | `https://cast1.my-control-panel.com/proxy/netmoham/radioislam.mp3` | MP3 | 48 kHz | Mono | 64 kbps | MP3 | PASS |
| Radio 786 | Western Cape | `https://stream.krypton.co.za/proxy/radio786?mp=/stream` | MP3 | 44.1 kHz | Mono | 56 kbps | MP3 | PASS |
| Voice of the Cape (VOC FM) | Western Cape | `https://streaming.fabrik.fm/voc/echocast/audio/low/index.m3u8` | AAC | 44.1 kHz | Stereo | Not reported | HLS | PASS |
| Radio Al Ansaar | KwaZulu-Natal | `https://edge.iono.fm/xice/467_medium.aac` | AAC | 44.1 kHz | Stereo | ~80.1 kbps | AAC | PASS |
| Markaz Sahaba Online Radio | KwaZulu-Natal | `http://zas4.ndx.co.za:9088/stream` | MP3 | 44.1 kHz | Mono | 160 kbps | MP3 | PASS |
| Sirius FM 105.7 | Gauteng | `http://s8.voscast.com:7112/;` | MP3 | 44.1 kHz | Mono | 64 kbps | MP3 | PASS |

The remaining candidates are not yet ready for production catalogue inclusion:

| Station | Status |
|---|---|
| Salaamedia | Current first-party player confirmed; current raw stream still unresolved |
| Radio 786 fallback | Primary stream resolved and validated; genuine independent Stream 2/fallback remains unresolved |
| IFM 88.3 | Current online player confirmed; raw stream unresolved |
| Islam Alive Radio | Lower priority; not proposed for initial catalogue |

## Validation Environments

### MasjidPi-Dev

The seven resolved streams were tested using `mpv` for 30 seconds with audio output disabled (`--ao=null`) to isolate network, transport, demuxing and decoding behaviour.

All seven streams passed.

`ffprobe` was then used against the same endpoints to measure codec, sample rate, channel layout, bitrate and container information. The measured properties are recorded in the validation table above.

### Raspberry Pi 4 — MasjidPi-Test

The same seven streams were tested on the Raspberry Pi test platform:

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

This confirms successful operation through the Raspberry Pi network, TLS/HTTP/HLS transport, `mpv`/FFmpeg demuxing and audio decoding stack.

Radio 786 emitted the following message when the timed test ended:

```text
[ffmpeg] tls: Error decoding the received TLS packet.
```

The stream nevertheless remained operational for the complete 300-second test. Similar TLS teardown output was observed during forced termination of another persistent stream during development testing. At present this is classified as non-fatal connection teardown noise rather than a playback failure. It should be monitored during longer uninterrupted testing.

## Channel Islam International

Channel Islam International currently uses iono.fm for live streaming.

The preferred direct endpoint is:

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

The endpoint passed both the x86_64 development test and the Raspberry Pi 4 five-minute soak test.

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

The endpoint passed both the x86_64 development test and Raspberry Pi 4 five-minute soak test.

**Assessment:** Validated candidate.

## Radio 786

Radio 786's official player provides multiple streams and advises listeners to switch streams when one is buffering or unavailable.

The validated primary endpoint is:

```text
https://stream.krypton.co.za/proxy/radio786?mp=/stream
```

The underlying HTTP stream is also represented as:

```text
http://stream.krypton.co.za:8040/stream/
```

These appear to be two access paths to the same Krypton source and should not be treated as independent fallback services.

Measured properties of the validated HTTPS endpoint:

```text
Codec:       MP3
Sample rate: 44100 Hz
Channels:    Mono
Bitrate:     56 kbps
Container:   MP3
```

The HTTPS endpoint passed the Raspberry Pi 4 five-minute soak test. A non-fatal TLS message was emitted when the timed test terminated.

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

VOC is the current architectural outlier because it uses HLS rather than a continuous MP3/AAC HTTP stream. `mpv` successfully handled the HLS stream on both development and Raspberry Pi 4 environments, including a full five-minute soak test.

This confirms that the proposed radio catalogue does not need to be restricted to Icecast/SHOUTcast-style streams.

**Assessment:** Validated candidate.

## Radio Al Ansaar

Radio Al Ansaar's current first-party service uses iono.fm. The preferred direct endpoint has now been experimentally confirmed as:

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

The endpoint passed both development and Raspberry Pi 4 playback testing.

The older SimpleStreaming endpoint:

```text
https://al-ansaar.simplestreaming.co.za/listen/al-ansaar_radio/radio.mp3
```

should no longer be treated as the preferred long-term endpoint because the station's current first-party web presence uses iono.fm. It may still be investigated as a possible fallback if it remains operational and independent.

**Assessment:** Preferred iono.fm endpoint validated.

## Markaz Sahaba Online Radio

Markaz Sahaba's official web presence points listeners through NetDynamix. The underlying direct stream has been resolved as:

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

The endpoint passed both development and Raspberry Pi 4 testing.

At 160 kbps, Markaz Sahaba currently has the highest measured bitrate among the validated stations and therefore consumes substantially more bandwidth than the 56–80 kbps streams in the rest of the group.

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

This is a SHOUTcast-style stream. It passed both development and Raspberry Pi 4 playback testing.

**Assessment:** Validated candidate.

## Salaamedia

Salaamedia's current official Listen Live service remains active.

The historically indexed endpoint is:

```text
http://capeant.antfarm.co.za:1935/salaam/salaam.stream/playlist.m3u8
```

This should currently be treated as a legacy candidate rather than a production endpoint. It has not been sufficiently tied to Salaamedia's current first-party player.

The raw stream requested by the current first-party player still needs to be extracted and validated.

**Assessment:** Raw current endpoint unresolved.

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

The current target catalogue is:

### Validated candidates

1. Channel Islam International
2. Radio Islam International
3. Voice of the Cape
4. Radio 786 — primary stream
5. Radio Al Ansaar
6. Markaz Sahaba Online Radio
7. Sirius FM 105.7

### Pending endpoint discovery

8. Salaamedia
9. IFM 88.3

Radio 786's genuine independent fallback stream also remains to be resolved.

Islam Alive Radio should remain outside the initial catalogue pending further investigation.

## Proposed MasjidPi Architecture

Radio stations should not be represented as ordinary masjid catalogue entries.

MasjidPi Listen should distinguish between:

- **Masjids** — individual masjid live streams
- **Radio Stations** — continuous Islamic radio broadcasters

The validation work also shows that radio streams may use different transports and codecs, including MP3, raw AAC and HLS/AAC. MasjidPi should therefore continue to delegate stream handling to `mpv` rather than attempting to restrict catalogue entries to one transport type.

## Fallback Stream Support

The radio catalogue should support more than one endpoint per station.

Radio 786 provides the clearest first-party justification: its own player exposes multiple streams specifically so listeners can switch when one is unavailable or buffering. CII also has multiple provider variants, and Radio Al Ansaar has changed preferred streaming infrastructure over time.

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

Expected player behaviour:

```text
Try streams[0]
      |
      +-- playback starts --> use primary
      |
      +-- connection/playback failure
                  |
                  v
            Try streams[1]
                  |
                  +-- playback starts --> use fallback
                  |
                  +-- failure --> station unavailable
```

Fallback should be triggered by a genuine connection or playback failure, not by a brief pause or transient buffering event. Aggressive switching could otherwise cause unnecessary movement between healthy sources.

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
- Verify Salaamedia's current first-party raw stream.
- Identify IFM 88.3's current raw stream.
- Test the seven validated streams on Raspberry Pi 3 B hardware.
- Perform longer soak testing, prioritising VOC, Radio 786, CII and Markaz Sahaba.
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

Seven South African Islamic radio stations now have direct endpoints that have successfully played through MasjidPi's `mpv`/FFmpeg stack on both the x86_64 development environment and an aarch64 Raspberry Pi 4. The test set includes MP3, raw AAC and HLS/AAC streams.

This provides strong evidence that a dedicated Radio Stations catalogue is technically viable in MasjidPi Listen.

Production implementation should wait until the remaining endpoint discovery is completed and the validated streams have been tested on Raspberry Pi 3 B hardware, through the production ALSA output path, and under longer/reconnection test conditions.
