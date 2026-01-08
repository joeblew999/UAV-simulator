# nats-server

NATS server binary.

## Tasks

```sh
task nats-server:start         # Start server
task nats-server:stop          # Stop server
task nats-server:status        # Check status
task nats-server:deps:install  # Install binary
task nats-server:deps:clean    # Remove binary and data
task nats-server:debug:self    # Print debug info
```

## Config

- `NATS_PORT` - Server port (default: 4222)

Data stored in `.data/nats-server/` with JetStream enabled.
