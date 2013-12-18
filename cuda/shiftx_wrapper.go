package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"unsafe"
)

var shiftx_code cu.Function

type shiftx_args struct {
	arg_dst    unsafe.Pointer
	arg_src    unsafe.Pointer
	arg_Nx     int
	arg_Ny     int
	arg_Nz     int
	arg_shx    int
	arg_clampL float32
	arg_clampR float32
	argptr     [8]unsafe.Pointer
}

// Wrapper for shiftx CUDA kernel, asynchronous.
func k_shiftx_async(dst unsafe.Pointer, src unsafe.Pointer, Nx int, Ny int, Nz int, shx int, clampL float32, clampR float32, cfg *config) {
	if Synchronous { // debug
		Sync()
	}

	if shiftx_code == 0 {
		shiftx_code = fatbinLoad(shiftx_map, "shiftx")
	}

	var _a_ shiftx_args

	_a_.arg_dst = dst
	_a_.argptr[0] = unsafe.Pointer(&_a_.arg_dst)
	_a_.arg_src = src
	_a_.argptr[1] = unsafe.Pointer(&_a_.arg_src)
	_a_.arg_Nx = Nx
	_a_.argptr[2] = unsafe.Pointer(&_a_.arg_Nx)
	_a_.arg_Ny = Ny
	_a_.argptr[3] = unsafe.Pointer(&_a_.arg_Ny)
	_a_.arg_Nz = Nz
	_a_.argptr[4] = unsafe.Pointer(&_a_.arg_Nz)
	_a_.arg_shx = shx
	_a_.argptr[5] = unsafe.Pointer(&_a_.arg_shx)
	_a_.arg_clampL = clampL
	_a_.argptr[6] = unsafe.Pointer(&_a_.arg_clampL)
	_a_.arg_clampR = clampR
	_a_.argptr[7] = unsafe.Pointer(&_a_.arg_clampR)

	args := _a_.argptr[:]
	cu.LaunchKernel(shiftx_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
	}
}

var shiftx_map = map[int]string{0: "",
	20: shiftx_ptx_20,
	30: shiftx_ptx_30,
	35: shiftx_ptx_35}

const (
	shiftx_ptx_20 = `
.version 3.2
.target sm_20
.address_size 64


.visible .entry shiftx(
	.param .u64 shiftx_param_0,
	.param .u64 shiftx_param_1,
	.param .u32 shiftx_param_2,
	.param .u32 shiftx_param_3,
	.param .u32 shiftx_param_4,
	.param .u32 shiftx_param_5,
	.param .f32 shiftx_param_6,
	.param .f32 shiftx_param_7
)
{
	.reg .pred 	%p<8>;
	.reg .s32 	%r<22>;
	.reg .f32 	%f<6>;
	.reg .s64 	%rd<9>;


	ld.param.u64 	%rd3, [shiftx_param_0];
	ld.param.u64 	%rd4, [shiftx_param_1];
	ld.param.u32 	%r5, [shiftx_param_2];
	ld.param.u32 	%r6, [shiftx_param_3];
	ld.param.u32 	%r8, [shiftx_param_4];
	ld.param.u32 	%r7, [shiftx_param_5];
	ld.param.f32 	%f3, [shiftx_param_6];
	ld.param.f32 	%f4, [shiftx_param_7];
	cvta.to.global.u64 	%rd1, %rd3;
	cvta.to.global.u64 	%rd2, %rd4;
	.loc 1 9 1
	mov.u32 	%r9, %ntid.x;
	mov.u32 	%r10, %ctaid.x;
	mov.u32 	%r11, %tid.x;
	mad.lo.s32 	%r1, %r9, %r10, %r11;
	.loc 1 10 1
	mov.u32 	%r12, %ntid.y;
	mov.u32 	%r13, %ctaid.y;
	mov.u32 	%r14, %tid.y;
	mad.lo.s32 	%r2, %r12, %r13, %r14;
	.loc 1 11 1
	mov.u32 	%r15, %ntid.z;
	mov.u32 	%r16, %ctaid.z;
	mov.u32 	%r17, %tid.z;
	mad.lo.s32 	%r3, %r15, %r16, %r17;
	.loc 1 13 1
	setp.lt.s32	%p1, %r1, %r5;
	setp.lt.s32	%p2, %r2, %r6;
	and.pred  	%p3, %p1, %p2;
	.loc 1 13 1
	setp.lt.s32	%p4, %r3, %r8;
	and.pred  	%p5, %p3, %p4;
	.loc 1 13 1
	@!%p5 bra 	BB0_5;
	bra.uni 	BB0_1;

BB0_1:
	.loc 1 14 1
	sub.s32 	%r4, %r1, %r7;
	.loc 1 16 1
	setp.lt.s32	%p6, %r4, 0;
	mov.f32 	%f5, %f3;
	@%p6 bra 	BB0_4;

	.loc 1 18 1
	setp.ge.s32	%p7, %r4, %r5;
	mov.f32 	%f5, %f4;
	@%p7 bra 	BB0_4;

	.loc 1 21 1
	mad.lo.s32 	%r18, %r3, %r6, %r2;
	mad.lo.s32 	%r19, %r18, %r5, %r4;
	mul.wide.s32 	%rd5, %r19, 4;
	add.s64 	%rd6, %rd2, %rd5;
	.loc 1 21 1
	ld.global.f32 	%f5, [%rd6];

BB0_4:
	.loc 1 23 1
	mad.lo.s32 	%r20, %r3, %r6, %r2;
	mad.lo.s32 	%r21, %r20, %r5, %r1;
	mul.wide.s32 	%rd7, %r21, 4;
	add.s64 	%rd8, %rd1, %rd7;
	.loc 1 23 1
	st.global.f32 	[%rd8], %f5;

BB0_5:
	.loc 1 25 2
	ret;
}


`
	shiftx_ptx_30 = `
.version 3.2
.target sm_30
.address_size 64


.visible .entry shiftx(
	.param .u64 shiftx_param_0,
	.param .u64 shiftx_param_1,
	.param .u32 shiftx_param_2,
	.param .u32 shiftx_param_3,
	.param .u32 shiftx_param_4,
	.param .u32 shiftx_param_5,
	.param .f32 shiftx_param_6,
	.param .f32 shiftx_param_7
)
{
	.reg .pred 	%p<8>;
	.reg .s32 	%r<22>;
	.reg .f32 	%f<6>;
	.reg .s64 	%rd<9>;


	ld.param.u64 	%rd3, [shiftx_param_0];
	ld.param.u64 	%rd4, [shiftx_param_1];
	ld.param.u32 	%r5, [shiftx_param_2];
	ld.param.u32 	%r6, [shiftx_param_3];
	ld.param.u32 	%r8, [shiftx_param_4];
	ld.param.u32 	%r7, [shiftx_param_5];
	ld.param.f32 	%f3, [shiftx_param_6];
	ld.param.f32 	%f4, [shiftx_param_7];
	cvta.to.global.u64 	%rd1, %rd3;
	cvta.to.global.u64 	%rd2, %rd4;
	.loc 1 9 1
	mov.u32 	%r9, %ntid.x;
	mov.u32 	%r10, %ctaid.x;
	mov.u32 	%r11, %tid.x;
	mad.lo.s32 	%r1, %r9, %r10, %r11;
	.loc 1 10 1
	mov.u32 	%r12, %ntid.y;
	mov.u32 	%r13, %ctaid.y;
	mov.u32 	%r14, %tid.y;
	mad.lo.s32 	%r2, %r12, %r13, %r14;
	.loc 1 11 1
	mov.u32 	%r15, %ntid.z;
	mov.u32 	%r16, %ctaid.z;
	mov.u32 	%r17, %tid.z;
	mad.lo.s32 	%r3, %r15, %r16, %r17;
	.loc 1 13 1
	setp.lt.s32	%p1, %r1, %r5;
	setp.lt.s32	%p2, %r2, %r6;
	and.pred  	%p3, %p1, %p2;
	.loc 1 13 1
	setp.lt.s32	%p4, %r3, %r8;
	and.pred  	%p5, %p3, %p4;
	.loc 1 13 1
	@!%p5 bra 	BB0_5;
	bra.uni 	BB0_1;

BB0_1:
	.loc 1 14 1
	sub.s32 	%r4, %r1, %r7;
	.loc 1 16 1
	setp.lt.s32	%p6, %r4, 0;
	mov.f32 	%f5, %f3;
	@%p6 bra 	BB0_4;

	.loc 1 18 1
	setp.ge.s32	%p7, %r4, %r5;
	mov.f32 	%f5, %f4;
	@%p7 bra 	BB0_4;

	.loc 1 21 1
	mad.lo.s32 	%r18, %r3, %r6, %r2;
	mad.lo.s32 	%r19, %r18, %r5, %r4;
	mul.wide.s32 	%rd5, %r19, 4;
	add.s64 	%rd6, %rd2, %rd5;
	.loc 1 21 1
	ld.global.f32 	%f5, [%rd6];

BB0_4:
	.loc 1 23 1
	mad.lo.s32 	%r20, %r3, %r6, %r2;
	mad.lo.s32 	%r21, %r20, %r5, %r1;
	mul.wide.s32 	%rd7, %r21, 4;
	add.s64 	%rd8, %rd1, %rd7;
	.loc 1 23 1
	st.global.f32 	[%rd8], %f5;

BB0_5:
	.loc 1 25 2
	ret;
}


`
	shiftx_ptx_35 = `
.version 3.2
.target sm_35
.address_size 64


.weak .func  (.param .b32 func_retval0) cudaMalloc(
	.param .b64 cudaMalloc_param_0,
	.param .b64 cudaMalloc_param_1
)
{
	.reg .s32 	%r<2>;


	mov.u32 	%r1, 30;
	st.param.b32	[func_retval0+0], %r1;
	.loc 2 66 3
	ret;
}

.weak .func  (.param .b32 func_retval0) cudaFuncGetAttributes(
	.param .b64 cudaFuncGetAttributes_param_0,
	.param .b64 cudaFuncGetAttributes_param_1
)
{
	.reg .s32 	%r<2>;


	mov.u32 	%r1, 30;
	st.param.b32	[func_retval0+0], %r1;
	.loc 2 71 3
	ret;
}

.visible .entry shiftx(
	.param .u64 shiftx_param_0,
	.param .u64 shiftx_param_1,
	.param .u32 shiftx_param_2,
	.param .u32 shiftx_param_3,
	.param .u32 shiftx_param_4,
	.param .u32 shiftx_param_5,
	.param .f32 shiftx_param_6,
	.param .f32 shiftx_param_7
)
{
	.reg .pred 	%p<8>;
	.reg .s32 	%r<22>;
	.reg .f32 	%f<6>;
	.reg .s64 	%rd<9>;


	ld.param.u64 	%rd3, [shiftx_param_0];
	ld.param.u64 	%rd4, [shiftx_param_1];
	ld.param.u32 	%r5, [shiftx_param_2];
	ld.param.u32 	%r6, [shiftx_param_3];
	ld.param.u32 	%r8, [shiftx_param_4];
	ld.param.u32 	%r7, [shiftx_param_5];
	ld.param.f32 	%f3, [shiftx_param_6];
	ld.param.f32 	%f4, [shiftx_param_7];
	cvta.to.global.u64 	%rd1, %rd3;
	cvta.to.global.u64 	%rd2, %rd4;
	.loc 1 9 1
	mov.u32 	%r9, %ntid.x;
	mov.u32 	%r10, %ctaid.x;
	mov.u32 	%r11, %tid.x;
	mad.lo.s32 	%r1, %r9, %r10, %r11;
	.loc 1 10 1
	mov.u32 	%r12, %ntid.y;
	mov.u32 	%r13, %ctaid.y;
	mov.u32 	%r14, %tid.y;
	mad.lo.s32 	%r2, %r12, %r13, %r14;
	.loc 1 11 1
	mov.u32 	%r15, %ntid.z;
	mov.u32 	%r16, %ctaid.z;
	mov.u32 	%r17, %tid.z;
	mad.lo.s32 	%r3, %r15, %r16, %r17;
	.loc 1 13 1
	setp.lt.s32	%p1, %r1, %r5;
	setp.lt.s32	%p2, %r2, %r6;
	and.pred  	%p3, %p1, %p2;
	.loc 1 13 1
	setp.lt.s32	%p4, %r3, %r8;
	and.pred  	%p5, %p3, %p4;
	.loc 1 13 1
	@!%p5 bra 	BB2_5;
	bra.uni 	BB2_1;

BB2_1:
	.loc 1 14 1
	sub.s32 	%r4, %r1, %r7;
	.loc 1 16 1
	setp.lt.s32	%p6, %r4, 0;
	mov.f32 	%f5, %f3;
	@%p6 bra 	BB2_4;

	.loc 1 18 1
	setp.ge.s32	%p7, %r4, %r5;
	mov.f32 	%f5, %f4;
	@%p7 bra 	BB2_4;

	.loc 1 21 1
	mad.lo.s32 	%r18, %r3, %r6, %r2;
	mad.lo.s32 	%r19, %r18, %r5, %r4;
	mul.wide.s32 	%rd5, %r19, 4;
	add.s64 	%rd6, %rd2, %rd5;
	.loc 1 21 1
	ld.global.nc.f32 	%f5, [%rd6];

BB2_4:
	.loc 1 23 1
	mad.lo.s32 	%r20, %r3, %r6, %r2;
	mad.lo.s32 	%r21, %r20, %r5, %r1;
	mul.wide.s32 	%rd7, %r21, 4;
	add.s64 	%rd8, %rd1, %rd7;
	.loc 1 23 1
	st.global.f32 	[%rd8], %f5;

BB2_5:
	.loc 1 25 2
	ret;
}


`
)
