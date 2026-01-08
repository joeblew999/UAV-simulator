# NATS Integration

https://github.com/pawel-dubiel/UAV-simulator is the origin

Adding nats to:
https://github.com/joeblew999/UAV-simulator

## Why

Integrate NATS messaging so any external agent (including Claude) can control the drone simulator via NATS CLI. This enables AI-driven flight control, automated testing, and remote operation.

## How

You only need Taskfile to run everything.

https://taskfile.dev/installation

We use Taskfile and Process Compose to run everything.
This keeps development localised fully with idempotency and single source of truth.

**Architecture:** Task → Process Compose → Task

We use `.env` and expose that to Task and Process Compose for config like NATS_PORT, NATS_URL, etc.

### Quick Start

```sh
# Start everything (NATS + Simulator with GUI)
task pc:up

# Stop everything
task pc:down

# Start in background
task pc:up:bg

# Attach to running TUI
task pc:attach

# Start without TUI (for CI/testing)
task pc:up:headless
```

### Control Drones via NATS

```sh
# Arm drone 0
nats pub drone.0.arm ''

# Take off to 5 meters
nats pub drone.0.takeoff '{"altitude": 5}'

# Fly to position (x, y, z)
nats pub drone.0.goto '{"x": 10, "y": 5, "z": 10}'

# Set throttle/yaw/pitch/roll directly
nats pub drone.0.input '{"throttle": 0.6, "yaw": 0, "pitch": 0, "roll": 0}'

# Change flight mode (Manual, AltitudeHold, Hover)
nats pub drone.0.mode '{"mode": "Hover"}'

# Land
nats pub drone.0.land ''

# Disarm
nats pub drone.0.disarm ''

# Emergency stop
nats pub drone.0.stop ''
```

### Monitor Telemetry

```sh
# Watch single drone
nats sub drone.0.telemetry

# Watch all drones
nats sub "drone.*.telemetry"

# Using task helper
task nats:sub
```

### Other Useful Commands

```sh
# Debug - print all vars and env
task debug

# Check NATS server status
task nats:status

# Run simulator standalone (no NATS)
task sim:run:standalone

# Run tests
task go:test
```

### Dependencies

Dependencies are installed automatically and idempotently to `.bin/`:

```sh
# Install all deps (process-compose, nats-server, nats-cli)
task deps:install

# Clean all deps
task deps:clean
```

## Telemetry Format

Telemetry is published at 10Hz on `drone.<id>.telemetry`:

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

## Command Reference

| Subject | Payload | Description |
|---------|---------|-------------|
| `drone.<id>.arm` | `''` | Arm the drone |
| `drone.<id>.disarm` | `''` | Disarm the drone |
| `drone.<id>.takeoff` | `{"altitude": 5}` | Take off to altitude (m) |
| `drone.<id>.land` | `''` | Land at current position |
| `drone.<id>.goto` | `{"x": 0, "y": 10, "z": 0}` | Fly to position |
| `drone.<id>.input` | `{"throttle": 0.5, "yaw": 0, "pitch": 0, "roll": 0}` | Direct control |
| `drone.<id>.mode` | `{"mode": "Hover"}` | Set flight mode |
| `drone.<id>.stop` | `''` | Emergency stop (disarm + zero throttle) |

## Example: Claude Controlling Drones

```bash
# 1. Start simulator
task pc:up

# 2. In separate terminal, Claude controls drones:
nats pub drone.0.arm ''
nats pub drone.0.takeoff '{"altitude": 5}'

# 3. Wait for altitude, then navigate
nats pub drone.0.goto '{"x": 10, "y": 5, "z": 10}'

# 4. Monitor telemetry
nats sub drone.0.telemetry --count=1

# 5. Land when done
nats pub drone.0.land ''
```

## Implementation

- **Client**: `internal/nats/client.go`
- **Wired in**: `main.go` (via `-nats-url` flag)
- **Thread-safe**: Uses `Simulator.Lock()/Unlock()` for command execution
- **Telemetry rate**: 10Hz (configurable in client)
- **Graceful degradation**: Simulator works without NATS if not configured

## Future: Swarm Commands

- [ ] `swarm.arm` - Arm all drones
- [ ] `swarm.disarm` - Disarm all
- [ ] `swarm.takeoff` - `{"altitude": 10}`
- [ ] `swarm.land` - Land all
- [ ] `swarm.formation` - `{"type": "circle", "radius": 10}`
- [ ] `swarm.goto` - Move entire swarm
