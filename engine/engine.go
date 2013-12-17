package engine

import (
	"github.com/mumax/3/cuda"
	"github.com/mumax/3/data"
	"github.com/mumax/3/mag"
	"log"
	"runtime"
)

const VERSION = "mumax3.3"

var UNAME = VERSION + " " + runtime.GOOS + "_" + runtime.GOARCH + " " + runtime.Version() + " (" + runtime.Compiler + ")"

var (
	globalmesh    data.Mesh     // mesh for m and everything that has the same size
	M             magnetization // reduced magnetization (unit length)
	B_eff, Torque setter
)

func init() {
	DeclFunc("SetGridSize", setGridSize, `Sets the number of cells for X,Y,Z`)
	DeclFunc("SetCellSize", setCellSize, `Sets the X,Y,Z cell size in meters`)
	DeclFunc("SetPBC", setPBC, `Sets number of repetitions in X,Y,Z`)

	// magnetization
	M.init(3, "m", "", `Reduced magnetization (unit length)`, &globalmesh)
	DeclLValue("m", &M, `Reduced magnetization (unit length)`)

	// effective field
	B_eff.init(3, &globalmesh, "B_eff", "T", "Effective field", func(dst *data.Slice) {
		B_demag.set(dst)
		B_exch.addTo(dst)
		B_anis.addTo(dst)
		B_ext.addTo(dst)
		B_therm.addTo(dst)
		cuda.Sync(0)
	})

	// torque inited in torque.go
}

func Mesh() *data.Mesh {
	checkMesh()
	return &globalmesh
}

func WorldSize() [3]float64 {
	w := Mesh().WorldSize()
	return [3]float64{w[2], w[1], w[0]} // swaps XYZ
}

// Set the simulation mesh to Nx x Ny x Nz cells of given size.
// Can be set only once at the beginning of the simulation.
func SetMesh(Nx, Ny, Nz int, cellSizeX, cellSizeY, cellSizeZ float64) {
	if Nx <= 1 {
		log.Fatal("mesh size X should be > 1, have: ", Nx)
	}
	globalmesh = *data.NewMesh(Nz, Ny, Nx, cellSizeZ, cellSizeY, cellSizeX, pbczyx...)
	log.Println("set mesh:", Mesh().UserString())
	alloc()
}

// allocate m and regions buffer (after mesh is set)
func alloc() {
	M.alloc()
	regions.alloc()

	Solver = NewSolver(M.buffer, Torque.set, normalize, 1e-15, mag.Gamma0, HeunStep)
	solvertype = 2 // HeunStep

	Table.Add(&M)
}

func normalize(m *data.Slice) {
	cuda.Normalize(m, nil)
}

// for lazy setmesh: set gridsize and cellsize in separate calls
var (
	gridsize []int
	cellsize []float64
	pbczyx   []int
)

func setGridSize(Nx, Ny, Nz int) {
	gridsize = []int{Nx, Ny, Nz}
	if cellsize != nil {
		SetMesh(Nx, Ny, Nz, cellsize[0], cellsize[1], cellsize[2])
	}
}

func setCellSize(cx, cy, cz float64) {
	cellsize = []float64{cx, cy, cz}
	if gridsize != nil {
		SetMesh(gridsize[0], gridsize[1], gridsize[2], cx, cy, cz)
	}
}

func setPBC(nx, ny, nz int) {
	if pbczyx != nil {
		log.Panicf("PBC alread set")
	}
	if globalmesh.Size() != [3]int{0, 0, 0} {
		log.Panicf("PBC must be set before MeshSize and GridSize")
	}
	pbczyx = []int{nz, ny, nx}
}

// check if mesh is set
func checkMesh() {
	if globalmesh.Size() == [3]int{0, 0, 0} {
		log.Panic("need to set mesh first")
	}
}

// check if m is set
func checkM() {
	checkMesh()
	if M.buffer.DevPtr(0) == nil {
		log.Fatal("need to initialize magnetization first")
	}
	if cuda.MaxVecNorm(M.buffer) == 0 {
		log.Fatal("need to initialize magnetization first")
	}
}

// Cleanly exits the simulation, assuring all output is flushed.
func Close() {
	log.Println("shutting down")
	drainOutput()
	Table.flush()
	var memstats runtime.MemStats
	runtime.ReadMemStats(&memstats)
	log.Println("Total memory allocation", memstats.TotalAlloc/(1024), "KiB")

	// debug. TODO: rm
	//	for n, p := range params {
	//		if u, ok := p.(interface {
	//			nUpload() int
	//		}); ok {
	//			log.Println(n, "\t:\t", u.nUpload(), "uploads")
	//		}
	//	}
}

//func sanitycheck() {
//	if Msat() == 0 {
//		log.Fatal("Msat should be nonzero")
//	}
//}
