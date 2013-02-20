package cuda

// common code for all convolutions.

import (
	"code.google.com/p/mx3/data"
	"code.google.com/p/mx3/util"
	"github.com/barnex/fmath"
	"log"
)

// Output size of R2C FFT with given logic size, expressed in floats.
func fftR2COutputSizeFloats(logicSize [3]int) [3]int {
	return [3]int{logicSize[0], logicSize[1], 2 * (logicSize[2]/2 + 1)}
}

func prod(size [3]int) int {
	return size[0] * size[1] * size[2]
}

// Extract real parts, copy them from src to dst.
// In the meanwhile, check if imaginary parts are nearly zero
// and scale the kernel to compensate for unnormalized FFTs.
func scaleRealParts(dst, src *data.Slice, scale float32) {
	util.Argument(2*dst.Len() == src.Len())
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)

	srcList := src.HostCopy().Host()[0]
	dstList := dst.Host()[0]

	// Normally, the FFT'ed kernel is purely real because of symmetry,
	// so we only store the real parts...
	maximg := float32(0.)
	maxreal := float32(0.)
	for i := 0; i < src.Len()/2; i++ {
		dstList[i] = srcList[2*i] * scale
		if fmath.Abs(srcList[2*i+0]) > maxreal {
			maxreal = fmath.Abs(srcList[2*i+0])
		}
		if fmath.Abs(srcList[2*i+1]) > maximg {
			maximg = fmath.Abs(srcList[2*i+1])
		}
	}
	// ...however, we check that the imaginary parts are nearly zero,
	// just to be sure we did not make a mistake during kernel creation.
	if maximg/maxreal > FFT_IMAG_TOLERANCE {
		log.Fatalf("Too large FFT kernel imaginary/real part: %v", maximg/maxreal)
	}
}

const FFT_IMAG_TOLERANCE = 1e-5
