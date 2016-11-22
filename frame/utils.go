package frame

import "sync"
import "bytes"

const (
	chunkSize = 1 << 15
)

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	bytesPool = sync.Pool{
		New: func() interface{} {
			return &[chunkSize]byte{}
		},
	}
	voidBytes = make([]byte, chunkSize)
)

func getBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func returnBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

func getBytes() *[chunkSize]byte {
	return bytesPool.Get().(*[chunkSize]byte)
}

func returnBytes(p *[chunkSize]byte) {
	copy((*p)[:], voidBytes)
	bytesPool.Put(p)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
