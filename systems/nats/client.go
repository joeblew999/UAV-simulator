package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	sim "drone-simulator/internal/sim"

	"github.com/nats-io/nats.go"
)

// Client handles NATS connection and message routing for the simulator.
type Client struct {
	nc        *nats.Conn
	simulator *sim.Simulator
	subs      []*nats.Subscription
	stopCh    chan struct{}
	wg        sync.WaitGroup

	// Telemetry config
	telemetryHz float64
}

// TelemetryMsg is published to drone.<id>.telemetry
type TelemetryMsg struct {
	ID         int       `json:"id"`
	Timestamp  int64     `json:"timestamp"`
	Position   Vec3Msg   `json:"position"`
	Velocity   Vec3Msg   `json:"velocity"`
	Rotation   Vec3Msg   `json:"rotation"`
	Battery    float64   `json:"battery"`
	FlightMode string    `json:"flightMode"`
	Throttle   float64   `json:"throttle"`
	Armed      bool      `json:"armed"`
	OnGround   bool      `json:"onGround"`
	Destroyed  bool      `json:"destroyed"`
}

type Vec3Msg struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// GotoCmd is received on drone.<id>.goto
type GotoCmd struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// TakeoffCmd is received on drone.<id>.takeoff
type TakeoffCmd struct {
	Altitude float64 `json:"altitude"`
}

// InputCmd is received on drone.<id>.input
type InputCmd struct {
	Throttle float64 `json:"throttle"`
	Yaw      float64 `json:"yaw"`
	Pitch    float64 `json:"pitch"`
	Roll     float64 `json:"roll"`
}

// ModeCmd is received on drone.<id>.mode
type ModeCmd struct {
	Mode string `json:"mode"`
}

// New creates a new NATS client connected to the given URL.
func New(url string, simulator *sim.Simulator) (*Client, error) {
	nc, err := nats.Connect(url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Println("NATS reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	c := &Client{
		nc:          nc,
		simulator:   simulator,
		stopCh:      make(chan struct{}),
		telemetryHz: 10, // Default 10 Hz
	}

	return c, nil
}

// Start begins listening for commands and publishing telemetry.
func (c *Client) Start() error {
	// Subscribe to command topics
	if err := c.subscribeCommands(); err != nil {
		return err
	}

	// Start telemetry publisher
	c.wg.Add(1)
	go c.publishTelemetryLoop()

	log.Printf("NATS client started (telemetry @ %.0f Hz)", c.telemetryHz)
	return nil
}

// Stop gracefully shuts down the NATS client.
func (c *Client) Stop() {
	close(c.stopCh)
	c.wg.Wait()

	for _, sub := range c.subs {
		sub.Unsubscribe()
	}

	c.nc.Drain()
	log.Println("NATS client stopped")
}

// Conn returns the underlying NATS connection for use by other services.
func (c *Client) Conn() *nats.Conn {
	return c.nc
}

func (c *Client) subscribeCommands() error {
	// drone.<id>.arm
	sub, err := c.nc.Subscribe("drone.*.arm", c.handleArm)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.disarm
	sub, err = c.nc.Subscribe("drone.*.disarm", c.handleDisarm)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.takeoff
	sub, err = c.nc.Subscribe("drone.*.takeoff", c.handleTakeoff)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.land
	sub, err = c.nc.Subscribe("drone.*.land", c.handleLand)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.goto
	sub, err = c.nc.Subscribe("drone.*.goto", c.handleGoto)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.input
	sub, err = c.nc.Subscribe("drone.*.input", c.handleInput)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.mode
	sub, err = c.nc.Subscribe("drone.*.mode", c.handleMode)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	// drone.<id>.stop (emergency stop)
	sub, err = c.nc.Subscribe("drone.*.stop", c.handleStop)
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	return nil
}

func (c *Client) parseDroneID(subject string) (int, error) {
	// subject format: drone.<id>.command
	parts := strings.Split(subject, ".")
	if len(parts) < 2 {
		return -1, fmt.Errorf("invalid subject: %s", subject)
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1, fmt.Errorf("invalid drone ID: %s", parts[1])
	}
	return id, nil
}

func (c *Client) getDrone(id int) *sim.Drone {
	drones := c.simulator.Drones()
	if id < 0 || id >= len(drones) {
		return nil
	}
	return drones[id]
}

func (c *Client) handleArm(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("arm: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("arm: drone %d not found", id)
		return
	}

	c.simulator.Lock()
	drone.Arm()
	c.simulator.Unlock()
	log.Printf("drone %d armed", id)
}

func (c *Client) handleDisarm(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("disarm: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("disarm: drone %d not found", id)
		return
	}

	c.simulator.Lock()
	drone.Disarm()
	c.simulator.Unlock()
	log.Printf("drone %d disarmed", id)
}

func (c *Client) handleTakeoff(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("takeoff: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("takeoff: drone %d not found", id)
		return
	}

	var cmd TakeoffCmd
	if len(msg.Data) > 0 {
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			log.Printf("takeoff: invalid payload: %v", err)
			return
		}
	}
	if cmd.Altitude <= 0 {
		cmd.Altitude = 10 // Default altitude
	}

	c.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeAltitudeHold)
	drone.AltitudeHold = cmd.Altitude
	drone.SetThrottle(drone.HoverThrottlePercent() * 1.2) // Slight boost for takeoff
	c.simulator.Unlock()
	log.Printf("drone %d taking off to %.1fm", id, cmd.Altitude)
}

func (c *Client) handleLand(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("land: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("land: drone %d not found", id)
		return
	}

	c.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeAltitudeHold)
	drone.AltitudeHold = 0 // Descend to ground
	drone.SetThrottle(drone.HoverThrottlePercent() * 0.8) // Reduce for descent
	c.simulator.Unlock()
	log.Printf("drone %d landing", id)
}

func (c *Client) handleGoto(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("goto: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("goto: drone %d not found", id)
		return
	}

	var cmd GotoCmd
	if err := json.Unmarshal(msg.Data, &cmd); err != nil {
		log.Printf("goto: invalid payload: %v", err)
		return
	}

	// GOTO currently only supports altitude control.
	// Lateral positioning (X, Z) is not implemented because the simulator's
	// internal physics (stability damping, motor torques) conflicts with
	// external torque control. Full 3D positioning would require either:
	// - Modifications to the simulator's flight controller
	// - Use of swarm mode (which has its own lateral control)
	c.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeAltitudeHold)
	drone.AltitudeHold = cmd.Y
	c.simulator.Unlock()

	log.Printf("drone %d goto altitude %.1f (lateral X=%.1f Z=%.1f not yet supported)", id, cmd.Y, cmd.X, cmd.Z)
}

func (c *Client) handleInput(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("input: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("input: drone %d not found", id)
		return
	}

	var cmd InputCmd
	if err := json.Unmarshal(msg.Data, &cmd); err != nil {
		log.Printf("input: invalid payload: %v", err)
		return
	}

	c.simulator.Lock()
	drone.SetFlightMode(sim.FlightModeManual)
	drone.SetThrottle(cmd.Throttle * 100) // Convert 0-1 to 0-100%
	// Yaw/Pitch/Roll would need additional methods on Drone
	c.simulator.Unlock()
}

func (c *Client) handleMode(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("mode: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("mode: drone %d not found", id)
		return
	}

	var cmd ModeCmd
	if err := json.Unmarshal(msg.Data, &cmd); err != nil {
		log.Printf("mode: invalid payload: %v", err)
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
		log.Printf("mode: unknown mode: %s", cmd.Mode)
		return
	}

	c.simulator.Lock()
	drone.SetFlightMode(mode)
	c.simulator.Unlock()
	log.Printf("drone %d mode set to %s", id, cmd.Mode)
}

func (c *Client) handleStop(msg *nats.Msg) {
	id, err := c.parseDroneID(msg.Subject)
	if err != nil {
		log.Printf("stop: %v", err)
		return
	}
	drone := c.getDrone(id)
	if drone == nil {
		log.Printf("stop: drone %d not found", id)
		return
	}

	c.simulator.Lock()
	drone.Disarm()
	c.simulator.Unlock()
	log.Printf("drone %d emergency stop", id)
}

func (c *Client) publishTelemetryLoop() {
	defer c.wg.Done()

	interval := time.Duration(float64(time.Second) / c.telemetryHz)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.publishTelemetry()
		}
	}
}

func (c *Client) publishTelemetry() {
	c.simulator.RLock()
	drones := c.simulator.Drones()

	for i, d := range drones {
		msg := TelemetryMsg{
			ID:        i,
			Timestamp: time.Now().UnixMilli(),
			Position:  Vec3Msg{X: d.Position.X, Y: d.Position.Y, Z: d.Position.Z},
			Velocity:  Vec3Msg{X: d.Velocity.X, Y: d.Velocity.Y, Z: d.Velocity.Z},
			Rotation:  Vec3Msg{X: d.Rotation.X, Y: d.Rotation.Y, Z: d.Rotation.Z},
			Battery:   d.BatteryPercent,
			FlightMode: flightModeString(d.FlightMode),
			Throttle:  d.ThrottlePercent,
			Armed:     d.IsArmed,
			OnGround:  d.OnGround,
			Destroyed: d.Destroyed,
		}

		data, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		c.nc.Publish(TelemetrySubject(i), data)
	}
	c.simulator.RUnlock()
}

func flightModeString(mode sim.FlightMode) string {
	switch mode {
	case sim.FlightModeManual:
		return "Manual"
	case sim.FlightModeAltitudeHold:
		return "AltitudeHold"
	case sim.FlightModeHover:
		return "Hover"
	default:
		return "Unknown"
	}
}

