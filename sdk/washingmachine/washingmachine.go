package washingmachine

import (
	"errors"
	"math"
)

const (
	MinTemperature              = 30
	MaxTemperature              = 90
	DefaultOperationTemperature = 45 // Only needed in exercise 3
)

var (
	ErrTemperatureOutOfRange      = errors.New("temperature out of range")
	ErrWashingMachineNotPoweredOn = errors.New("washing machine is not powered on")
)

type WashingMachine struct {
	PoweredOn   bool
	Temperature int
	Running     bool
	WattUse     int
}

func NewWashingMachine() *WashingMachine {
	return &WashingMachine{
		PoweredOn:   false,
		Temperature: 0,
		Running:     false,
		WattUse:     900,
	}
}

func (d *WashingMachine) PowerOn() {
	d.PoweredOn = true
}

func (d *WashingMachine) PowerOff() {
	d.Temperature = 0
	d.Running = false
	d.PoweredOn = false
}

func (d *WashingMachine) Start() error {
	if !d.PoweredOn {
		return ErrWashingMachineNotPoweredOn
	}

	if isTemperatureOutOfRange(d.Temperature) {
		return ErrTemperatureOutOfRange
	}

	d.Running = true

	return nil
}

func (d *WashingMachine) Stop() error {
	if !d.PoweredOn {
		return ErrWashingMachineNotPoweredOn
	}

	if d.Running {
		d.Running = false
	}

	d.Temperature = 0

	return nil
}

func (d *WashingMachine) SetTemperature(temp int) error {
	if !d.PoweredOn {
		return ErrWashingMachineNotPoweredOn
	}
	if isTemperatureOutOfRange(temp) {
		return ErrTemperatureOutOfRange
	}

	// Set temperature
	d.Temperature = temp

	return nil
}

// Only needed in exercise 3
func (d *WashingMachine) GetWattUseage() int {
	// Naive way of calculating watt usage based on temperature assuming the smarthome will want to start with defaults
	if d.Temperature == 0 {
		return d.WattUse
	}

	multiplication := (float64(d.Temperature) / float64(DefaultOperationTemperature)) * float64(d.WattUse)
	return int(math.Round(multiplication))
}

// Only needed in exercise 3
func (d *WashingMachine) IsRunning() bool {
	return d.Running
}

// Only needed in exercise 3
func (d *WashingMachine) StartWithDefaults() error {
	if !d.PoweredOn {
		d.PowerOn()
	}

	d.SetTemperature(DefaultOperationTemperature)

	return d.Start()
}

func isTemperatureOutOfRange(temp int) bool {
	return temp < MinTemperature || temp > MaxTemperature
}
