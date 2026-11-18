package dishwasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetWattUseage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		temp     int
		expected int
	}{
		"default": {
			temp:     DefaultOperationTemperature,
			expected: 600,
		},
		"30_degrees": {
			temp:     30,
			expected: 450,
		},
		"45_degrees": {
			temp:     45,
			expected: 675,
		},
		"70_degrees": {
			temp:     70,
			expected: 1050,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			d := NewDishwasher()
			d.PowerOn()
			d.SetTemperature(tc.temp)

			// Act
			got := d.GetWattUseage()

			// Assert
			require.Equal(t, tc.expected, got)
		})
	}
}
