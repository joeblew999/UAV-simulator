# nats

Go NATS client for the simulator. Part of the main module.

```go
import natsclient "drone-simulator/systems/nats"
```

## Commands

| Subject | Payload | Description |
|---------|---------|-------------|
| `drone.<id>.arm` | `''` | Arm drone |
| `drone.<id>.disarm` | `''` | Disarm drone |
| `drone.<id>.takeoff` | `{"altitude": 5}` | Take off |
| `drone.<id>.land` | `''` | Land |
| `drone.<id>.goto` | `{"x": 0, "y": 10, "z": 0}` | Fly to position |
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
- **Telemetry rate**: 10Hz
- **Graceful degradation**: Simulator works without NATS
