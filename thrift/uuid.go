package thrift

import (
	"encoding/hex"
	"fmt"
)

type TUUID [16]byte

// String 生成规范形式的字符串
// 该字符串适合用 TJSONProtocol 写入.
func (u TUUID) String() string {
	var buf [36]byte
	hex.Encode(buf[0:], u[:4])
	buf[8] = '-'
	hex.Encode(buf[9:], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])
	return string(buf[:])
}

func hexToDec(b byte) (byte, bool) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', true
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, true
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, true
	default:
		return 0, false
	}
}

func hexToByte(b1, b2 byte) (b byte, ok bool) {
	b1, ok = hexToDec(b1)
	if !ok {
		return 0, ok
	}
	b2, ok = hexToDec(b2)
	if !ok {
		return 0, ok
	}
	return b1<<4 + b2, true
}

// ParseTuuid parses a canonical form UUID string into TUUID.
//
// Note that this function only supports case insensitive canonical form
// (8-4-4-4-12/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx), and rejects any other forms.
// For a more flexible UUID string parser, please use third party UUID libraries.
//
// This function is suitable for reading with TJSONProtocol.
func ParseTuuid(s string) (u TUUID, err error) {
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return u, fmt.Errorf("malformed TUUID string: %q", s)
	}
	var ok bool
	for i, j := range []int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	} {
		u[i], ok = hexToByte(s[j], s[j+1])
		if !ok {
			return u, fmt.Errorf("malformed Tuuid string: %q", s)
		}
	}
	return u, nil
}

// Must is a sugar to be used in places that error handling is impossible (for
// example, global variable declarations) and also errors are not in general
// expected.
//
// This is an example to use Must with ParseTuuid to declare a global special
// uuid:
//
//	var NameSpaceDNSUUID = thrift.Must(thrift.ParseTuuid("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
