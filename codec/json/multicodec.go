package json

import (
	"io"

	rfmtjson "github.com/polydawn/refmt/json"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/multicodec"
)

var (
	_ ipld.Decoder = Decode
	_ ipld.Encoder = Encode
)

func init() {
	multicodec.RegisterEncoder(0x0200, Encode)
	multicodec.RegisterDecoder(0x0200, Decode)
}

// Decode deserializes data from the given io.Reader and feeds it into the given ipld.NodeAssembler.
// Decode fits the ipld.Decoder function interface.
//
// This is the function that will be registered in the default multicodec registry during package init time.
func Decode(na ipld.NodeAssembler, r io.Reader) error {
	return dagjson.DecodeOptions{
		ParseLinks: false,
		ParseBytes: false,
	}.Decode(na, r)
}

// Encode walks the given ipld.Node and serializes it to the given io.Writer.
// Encode fits the ipld.Encoder function interface.
//
// This is the function that will be registered in the default multicodec registry during package init time.
func Encode(n ipld.Node, w io.Writer) error {
	// Shell out directly to generic inspection path.
	//  (There's not really any fastpaths of note for json.)
	// Write another function if you need to tune encoding options about whitespace.
	return dagjson.Marshal(n, rfmtjson.NewEncoder(w, rfmtjson.EncodeOptions{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
	}), dagjson.EncodeOptions{
		EncodeLinks: false,
		EncodeBytes: false,
		MapSortMode: codec.MapSortMode_None,
	})
}
