package main

import (
	. "nimble-cube/core"
	"nimble-cube/gpu/conv"
	"nimble-cube/mag"
	"os"
	"time"
	"fmt"
)

func main() {
	N0, N1, N2 := IntArg(0), IntArg(1), IntArg(2)
	cx, cy, cz := 3e-9, 3.125e-9, 3.125e-9
	mesh := NewMesh(N0, N1, N2, cx, cy, cz)
	size := mesh.GridSize()
	Log("mesh:", mesh)
	Log("block:", BlockSize(mesh.GridSize()))

	m := MakeChan3(size, "m")
	hd := MakeChan3(size, "Hd")

	acc := 1
	kernel := mag.BruteKernel(mesh.ZeroPadded(), acc)
	Stack(conv.NewSymmetricHtoD(size, kernel, m.MakeRChan3(), hd))

	Msat := 1.0053
	aex := Mu0 * 13e-12 / Msat
	hex := MakeChan3(size, "Hex")
	Stack(mag.NewExchange2D(m.MakeRChan3(), hex, mesh, aex))

	heff := MakeChan3(size, "Heff")
	Stack(NewAdder3(heff, hd.MakeRChan3(), hex.MakeRChan3()))

	const alpha = 1
	torque := MakeChan3(size, "τ")
	Stack(mag.NewLLGTorque(torque, m.MakeRChan3(), heff.MakeRChan3(), alpha))

	const dt = 100e-15
	solver := mag.NewEuler(m, torque.MakeRChan3(), dt)
	mag.SetAll(m.UnsafeArray(), mag.Uniform(0, 0.1, 1))

	RunStack()

	start := time.Now()
	duration := time.Since(start)
	for duration < 10*time.Second{
		solver.Steps(100)
		duration = time.Since(start)
	}

	fmt.Println(N0, N1, N2, *Flag_maxblocklen, duration.Nanoseconds()/1e6)
	
	ProfDump(os.Stdout)
	Cleanup()
}
