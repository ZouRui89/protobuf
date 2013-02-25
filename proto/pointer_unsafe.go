// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2013 Vastech SA (PTY) LTD. All rights reserved.
//
// Copyright 2012 The Go Authors.  All rights reserved.
// http://code.google.com/p/goprotobuf/
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// +build !appengine

// This file contains the implementation of the proto field accesses using package unsafe.

package proto

import (
	"reflect"
	"unsafe"
)

// NOTE: These type_Foo functions would more idiomatically be methods,
// but Go does not allow methods on pointer types, and we must preserve
// some pointer type for the garbage collector. We use these
// funcs with clunky names as our poor approximation to methods.
//
// An alternative would be
//	type structPointer struct { p unsafe.Pointer }
// but that does not registerize as well.

// A structPointer is a pointer to a struct.
type structPointer unsafe.Pointer

// toStructPointer returns a structPointer equivalent to the given reflect value.
func toStructPointer(v reflect.Value) structPointer {
	return structPointer(unsafe.Pointer(v.Pointer()))
}

// IsNil reports whether p is nil.
func structPointer_IsNil(p structPointer) bool {
	return p == nil
}

// Interface returns the struct pointer, assumed to have element type t,
// as an interface value.
func structPointer_Interface(p structPointer, t reflect.Type) interface{} {
	return reflect.NewAt(t, unsafe.Pointer(p)).Interface()
}

func structPointer_InterfaceAt(p structPointer, f field, t reflect.Type) interface{} {
	point := unsafe.Pointer(uintptr(p) + uintptr(f))
	r := reflect.NewAt(t, point)
	return r.Interface()
}

func structPointer_InterfaceRef(p structPointer, f field, t reflect.Type) interface{} {
	point := unsafe.Pointer(uintptr(p) + uintptr(f))
	r := reflect.NewAt(t, point)
	if r.Elem().IsNil() {
		return nil
	}
	return r.Elem().Interface()
}

func copyUintPtr(oldptr, newptr uintptr, size int) {
	for j := 0; j < size; j++ {
		oldb := (*byte)(unsafe.Pointer(oldptr + uintptr(j)))
		*(*byte)(unsafe.Pointer(newptr + uintptr(j))) = *oldb
	}
}

func structPointer_Copy(oldptr structPointer, newptr structPointer, size int) {
	copyUintPtr(uintptr(oldptr), uintptr(newptr), size)
}

func appendStructPointer(base structPointer, f field, typ reflect.Type) structPointer {
	size := typ.Elem().Size()
	oldHeader := structPointer_GetSliceHeader(base, f)
	newLen := oldHeader.Len + 1
	slice := reflect.MakeSlice(typ, newLen, newLen)
	bas := toStructPointer(slice)
	for i := 0; i < oldHeader.Len; i++ {
		newElemptr := uintptr(bas) + uintptr(i)*size
		oldElemptr := oldHeader.Data + uintptr(i)*size
		copyUintPtr(oldElemptr, newElemptr, int(size))
	}

	oldHeader.Data = uintptr(bas)
	oldHeader.Len = newLen
	oldHeader.Cap = newLen

	return structPointer(unsafe.Pointer(uintptr(unsafe.Pointer(bas)) + uintptr(uintptr(newLen-1)*size)))
}

// A field identifies a field in a struct, accessible from a structPointer.
// In this implementation, a field is identified by its byte offset from the start of the struct.
type field uintptr

// toField returns a field equivalent to the given reflect field.
func toField(f *reflect.StructField) field {
	return field(f.Offset)
}

// invalidField is an invalid field identifier.
const invalidField = ^field(0)

// IsValid reports whether the field identifier is valid.
func (f field) IsValid() bool {
	return f != ^field(0)
}

// Bytes returns the address of a []byte field in the struct.
func structPointer_Bytes(p structPointer, f field) *[]byte {
	return (*[]byte)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// BytesSlice returns the address of a [][]byte field in the struct.
func structPointer_BytesSlice(p structPointer, f field) *[][]byte {
	return (*[][]byte)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// Bool returns the address of a *bool field in the struct.
func structPointer_Bool(p structPointer, f field) **bool {
	return (**bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// RefBool returns a *bool field in the struct.
func structPointer_RefBool(p structPointer, f field) *bool {
	return (*bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// BoolSlice returns the address of a []bool field in the struct.
func structPointer_BoolSlice(p structPointer, f field) *[]bool {
	return (*[]bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// String returns the address of a *string field in the struct.
func structPointer_String(p structPointer, f field) **string {
	return (**string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// RefString returns the address of a string field in the struct.
func structPointer_RefString(p structPointer, f field) *string {
	return (*string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// StringSlice returns the address of a []string field in the struct.
func structPointer_StringSlice(p structPointer, f field) *[]string {
	return (*[]string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// ExtMap returns the address of an extension map field in the struct.
func structPointer_ExtMap(p structPointer, f field) *map[int32]Extension {
	return (*map[int32]Extension)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// SetStructPointer writes a *struct field in the struct.
func structPointer_SetStructPointer(p structPointer, f field, q structPointer) {
	*(*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f))) = q
}

func structPointer_FieldPointer(p structPointer, f field) structPointer {
	return structPointer(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// GetStructPointer reads a *struct field in the struct.
func structPointer_GetStructPointer(p structPointer, f field) structPointer {
	return *(*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

func structPointer_GetRefStructPointer(p structPointer, f field) structPointer {
	return structPointer((*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

func structPointer_GetSliceHeader(p structPointer, f field) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

func structPointer_Add(p structPointer, size field) structPointer {
	return structPointer(unsafe.Pointer(uintptr(p) + uintptr(size)))
}

func structPointer_Len(p structPointer, f field) int {
	return len(*(*[]interface{})(unsafe.Pointer(structPointer_GetRefStructPointer(p, f))))
}

// StructPointerSlice the address of a []*struct field in the struct.
func structPointer_StructPointerSlice(p structPointer, f field) *structPointerSlice {
	return (*structPointerSlice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// A structPointerSlice represents a slice of pointers to structs (themselves submessages or groups).
type structPointerSlice []structPointer

func (v *structPointerSlice) Len() int                  { return len(*v) }
func (v *structPointerSlice) Index(i int) structPointer { return (*v)[i] }
func (v *structPointerSlice) Append(p structPointer)    { *v = append(*v, p) }

// A word32 is the address of a "pointer to 32-bit value" field.
type word32 **uint32

// IsNil reports whether *v is nil.
func word32_IsNil(p word32) bool {
	return *p == nil
}

// Set sets *v to point at a newly allocated word set to x.
func word32_Set(p word32, o *Buffer, x uint32) {
	if len(o.uint32s) == 0 {
		o.uint32s = make([]uint32, uint32PoolSize)
	}
	o.uint32s[0] = x
	*p = &o.uint32s[0]
	o.uint32s = o.uint32s[1:]
}

// Get gets the value pointed at by *v.
func word32_Get(p word32) uint32 {
	return **p
}

// Word32 returns the address of a *int32, *uint32, *float32, or *enum field in the struct.
func structPointer_Word32(p structPointer, f field) word32 {
	return word32((**uint32)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

// refWord32 is the address of a 32-bit value field.
type refWord32 *uint32

func refWord32_IsNil(p refWord32) bool {
	return p == nil
}

func refWord32_Set(p refWord32, o *Buffer, x uint32) {
	if len(o.uint32s) == 0 {
		o.uint32s = make([]uint32, uint32PoolSize)
	}
	o.uint32s[0] = x
	*p = o.uint32s[0]
	o.uint32s = o.uint32s[1:]
}

func refWord32_Get(p refWord32) uint32 {
	return *p
}

func structPointer_RefWord32(p structPointer, f field) refWord32 {
	return refWord32((*uint32)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

// A word32Slice is a slice of 32-bit values.
type word32Slice []uint32

func (v *word32Slice) Append(x uint32)    { *v = append(*v, x) }
func (v *word32Slice) Len() int           { return len(*v) }
func (v *word32Slice) Index(i int) uint32 { return (*v)[i] }

// Word32Slice returns the address of a []int32, []uint32, []float32, or []enum field in the struct.
func structPointer_Word32Slice(p structPointer, f field) *word32Slice {
	return (*word32Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// word64 is like word32 but for 64-bit values.
type word64 **uint64

func word64_Set(p word64, o *Buffer, x uint64) {
	if len(o.uint64s) == 0 {
		o.uint64s = make([]uint64, uint64PoolSize)
	}
	o.uint64s[0] = x
	*p = &o.uint64s[0]
	o.uint64s = o.uint64s[1:]
}

func word64_IsNil(p word64) bool {
	return *p == nil
}

func word64_Get(p word64) uint64 {
	return **p
}

func structPointer_Word64(p structPointer, f field) word64 {
	return word64((**uint64)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

// refWord64 is like refWord32 but for 32-bit values.
type refWord64 *uint64

func refWord64_Set(p refWord64, o *Buffer, x uint64) {
	if len(o.uint64s) == 0 {
		o.uint64s = make([]uint64, uint64PoolSize)
	}
	o.uint64s[0] = x
	*p = o.uint64s[0]
	o.uint64s = o.uint64s[1:]
}

func refWord64_IsNil(p refWord64) bool {
	return p == nil
}

func refWord64_Get(p refWord64) uint64 {
	return *p
}

func structPointer_RefWord64(p structPointer, f field) refWord64 {
	return refWord64((*uint64)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

// word64Slice is like word32Slice but for 64-bit values.
type word64Slice []uint64

func (v *word64Slice) Append(x uint64)    { *v = append(*v, x) }
func (v *word64Slice) Len() int           { return len(*v) }
func (v *word64Slice) Index(i int) uint64 { return (*v)[i] }

func structPointer_Word64Slice(p structPointer, f field) *word64Slice {
	return (*word64Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}

// A word16 is the address of a "pointer to 16-bit value" field.
type word16 **uint16

func word16_IsNil(p word16) bool {
	return *p == nil
}

func word16_Set(p word16, o *Buffer, x uint16) {
	if len(o.uint16s) == 0 {
		o.uint16s = make([]uint16, uint16PoolSize)
	}
	o.uint16s[0] = x
	*p = &o.uint16s[0]
	o.uint16s = o.uint16s[1:]
}

func word16_Get(p word16) uint16 {
	return **p
}

// Word16 returns the address of a *int16, *uint16, or *enum field in the struct.
func structPointer_Word16(p structPointer, f field) word16 {
	return word16((**uint16)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

type refWord16 *uint16

func refWord16_IsNil(p refWord16) bool {
	return p == nil
}

func refWord16_Set(p refWord16, o *Buffer, x uint16) {
	if len(o.uint16s) == 0 {
		o.uint16s = make([]uint16, uint16PoolSize)
	}
	o.uint16s[0] = x
	*p = o.uint16s[0]
	o.uint16s = o.uint16s[1:]
}

func refWord16_Get(p refWord16) uint16 {
	return *p
}

// RefWord16 returns the address of a int16, uint16 of enum field in the struct.
func structPointer_RefWord16(p structPointer, f field) refWord16 {
	return refWord16((*uint16)(unsafe.Pointer(uintptr(p) + uintptr(f))))
}

// A word16Slice is a slice of 16-bit values.
type word16Slice []uint16

func (v *word16Slice) Append(x uint16)    { *v = append(*v, x) }
func (v *word16Slice) Len() int           { return len(*v) }
func (v *word16Slice) Index(i int) uint16 { return (*v)[i] }

// Word16Slice returns the address of a []int16, []uint16 or []enum field in the struct.
func structPointer_Word16Slice(p structPointer, f field) *word16Slice {
	return (*word16Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
}
