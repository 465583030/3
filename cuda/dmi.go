package cuda

import (
	"code.google.com/p/mx3/data"
	"code.google.com/p/mx3/kernel"
)

// Dzyaloshinskii-Moriya interaction
func DMI(Hdm *data.Slice, m *data.Slice, Dx, Dy, Dz float64) {
	mesh := Hdm.Mesh()
	N := mesh.Size()
	c := mesh.CellSize()

	dx := float32(Dx)
	dy := float32(Dy)
	dz := float32(Dz)
	cx := float32(c[0])
	cy := float32(c[1])
	cz := float32(c[2])

	gr, bl := Make2DConf(N[2], N[1])
	kernel.K_dmi(Hdm.DevPtr(0), Hdm.DevPtr(1), Hdm.DevPtr(2),
		m.DevPtr(0), m.DevPtr(1), m.DevPtr(2),
		dx, dy, dz, cx, cy, cz,
		N[0], N[1], N[2], gr, bl)
}
