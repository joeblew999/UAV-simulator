# Systems

This project uses **Taskfile** and **Process Compose** to orchestrate all services.

## Architecture

```
Task → Process Compose → Task
```

- **Taskfile** (`Taskfile.yml`) - Main entry point for all commands
- **Process Compose** (`process-compose.yml`) - Orchestrates services with dependencies
- **Systems** (`systems/`) - Each subsystem has its own folder with Taskfile and README

## Quick Start

```sh
task pc:up          # Start everything (NATS, simulator, voxel, narun)
task pc:down        # Stop everything
task pc:up:bg       # Start in background
task pc:attach      # Attach to running TUI
```

## Systems Folder

Each system in `systems/` follows the same pattern:

| System | Description |
|--------|-------------|
| `nats-server/` | NATS server binary |
| `nats-cli/` | NATS CLI tool |
| `nats/` | Go NATS client (part of main module) |
| `narun/` | HTTP/gRPC gateway to NATS |
| `voxel/` | Voxel-fun 3D visualizer |
| `pc/` | Process Compose orchestration |

### Standard Tasks

Each system provides:

```sh
task <system>:start        # Start the service
task <system>:stop         # Stop the service
task <system>:deps:install # Install dependencies/binaries
task <system>:deps:clean   # Remove dependencies
task <system>:debug        # Print debug info
```

## Directory Structure

```
.bin/           # Compiled binaries (nats-server, narun, etc.)
.data/          # Runtime data directories
  nats-server/  # NATS JetStream data
  narun/        # Narun data
.src/           # Cloned source repos
  narun/        # Narun source
  voxel-fun/    # Voxel visualizer source
.env            # Environment variables (NATS_PORT, NATS_URL, etc.)
```

## Configuration

All services use `.env` for configuration:

```env
NATS_PORT=4222
NATS_URL=nats://localhost:4222
VOXEL_PORT=5173
NARUN_PORT=8080
```

Both Taskfile (`dotenv: ['.env']`) and Process Compose (`env_file: .env`) load this file.
