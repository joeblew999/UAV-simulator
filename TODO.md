# TODO

## NATS

### Completed
- [x] Basic NATS client
- [x] Telemetry publishing at 10Hz
- [x] Command subscriptions (arm, disarm, takeoff, land, goto, input, mode, stop)
- [x] Thread-safe simulator access
- [x] Process Compose orchestration

### In Progress
- [ ] GOTO lateral control (X/Z) - currently altitude (Y) only
- [ ] Fix crash when drone destroyed
- [ ] Reconnection handling

### Future
- [ ] JetStream persistence (replay, save/load)
- [ ] Swarm commands (arm all, formation, etc.)
- [ ] Request-reply (status queries)

## Swarm

### Current Implementation (Centralized)
Leader-follower with single point of failure.

### Future: Distributed Swarm
Boids algorithm with local rules:
- Separation (avoid collision)
- Alignment (match neighbors' velocity)
- Cohesion (move toward group center)

Each drone only knows about nearby neighbors. No single leader.

## Physics

- [ ] Proper rotor dynamics (blade element theory)
- [ ] Flight envelope protection (vortex ring state)
- [ ] Cascaded control loops (position → velocity → attitude → rates)
- [ ] Battery modeling with voltage curves
- [ ] Sensor modeling (IMU noise, GPS errors)
- [ ] Wind field with turbulence
