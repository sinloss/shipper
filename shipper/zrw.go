package shipper

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"

	"github.com/sinloss/shipper/util"
)

type w struct {
	f *os.File
}

// Gzip compresses the given bytes
func Gzip(f *os.File, p []byte) (n int, err error) {
	zw := gzip.NewWriter(&w{f})
	defer zw.Close()
	return zw.Write(p)
}

func (w *w) Write(p []byte) (n int, err error) {
	hex := make([]byte, len(p)*4)
	j := 0
	for _, b := range p {
		hex[j], hex[j+1] = '\\', 'x'
		hex[j+2], hex[j+3] = util.Hexchar(b)
		j += 4
	}
	return w.f.Write(hex)
}

// UnGzip uncompresses the given gz format bytes
func UnGzip(p []byte) (data []byte, err error) {
	zr, err := gzip.NewReader(bytes.NewBuffer(p))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return ioutil.ReadAll(zr)
}
