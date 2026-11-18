package dishwasher

import (
	"errors"
	"math"
)

const (
	MinTemperature              = 30
	MaxTemperature              = 100
	DefaultOperationTemperature = 40 // Only needed in exercise 3
)

var (
	ErrTemperatureOutOfRange  = errors.New("temperature out of range")
	ErrDishwasherNotPoweredOn = errors.New("dishwasher is not powered on")
)

type Dishwasher struct {
	PoweredOn   bool
	Temperature int
	Running     bool
	WattUse     int
}

func NewDishwasher() *Dishwasher {
	return &Dishwasher{
		PoweredOn:   false,
		Temperature: 0,
		Running:     false,
		WattUse:     600,
	}
}

func (d *Dishwasher) PowerOn() {
	d.PoweredOn = true
}

func (d *Dishwasher) PowerOff() {
	d.Temperature = 0
	d.Running = false
	d.PoweredOn = false
}

func (d *Dishwasher) Start() error {
	if !d.PoweredOn {
		return ErrDishwasherNotPoweredOn
	}

	if isTemperatureOutOfRange(d.Temperature) {
		return ErrTemperatureOutOfRange
	}

	d.Running = true

	return nil
}

func (d *Dishwasher) Stop() error {
	if !d.PoweredOn {
		return ErrDishwasherNotPoweredOn
	}

	if d.Running {
		d.Running = false
	}

	d.Temperature = 0

	return nil
}

func (d *Dishwasher) SetTemperature(temp int) error {
	if !d.PoweredOn {
		return ErrDishwasherNotPoweredOn
	}
	if isTemperatureOutOfRange(temp) {
		return ErrTemperatureOutOfRange
	}

	// Set temperature
	d.Temperature = temp

	return nil
}

func (d *Dishwasher) GetWattUseage() int {
	// Naive way of calculating watt usage based on temperature assuming the smarthome will want to start with defaults
	if d.Temperature == 0 {
		return d.WattUse
	}

	multiplication := (float64(d.Temperature) / float64(DefaultOperationTemperature)) * float64(d.WattUse)
	return int(math.Round(multiplication))
}

// Only needed in exercise 3
func (d *Dishwasher) IsRunning() bool {
	return d.Running
}

func (d *Dishwasher) StartWithDefaults() error {
	if !d.PoweredOn {
		d.PowerOn()
	}

	d.SetTemperature(DefaultOperationTemperature)

	return d.Start()
}

func isTemperatureOutOfRange(temp int) bool {
	return temp < MinTemperature || temp > MaxTemperature
}
