package smarthome

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/interfaces"
	"github.com/WilliamDeKeyzer/bdd-workshop-sdk/sdk/solarpanel"
)

var _ interfaces.SmartHomeManual = (*SmartHome)(nil)

type SmartHome struct {
	Heartbeat    int
	SolarPanels  []*solarpanel.SolarPanel
	SmartDevices []*interfaces.SmartDevice

	mutex          sync.Mutex
	RunningDevices []*interfaces.SmartDevice

	currentWeather string
}

var (
	weatherMapping = map[string]int{
		"reallysunny": 1000,
		"sunny":       700,
		"cloudy":      300,
		"rainy":       150,
		"twilight":    100,
		"night":       0,
	}
)

func NewSmartHome(ctx context.Context) *SmartHome {
	smartHome := &SmartHome{
		Heartbeat:    1,
		SolarPanels:  []*solarpanel.SolarPanel{},
		SmartDevices: []*interfaces.SmartDevice{},
	}

	// start a goroutine that runs at heartbeat interval to monitor and optimize the household's power usage
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(smartHome.Heartbeat) * time.Second):
				smartHome.Run()
			}
		}
	}()

	return smartHome
}

func (sh *SmartHome) Run() {
	fmt.Printf("SmartHome brain 🧠 is running with %d solar panels and %d smart devices\n", len(sh.SolarPanels), len(sh.SmartDevices))

	// This is the brain func that will attempt to power on devices in a smart fashion whenever enough electricity is available
	// It will check the total watt usage of all smart devices and compare it with the available power from solar panels.

	availablePower := sh.GetInjection()
	currentWattUseage := sh.GetConsumption()

	if availablePower >= currentWattUseage {
		for _, device := range sh.SmartDevices {
			if (*device).GetWattUseage() <= (availablePower - currentWattUseage) {
				_ = (*device).StartWithDefaults()
				sh.RunningDevices = append(sh.RunningDevices, device)

				// Start a timer to turn it off automatically
				go func(d *interfaces.SmartDevice) {
					// Wait by sleeping
					time.Sleep(3 * time.Second) // Yes we have very fast smart devices

					_ = (*d).Stop()

					// Remove the device from `RunningDevices`` after stopping
					sh.mutex.Lock()
					for i, runningDevice := range sh.RunningDevices {
						if runningDevice == d {
							sh.RunningDevices = append(sh.RunningDevices[:i], sh.RunningDevices[i+1:]...)
							break
						}
					}
					sh.mutex.Unlock()
				}(device)

				currentWattUseage += (*device).GetWattUseage()
			}
		}
	}
}

func (sh *SmartHome) AddSolarPanel(panel *solarpanel.SolarPanel) {
	wattUseage, ok := weatherMapping[sh.currentWeather]
	if !ok {
		// Default to 0 for unknown weather
		wattUseage = 0
	}
	panel.WattGeneration = wattUseage

	sh.SolarPanels = append(sh.SolarPanels, panel)
}

func (sh *SmartHome) AddSmartDevice(device *interfaces.SmartDevice) {
	sh.SmartDevices = append(sh.SmartDevices, device)
}

func (sh *SmartHome) GetInjection() int {
	injection := 0
	for _, panel := range sh.SolarPanels {
		injection += panel.GetInjection()
	}
	return injection
}

func (sh *SmartHome) GetConsumption() int {
	currentWattUseage := 0
	for _, device := range sh.SmartDevices {
		if (*device).IsRunning() {
			currentWattUseage += (*device).GetWattUseage()
		}
	}
	return currentWattUseage
}

func (sh *SmartHome) GetRunningDevices() []*interfaces.SmartDevice {
	return sh.RunningDevices
}

func (sh *SmartHome) GetSmartDevices() []*interfaces.SmartDevice {
	return sh.SmartDevices
}

func (sh *SmartHome) GetSolarPanels() []*solarpanel.SolarPanel {
	return sh.SolarPanels
}

// Normally a solar panel is smart enough to handle this internally, but we simulate it here for testing purposes
func (sh *SmartHome) HandleWeatherLogic(weather string) {
	sh.currentWeather = weather

	// weather change logic: we adjust solar panel output
	newWattUseage, ok := weatherMapping[weather]
	if !ok {
		// Default to 0 for unknown weather
		newWattUseage = 0
	}

	for _, panel := range sh.SolarPanels {
		panel.SetNewInjection(newWattUseage)
	}
}
