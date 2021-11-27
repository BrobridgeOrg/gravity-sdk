package main

/*
#include "./error.h"
*/
import "C"
import "unsafe"

func NewError(message string) *C.GravityError {
	e := (*C.GravityError)(C.malloc(C.size_t(unsafe.Sizeof(C.GravityError{}))))
	e.message = C.CString(message)

	return e
}
