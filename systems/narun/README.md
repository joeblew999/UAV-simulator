# narun

HTTP/gRPC gateway to NATS. https://github.com/akhenakh/narun

## Tasks

```sh
task narun:start         # Start gateway
task narun:stop          # Stop gateway
task narun:deps:install  # Clone and build
task narun:deps:clean    # Remove source and binaries
task narun:debug:self    # Print debug info
```

## Config

- `NATS_URL` - NATS server URL
- `NARUN_PORT` - Gateway port (default: 8081)
- `NARUN_METRICS_PORT` - Metrics port (default: 9091)

Source cloned to `.src/narun/`, binaries to `.bin/`.

## macOS Note

The `node-runner` component doesn't build on macOS (missing `runLauncher()` stub). Only `narun` CLI and `narun-gw` are built.
