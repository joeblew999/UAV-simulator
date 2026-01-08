# CLAUDE.md

Go-based 3D quadcopter physics simulator with OpenGL 4.1 rendering, realistic flight dynamics, and 20-drone swarm coordination.

## Quick Start

```sh
task pc:up      # Start everything (NATS + Simulator + Voxel + Narun)
task pc:down    # Stop everything
task pc:attach  # Attach to TUI
```

## CLI Flags

- `-headless` - Run without window
- `-steps=N` - Fixed update count (headless)
- `-duration=5s` - Run duration (headless)
- `-ups=120` - Updates per second

## Project Structure

```
main.go                  # Entry point, flags
internal/sim/
├── drone.go             # Physics, PID control, engines
├── simulator.go         # Main loop
├── swarm.go             # Formation control
├── camera.go            # Camera modes
├── renderer.go          # OpenGL
├── input.go             # Keyboard/mouse
├── audio.go             # Rotor sound
├── math.go              # Vec3, Mat4
└── ui.go                # HUD overlay
systems/                 # See SYSTEM.md
test/                    # Unit tests
```

## Key Types

```go
type Drone struct {
    Position, Velocity, Rotation Vec3
    PropSpeeds [4]float64
    Engines []Engine
    FlightMode FlightMode
    PitchPID, RollPID, YawPID, AltitudePID PIDController
    IsArmed, OnGround, Destroyed bool
}
```

## Controls

- **Throttle**: W/S
- **Yaw**: A/D
- **Roll**: Q/E
- **Pitch**: Up/Down arrows
- **Arm**: Hold SPACE 2s
- **Disarm**: ESC
- **Camera**: Right-drag orbit, scroll zoom
- **Drone select**: [ / ]
- **Flight modes**: 1/2/3

## Commands

```sh
task main:test      # Run tests
task main:lint      # fmt + vet
task ci             # Full CI pipeline
task debug:all      # Print all debug info
```

## Code Style

- Conventional Commits: `feat:`, `fix:`, `refactor:`
- See AGENTS.md for guidelines
