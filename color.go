package spectogram

import (
	"image/color"
	"log"
	"strconv"
)

// ParseColor interpretates a string representation of a color as RGBA
func ParseColor(text string) color.RGBA {
	if text == "transparent" {
		return color.RGBA{}
	}
	var err error
	r, g, b, a := uint64(0), uint64(0), uint64(0), uint64(0xffffffff)
	switch len(text) {
	case 4:
		a, err = strconv.ParseUint(text[3:4], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		a |= a << 4
		a |= a << 8
		a |= a << 16
		fallthrough

	case 3:
		r, err = strconv.ParseUint(text[0:1], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		r |= r << 4
		r |= r << 8
		r |= r << 16

		g, err = strconv.ParseUint(text[1:2], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		g |= g << 4
		g |= g << 8
		g |= g << 16

		b, err = strconv.ParseUint(text[2:3], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		b |= b << 4
		b |= b << 8
		b |= b << 16

	case 8:
		a, err = strconv.ParseUint(text[6:8], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}
		fallthrough

	case 6:
		r, err = strconv.ParseUint(text[0:2], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}

		g, err = strconv.ParseUint(text[2:4], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}

		b, err = strconv.ParseUint(text[4:6], 16, 8)
		if err != nil {
			log.Fatalf("invalid color %q", text)
		}

	default:
		log.Fatalf("invalid color %q", text)
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
