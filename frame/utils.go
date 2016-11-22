package frame

import "sync"
import "bytes"

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	bytesPool = sync.Pool{
		New: func() interface{} {
			return &[8]byte{}
		},
	}
	voidBytes = make([]byte, 8)
)

func getBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func returnBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

func getBytes() *[8]byte {
	return bytesPool.Get().(*[8]byte)
}

func returnBytes(p *[8]byte) {
	copy((*p)[:], voidBytes)
	bytesPool.Put(p)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
