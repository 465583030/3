#include "stencil.h"
#include "float3.h"

// Exchange + Dzyaloshinskii-Moriya interaction according to
// Bagdanov and Röβler, PRL 87, 3, 2001. eq.8 (out-of-plane symmetry breaking).
// Taking into account proper boundary conditions.
// m: normalized magnetization
// H: effective field in Tesla
// D: dmi strength / Msat, in Tesla*m
// A: Aex/Msat
extern "C" __global__ void
adddmi(float* __restrict__ Hx, float* __restrict__ Hy, float* __restrict__ Hz,
       float* __restrict__ mx, float* __restrict__ my, float* __restrict__ mz,
       float cx, float cy, float cz, float DL, float DH, float A, int N0, int N1, int N2) {

    int i = blockIdx.z * blockDim.z + threadIdx.z;
    int j = blockIdx.y * blockDim.y + threadIdx.y;
    int k = blockIdx.x * blockDim.x + threadIdx.x;

    if (i >= N0 || j >= N1 || k >= N2) {
        return;
    }

    int I = idx(i, j, k);                        // central cell index
    float3 h = make_float3(Hx[I], Hy[I], Hz[I]); // add to H
    float3 m = make_float3(mx[I], my[I], mz[I]); // central m
    float DL_2A = (DL/(2.0f*A));
    float DH_2A = (DH/(2.0f*A));

    // z derivatives (along length)
    {
        int I1 = idx(i, j, hclamp(k+1, N2));  // right index, clamped
        int I2 = idx(i, j, lclamp(k-1));      // left index, clamped

        // DMI
        float mz1 = (k+1<N2)? mz[I1] : (m.z + (cz * DL_2A * m.x)); // right neighbor
        float mz2 = (k-1>=0)? mz[I2] : (m.z - (cz * DL_2A * m.x)); // left neighbor
        h.x -= DL*(mz1-mz2)/cz;
        // note: actually 2*D * delta / (2*c)

        float mx1 = (k+1<N2)? mx[I1] : (m.x - (cz * DL_2A * m.z));
        float mx2 = (k-1>=0)? mx[I2] : (m.x + (cz * DL_2A * m.z));
        h.z += DL*(mx1-mx2)/cz;

        // Exchange
        float3 m1 = make_float3(mx1, my[I1], mz1); // right neighbor
        float3 m2 = make_float3(mx2, my[I2], mz2); // left neighbor
        h +=  (2.0f*A/(cz*cz)) * ((m1 - m) + (m2 - m));
    }

    // y derivatives (along height)
    {
        int I1 = idx(i, hclamp(j+1, N1), k);
        int I2 = idx(i, lclamp(j-1), k);

        // DMI
        float my1 = (j+1<N1)? my[I1] : (m.y + (cy * DH_2A * m.x));
        float my2 = (j-1>=0)? my[I2] : (m.y - (cy * DH_2A * m.x));
        h.x -= DH*(my1-my2)/cy;

        float mx1 = (j+1<N1)? mx[I1] : (m.x - (cy * DH_2A * m.y));
        float mx2 = (j-1>=0)? mx[I2] : (m.x + (cy * DH_2A * m.y));
        h.y += DH*(mx1-mx2)/cy;

        // Exchange
        float3 m1 = make_float3(mx1, my1, mz[I1]);
        float3 m2 = make_float3(mx2, my2, mz[I2]);
        h +=  (2.0f*A/(cy*cy)) * ((m1 - m) + (m2 - m));
    }

    // write back, result is H + Hdmi + Hex
    Hx[I] = h.x;
    Hy[I] = h.y;
    Hz[I] = h.z;
}

// Note on boundary conditions.
//
// We need the derivative and laplacian of m in point A, but e.g. C lies out of the boundaries.
// We use the boundary condition in B (derivative of the magnetization) to extrapolate m to point C:
// 	m_C = m_A + (dm/dx)|_B * cellsize
//
// When point C is inside the boundary, we just use its actual value.
//
// Then we can take the central derivative in A:
// 	(dm/dx)|_A = (m_C - m_D) / (2*cellsize)
// And the laplacian:
// 	lapl(m)|_A = (m_C + m_D - 2*m_A) / (cellsize^2)
//
// All these operations should be second order as they involve only central derivatives.
//
//    ------------------------------------------------------------------ *
//   |                                                   |             C |
//   |                                                   |          **   |
//   |                                                   |        ***    |
//   |                                                   |     ***       |
//   |                                                   |   ***         |
//   |                                                   | ***           |
//   |                                                   B               |
//   |                                               *** |               |
//   |                                            ***    |               |
//   |                                         ****      |               |
//   |                                     ****          |               |
//   |                                  ****             |               |
//   |                              ** A                 |               |
//   |                         *****                     |               |
//   |                   ******                          |               |
//   |          *********                                |               |
//   |D ********                                         |               |
//   |                                                   |               |
//   +----------------+----------------+-----------------+---------------+
//  -1              -0.5               0               0.5               1
//                                 x
