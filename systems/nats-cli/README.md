# nats-cli

NATS CLI tool for interacting with the NATS server.

## Why

Provides command-line access to NATS for debugging, monitoring telemetry, and sending commands to drones.

## How

```sh
task nats-cli:sub             # Subscribe to all drone telemetry
task nats-cli:pub             # Publish a message
task nats-cli:rtt             # Check server round-trip time
task nats-cli:deps:install    # Install NATS CLI binary
task nats-cli:deps:clean      # Remove CLI binary
task nats-cli:debug           # Print debug info
```

## Examples

```sh
# Subscribe to all drone telemetry
task nats-cli:sub

# Publish a command
task nats-cli:pub -- drone.0.arm ''
task nats-cli:pub -- drone.0.takeoff '{"altitude": 5}'
```
