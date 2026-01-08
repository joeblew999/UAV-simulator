# pc

Process Compose orchestration for the UAV simulator.

## Why

Runs all services (NATS server, simulator, voxel visualizer) together with proper dependency ordering and health checks.

## How

```sh
task pc:up            # Start all services with TUI
task pc:up:bg         # Start all services in background
task pc:up:headless   # Start without TUI (for CI)
task pc:down          # Stop all services
task pc:attach        # Attach to running TUI
task pc:deps:install  # Install process-compose binary
task pc:deps:clean    # Remove binary
task pc:debug         # Print debug info
```

## Configuration

Services are defined in `process-compose.yml` at the project root.
Environment variables are loaded from `.env`.
