package mailer

import (
	"io"
	"sync"
)

var (
	chunkPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, CHUNK_SIZE)
		},
	}
	voidChunk = make([]byte, CHUNK_SIZE)
)

func getChunk() []byte {
	return chunkPool.Get().([]byte)
}

func returnChunk(p []byte) {
	copy(p, voidChunk)
	chunkPool.Put(p)
}

func copyBuf(w io.Writer, r io.Reader, buf []byte) (int, error) {
	n, err := r.Read(buf)
	if err != nil {
		return 0, err
	}
	n, err = w.Write(buf[:n])
	return n, err
}
