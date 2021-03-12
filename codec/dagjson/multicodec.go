package dagjson

import (
	"fmt"
	"io"

	"github.com/polydawn/refmt/json"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
)

var (
	_ ipld.Decoder = Decode
	_ ipld.Encoder = Encode
)

func init() {
	multicodec.RegisterEncoder(0x0129, Encode)
	multicodec.RegisterDecoder(0x0129, Decode)
}

func Decode(na ipld.NodeAssembler, r io.Reader) error {
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

func Encode(n ipld.Node, w io.Writer) error {
	// Shell out directly to generic inspection path.
	//  (There's not really any fastpaths of note for json.)
	// Write another function if you need to tune encoding options about whitespace.
	return Marshal(n, json.NewEncoder(w, json.EncodeOptions{
		Line:   []byte{'\n'},
		Indent: []byte{'\t'},
	}))
}
