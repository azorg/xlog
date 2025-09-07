// File: "buffer.go"
// Код на основе https://github.com/lmittmann/tint/blob/main/buffer.go

package xlog

import "sync"

const (
	bufferDefaultCap = 1 * 1024  // 1 KB
	bufferMaxSize    = 16 * 1024 // 16 KB
)

type buffer []byte

var bufPool = sync.Pool{
	New: func() any {
		b := make(buffer, 0, bufferDefaultCap)
		return (*buffer)(&b)
	},
}

func newBuffer() *buffer {
	return bufPool.Get().(*buffer)
}

func (b *buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool
	if cap(*b) <= bufferMaxSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *buffer) Write(bytes []byte) int {
	*b = append(*b, bytes...)
	return len(bytes)
}

func (b *buffer) WriteByte(char byte) error {
	*b = append(*b, char)
	return nil
}

func (b *buffer) WriteString(str string) int {
	*b = append(*b, str...)
	return len(str)
}

func (b *buffer) WriteStringIf(ok bool, str string) int {
	if !ok {
		return 0
	}
	return b.WriteString(str)
}

func (b *buffer) String() string {
	return string(*b)
}

// EOF: "buffer.go"
