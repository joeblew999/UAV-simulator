# Narun Integration TODO

## Why Narun?

Narun provides an HTTP/gRPC gateway to NATS. This gives us:

1. **Web GUI Access** - Control drones from any browser via REST API
2. **Datastar + HTMX** - Reactive hypermedia for real-time drone control
3. **Voxel WebGL** - Connect voxel-fun (Three.js) to the simulator via HTTP
4. **Mobile Apps** - iOS/Android can call REST endpoints
5. **External Tools** - curl, Postman, any HTTP client can interact with drones

## Architecture

```
Browser (Datastar/HTMX)
        ↓ HTTP + SSE
    narun-gw (:8081)
        ↓ NATS Micro
    Drone Simulator (Go)
        ↓ NATS pub/sub
    voxel-fun (WebGL viz)
    selfhostmap (future)
```

## HTTP API (via narun-gw)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/drone/` | List all drones with telemetry |
| GET | `/drone/{id}` | Get single drone status |
| POST | `/drone/{id}/arm` | Arm drone |
| POST | `/drone/{id}/disarm` | Disarm drone |
| POST | `/drone/{id}/takeoff` | Take off `{"altitude": 10}` |
| POST | `/drone/{id}/land` | Land drone |
| POST | `/drone/{id}/goto` | Go to position `{"x":0, "y":10, "z":0}` |
| POST | `/drone/{id}/mode` | Set flight mode `{"mode": "Hover"}` |
| POST | `/drone/{id}/stop` | Emergency stop |

## Completed

- [x] NATS Micro service (`systems/nats/micro.go`)
- [x] narun-gw config (`systems/narun/config.yaml`)
- [x] Wire up MicroService in main.go

## TODO

### Phase 1: Basic HTTP API
- [ ] Test with curl
- [ ] Verify narun-gw routing works
- [ ] Add SSE endpoint for telemetry streaming

### Phase 2: Via + Datastar Dashboard
- [ ] Add Via web server system (`systems/via/`)
- [ ] Create Go templates with Datastar attributes
- [ ] Real-time telemetry updates via SSE
- [ ] Control buttons with `data-on-click="@post('/drone/0/arm')"`
- [ ] Drone list with reactive updates
- [ ] Battery/altitude gauges
- [ ] Flight mode selector

**Why Via + Datastar?**
- Via: Go web framework with first-class Datastar support
- Datastar: Hypermedia framework with SSE for real-time updates
- No build step, just HTML + Go templates
- Perfect for drone telemetry dashboards
- https://github.com/go-via/via
- https://data-star.dev/

**Architecture:**
```
┌─────────────────┐     ┌─────────────────┐
│  Via Server     │     │  nats2sse       │
│  (:8084)        │     │  (:8083)        │
│  - Templates    │     │  - NATS → SSE   │
│  - Datastar     │     │                 │
└────────┬────────┘     └────────┬────────┘
         │ HTTP POST             │ SSE (telemetry)
         ▼                       ▼
    ┌─────────────────────────────────────┐
    │           Browser (Datastar)        │
    │  - Control buttons → Via → narun-gw │
    │  - Telemetry ← nats2sse SSE stream  │
    └─────────────────────────────────────┘
         │                       ▲
         ▼ NATS Micro            │ NATS pub
    ┌─────────────────┐     ┌────┴────────┐
    │  narun-gw       │     │  Simulator  │
    │  (:8081)        │────▶│  (Go)       │
    └─────────────────┘     └─────────────┘
```

- **Via** serves the dashboard HTML + handles control actions
- **nats2sse** streams real-time telemetry to browser via SSE
- **narun-gw** routes HTTP commands to NATS Micro service
- https://github.com/go-via/via
- https://github.com/akhenakh/nats2sse

### Phase 3: Voxel Integration
- [ ] Connect voxel-fun to narun-gw
- [ ] Send drone positions to voxel via HTTP/SSE
- [ ] 3D visualization in browser
- [ ] Sync camera with selected drone

### Phase 4: Mapping
- [ ] Integrate selfhostmap (https://github.com/akhenakh/selfhostmap)
- [ ] Display drone positions on map
- [ ] Geofencing alerts
- [ ] Mission waypoints

## Ideas

- SSE endpoint: `GET /drone/stream` - push telemetry every 100ms
- Datastar fragments for partial updates
- Go templates served by narun-gw or separate static server
- Authentication via NATS credentials
- OpenAPI/Swagger docs
- Record/replay flight sessions via NATS JetStream
