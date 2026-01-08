# nats

Go NATS client for the UAV simulator.

## Why

Connects the simulator to NATS for pub/sub messaging. Publishes drone telemetry and subscribes to control commands, enabling AI-driven flight control and remote operation.

## How

This is Go code, part of the main module. Imported as:

```go
import natsclient "drone-simulator/systems/nats"
```

Wired in via `main.go` with the `-nats-url` flag.

## Commands

| Subject | Payload | Description |
|---------|---------|-------------|
| `drone.<id>.arm` | `''` | Arm the drone |
| `drone.<id>.disarm` | `''` | Disarm the drone |
| `drone.<id>.takeoff` | `{"altitude": 5}` | Take off to altitude (m) |
| `drone.<id>.land` | `''` | Land at current position |
| `drone.<id>.goto` | `{"x": 0, "y": 10, "z": 0}` | Fly to altitude (lateral X/Z not yet supported) |
| `drone.<id>.input` | `{"throttle": 0.5, ...}` | Direct control |
| `drone.<id>.mode` | `{"mode": "Hover"}` | Set flight mode |
| `drone.<id>.stop` | `''` | Emergency stop |

## Telemetry

Published at 10Hz on `drone.<id>.telemetry`:

```json
{
  "id": 0,
  "timestamp": 1767853375086,
  "position": {"x": 0, "y": 5.48, "z": 0},
  "velocity": {"x": 0, "y": -1.36, "z": 0},
  "rotation": {"x": 0, "y": 0, "z": 0},
  "battery": 99.64,
  "flightMode": "AltitudeHold",
  "throttle": 75.89,
  "armed": true,
  "onGround": false,
  "destroyed": false
}
```

## Implementation

- **File**: `systems/nats/client.go`
- **Thread-safe**: Uses `Simulator.Lock()/Unlock()`
- **Telemetry rate**: 10Hz (configurable)
- **Graceful degradation**: Simulator works without NATS

See [TODO.md](TODO.md) for roadmap.
