package spectogram

import "image/color"

type Gradient []color.RGBA

func NewGradient() Gradient {
	return make(Gradient, 0)
}

func (grad *Gradient) Append(c ...color.RGBA) {
	*grad = append(*grad, c...)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (grad Gradient) Index(i float64) (int, float64) {
	n := len(grad)
	if n == 0 {
		return -1, -1
	}
	c := clamp(i, 0, 1)
	if c == 1 {
		return n - 1, 0.0
	}
	x := int(c * float64(n-1))
	return x, (c - float64(x)/float64(n-1)) * float64(n-1)
}

func (grad Gradient) ColorAt(i float64) color.RGBA {
	n := len(grad)
	if n == 0 {
		return color.RGBA{}
	}
	if n == 1 {
		return grad[0]
	}
	nx, ni := grad.Index(i)
	if nx >= n-1 {
		return grad[nx]
	}

	return color.RGBA{
		R: uint8((float64(grad[nx+1].R)-float64(grad[nx].R))*ni + float64(grad[nx].R)),
		G: uint8((float64(grad[nx+1].G)-float64(grad[nx].G))*ni + float64(grad[nx].G)),
		B: uint8((float64(grad[nx+1].B)-float64(grad[nx].B))*ni + float64(grad[nx].B)),
		A: uint8((float64(grad[nx+1].A)-float64(grad[nx].A))*ni + float64(grad[nx].A)),
	}
}
