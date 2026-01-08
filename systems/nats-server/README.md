# nats-server

NATS server binary management for the UAV simulator.

## Why

The simulator uses NATS for pub/sub messaging between components. This system manages the NATS server binary installation and lifecycle.

## How

```sh
task nats-server:start        # Start NATS server
task nats-server:stop         # Stop NATS server
task nats-server:status       # Check server status
task nats-server:deps:install # Install NATS server binary
task nats-server:deps:clean   # Remove binary and data
task nats-server:debug        # Print debug info
```

## Configuration

Uses environment variables from `.env`:
- `NATS_PORT` - Server port (default: 4222)
- `NATS_URL` - Full URL (default: nats://localhost:4222)

Data is stored in `.data/nats/` with JetStream enabled.
