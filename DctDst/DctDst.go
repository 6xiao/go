package DctDst

import "math"

func DST(in []float64) []float64 {
	out := make([]float64, len(in))
	nr := len(in)

	for i := 0; i < nr; i++ {
		for j := 0; j < nr; j++ {
			tmp := (float64(i) + 0.5) * (float64(j) + 0.5) / float64(nr)
			out[i] += in[j] * math.Sin(tmp*math.Pi)
		}
	}

	return out
}

func IDST(in []float64) []float64 {
	out := DST(in)

	for i, v := range out {
		out[i] = v * 2 / float64(len(in))
	}

	return out
}

func DCT(in []float64) []float64 {
	out := make([]float64, len(in))
	nr := len(in)

	for i := 0; i < nr; i++ {
		for j := 0; j < nr; j++ {
			tmp := float64(i) * (float64(j) + 0.5) / float64(nr)
			out[i] += in[j] * math.Cos(tmp*math.Pi)
		}
	}

	return out
}

func IDCT(in []float64) []float64 {
	out := make([]float64, len(in))
	nr := len(in)

	for i := 0; i < nr; i++ {
		out[i] = in[i] / 2

		for j := 0; j < nr; j++ {
			tmp := float64(j) * (float64(i) + 0.5) / float64(nr)
			out[i] += in[j] * math.Cos(tmp*math.Pi)
		}

		out[i] = out[i] * 2 / float64(nr)
	}

	return out
}
