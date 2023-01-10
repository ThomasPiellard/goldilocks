package fft

import (
	"fmt"
	"testing"

	"github.com/ThomasPiellard/goldilocks/fr"
)

func BenchmarkFFT(b *testing.B) {

	sizes := []int{32, 64, 128, 256}

	for _, v := range sizes {

		b.Run(fmt.Sprintf("[FFT SIZE %d]", v), func(b *testing.B) {

			a := make([]fr.Element, v)
			d := NewDomain(uint64(v))

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				d.FFT(a, DIF, true)
			}

		})

	}

}
