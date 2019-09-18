package compress

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func ZlibCompress(bs []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer, err := zlib.NewWriterLevel(buf, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err = writer.Write(bs); err != nil {
		return nil, err
	}
	if err = writer.Flush(); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ZlibDecompress(bs []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
