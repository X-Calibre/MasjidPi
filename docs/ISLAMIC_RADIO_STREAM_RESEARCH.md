# South African Islamic Radio Stream Research

**Status:** Research / preliminary  
**Purpose:** Evaluate South African Islamic radio stations for inclusion in MasjidPi Listen as secondary audio sources.

## Objective

MasjidPi Listen currently focuses on live audio streams from individual masjids. A secondary catalogue of continuous Islamic radio stations would provide users with an alternative source of Islamic programming.

This investigation identifies suitable South African Islamic radio stations and evaluates their potential compatibility with MasjidPi's `mpv`-based playback architecture.

Stream endpoints documented here should be considered preliminary until they have been validated on supported Raspberry Pi platforms.

## Candidate Stations

| Station | Region | Stream Status | Likely Format | Assessment |
|---|---|---|---|---|
| Channel Islam International | Gauteng | Online stream exists; raw endpoint still required | Unknown | Include after validation |
| Radio Islam International | Gauteng | Direct stream identified | MP3 | Strong candidate |
| Salaamedia | Gauteng | Stream identified; requires validation | HLS | Strong candidate |
| Voice of the Cape (VOC FM) | Western Cape | Direct stream identified | HLS/AAC | Strong candidate |
| Radio 786 | Western Cape | Online streams available; current raw endpoint requires validation | Icecast/unknown | Strong candidate |
| Radio Al Ansaar | KwaZulu-Natal | Direct stream identified | MP3 | Strong candidate |
| Markaz Sahaba Online Radio | KwaZulu-Natal | Online player identified; raw endpoint required | Unknown | Promising |
| Sirius FM 105.7 | Gauteng | Online stream exists; raw endpoint required | Unknown | Promising |
| IFM 88.3 | Eastern Cape | Online player exists; raw endpoint required | Unknown | Promising |
| Islam Alive Radio | South Africa | Third-party online stream references found | Unknown | Lower priority |

## Radio Islam International

Radio Islam International is an established Islamic broadcaster based in Gauteng.

A direct MP3 stream has been identified:

```text
https://cast1.my-control-panel.com/proxy/netmoham/radioislam.mp3
```

The endpoint serves `audio/mpeg`, making it well suited to direct playback using `mpv`.

**Assessment:** Strong candidate for inclusion.

## Radio Al Ansaar

Radio Al Ansaar is an established Muslim community broadcaster serving KwaZulu-Natal, including Durban and Pietermaritzburg.

A direct MP3 endpoint has been identified:

```text
https://al-ansaar.simplestreaming.co.za/listen/al-ansaar_radio/radio.mp3
```

The station's official web presence also uses iono.fm infrastructure. Both sources should be tested to determine which provides the most stable long-term endpoint.

**Assessment:** Strong candidate for inclusion.

## Voice of the Cape (VOC FM)

Voice of the Cape is an established Muslim community broadcaster based in Cape Town.

A current HLS endpoint has been identified:

```text
https://streaming.fabrik.fm/voc/echocast/audio/low/index.m3u8
```

VOC has used different streaming providers previously. This indicates that radio catalogue endpoints should be designed so that they can be updated independently when providers change.

`mpv` supports HLS playback.

**Assessment:** Strong candidate, subject to Raspberry Pi playback validation.

## Radio 786

Radio 786 is an established Cape Town community broadcaster.

Its official website currently provides multiple online streams, including a fallback stream.

A historically/currently indexed endpoint is:

```text
https://stream.krypton.co.za/proxy/radio786
```

This endpoint should not yet be considered production-ready because its current status has not been sufficiently verified.

The active endpoints used by Radio 786's current web player should be extracted and tested.

**Assessment:** Strong candidate; endpoint validation required.

## Salaamedia

Salaamedia is an established South African Islamic media organisation providing continuous online audio programming.

An HLS endpoint has been identified:

```text
http://capeant.antfarm.co.za:1935/salaam/salaam.stream/playlist.m3u8
```

The endpoint was obtained through external stream listings rather than directly from the current Salaamedia player and may therefore be stale.

The current first-party player endpoint should be identified before inclusion.

**Assessment:** Strong candidate; endpoint validation required.

## Markaz Sahaba Online Radio

Markaz Sahaba operates a digital Islamic radio service from Durban with worldwide streaming.

Its official live page directs listeners to:

```text
https://ndstream.net/markazsahaba
```

This appears to be a player or redirect rather than necessarily the underlying raw audio stream.

The actual audio endpoint needs to be extracted.

**Assessment:** Promising candidate.

## Sirius FM 105.7

Sirius FM broadcasts from Springs, Gauteng, on 105.7 FM and also provides internet streaming.

Programming includes Qur'an recitation, Islamic lectures, motivational programming, naats/nasheeds and community content.

A stable raw audio endpoint has not yet been identified.

**Assessment:** Promising candidate.

## IFM 88.3

IFM 88.3 is a community broadcaster in the Eastern Cape with faith, community, educational, news and current-affairs programming.

The station provides an online Listen Live facility, but the underlying raw stream still needs to be identified.

Its inclusion would also improve the geographical coverage of the MasjidPi radio catalogue.

**Assessment:** Promising candidate.

## Channel Islam International

Channel Islam International (CII) is an established South African Islamic broadcaster and should be considered a core candidate.

The station continues to be distributed online, but its web and streaming infrastructure has changed over time.

A sufficiently reliable current raw stream endpoint has not yet been identified.

**Assessment:** Intended for inclusion once the current stream is identified and validated.

## Islam Alive Radio

Islam Alive Radio appears in current South African internet-radio listings and is reported to have operated since approximately 2023.

At present there is less first-party evidence and infrastructure available than for the other candidate stations.

**Assessment:** Do not include in the initial catalogue. Reassess later.

## Proposed Initial Catalogue

Subject to technical validation, the target South African Islamic radio catalogue is:

1. Channel Islam International
2. Radio Islam International
3. Salaamedia
4. Voice of the Cape
5. Radio 786
6. Radio Al Ansaar
7. Markaz Sahaba Online Radio
8. Sirius FM 105.7
9. IFM 88.3

Islam Alive Radio should remain outside the initial catalogue pending further investigation.

## Proposed MasjidPi Architecture

Radio stations should not be represented as ordinary masjid catalogue entries.

MasjidPi Listen should distinguish between:

- **Masjids** — individual masjid live streams
- **Radio Stations** — continuous Islamic radio broadcasters

Internally, a radio entry may eventually contain metadata similar to:

```yaml
type: radio
name: Radio Islam International
country: ZA
region: Gauteng
stream:
  primary: https://example.org/live.mp3
  fallback: https://example.org/backup.mp3
format: mp3
```

The exact schema should only be decided after the stream-validation work is complete.

The UI does not need to expose this implementation detail. Users could simply select a **Radio Stations** section and play a station using the existing MasjidPi Listen playback controls.

## Potential Secondary-Stream Behaviour

A future enhancement could allow a selected radio station to operate as a secondary or fallback audio source.

For example:

```text
Selected Masjid
      |
      +-- Broadcasting ----> Play masjid
      |
      +-- Not broadcasting -> Radio station
```

This behaviour should **not** be coupled to the initial radio catalogue implementation.

The recommended implementation sequence is:

1. Establish a validated radio catalogue.
2. Add Radio Stations as a separate selectable category.
3. Reuse the existing MasjidPi playback architecture where practical.
4. Validate normal manual radio playback.
5. Consider automatic fallback/secondary-stream behaviour as a separate enhancement.

## Technical Validation Required

Before any stream is added to the production catalogue, it should be tested using the same `mpv` stack used by MasjidPi.

Validation should cover:

- DNS resolution
- HTTPS/TLS compatibility
- HTTP redirects
- Direct versus player-only URLs
- Stream format/container
- Audio codec
- Bitrate
- Time to first audio
- Continuous playback stability
- Reconnection behaviour
- Raspberry Pi CPU and memory impact
- Raspberry Pi 3 B compatibility
- Raspberry Pi 4 compatibility

Particular attention should be paid to HLS streams because they behave differently from continuous MP3/Icecast streams.

## Outstanding Investigation

- Identify the current direct Channel Islam International stream.
- Extract Radio 786's current primary and fallback streams.
- Verify Salaamedia's current first-party stream.
- Extract the underlying Markaz Sahaba audio stream.
- Identify the Sirius FM raw stream.
- Identify the IFM 88.3 raw stream.
- Verify Radio Al Ansaar's preferred long-term endpoint.
- Test all candidate streams with `mpv`.
- Test validated streams on Raspberry Pi hardware.
- Record codec, format, bitrate and connection behaviour.
- Determine whether fallback URLs should be supported by the radio catalogue schema.

## Conclusion

There is sufficient evidence to justify adding a dedicated South African Islamic radio catalogue to MasjidPi Listen.

Nine stations currently warrant further technical validation, with Radio Islam International, Radio Al Ansaar and Voice of the Cape already presenting particularly strong candidates.

No production catalogue changes should be made until direct stream endpoints have been identified and playback-tested on supported MasjidPi hardware.
