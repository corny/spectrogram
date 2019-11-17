package spectogram

import (
	"image/color"
	"image/draw"
	"math"
)

type WaveformParams struct {
	AvgColor color.Color
	MaxColor color.Color
	Draw     int
}

type Line struct {
	Min float64
	Max float64
	Avg float64
}

const (
	FlagDrawAvg = 1 // Draws the average
	FlagDrawMax = 2 // Draws the maximum
)

// DrawWaveform calculates and draws the waveform
func DrawWaveform(params WaveformParams, img draw.Image, samples []float64) {
	bn := img.Bounds()

	var gmin, gmax, gabs float64

	// -------------------------------------------

	middle := bn.Dy() / 2
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
	}

	// Calculate the boundaries
	if params.Draw&FlagDrawMax != 0 {
		// use min/max as boundaries
		for i := 0; i < len(samples); i++ {
			gsmp := samples[i]
			gmax = math.Max(gsmp, gmax)
			gmin = math.Min(gsmp, gmin)
		}
		gabs = math.Max(math.Abs(gmin), math.Abs(gmax))
	} else {
		// use the average values as boundaries
		for i := 0; i < len(lines); i++ {
			gmax = math.Max(lines[i].Avg, gmax)
		}
		gmin = -gmax
	}

	for i := range lines {
		line := &lines[i]

		if params.Draw&FlagDrawMax != 0 {
			s0 := int(mapRange(line.Min, -gabs, gabs, -float64(middle), float64(middle)))
			s1 := int(mapRange(line.Max, -gabs, gabs, -float64(middle), float64(middle)))
			if s0 != 0 || s1 != 0 {
				drawLine(img, i, middle-s0, i, middle-s1, params.MaxColor)
			}
		}

		if params.Draw&FlagDrawAvg != 0 {
			val := int(mapRange(line.Avg, gmin, gmax, -float64(middle), float64(middle)))
			if val != 0 {
				drawLine(img, i, middle-val, i, middle+val, params.AvgColor)
			}
		}
	}
}
