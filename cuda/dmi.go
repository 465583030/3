package cuda

import (
	"github.com/mumax/3/data"
	"github.com/mumax/3/util"
)

// Add effective field of Dzyaloshinskii-Moriya interaction to Beff (Tesla).
// According to Bagdanov and Röβler, PRL 87, 3, 2001. eq.8 (out-of-plane symmetry breaking).
// See dmi.cu
func AddDMI(Beff *data.Slice, m *data.Slice, DL_red, DH_red, A_red float32, str int) {
	mesh := Beff.Mesh()
	N := mesh.Size()
	c := mesh.CellSize()

	util.Argument(m.Mesh().Size() == mesh.Size())
	util.AssertMsg(N[0] == 1, "DMI available in 2D only")
	util.AssertMsg(mesh.PBC_code() == 0, "DMI not available with PBC")

	cfg := make3DConf(N)

	k_adddmi_async(Beff.DevPtr(0), Beff.DevPtr(1), Beff.DevPtr(2),
		m.DevPtr(0), m.DevPtr(1), m.DevPtr(2),
		float32(c[0]), float32(c[1]), float32(c[2]),
		DL_red, DH_red, A_red, N[0], N[1], N[2], cfg, str)
}
