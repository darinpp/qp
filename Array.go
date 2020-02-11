/* ----------------------------------------------------------------------------
 * This file was automatically generated by SWIG (http://www.swig.org).
 * Version 4.0.1
 *
 * This file is not intended to be easily readable and contains a number of
 * coding conventions designed to improve portability and efficiency. Do not make
 * changes to this file unless you know what you are doing--modify the SWIG
 * interface file instead.
 * ----------------------------------------------------------------------------- */

// source: Array.i

package qp

/*
#define intgo swig_intgo
typedef void *swig_voidp;

#include <stdint.h>


typedef long long intgo;
typedef unsigned long long uintgo;



typedef struct { char *p; intgo n; } _gostring_;
typedef struct { void* array; intgo len; intgo cap; } _goslice_;


typedef _goslice_ swig_type_1;
typedef _goslice_ swig_type_2;
extern void _wrap_Swig_free_qp_81cbf099b3eba5bb(uintptr_t arg1);
extern uintptr_t _wrap_Swig_malloc_qp_81cbf099b3eba5bb(swig_intgo arg1);
extern double _wrap_solve_quadprog_qp_81cbf099b3eba5bb(uintptr_t arg1, uintptr_t arg2, uintptr_t arg3, uintptr_t arg4, uintptr_t arg5, uintptr_t arg6, uintptr_t arg7);
extern uintptr_t _wrap_new_matrix_qp_81cbf099b3eba5bb(void);
extern void _wrap_delete_matrix_qp_81cbf099b3eba5bb(uintptr_t arg1);
extern void _wrap_matrix_set_qp_81cbf099b3eba5bb(uintptr_t arg1, swig_type_1 arg2, swig_intgo arg3, swig_intgo arg4);
extern swig_intgo _wrap_matrix_nrows_qp_81cbf099b3eba5bb(uintptr_t arg1);
extern swig_intgo _wrap_matrix_ncols_qp_81cbf099b3eba5bb(uintptr_t arg1);
extern uintptr_t _wrap_new_vector_qp_81cbf099b3eba5bb(void);
extern void _wrap_delete_vector_qp_81cbf099b3eba5bb(uintptr_t arg1);
extern void _wrap_vector_set_qp_81cbf099b3eba5bb(uintptr_t arg1, swig_type_2 arg2, swig_intgo arg3);
extern swig_intgo _wrap_vector_size_qp_81cbf099b3eba5bb(uintptr_t arg1);
extern double _wrap_vector_at_qp_81cbf099b3eba5bb(uintptr_t arg1, swig_intgo arg2);
#undef intgo
*/
import "C"

import "unsafe"
import _ "runtime/cgo"
import "sync"


type _ unsafe.Pointer



var Swig_escape_always_false bool
var Swig_escape_val interface{}


type _swig_fnptr *byte
type _swig_memberptr *byte


type _ sync.Mutex

func Swig_free(arg1 uintptr) {
	_swig_i_0 := arg1
	C._wrap_Swig_free_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0))
}

func Swig_malloc(arg1 int) (_swig_ret uintptr) {
	var swig_r uintptr
	_swig_i_0 := arg1
	swig_r = (uintptr)(C._wrap_Swig_malloc_qp_81cbf099b3eba5bb(C.swig_intgo(_swig_i_0)))
	return swig_r
}

func Solve_quadprog(arg1 Matrix, arg2 Vector, arg3 Matrix, arg4 Vector, arg5 Matrix, arg6 Vector, arg7 Vector) (_swig_ret float64) {
	var swig_r float64
	_swig_i_0 := arg1.Swigcptr()
	_swig_i_1 := arg2.Swigcptr()
	_swig_i_2 := arg3.Swigcptr()
	_swig_i_3 := arg4.Swigcptr()
	_swig_i_4 := arg5.Swigcptr()
	_swig_i_5 := arg6.Swigcptr()
	_swig_i_6 := arg7.Swigcptr()
	swig_r = (float64)(C._wrap_solve_quadprog_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0), C.uintptr_t(_swig_i_1), C.uintptr_t(_swig_i_2), C.uintptr_t(_swig_i_3), C.uintptr_t(_swig_i_4), C.uintptr_t(_swig_i_5), C.uintptr_t(_swig_i_6)))
	return swig_r
}

type SwigcptrMatrix uintptr

func (p SwigcptrMatrix) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrMatrix) SwigIsMatrix() {
}

func NewMatrix() (_swig_ret Matrix) {
	var swig_r Matrix
	swig_r = (Matrix)(SwigcptrMatrix(C._wrap_new_matrix_qp_81cbf099b3eba5bb()))
	return swig_r
}

func DeleteMatrix(arg1 Matrix) {
	_swig_i_0 := arg1.Swigcptr()
	C._wrap_delete_matrix_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0))
}

func (arg1 SwigcptrMatrix) Set(arg2 []float64, arg3 uint, arg4 uint) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	_swig_i_2 := arg3
	_swig_i_3 := arg4
	C._wrap_matrix_set_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0), *(*C.swig_type_1)(unsafe.Pointer(&_swig_i_1)), C.swig_intgo(_swig_i_2), C.swig_intgo(_swig_i_3))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrMatrix) Nrows() (_swig_ret uint) {
	var swig_r uint
	_swig_i_0 := arg1
	swig_r = (uint)(C._wrap_matrix_nrows_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0)))
	return swig_r
}

func (arg1 SwigcptrMatrix) Ncols() (_swig_ret uint) {
	var swig_r uint
	_swig_i_0 := arg1
	swig_r = (uint)(C._wrap_matrix_ncols_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0)))
	return swig_r
}

type Matrix interface {
	Swigcptr() uintptr
	SwigIsMatrix()
	Set(arg2 []float64, arg3 uint, arg4 uint)
	Nrows() (_swig_ret uint)
	Ncols() (_swig_ret uint)
}

type SwigcptrVector uintptr

func (p SwigcptrVector) Swigcptr() uintptr {
	return (uintptr)(p)
}

func (p SwigcptrVector) SwigIsVector() {
}

func NewVector() (_swig_ret Vector) {
	var swig_r Vector
	swig_r = (Vector)(SwigcptrVector(C._wrap_new_vector_qp_81cbf099b3eba5bb()))
	return swig_r
}

func DeleteVector(arg1 Vector) {
	_swig_i_0 := arg1.Swigcptr()
	C._wrap_delete_vector_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0))
}

func (arg1 SwigcptrVector) Set(arg2 []float64, arg3 uint) {
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	_swig_i_2 := arg3
	C._wrap_vector_set_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0), *(*C.swig_type_2)(unsafe.Pointer(&_swig_i_1)), C.swig_intgo(_swig_i_2))
	if Swig_escape_always_false {
		Swig_escape_val = arg2
	}
}

func (arg1 SwigcptrVector) Size() (_swig_ret uint) {
	var swig_r uint
	_swig_i_0 := arg1
	swig_r = (uint)(C._wrap_vector_size_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0)))
	return swig_r
}

func (arg1 SwigcptrVector) At(arg2 uint) (_swig_ret float64) {
	var swig_r float64
	_swig_i_0 := arg1
	_swig_i_1 := arg2
	swig_r = (float64)(C._wrap_vector_at_qp_81cbf099b3eba5bb(C.uintptr_t(_swig_i_0), C.swig_intgo(_swig_i_1)))
	return swig_r
}

type Vector interface {
	Swigcptr() uintptr
	SwigIsVector()
	Set(arg2 []float64, arg3 uint)
	Size() (_swig_ret uint)
	At(arg2 uint) (_swig_ret float64)
}

