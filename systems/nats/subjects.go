package nats

import "fmt"

// NATS subject patterns for the simulator.
// These are shared across all components that interact with NATS.
const (
	// Telemetry subjects - published by simulator, consumed by nats2sse for SSE
	SubjectTelemetryPattern = "telemetry.>"     // For nats2sse subscription
	SubjectTelemetryFmt     = "telemetry.%d"    // telemetry.<droneID>

	// Micro service subjects - request/reply via narun-gw
	SubjectDroneList    = "drone.list"
	SubjectDroneStatus  = "drone.status"
	SubjectDroneArm     = "drone.arm"
	SubjectDroneDisarm  = "drone.disarm"
	SubjectDroneTakeoff = "drone.takeoff"
	SubjectDroneLand    = "drone.land"
	SubjectDroneGoto    = "drone.goto"
	SubjectDroneMode    = "drone.mode"
	SubjectDroneStop    = "drone.stop"

	// Legacy command subjects (direct pub/sub) - deprecated, use micro service
	SubjectCommandFmt = "drone.%d.%s" // drone.<droneID>.<command>
)

// TelemetrySubject returns the telemetry subject for a specific drone.
func TelemetrySubject(droneID int) string {
	return fmt.Sprintf(SubjectTelemetryFmt, droneID)
}

// CommandSubject returns a command subject for a specific drone (legacy).
func CommandSubject(droneID int, command string) string {
	return fmt.Sprintf(SubjectCommandFmt, droneID, command)
}
