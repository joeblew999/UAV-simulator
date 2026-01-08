# TODO: NATS Integration

## Overview

Integrate NATS messaging so or any external agent can control the drone simulator via NATS CLI. This enables AI-driven flight control, automated testing, and remote operation.

## Dev

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

### Other Useful Commands

```sh
# Debug - print all vars and env
task debug

# Subscribe to all drone telemetry (separate terminal)
task nats:sub

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



## Use Case: Claude Controlling Drones via NATS CLI

```bash
# Start everything (NATS server + Simulator with GUI)
task pc:up

# In a separate terminal, Claude can control drones:
nats pub drone.0.arm ''                              # Arm drone 0
nats pub drone.0.takeoff '{"altitude": 10}'          # Take off to 10m
nats pub drone.0.goto '{"x": 5, "y": 10, "z": 5}'    # Fly to position
nats pub drone.0.land ''                             # Land
nats pub swarm.formation '{"type": "circle", "radius": 10}'

# Monitor state:
nats sub "drone.*.telemetry"                         # Watch all drones
nats sub "simulator.status"                          # Simulator state
```

## Tasks

### Phase 1: Basic NATS Setup
- [ ] Add NATS client dependency (`github.com/nats-io/nats.go`)
- [ ] Create `internal/sim/nats.go` for NATS connection management
- [ ] Add CLI flags: `-nats-url`, `-nats-enabled`
- [ ] Implement connection lifecycle (connect, reconnect, disconnect)
- [ ] Graceful degradation if NATS unavailable

### Phase 2: Telemetry Publishing
- [ ] Publish drone state on `drone.<id>.telemetry` (JSON)
  - Position (x, y, z)
  - Velocity
  - Rotation (pitch, roll, yaw)
  - Battery percentage
  - Flight mode
  - Armed/OnGround/Destroyed status
- [ ] Publish at configurable rate (default: 10 Hz, flag: `-nats-telemetry-hz`)
- [ ] Publish simulator status on `simulator.status`

### Phase 3: Command Subscription (Core for Claude Control)
- [ ] Subscribe to `drone.<id>.arm` - Arm specific drone
- [ ] Subscribe to `drone.<id>.disarm` - Disarm specific drone
- [ ] Subscribe to `drone.<id>.takeoff` - `{"altitude": 10}`
- [ ] Subscribe to `drone.<id>.land` - Land at current position
- [ ] Subscribe to `drone.<id>.goto` - `{"x": 0, "y": 10, "z": 0}`
- [ ] Subscribe to `drone.<id>.input` - Direct control `{"throttle": 0.5, "yaw": 0, "pitch": 0, "roll": 0}`
- [ ] Subscribe to `drone.<id>.mode` - `{"mode": "Hover"}` (Manual/AltitudeHold/Hover)
- [ ] Subscribe to `drone.<id>.stop` - Emergency stop

### Phase 4: Swarm Commands
- [ ] Subscribe to `swarm.arm` - Arm all drones
- [ ] Subscribe to `swarm.disarm` - Disarm all
- [ ] Subscribe to `swarm.takeoff` - `{"altitude": 10}`
- [ ] Subscribe to `swarm.land` - Land all
- [ ] Subscribe to `swarm.formation` - `{"type": "line|circle|grid", ...params}`
- [ ] Subscribe to `swarm.goto` - Move entire swarm

### Phase 5: Request-Reply (for Claude to query state)
- [ ] `drone.<id>.status` - Request current drone state, get reply
- [ ] `swarm.status` - Request all drone states
- [ ] `simulator.info` - Get simulator config/state

## Message Schemas

### Telemetry (drone.<id>.telemetry) - Published by Simulator
```json
{
  "id": 0,
  "timestamp": 1704672000000,
  "position": {"x": 0.0, "y": 10.0, "z": 0.0},
  "velocity": {"x": 0.0, "y": 0.0, "z": 0.0},
  "rotation": {"pitch": 0.0, "roll": 0.0, "yaw": 0.0},
  "battery": 0.95,
  "flightMode": "Hover",
  "armed": true,
  "onGround": false,
  "destroyed": false
}
```

### Goto Command (drone.<id>.goto)
```json
{"x": 5.0, "y": 10.0, "z": 5.0}
```

### Takeoff Command (drone.<id>.takeoff)
```json
{"altitude": 10.0}
```

### Formation Command (swarm.formation)
```json
{"type": "circle", "radius": 10.0, "altitude": 15.0}
```

### Direct Input (drone.<id>.input)
```json
{"throttle": 0.5, "yaw": 0.1, "pitch": 0.0, "roll": 0.0}
```

## Example Claude Workflow

1. **Start simulator**: User runs `task pc:up`
2. **Claude subscribes to telemetry**: `nats sub drone.0.telemetry`
3. **Claude arms drone**: `nats pub drone.0.arm ''`
4. **Claude commands takeoff**: `nats pub drone.0.takeoff '{"altitude":5}'`
5. **Claude monitors altitude** via telemetry stream
6. **Claude navigates**: `nats pub drone.0.goto '{"x":10,"y":5,"z":0}'`
7. **Claude lands**: `nats pub drone.0.land ''`

## Dependencies

```bash
go get github.com/nats-io/nats.go
```

## Implementation Notes

- All commands are fire-and-forget (pub) except status queries (request-reply)
- Thread-safe command queue in simulator to process NATS messages
- JSON for simplicity (Claude can easily construct messages)
- Empty payload `''` for simple commands (arm, disarm, land, stop)
- Waypoint/goto commands should use the existing PID controllers
