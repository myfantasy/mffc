package compress

import (
	"bytes"
	"compress/flate"
	"io"
	"strings"
	"sync"
)

// DeflateCompressor - compressor
//
type DeflateCompressor struct {
	CompressLevel int

	mx  sync.Mutex
	r   strings.Reader
	b   bytes.Buffer
	buf []byte
	zw  *flate.Writer
	zr  io.ReadCloser
}

// DeflateCompressorCreate creates deflate compressor
func DeflateCompressorCreate(level int) (*DeflateCompressor, error) {

	zw, err := flate.NewWriter(nil, level)

	if err != nil {
		return nil, err
	}

	return &DeflateCompressor{
		CompressLevel: level,
		buf:           make([]byte, 32<<10),
		zw:            zw,
		zr:            flate.NewReader(nil),
	}, nil
}

// Clone creates deflate compressor
// if it throw error then panic
func (dc *DeflateCompressor) Clone() *DeflateCompressor {

	r, err := DeflateCompressorCreate(dc.CompressLevel)

	if err != nil {
		panic(err)
	}

	return r

}

// Compress bytes data
func (dc *DeflateCompressor) Compress(data []byte) (res []byte, err error) {

	dc.mx.Lock()
	defer dc.mx.Unlock()

	dc.b.Reset()
	dc.zw.Reset(&dc.b)

	r := bytes.NewReader(data)

	if _, err = io.CopyBuffer(dc.zw, r, dc.buf); err != nil {
		return res, err
	}
	if err = dc.zw.Close(); err != nil {
		return res, err
	}

	return dc.b.Bytes(), nil
}

// Restore data into bytes
func (dc *DeflateCompressor) Restore(data []byte) (res []byte, err error) {

	dc.mx.Lock()
	defer dc.mx.Unlock()

	dc.b.Reset()

	dc.b.Write(data)

	// Reset the decompressor and decode to some output stream.
	if err := (dc.zr).(flate.Resetter).Reset(&dc.b, nil); err != nil {
		return res, err
	}
	var b bytes.Buffer
	if _, err := io.CopyBuffer(&b, dc.zr, dc.buf); err != nil {
		return res, err
	}
	if err := dc.zr.Close(); err != nil {
		return res, err
	}

	return b.Bytes(), nil
}
