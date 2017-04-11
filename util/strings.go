package util

import (
	"reflect"
	"strings"
	"unsafe"
)

func Extract2(s string, delim string) (string, string) {
	a := strings.Split(s, delim)
	if len(a) == 1 {
		return s, ""
	} else {
		return a[0], strings.Join(a[1:], delim)
	}
}

func ExtractFirst(s string, delim string) string {
	i := strings.Index(s, delim)
	if i >= 0 {
		return s[0:i]
	} else {
		return s
	}
}

func ExtractLast(s string, delim string) string {
	i := strings.LastIndex(s, delim)
	if i >= 0 {
		return s[i+1:]
	} else {
		return s
	}
}

func CToGoString(c []byte) string {
	for i, b := range c {
		if b == 0 {
			return string(c[:i+1])
		}
	}
	return string(c[:])
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
