/*
@Author : hrbc
@Time : 2020/3/13 4:34 PM
*/
package geecache

import (
	"geecache/lru"
)

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
func (v ByteView) String() string {
	return string(v.b)
}

var _ lru.Value = (*ByteView)(nil)
