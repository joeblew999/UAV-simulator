# Narun Integration TODO

## Why Narun?

Narun provides an HTTP/gRPC gateway to NATS. This gives us:

1. **Web GUI Access** - Control drones from any browser via REST API
2. **VIA Integration** - Use VIA or any web framework to build a control dashboard
3. **Voxel WebGL** - Connect voxel-fun (Three.js) to the simulator via HTTP
4. **Mobile Apps** - iOS/Android can call REST endpoints
5. **External Tools** - curl, Postman, any HTTP client can interact with drones

## Architecture

```
Browser / Mobile / VIA
        ↓ HTTP
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

## TODO

### Phase 1: Basic HTTP API
- [ ] Wire up MicroService in main.go
- [ ] Test with curl
- [ ] Verify narun-gw routing works

### Phase 2: Web Dashboard
- [ ] Create simple HTML/JS dashboard
- [ ] Display drone positions in real-time
- [ ] Add control buttons (arm, takeoff, land)
- [ ] Consider VIA framework

### Phase 3: Voxel Integration
- [ ] Connect voxel-fun to narun-gw
- [ ] Send drone positions to voxel via HTTP/WebSocket
- [ ] 3D visualization in browser

### Phase 4: Mapping
- [ ] Integrate selfhostmap (https://github.com/akhenakh/selfhostmap)
- [ ] Display drone positions on map
- [ ] Geofencing alerts

## Ideas

- WebSocket endpoint for real-time telemetry streaming
- SSE (Server-Sent Events) for push updates
- Authentication/API keys
- Rate limiting
- OpenAPI/Swagger docs
