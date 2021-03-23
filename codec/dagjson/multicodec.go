package dagjson

import (
	"fmt"
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

func Decoder(na ipld.NodeAssembler, r io.Reader) error {
	// Shell out directly to generic builder path.
	//  (There's not really any fastpaths of note for json.)
	err := Unmarshal(na, json.NewDecoder(r))
	if err != nil {
		return err
	}
	// Slurp any remaining whitespace.
	//  (This is relevant if our reader is tee'ing bytes to a hasher, and
	//   the json contained any trailing whitespace.)
	//  (We can't actually support multiple objects per reader from here;
	//   we can't unpeek if we find a non-whitespace token, so our only
	//    option is to error if this reader seems to contain more content.)
	var buf [1]byte
	for {
		_, err := r.Read(buf[:])
		switch buf[0] {
		case ' ', 0x0, '\t', '\r', '\n': // continue
		default:
			return fmt.Errorf("unexpected content after end of json object")
		}
		if err == nil {
			continue
		} else if err == io.EOF {
			return nil
		} else {
			return err
		}
	}
	return err
}

func Encoder(n ipld.Node, w io.Writer) error {
	// Shell out directly to generic inspection path.
	//  (There's not really any fastpaths of note for json.)
	// Write another function if you need to tune encoding options about whitespace.
	return Marshal(n, json.NewEncoder(w, json.EncodeOptions{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
	}), true)
}
