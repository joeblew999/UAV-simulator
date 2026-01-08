# nats2sse

NATS to Server-Sent Events bridge for real-time web updates.

https://github.com/akhenakh/nats2sse

## Why

Bridges NATS pub/sub to browser SSE streams. Perfect for:
- Datastar reactive updates
- Real-time telemetry in browser
- No WebSocket complexity
- Works with any SSE-capable frontend

## Tasks

```sh
task nats2sse:start         # Start SSE bridge
task nats2sse:stop          # Stop bridge
task nats2sse:deps:install  # Clone and build
task nats2sse:deps:clean    # Remove source and binary
task nats2sse:debug:self    # Print debug info
```

## Config

- `NATS2SSE_PORT` - HTTP port (default: 8083)
- `NATS_URL` - NATS server URL

Source cloned to `.src/nats2sse/`, binary to `.bin/`.

## Usage

Subscribe to NATS subject via SSE:
```
GET http://localhost:8083/sse?subject=drone.0.telemetry
```

In Datastar:
```html
<div data-on-sse="/sse?subject=drone.0.telemetry"
     data-on-message="updateDrone(event.data)">
</div>
```
