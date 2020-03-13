package compress

import (
	"compress/flate"
	"testing"
)

func TestCompressRestore(t *testing.T) {

	text := "some easy text"

	ddc, err := DeflateCompressorCreate(flate.BestCompression)
	if err != nil {
		t.Fatal("DeflateCompressorCreate", err)
	}

	b, err := ddc.Compress([]byte(text))
	if err != nil {
		t.Fatal("ddc.Compress", err)
	}

	b2, err := ddc.Restore(b)
	if err != nil {
		t.Fatal("ddc.Restore", err)
	}

	if string(b2) != text {
		t.Fatal("Compress and Restore fail ", string(b2), " ", text)
	}
}

func BenchmarkCompress(b *testing.B) {
	text := "some easy text"
	for k := 0; k < 10; k++ {
		text = text + text
	}

	ddc, err := DeflateCompressorCreate(flate.BestCompression)
	if err != nil {
		b.Fatal("DeflateCompressorCreate", err)
	}

	for i := 0; i < b.N; i++ {
		_, err = ddc.Compress([]byte(text))
		if err != nil {
			b.Fatal("ddc.Compress", err)
		}
	}
}
