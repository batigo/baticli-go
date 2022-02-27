package baticli

import (
	"bytes"
	"io/ioutil"

	"github.com/klauspost/compress/flate"
)

func newCompressor(typ CompressorType) Compressor {
	switch typ {
	case CompressorTypeDeflate:
		return DeflateCompressor{}
	default:
		return NullCompressor{}
	}
}

type CompressorType string

const (
	CompressorTypeNull    CompressorType = "null"
	CompressorTypeDeflate CompressorType = "deflate"
)

type Compressor interface {
	Compress([]byte) ([]byte, error)
	Uncompress([]byte) ([]byte, error)
	String() CompressorType
}

type DeflateCompressor struct{}

func (c DeflateCompressor) Compress(i []byte) (o []byte, err error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, 6)
	if err != nil {
		return
	}

	defer w.Close()
	_, err = w.Write(i)
	if err != nil {
		return
	}

	if err = w.Flush(); err != nil {
		return
	}

	o = buf.Bytes()
	return
}

func (c DeflateCompressor) Uncompress(i []byte) (o []byte, err error) {
	r := flate.NewReader(bytes.NewReader(i))
	defer r.Close()
	o, _ = ioutil.ReadAll(r)
	return
}

func (c DeflateCompressor) String() CompressorType {
	return CompressorTypeDeflate
}

type NullCompressor struct{}

func (c NullCompressor) Compress(i []byte) (o []byte, err error) {
	return i, nil
}

func (c NullCompressor) Uncompress(i []byte) (o []byte, err error) {
	return i, nil
}

func (c NullCompressor) String() CompressorType {
	return CompressorTypeNull
}
