package dagjson

import (
	"io"

	"github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

var (
	_ cidlink.MulticodecDecoder = Decoder
	_ cidlink.MulticodecEncoder = Encoder
)

func init() {
	cidlink.RegisterMulticodecDecoder(0x0129, Decoder)
	cidlink.RegisterMulticodecEncoder(0x0129, Encoder)
}

func Decoder(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Shell out directly to generic builder path.
	//  (There's not really any fastpaths of note for json.)
	return Unmarshal(nb, json.NewDecoder(r))
}

func Encoder(n ipld.Node, w io.Writer) error {
	// Shell out directly to generic inspection path.
	//  (There's not really any fastpaths of note for json.)
	// Write another function if you need to tune encoding options about whitespace.
	return Marshal(n, json.NewEncoder(w, json.EncodeOptions{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
	}))
}
