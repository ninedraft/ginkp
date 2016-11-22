package frame

import (
	"bytes"
	"testing"
)

func TestFrameDataSize(t *testing.T) {
	fr := NewFrame()
	str := bytes.NewBufferString("lorem")
	_, err := fr.ReadDataFrom(str)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("frame %d", fr.Bytes())
	t.Logf("frame %q", string(fr.Bytes()))
	t.Logf("frame len %d", fr.Len())
	t.Logf("data size %d", fr.DataSize())
	if fr.Len()-2 != fr.DataSize() {
		t.Fatal("frame invalid size field!")
	}
}
