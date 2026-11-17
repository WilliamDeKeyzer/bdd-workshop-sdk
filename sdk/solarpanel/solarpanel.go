package solarpanel

import "errors"

var (
	ErrInvalidWattGeneration = errors.New("invalid watt generation")
)

type SolarPanel struct {
	WattGeneration int
}

func NewSolarPanel() *SolarPanel {
	return &SolarPanel{
		WattGeneration: 0,
	}
}

func (sp *SolarPanel) SetNewInjection(watt int) error {
	if watt <= 0 {
		return ErrInvalidWattGeneration
	}

	sp.WattGeneration = watt

	return nil
}

func (sp *SolarPanel) GetInjection() int {
	return sp.WattGeneration
}
