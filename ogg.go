package spectogram

import (
	"io"
	"os"

	"github.com/jfreymuth/oggvorbis"
)

// ReadOggFile reads a ogg/vorbis file
func ReadOggFile(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r, err := oggvorbis.NewReader(file)
	if err != nil {
		return nil, err
	}

	return ReadOgg(r)
}

// ReadOgg reads a ogg/vorbis stream
func ReadOgg(r *oggvorbis.Reader) ([]float64, error) {
	numChannels := r.Channels()
	samples := make([]float64, r.Length())
	read := 0

	for {
		buffer := make([]float32, 8192)
		n, err := r.Read(buffer)

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return samples, err
		}

		for i := 0; i < n; i += numChannels {
			read++
			if read >= len(samples) {
				// might happen if one channel has more samples than another
				break
			}

			// reduce to a single channel
			// out of phase channels will cancel out!
			for j := 0; j < numChannels; j++ {
				samples[read] += float64(buffer[i+j])
			}
			samples[read] /= float64(numChannels)
		}
	}
}
