// File: "buffer.go"
// Code based on https://github.com/lmittmann/tint/blob/main/buffer.go

package xlog

import "sync"

const BUFFER_DEFAULT_CAP = 1024

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
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *Buffer) Write(bytes []byte) (int, error) {
	*b = append(*b, bytes...)
	return len(bytes), nil
}

func (b *Buffer) WriteByte(char byte) error {
	*b = append(*b, char)
	return nil
}

func (b *Buffer) WriteString(str string) (int, error) {
	*b = append(*b, str...)
	return len(str), nil
}

func (b *Buffer) WriteStringIf(ok bool, str string) (int, error) {
	if !ok {
		return 0, nil
	}
	return b.WriteString(str)
}

// EOF: "buffer.go"
