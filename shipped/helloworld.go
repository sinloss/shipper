// Code generated by shipper; DO NOT EDIT.

package shipped

import (
	"github.com/sinloss/shipper/shipper"
)

// A is the Asset
var A = &shipper.Assets{

	"hello": shipper.Content{
		Gziped: false,
		Bytes:  []byte("\x68\x0d\x0a"),
	},
	"world/bar.foo": shipper.Content{
		Gziped: false,
		Bytes:  []byte("\x62\x0d\x0a"),
	},
	"world/foo.bar": shipper.Content{
		Gziped: true,
		Bytes:  []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\x53\xe0\xe5\x02\x04\x00\x00\xff\xff\x70\xa6\x3f\x52\x04\x00\x00\x00"),
	},
}
