package spectogram

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ParseColor interpretates a string representation of a color as RGBA
func TestParseHex(t *testing.T) {

	tests := []struct {
		input  string
		result color.RGBA
	}{
		{"000", color.RGBA{0, 0, 0, 255}},
		{"f93", color.RGBA{0xff, 0x99, 0x33, 255}},
		{"f938", color.RGBA{0xff, 0x99, 0x33, 0x88}},
		{"0609f9", color.RGBA{0x06, 0x09, 0xf9, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

			assert.Equal(t, tt.result, ParseColor(tt.input))
		})
	}
}
