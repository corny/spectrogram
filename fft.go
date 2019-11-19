package spectogram

import (
	"image/draw"
	"math"
	"math/cmplx"
)

// DrawFFT calculates and draws the fast Fourier transform (FFT)
func DrawFFT(img draw.Image, gr Gradient, samples []float64, bins int) {
	bn := img.Bounds()
	sub := make([]float64, bins*2)

	for x := 0; x < bn.Dx(); x++ {
		n0 := int64(mapRange(float64(x+0), 0, float64(bn.Dx()), 0, float64(len(samples))))
		n1 := int64(mapRange(float64(x+1), 0, float64(bn.Dx()), 0, float64(len(samples))))
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

		for y := 0; y < bins; y++ {
			r := cmplx.Abs(freqs[y])
			img.Set(x, y+bn.Min.Y, gr.ColorAt(r))
		}
	}
}

func dft(input []float64) []complex128 {
	output := make([]complex128, len(input))

	arg := -2.0 * math.Pi / float64(len(input))
	for k := range input {
		r, i := 0.0, 0.0
		for n := 0; n < len(input); n++ {
			r += input[n] * math.Cos(arg*float64(n)*float64(k))
			i += input[n] * math.Sin(arg*float64(n)*float64(k))
		}
		output[k] = complex(r, i)
	}
	return output
}

func hfft(samples []float64, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(samples[0], 0)
		return
	}

	half := n / 2

	hfft(samples, freqs, half, 2*step)
	hfft(samples[step:], freqs[half:], half, 2*step)

	for k := 0; k < half; k++ {
		a := -2 * math.Pi * float64(k) / float64(n)
		e := cmplx.Rect(1, a) * freqs[k+half]

		freqs[k], freqs[k+half] = freqs[k]+e, freqs[k]-e
	}
}

func fft(samples []float64) []complex128 {
	n := len(samples)
	freqs := make([]complex128, n)
	hfft(samples, freqs, n, 1)
	return freqs
}
