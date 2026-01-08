# NATS TODO

## Completed

- [x] Basic NATS client (`internal/nats/client.go`)
- [x] Telemetry publishing at 10Hz on `drone.<id>.telemetry`
- [x] Command subscriptions: arm, disarm, takeoff, land, goto, input, mode, stop
- [x] Thread-safe simulator access via Lock/Unlock methods
- [x] Wire up in main.go with `-nats-url` flag
- [x] Process Compose orchestration (NATS + Simulator)
- [x] Taskfile integration

## TODO

### GOTO Lateral Control

The GOTO command currently only supports altitude (Y) control. Lateral positioning (X, Z) is not implemented because:

1. The simulator's internal physics (stability damping, motor torques) conflicts with external torque control
2. Attempting to apply AddTorque from the NATS client results in oscillation and movement in the wrong direction
3. The swarm.go has lateral control that works because it runs inside the simulator's physics loop

**Possible solutions:**
- Modify simulator to expose a position controller API
- Add a "target position" field that the simulator's internal control loop handles
- Use swarm mode for multi-drone positioning

### Stability

- [ ] Fix simulator crash when drone crashes/destroyed
- [ ] Graceful handling of all drone states
- [ ] Reconnection handling if NATS disconnects

### JetStream Persistence

- [ ] Enable JetStream on NATS server
- [ ] Create `DRONE_TELEMETRY` stream for all telemetry
- [ ] Create `DRONE_COMMANDS` stream for all commands
- [ ] Add replay capability (replay a flight session)
- [ ] Add scene save/load (save current swarm state, reload later)

### Swarm Commands

- [ ] `swarm.arm` - Arm all drones
- [ ] `swarm.disarm` - Disarm all drones
- [ ] `swarm.takeoff` - `{"altitude": 10}`
- [ ] `swarm.land` - Land all drones
- [ ] `swarm.formation` - `{"type": "circle", "radius": 10}`
- [ ] `swarm.goto` - Move entire swarm to position

### Request-Reply

- [ ] `drone.<id>.status` - Request current drone state, get reply
- [ ] `swarm.status` - Request all drone states
- [ ] `simulator.info` - Get simulator config/state

### Tooling

- [ ] Add `gojq` or native Go JSON parsing for testing
- [ ] Create test script that doesn't rely on `jq`
- [ ] Add NATS CLI wrapper tasks for common operations

### Documentation

- [ ] Add JetStream usage examples
- [ ] Document replay workflow
- [ ] Add troubleshooting section

## Ideas

- [ ] WebSocket bridge for browser-based control
- [ ] Mission scripting (sequence of commands)
- [ ] Geofencing alerts via NATS
- [ ] Battery low warnings published to NATS
- [ ] Collision detection alerts
