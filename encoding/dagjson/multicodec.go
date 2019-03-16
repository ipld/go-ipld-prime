package dagjson

import (
	"io"

	"github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
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

// FIXME: technically these are NOT dag-json; they're just regular json.
// We need to get encoder logic that handles the special links cases.

func Decoder(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Shell out directly to generic builder path.
	//  (There's not really any fastpaths of note for json.)
	return encoding.Unmarshal(nb, json.NewDecoder(r))
}

func Encoder(n ipld.Node, w io.Writer) error {
	// Shell out directly to generic inspection path.
	//  (There's not really any fastpaths of note for json.)
	// Write another function if you need to tune encoding options about whitespace.
	return encoding.Marshal(n, json.NewEncoder(w, json.EncodeOptions{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
	}))
}
