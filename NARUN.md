# narun

https://github.com/akhenakh/narun

## Why?

Narun is a NATS-based microservice orchestration platform. It could potentially be used to:
- Execute functions triggered by NATS messages
- Orchestrate multi-step drone missions
- Run WebAssembly workloads triggered by drone events

Currently exploring whether it's useful for the simulator or overkill for simple pub/sub.

## How

Taskfile include with Process Compose running it.

### Quick Start

```sh
# Clone and build narun
task narun:deps:install

# Run narun (requires NATS server)
task narun:run

# Or start everything with process-compose
task pc:up
```

### Available Commands

```sh
task narun:clone         # Clone repository to .src/narun
task narun:pull          # Pull latest changes
task narun:deps:install  # Build narun binary
task narun:deps:clean    # Remove narun source directory
task narun:run           # Run narun (connects to NATS)
```

## Architecture

Narun connects to NATS and listens for function invocation requests. It can:
- Run Go plugins
- Execute WebAssembly modules
- Schedule recurring tasks

See the [narun repository](https://github.com/akhenakh/narun) for full documentation.
