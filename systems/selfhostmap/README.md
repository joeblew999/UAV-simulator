# selfhostmap

Self-hosted map tile server for drone position visualization.

https://github.com/akhenakh/selfhostmap

## Why

Display drone positions on a real map without external tile services. Perfect for:
- Offline/air-gapped environments
- Low-latency local rendering
- Privacy (no external API calls)
- Geofencing visualization

## Tasks

```sh
task selfhostmap:start         # Start map server
task selfhostmap:stop          # Stop server
task selfhostmap:deps:install  # Clone and build
task selfhostmap:deps:clean    # Remove source and binary
task selfhostmap:debug:self    # Print debug info
```

## Config

- `SELFHOSTMAP_PORT` - HTTP port (default: 8082)

Source cloned to `.src/selfhostmap/`, binary to `.bin/`.
