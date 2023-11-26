// File: "buffer.go"
// Code based on https://github.com/lmittmann/tint/blob/main/buffer.go

package xlog

import "sync"

const BUFFER_DEFAULT_CAP = 1 << 10 // 1K

const BUFFER_MAX_SIZE = 16 << 10 // 16K

type Buffer []byte

var bufPool = sync.Pool{
	New: func() any {
		b := make(Buffer, 0, BUFFER_DEFAULT_CAP)
		return (*Buffer)(&b)
	},
}

func NewBuffer() *Buffer {
	return bufPool.Get().(*Buffer)
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool
	if cap(*b) <= BUFFER_MAX_SIZE {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *Buffer) Write(bytes []byte) int {
	*b = append(*b, bytes...)
	return len(bytes)
}

func (b *Buffer) WriteByte(char byte) {
	*b = append(*b, char)
}

func (b *Buffer) WriteString(str string) int {
	*b = append(*b, str...)
	return len(str)
}

func (b *Buffer) WriteStringIf(ok bool, str string) int {
	if !ok {
		return 0
	}
	return b.WriteString(str)
}

// EOF: "buffer.go"
