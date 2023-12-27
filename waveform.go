package spectogram

import (
	"image/color"
	"image/draw"
	"math"
	"math/cmplx"
)

type WaveformParams struct {
	AvgColor color.Color
	MaxColor color.Color
	Draw     int
}

type Line struct {
	Min      float64
	Max      float64
	Avg      float64
	Strength float64
}

const (
	FlagDrawAvg = 1 // Draws the average
	FlagDrawMax = 2 // Draws the maximum
)

// DrawWaveform calculates and draws the waveform
func DrawWaveform(params WaveformParams, img draw.Image, samples []float64) {
	bn := img.Bounds()

	const bins = 32
	sub := make([]float64, bins*2)

	var gmax float64

	// -------------------------------------------

	height := bn.Dy()
	lines := make([]Line, bn.Dx())

	for i := 0; i < bn.Dx(); i++ {
		n0 := int64(mapRange(float64(i-0), 0, float64(bn.Dx()), 0, float64(len(samples))))
		n1 := int64(mapRange(float64(i+1), 0, float64(bn.Dx()), 0, float64(len(samples))))

		sum, min, max := 0.0, 0.0, 0.0
		for i := n0; i < n1; i++ {
			smp := samples[i]
			sum += math.Abs(smp)
			max = math.Max(smp, max)
			min = math.Min(smp, min)
		}

		lines[i].Min = min
		lines[i].Max = max
		lines[i].Avg = sum / float64(n1-n0)

		// FFT
		c := int(n0 + (n1-n0)/2)
		for i := range sub {
			s := 0.0
			n := c - bins + i
			if n >= 0 && n < len(samples) {
				s = samples[n]
			}

			// Apply Hamming window
			s *= 0.54 - 0.46*math.Cos(float64(i)*math.Pi*2/float64(len(sub)))
			sub[i] = s
		}

		freqs := fft(sub)

		for y := 0; y < bins/3; y++ {
			lines[i].Strength += cmplx.Abs(freqs[y])
		}
	}

	// use the average values as boundaries
	for i := 0; i < len(lines); i++ {
		gmax = math.Max(lines[i].Avg, gmax)
	}

	gr := NewGradient()
	gr.Append(
		ParseColor("0000cc"),
		ParseColor("00cc00"),
		ParseColor("cc0000"),
	)

	for i := range lines {
		line := &lines[i]

		if params.Draw&FlagDrawAvg != 0 {
			val := int(mapRange(line.Avg, 0, gmax, 0, float64(height)))
			if val != 0 {

				drawLine(img, i, height-val, i, height, gr.ColorAt(line.Strength/15))
			}
		}
	}
}
