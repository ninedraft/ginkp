package frame

import "sync"
import "bytes"

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

func getBuf() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func returnBuf(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
