// Copyright (c) 2013, Vastech SA (PTY) LTD. All rights reserved.
// http://code.google.com/p/gogoprotobuf/gogoproto
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

/*
	Package custom contains custom types for test and example purposes.
	These types are used by the test structures generated by gogoprotobuf.
*/
package custom

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type Uint128 [2]uint64

func ParseUint128(s string) (u Uint128, err error) {
	ss := strings.Split(s, "-")
	if len(ss) != 2 {
		return u, errors.New("Uint128 Parse Error")
	}
	u[0], err = strconv.ParseUint(ss[0], 10, 64)
	if err != nil {
		return u, err
	}
	u[1], err = strconv.ParseUint(ss[1], 10, 64)
	if err != nil {
		return u, err
	}
	return u, nil
}

//http://code.google.com/p/go/issues/detail?id=4389 waiting to use math/big to make a big number string
func (u Uint128) String() string {
	return fmt.Sprintf("%v-%v", strconv.FormatUint(u[0], 10), strconv.FormatUint(u[1], 10))
}

func (u Uint128) MarshalToString() (*string, error) {
	s := u.String()
	return &s, nil
}

func (u *Uint128) UnmarshalFromString(str *string) error {
	if str == nil {
		u = nil
		return nil
	}
	if len(*str) == 0 {
		u = &Uint128{}
		return nil
	}
	pu, err := ParseUint128(*str)
	if err != nil {
		return err
	}
	*u = pu
	return nil
}

func (u Uint128) MarshalToBytes() ([]byte, error) {
	buffer := make([]byte, 16)
	PutLittleEndianUint128(buffer, 0, u)
	return buffer, nil
}

func GetLittleEndianUint64(b []byte, offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(&b[offset]))
}

func PutLittleEndianUint64(b []byte, offset int, v uint64) {
	b[offset] = byte(v)
	b[offset+1] = byte(v >> 8)
	b[offset+2] = byte(v >> 16)
	b[offset+3] = byte(v >> 24)
	b[offset+4] = byte(v >> 32)
	b[offset+5] = byte(v >> 40)
	b[offset+6] = byte(v >> 48)
	b[offset+7] = byte(v >> 56)
}

func PutLittleEndianUint128(buffer []byte, offset int, v [2]uint64) {
	PutLittleEndianUint64(buffer, offset, v[0])
	PutLittleEndianUint64(buffer, offset+8, v[1])
}

func GetLittleEndianUint128(buffer []byte, offset int) (value [2]uint64) {
	value[0] = GetLittleEndianUint64(buffer, offset)
	value[1] = GetLittleEndianUint64(buffer, offset+8)
	return
}

func GetLittleEndianUint32(b []byte, offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(&b[offset]))
}

func PutLittleEndianUint32(b []byte, offset int, v uint32) {
	b[offset] = byte(v)
	b[offset+1] = byte(v >> 8)
	b[offset+2] = byte(v >> 16)
	b[offset+3] = byte(v >> 24)
}

func GetLittleEndianUint16(b []byte, offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(&b[offset]))
}

func PutLittleEndianUint16(b []byte, offset int, v uint16) {
	b[offset] = byte(v)
	b[offset+1] = byte(v >> 8)
}

func (u *Uint128) UnmarshalFromBytes(data []byte) error {
	if data == nil {
		u = nil
		return nil
	}
	if len(data) == 0 {
		pu := Uint128{}
		*u = pu
		return nil
	}
	pu := Uint128(GetLittleEndianUint128(data, 0))
	*u = pu
	return nil
}

func (this Uint128) Equal(that Uint128) bool {
	return this == that
}

type Uint16 uint16

func ParseUint16(s string) (Uint16, error) {
	u, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return Uint16(u), err
	}
	return Uint16(u), nil
}

func (u Uint16) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (u Uint16) MarshalToString() (*string, error) {
	s := u.String()
	return &s, nil
}

func (u *Uint16) UnmarshalFromString(str *string) error {
	if str == nil {
		u = nil
		return nil
	}
	if len(*str) == 0 {
		d := Uint16(0)
		*u = d
		return nil
	}
	pu, err := ParseUint16(*str)
	if err != nil {
		return err
	}
	*u = pu
	return nil
}

func (u Uint16) MarshalToBytes() ([]byte, error) {
	buffer := make([]byte, 2)
	PutLittleEndianUint16(buffer, 0, uint16(u))
	return buffer, nil
}

func (u *Uint16) UnmarshalFromBytes(data []byte) error {
	if data == nil {
		u = nil
		return nil
	}
	if len(data) == 0 {
		pu := Uint16(0)
		*u = pu
		return nil
	}
	pu := Uint16(GetLittleEndianUint16(data, 0))
	*u = pu
	return nil
}

func (u Uint16) MarshalToUint32() (*uint32, error) {
	u32 := uint32(u)
	return &u32, nil
}

var errTooLargeUint16 = errors.New("Too Large Uint16")

func (u *Uint16) UnmarshalFromUint32(data *uint32) error {
	if data == nil {
		return nil
	}
	if *data > math.MaxUint16 {
		return errTooLargeUint16
	}
	u16 := Uint16(*data)
	*u = u16
	return nil
}

func (u Uint16) Equal(that Uint16) bool {
	return u == that
}

type Uuid []byte

func Parse(str string) (Uuid, error) {
	if len(str) == 38 {
		if str[0] != '{' || str[37] != '}' {
			return nil, syscall.EINVAL
		}
		str = str[1:37]
	}
	if len(str) != 36 {
		return nil, syscall.EINVAL
	}
	uuid := make(Uuid, 16)
	j := 0
	for i, c := range str {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return nil, syscall.EINVAL
			}
			continue
		}
		var v byte
		if c >= '0' && c <= '9' {
			v = byte(c - '0')
		} else if c >= 'a' && c <= 'f' {
			v = 10 + byte(c-'a')
		} else if c >= 'A' && c <= 'F' {
			v = 10 + byte(c-'A')
		} else {
			return nil, syscall.EINVAL
		}
		if j&0x1 == 0 {
			uuid[j>>1] = v << 4
		} else {
			uuid[j>>1] |= v
		}
		j++
	}
	version := uuid.Version()
	if version < 1 || version > 5 {
		return nil, syscall.EINVAL
	}
	return uuid, nil
}

func (uuid Uuid) Version() int {
	if len(uuid) != 16 {
		panic("invalid uuid: not 16 bytes")
	}
	return int(uuid[6] >> 4)
}

func (uuid Uuid) String() string {
	if len(uuid) == 0 {
		return "<empty uuid>"
	}
	if len(uuid) != 16 {
		panic("invalid uuid: not 16 bytes")
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", []byte(uuid[0:4]), []byte(uuid[4:6]), []byte(uuid[6:8]), []byte(uuid[8:10]), []byte(uuid[10:]))
}

func (uuid Uuid) MarshalToString() (*string, error) {
	if uuid == nil {
		return nil, nil
	}
	if len(uuid) == 0 {
		s := ""
		return &s, nil
	}
	id := uuid.String()
	return &id, nil
}

func (uuid *Uuid) UnmarshalFromString(str *string) error {
	if str == nil {
		uuid = nil
		return nil
	}
	if len(*str) == 0 {
		uuid = &Uuid{}
		return nil
	}
	id, err := Parse(*str)
	if err != nil {
		return err
	}
	*uuid = id
	return nil
}

func (uuid Uuid) MarshalToBytes() ([]byte, error) {
	return []byte(uuid), nil
}

func (uuid *Uuid) UnmarshalFromBytes(data []byte) error {
	if data == nil {
		uuid = nil
		return nil
	}
	id := Uuid(data)
	*uuid = id
	return nil
}
