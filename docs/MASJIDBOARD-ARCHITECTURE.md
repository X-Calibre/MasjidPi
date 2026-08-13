# MasjidBoard Application Architecture

**Status:** Architectural decision — agreed
**Branch:** `research/masjidboard-live`

## Decision

MasjidBoard will be designed as a **standalone application within the MasjidPi repository**, rather than as a tightly coupled feature of the audio-streaming application.

An end user will be able to choose whether to run:

- Audio streaming only
- MasjidBoard display only
- Both audio streaming and MasjidBoard display

The two applications/subsystems must remain operationally independent.

## Intended Architecture

```text
                         MasjidPi
                            |
              +-------------+-------------+
              |                           |
              v                           v
      Audio Application            MasjidBoard Application
              |                           |
             MPV                    MasjidBoard Live
                                          |
                              +-----------+-----------+
                              |                       |
                         Data / Cache            Display
                              |                       |
                              +-----------+-----------+
                                          |
                                         HDMI
```

The MasjidBoard application should be independently executable and should not require the audio application to be running.

## Application Boundaries

### Audio application

Responsible for the existing MasjidPi audio-streaming functionality, including stream selection, playback, volume and related audio controls.

### MasjidBoard application

Responsible for:

- MasjidBoard Live data retrieval
- Normalisation of upstream data
- Persistent local caching
- Image/media caching
- Board state
- Content scheduling
- Display rendering
- HDMI output
- MasjidBoard-specific configuration and recovery

MasjidBoard must not depend on MPV or the audio playback subsystem.

## Repository Strategy

MasjidBoard remains in the **MasjidPi repository** so that shared infrastructure, configuration, documentation, testing and release management remain together.

It should nevertheless be implemented as a distinct application/package boundary rather than being embedded directly into the audio playback code.

The eventual implementation may expose separate executables/services, for example:

```text
/opt/masjidpi/bin/masjidpi
/opt/masjidpi/bin/masjidboard
```

and corresponding independently managed services:

```text
masjidpi.service
masjidboard.service
```

Exact executable/service names are implementation details and are not yet fixed.

## Configuration

The MasjidPi configuration/UI should allow the end user to enable either application independently or both together.

Conceptually:

```text
Audio Streaming:   [ enabled / disabled ]
MasjidBoard:       [ enabled / disabled ]
```

This configuration must not cause one subsystem to fail simply because the other is disabled or unavailable.

## Failure Isolation

The architecture must support independent recovery:

```text
MasjidBoard failure
    -> Audio continues
    -> MasjidBoard restarts/recover independently

MasjidBoard Live unavailable
    -> Audio continues
    -> MasjidBoard displays cached data where possible

Audio failure
    -> MasjidBoard continues
    -> Audio restarts/recover independently
```

## Display Strategy

MasjidBoard will **not** simply embed or mirror the MasjidBoard Live webpage.

MasjidBoard Live is the primary data/content provider. MasjidPi will implement its own presentation layer using the structured upstream data discovered during the research phase.

The initial goal is **functional and content parity** with MasjidBoard Live: all available information and relevant display behaviour should be represented. The implementation should not unnecessarily reproduce MasjidBoard Live's HTML, CSS or JavaScript architecture.

This gives MasjidPi control over:

- HDMI output
- 1920×1080 display layout
- Content scheduling
- Offline behaviour
- Local caching
- Image caching
- Rendering performance
- Recovery behaviour
- Future MasjidPi-specific display improvements

## Data Flow

```text
MasjidBoard Live
        |
        v
MasjidBoard Provider
        |
        v
Normalised MasjidBoard Model
        |
        +----> Persistent Cache
        |
        v
Board State / Scheduler
        |
        v
MasjidBoard Renderer
        |
        v
HDMI Display
```

The upstream MasjidBoard Live positional API schema must remain isolated inside the provider/normalisation layer.

## Future Extensibility

Although MasjidBoard Live is the primary data source, the provider boundary should allow another data source to be introduced later without redesigning the display system.

The architecture should therefore distinguish between:

```text
Data provider
        -> Normalised board data
        -> Display system
```

rather than coupling the display directly to MasjidBoard Live.

## Consequences

### Benefits

- Audio-only installations remain lightweight and unaffected by HDMI requirements.
- Display-only installations are possible.
- Audio and display can be enabled together.
- Failures are isolated between subsystems.
- MasjidBoard can have its own caching and recovery behaviour.
- The display is not dependent on the MasjidBoard Live website implementation.
- The MasjidBoard application can evolve independently while remaining part of the MasjidPi project.

### Trade-offs

- There is additional application/service architecture compared with simply displaying the website.
- MasjidPi must maintain its own renderer and display behaviour.
- Changes to the upstream MasjidBoard Live data model may require provider updates.

These trade-offs are accepted because reliability, independence and appliance-style operation are more important than minimising implementation effort.

## Implementation Guardrail

This architectural decision does **not** mean production code should be started immediately.

The current research phase must first establish a sufficiently complete and stable MasjidBoard Live data model. Only then should the standalone MasjidBoard application be implemented.
