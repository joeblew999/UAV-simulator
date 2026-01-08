# nats-cli

NATS CLI tool.

## Tasks

```sh
task nats-cli:sub              # Subscribe to drone telemetry
task nats-cli:pub              # Publish a message
task nats-cli:rtt              # Check server RTT
task nats-cli:deps:install     # Install CLI binary
task nats-cli:deps:clean       # Remove binary
task nats-cli:debug:self       # Print debug info
```

## Examples

```sh
# Subscribe to all drone telemetry
task nats-cli:sub

# Publish commands
task nats-cli:pub -- drone.0.arm ''
task nats-cli:pub -- drone.0.takeoff '{"altitude": 5}'
```
