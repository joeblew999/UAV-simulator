# narun

NATS-based microservice orchestration.

https://github.com/akhenakh/narun

## Why

Narun can execute functions triggered by NATS messages, orchestrate multi-step drone missions, and run WebAssembly workloads. Exploring whether it's useful for the simulator.

## How

```sh
task narun:start        # Start narun gateway
task narun:stop         # Stop narun
task narun:deps:install # Clone and build narun binaries
task narun:deps:clean   # Remove narun source and binaries
task narun:debug        # Print debug info
```

## Configuration

- `NATS_URL` - NATS server URL
- `NARUN_PORT` - Narun gateway port (default: 8080)

Source is cloned to `.src/narun/`, binaries go to `.bin/`.

## macOS Limitation

**Note:** The `node-runner` component does not build on macOS due to a missing `runLauncher()` function in `launcher_other.go`. Only the CLI (`narun`) and gateway (`narun-gw`) are built.

The issue is that `cmd/node-runner/main.go:59` calls `runLauncher()` which is only defined in:
- `launcher_linux.go`
- `launcher_freebsd.go`

But `launcher_other.go` (used for macOS/Windows) only defines `runLandlockLauncher()`.

**To file an issue:** https://github.com/akhenakh/narun/issues

Suggested fix: Add a stub `runLauncher()` function to `launcher_other.go` that prints an error message about unsupported platform.
