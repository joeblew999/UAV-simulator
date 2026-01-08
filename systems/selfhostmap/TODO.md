# Selfhostmap Integration TODO

## Why Selfhostmap?

Self-hosted map tiles for drone visualization without external dependencies.

https://github.com/akhenakh/selfhostmap

## Use Cases

1. **Drone Position Overlay** - Real-time drone markers on map
2. **Geofencing** - Visual boundaries for no-fly zones
3. **Mission Planning** - Draw waypoints on map
4. **Offline Operation** - Works without internet
5. **Privacy** - No data sent to external tile servers

## Architecture

```
Browser (Leaflet/MapLibre)
        ↓ tile requests
    selfhostmap (:8082)
        ↓ vector tiles
    PMTiles / MBTiles

Drone positions:
    Simulator → NATS → narun-gw → SSE → Browser → Map overlay
```

## TODO

### Phase 1: Basic Setup
- [ ] Add to process-compose.yml
- [ ] Download map tiles for test area
- [ ] Verify tile serving works

### Phase 2: Drone Overlay
- [ ] Create Leaflet/MapLibre page
- [ ] Subscribe to drone telemetry via SSE
- [ ] Render drone markers on map
- [ ] Update positions in real-time

### Phase 3: Geofencing
- [ ] Define geofence polygons
- [ ] Publish geofence alerts to NATS
- [ ] Visual warning when drone approaches boundary

### Phase 4: Mission Planning
- [ ] Draw waypoints on map
- [ ] Export waypoints as NATS commands
- [ ] Execute mission via `drone.*.goto`

## Map Data

Need to download PMTiles for your region:
- https://protomaps.com/downloads
- OpenStreetMap extracts

## Ideas

- 3D terrain with drone altitude
- Flight path history trails
- Multiple map layers (satellite, terrain, streets)
- Integration with voxel-fun for 3D + 2D views
