package cbor

import (
	"io"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/multicodec"
)

var (
	_ ipld.Decoder = Decode
	_ ipld.Encoder = Encode
)

func init() {
	multicodec.RegisterEncoder(0x51, Encode)
	multicodec.RegisterDecoder(0x51, Decode)
}

// Decode deserializes data from the given io.Reader and feeds it into the given ipld.NodeAssembler.
// Decode fits the ipld.Decoder function interface.
//
// This is the function that will be registered in the default multicodec registry during package init time.
func Decode(na ipld.NodeAssembler, r io.Reader) error {
	return dagcbor.DecodeOptions{
		AllowLinks: false,
	}.Decode(na, r)
}

// Encode walks the given ipld.Node and serializes it to the given io.Writer.
// Encode fits the ipld.Encoder function interface.
//
// This is the function that will be registered in the default multicodec registry during package init time.
func Encode(n ipld.Node, w io.Writer) error {
	return dagcbor.EncodeOptions{
		AllowLinks: false,
	}.Encode(n, w)
}
