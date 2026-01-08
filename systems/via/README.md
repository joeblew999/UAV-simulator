# Via Dashboard

Web dashboard for drone control using [Via](https://github.com/nicois/via) + [Datastar](https://data-star.dev/).

## Architecture

```
Via Server (:8084)          nats2sse (:8083)
    │                           │
    │ HTTP (control)            │ SSE (telemetry)
    ▼                           ▼
┌───────────────────────────────────────┐
│         Browser (Datastar)            │
│  - Arm/Disarm/Takeoff buttons         │
│  - Real-time telemetry gauges         │
│  - Drone list                         │
└───────────────────────────────────────┘
```

## Usage

```bash
# Start all services
task pc:up

# Or start individually
task nats-server:start
task nats2sse:start
task via:start
```

## Ports

| Service | Port | Purpose |
|---------|------|---------|
| Via | 8084 | Dashboard HTML + control actions |
| nats2sse | 8083 | NATS → SSE telemetry stream |
| narun-gw | 8081 | HTTP → NATS command routing |

## TODO

- [ ] Create dashboard HTML template
- [ ] Add Datastar attributes for real-time updates
- [ ] Connect to nats2sse for telemetry
- [ ] Route control actions via narun-gw
