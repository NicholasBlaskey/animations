package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// https://www.nayuki.io/res/how-to-implement-the-discrete-fourier-transform/dft.py
// https://towardsdatascience.com/understanding-audio-data-fourier-transform-fft-spectrogram-and-speech-recognition-a4072d228520
func getDFT(in []float32) []complex64 {
	n := len(in)
	output := []complex64{}
	for k := 0; k < n; k++ {
		s := complex64(0)
		for t := 0; t < n; t++ {
			angle := complex(2, 1) * complex(float32(math.Pi), 0) *
				complex(float32(t)*float32(k)/float32(n), 0)
			s += complex(in[t], 0) * complex64(cmplx.Exp(-complex128(angle)))
		}
		output = append(output, s)
	}
	return output
}

func main() {
	wave := []float32{}
	for i := float64(-2); i < 2; i += 0.1 {
		wave = append(wave, float32(math.Sin(i)))
	}

	fmt.Println(getDFT(wave))
}
