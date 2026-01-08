# CLAUDE.md

## Project Overview

Go-based 3D quadcopter physics simulator with OpenGL 4.1 rendering, realistic flight dynamics, and 20-drone swarm coordination. Models a DJI Mini 2-equivalent drone (249g) with cascaded PID control, battery simulation, and spatial audio.

## Build & Run

```bash
go build -o drone-simulator .    # Build
go run .                         # Run with GUI
go run . -headless -steps=1000   # Headless benchmark
```

### CLI Flags
- `-headless` - Run without window (benchmarking)
- `-steps=N` - Fixed update count (headless)
- `-duration=5s` - Run duration (headless)
- `-ups=120` - Updates per second (default: 120)
- `-decoupled=true` - Separate physics/render loops (default)

## Testing

```bash
go test ./...           # Run all tests
go test -cover ./...    # With coverage
go vet ./...            # Static analysis
```

Tests are in `test/` using `sim_test` package - headless compatible, no OpenGL deps.

## Project Structure

```
main.go                  # Entry point, window setup, flags
internal/sim/
├── drone.go             # Physics model, PID control, engines (851 LOC)
├── simulator.go         # Main loop, subsystem coordination (910 LOC)
├── swarm.go             # Leader-follower formation control
├── camera.go            # Follow/TopDown/FPV camera modes
├── renderer.go          # OpenGL shaders, mesh rendering
├── input.go             # Keyboard/mouse handling
├── audio.go             # Rotor sound synthesis (thrust-based RPM)
├── math.go              # Vec3, Mat4 utilities
├── ui.go                # 2D HUD overlay
└── avaudio/             # Platform audio (darwin=AVFoundation)
test/                    # Unit tests (drone, math, swarm)
cmd/headless/            # Headless benchmark binary
```

## Key Types

```go
type Drone struct {
    Position, Velocity, Rotation Vec3
    PropSpeeds [4]float64          // Motor RPMs
    Engines []Engine               // 4 rotors with thrust/torque
    FlightMode FlightMode          // Manual, AltitudeHold, Hover
    PitchPID, RollPID, YawPID, AltitudePID PIDController
    IsArmed, OnGround, Destroyed bool
}

type Swarm struct {
    drones []*Drone
    latency float64                // Simulated 100ms comm delay
    leaderIdx int
}
```

## Physics Model

- **Gravity**: 9.81 m/s²
- **Air density**: 1.225 kg/m³ (varies with altitude)
- **Max thrust**: 2.5× weight across 4 engines
- **Simulation**: 120 Hz fixed timestep, decoupled from render

Update sequence: ground contact → air density → forces (gravity, thrust, drag) → altitude PID → integration → collision → angular motion

## Controls (Runtime)

- **Throttle**: W/S
- **Yaw**: A/D
- **Roll**: Q/E
- **Pitch**: Up/Down arrows
- **Arm**: Hold SPACE 2s
- **Disarm**: ESC
- **Camera**: Right-drag orbit, scroll zoom
- **Drone select**: [ / ]
- **Flight modes**: 1/2/3

## Code Style

- Standard `go fmt`
- Conventional Commits: `feat(scope):`, `refactor:`, `chore:`
- Tests in separate `test/` directory with `sim_test` package
- Physics/rendering cleanly separated

## Active Development

Recent focus: audio system (thrust-based RPM), vertical thrust tracking, physics refinements. See commits for patterns.
