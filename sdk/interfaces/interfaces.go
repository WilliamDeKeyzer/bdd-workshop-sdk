package interfaces

import (
	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/dishwasher"
	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/solarpanel"
	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/washingmachine"
)

type SmartDevice interface {
	PowerOn()
	PowerOff()
	Start() error
	StartWithDefaults() error
	Stop() error
	GetWattUseage() int
	IsRunning() bool
}

type SolarPanelManual interface {
	SetNewInjection(watt int) error
	GetInjection() int
}

type SmartHomeManual interface {
	AddSolarPanel(panel *solarpanel.SolarPanel)
	AddSmartDevice(device *SmartDevice)
	GetInjection() int
	GetConsumption() int
	GetRunningDevices() []*SmartDevice
	GetSmartDevices() []*SmartDevice
	GetSolarPanels() []*solarpanel.SolarPanel
	HandleWeatherLogic(weather string)
}

// Compile-time checks to ensure valid interfaces
var _ SmartDevice = (*dishwasher.Dishwasher)(nil)
var _ SmartDevice = (*washingmachine.WashingMachine)(nil)
var _ SolarPanelManual = (*solarpanel.SolarPanel)(nil)
