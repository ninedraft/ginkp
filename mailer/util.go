package mailer

import (
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
