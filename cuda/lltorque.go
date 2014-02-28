package cuda

import (
	"github.com/mumax/3/data"
	"unsafe"
)

// Landau-Lifshitz torque divided by gamma0:
// 	- 1/(1+α²) [ m x B +  α m x (m x B) ]
// 	torque in Tesla
// 	m normalized
// 	B in Tesla
// see lltorque.cu
func LLTorque(torque, m, B *data.Slice, alpha data.LUTPtr, regions *Bytes) {
	N := torque.Len()
	cfg := make1DConf(N)

	k_lltorque_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		unsafe.Pointer(alpha), regions.Ptr, N, cfg)
}
