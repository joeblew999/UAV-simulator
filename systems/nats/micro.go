package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	sim "drone-simulator/internal/sim"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

// MicroService exposes the simulator as a NATS Micro service for narun-gw.
// This allows HTTP requests to be routed to the simulator via narun gateway.
type MicroService struct {
	nc        *nats.Conn
	simulator *sim.Simulator
	service   micro.Service
}

// Response types for HTTP API
type DroneResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type DroneListResponse struct {
	Count  int           `json:"count"`
	Drones []TelemetryMsg `json:"drones"`
}

type DroneStatusResponse struct {
	Success bool         `json:"success"`
	Drone   TelemetryMsg `json:"drone"`
}

// NewMicroService creates a NATS Micro service for HTTP API access.
func NewMicroService(nc *nats.Conn, simulator *sim.Simulator) (*MicroService, error) {
	ms := &MicroService{
		nc:        nc,
		simulator: simulator,
	}

	// Create the micro service
	srv, err := micro.AddService(nc, micro.Config{
		Name:        "drone",
		Version:     "1.0.0",
		Description: "UAV Simulator HTTP API",
	})
	if err != nil {
		return nil, fmt.Errorf("create micro service: %w", err)
	}
	ms.service = srv

	// Add endpoint group for drone operations
	droneGroup := srv.AddGroup("drone")

	// GET /drone/ - List all drones
	droneGroup.AddEndpoint("list", micro.HandlerFunc(ms.handleList),
		micro.WithEndpointMetadata(map[string]string{
			"method": "GET",
			"path":   "/drone/",
		}))

	// GET /drone/{id} - Get drone status
	droneGroup.AddEndpoint("status", micro.HandlerFunc(ms.handleStatus),
		micro.WithEndpointMetadata(map[string]string{
			"method": "GET",
			"path":   "/drone/{id}",
		}))

	// POST /drone/{id}/arm
	droneGroup.AddEndpoint("arm", micro.HandlerFunc(ms.handleArm),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/arm",
		}))

	// POST /drone/{id}/disarm
	droneGroup.AddEndpoint("disarm", micro.HandlerFunc(ms.handleDisarm),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/disarm",
		}))

	// POST /drone/{id}/takeoff
	droneGroup.AddEndpoint("takeoff", micro.HandlerFunc(ms.handleTakeoff),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/takeoff",
		}))

	// POST /drone/{id}/land
	droneGroup.AddEndpoint("land", micro.HandlerFunc(ms.handleLand),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/land",
		}))

	// POST /drone/{id}/goto
	droneGroup.AddEndpoint("goto", micro.HandlerFunc(ms.handleGoto),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/goto",
		}))

	// POST /drone/{id}/mode
	droneGroup.AddEndpoint("mode", micro.HandlerFunc(ms.handleMode),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/mode",
		}))

	// POST /drone/{id}/stop
	droneGroup.AddEndpoint("stop", micro.HandlerFunc(ms.handleStop),
		micro.WithEndpointMetadata(map[string]string{
			"method": "POST",
			"path":   "/drone/{id}/stop",
		}))

	log.Println("NATS Micro service 'drone' started")
	return ms, nil
}

// Stop shuts down the micro service.
func (ms *MicroService) Stop() error {
	if ms.service != nil {
		return ms.service.Stop()
	}
	return nil
}

// parseRequest extracts drone ID and body from the micro request.
// narun-gw sends the original HTTP path in X-Original-Path header.
func (ms *MicroService) parseRequest(req micro.Request) (droneID int, body []byte, err error) {
	headers := req.Headers()
	path := headers.Get("X-Original-Path")
	body = req.Data()

	// Parse drone ID from path: /drone/0/arm -> 0
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		droneID, err = strconv.Atoi(parts[1])
		if err != nil {
			return -1, body, fmt.Errorf("invalid drone ID: %s", parts[1])
		}
	} else {
		droneID = -1 // No drone ID in path (list operation)
	}

	return droneID, body, nil
}

func (ms *MicroService) getDrone(id int) *sim.Drone {
	drones := ms.simulator.Drones()
	if id < 0 || id >= len(drones) {
		return nil
	}
	return drones[id]
}

func (ms *MicroService) respondJSON(req micro.Request, status int, data interface{}) {
	body, _ := json.Marshal(data)
	headers := micro.Headers{
		"Content-Type":  []string{"application/json"},
		"X-Status-Code": []string{strconv.Itoa(status)},
	}
	req.Respond(body, micro.WithHeaders(headers))
}

func (ms *MicroService) respondError(req micro.Request, status int, message string) {
	ms.respondJSON(req, status, DroneResponse{
		Success: false,
		Error:   message,
	})
}

func (ms *MicroService) respondSuccess(req micro.Request, message string) {
	ms.respondJSON(req, http.StatusOK, DroneResponse{
		Success: true,
		Message: message,
	})
}

// Handlers

func (ms *MicroService) handleList(req micro.Request) {
	ms.simulator.RLock()
	drones := ms.simulator.Drones()

	resp := DroneListResponse{
		Count:  len(drones),
		Drones: make([]TelemetryMsg, len(drones)),
	}

	for i, d := range drones {
		resp.Drones[i] = TelemetryMsg{
			ID:         i,
			Timestamp:  time.Now().UnixMilli(),
			Position:   Vec3Msg{X: d.Position.X, Y: d.Position.Y, Z: d.Position.Z},
			Velocity:   Vec3Msg{X: d.Velocity.X, Y: d.Velocity.Y, Z: d.Velocity.Z},
			Rotation:   Vec3Msg{X: d.Rotation.X, Y: d.Rotation.Y, Z: d.Rotation.Z},
			Battery:    d.BatteryPercent,
			FlightMode: flightModeString(d.FlightMode),
			Throttle:   d.ThrottlePercent,
			Armed:      d.IsArmed,
			OnGround:   d.OnGround,
			Destroyed:  d.Destroyed,
		}
	}
	ms.simulator.RUnlock()

	ms.respondJSON(req, http.StatusOK, resp)
}

func (ms *MicroService) handleStatus(req micro.Request) {
	id, _, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	ms.simulator.RLock()
	resp := DroneStatusResponse{
		Success: true,
		Drone: TelemetryMsg{
			ID:         id,
			Timestamp:  time.Now().UnixMilli(),
			Position:   Vec3Msg{X: drone.Position.X, Y: drone.Position.Y, Z: drone.Position.Z},
			Velocity:   Vec3Msg{X: drone.Velocity.X, Y: drone.Velocity.Y, Z: drone.Velocity.Z},
			Rotation:   Vec3Msg{X: drone.Rotation.X, Y: drone.Rotation.Y, Z: drone.Rotation.Z},
			Battery:    drone.BatteryPercent,
			FlightMode: flightModeString(drone.FlightMode),
			Throttle:   drone.ThrottlePercent,
			Armed:      drone.IsArmed,
			OnGround:   drone.OnGround,
			Destroyed:  drone.Destroyed,
		},
	}
	ms.simulator.RUnlock()

	ms.respondJSON(req, http.StatusOK, resp)
}

func (ms *MicroService) handleArm(req micro.Request) {
	id, _, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	ms.simulator.Lock()
	drone.Arm()
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d armed", id)
	ms.respondSuccess(req, fmt.Sprintf("drone %d armed", id))
}

func (ms *MicroService) handleDisarm(req micro.Request) {
	id, _, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	ms.simulator.Lock()
	drone.Disarm()
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d disarmed", id)
	ms.respondSuccess(req, fmt.Sprintf("drone %d disarmed", id))
}

func (ms *MicroService) handleTakeoff(req micro.Request) {
	id, body, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	var cmd TakeoffCmd
	if len(body) > 0 {
		if err := json.Unmarshal(body, &cmd); err != nil {
			ms.respondError(req, http.StatusBadRequest, "invalid JSON payload")
			return
		}
	}
	if cmd.Altitude <= 0 {
		cmd.Altitude = 10
	}

	ms.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeAltitudeHold)
	drone.AltitudeHold = cmd.Altitude
	drone.SetThrottle(drone.HoverThrottlePercent() * 1.2)
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d takeoff to %.1fm", id, cmd.Altitude)
	ms.respondSuccess(req, fmt.Sprintf("drone %d taking off to %.1fm", id, cmd.Altitude))
}

func (ms *MicroService) handleLand(req micro.Request) {
	id, _, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	ms.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeAltitudeHold)
	drone.AltitudeHold = 0
	drone.SetThrottle(drone.HoverThrottlePercent() * 0.8)
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d landing", id)
	ms.respondSuccess(req, fmt.Sprintf("drone %d landing", id))
}

func (ms *MicroService) handleGoto(req micro.Request) {
	id, body, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	var cmd GotoCmd
	if err := json.Unmarshal(body, &cmd); err != nil {
		ms.respondError(req, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	ms.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeAltitudeHold)
	drone.AltitudeHold = cmd.Y
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d goto altitude %.1f", id, cmd.Y)
	ms.respondSuccess(req, fmt.Sprintf("drone %d going to altitude %.1f (lateral not yet supported)", id, cmd.Y))
}

func (ms *MicroService) handleMode(req micro.Request) {
	id, body, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	var cmd ModeCmd
	if err := json.Unmarshal(body, &cmd); err != nil {
		ms.respondError(req, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	var mode sim.FlightMode
	switch strings.ToLower(cmd.Mode) {
	case "manual":
		mode = sim.FlightModeManual
	case "altitudehold", "altitude":
		mode = sim.FlightModeAltitudeHold
	case "hover":
		mode = sim.FlightModeHover
	default:
		ms.respondError(req, http.StatusBadRequest, fmt.Sprintf("unknown mode: %s", cmd.Mode))
		return
	}

	ms.simulator.Lock()
	drone.SetFlightMode(mode)
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d mode set to %s", id, cmd.Mode)
	ms.respondSuccess(req, fmt.Sprintf("drone %d mode set to %s", id, cmd.Mode))
}

func (ms *MicroService) handleStop(req micro.Request) {
	id, _, err := ms.parseRequest(req)
	if err != nil {
		ms.respondError(req, http.StatusBadRequest, err.Error())
		return
	}

	drone := ms.getDrone(id)
	if drone == nil {
		ms.respondError(req, http.StatusNotFound, fmt.Sprintf("drone %d not found", id))
		return
	}

	ms.simulator.Lock()
	drone.Disarm()
	ms.simulator.Unlock()

	log.Printf("HTTP: drone %d emergency stop", id)
	ms.respondSuccess(req, fmt.Sprintf("drone %d emergency stopped", id))
}
