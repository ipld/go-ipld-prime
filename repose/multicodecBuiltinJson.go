package repose

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/polydawn/refmt/json"
)

var (
	_ MulticodecDecoder = DecoderDagJson
	_ MulticodecEncoder = EncoderDagJson
)

func DecoderDagJson(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Shell out directly to generic builder path.
	//  (There's not really any fastpaths of note for json.)
	return encoding.Unmarshal(nb, json.NewDecoder(r))
}

func EncoderDagJson(n ipld.Node, w io.Writer) error {
	// Shell out directly to generic inspection path.
	//  (There's not really any fastpaths of note for json.)
	// Write another function if you need to tune encoding options about whitespace.
	return encoding.Marshal(n, json.NewEncoder(w, json.EncodeOptions{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
	}))
}
