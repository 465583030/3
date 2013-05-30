package engine

import (
	"code.google.com/p/mx3/cuda"
	"code.google.com/p/mx3/util"
)

// average in userspace XYZ order
// does not yet take into account volume.
// pass volume parameter, possibly nil?
func Average(s GPU_Getter) []float64 {
	b, recycle := s.GetGPU()
	if recycle {
		defer cuda.RecycleBuffer(b)
	}
	nComp := b.NComp()
	avg := make([]float64, nComp)
	for i := range avg {
		I := util.SwapIndex(i, nComp)
		avg[i] = float64(cuda.Sum(b.Comp(I))) / float64(b.Mesh().NCell())
	}
	return avg
}
