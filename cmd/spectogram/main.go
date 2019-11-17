package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	spectogram "github.com/corny/spectrogram"
)

var (
	width   = flag.Int("width", 1024, "set width")
	height  = flag.Int("height", 256, "set height")
	hideavg = flag.Bool("hideavg", false, "hide average")
	out     = flag.String("out", "out.png", "set output filename")
	bins    = flag.Int("bins", 256, "set freq bins")

	bgWaveform = flag.String("bgWaveform", "000", "set background color for the waveform")
	bgFFT      = flag.String("bgFFT", "333", "set background color for FFT")
	maxColor   = flag.String("maxColor", "6b5f7e", "set waveform max color")
	avgColor   = flag.String("avgColor", "0972a2", "set waveform avg color")
)

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("usage: fft [options] file.wav")
		os.Exit(1)
	}

	path := flag.Arg(0)
	samples, err := spectogram.ReadOggFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	bounds := image.Rect(0, 0, *width, *height+*bins)
	img := image.NewRGBA(bounds)

	if *height > 0 {
		params := spectogram.WaveformParams{
			AvgColor: spectogram.ParseColor(*avgColor),
			MaxColor: spectogram.ParseColor(*maxColor),
			Draw:     spectogram.FlagDrawAvg | spectogram.FlagDrawMax,
		}

		sub := img.SubImage(image.Rect(0, 0, *width, *height)).(*image.RGBA)
		draw.Draw(sub, sub.Bounds(), image.NewUniform(spectogram.ParseColor(*bgWaveform)), image.ZP, draw.Src)
		spectogram.DrawWaveform(params, sub, samples)
	}

	if *bins > 0 {
		gr := spectogram.NewGradient()
		gr.Append(
			spectogram.ParseColor("000000"),
			spectogram.ParseColor("380F6D"),
			spectogram.ParseColor("B63679"),
			spectogram.ParseColor("FD9A69"),
			spectogram.ParseColor("FCF6B8"),
		)

		sub := img.SubImage(image.Rect(0, *height, *width, *height+*bins)).(*image.RGBA)
		draw.Draw(sub, sub.Bounds(), image.NewUniform(spectogram.ParseColor(*bgFFT)), image.ZP, draw.Src)
		spectogram.DrawFFT(sub, gr, samples, *bins)
	}

	err = saveImage(img, *out)
	if err != nil {
		log.Fatalf("savePng failed: %v", err)
	}
}

func saveImage(img image.Image, fileName string) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	bf := bufio.NewWriter(outFile)

	err = png.Encode(bf, img)
	if err != nil {
		return err
	}

	err = bf.Flush()
	if err != nil {
		return err
	}

	return nil
}
