package dagjson

import (
	"bytes"
	"io"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	cid "github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestRoundtripCidlink(t *testing.T) {
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x13,
		MhLength: 4,
	}}
	lsys := cidlink.DefaultLinkSystem()

	buf := bytes.Buffer{}
	lsys.StorageWriteOpener = func(lnkCtx linking.LinkContext) (io.Writer, linking.BlockWriteCommitter, error) {
		return &buf, func(lnk datamodel.Link) error { return nil }, nil
	}
	lsys.StorageReadOpener = func(lnkCtx linking.LinkContext, lnk datamodel.Link) (io.Reader, error) {
		return bytes.NewReader(buf.Bytes()), nil
	}

	lnk, err := lsys.Store(linking.LinkContext{}, lp, n)
	Require(t, err, ShouldEqual, nil)

	n2, err := lsys.Load(linking.LinkContext{}, lnk, basicnode.Prototype.Any)
	Require(t, err, ShouldEqual, nil)
	Wish(t, n2, ShouldEqual, nSorted)
}

// Make sure that a map that *almost* looks like a link is handled safely.
//
// This is aiming very specifically at the corner case where a minimal number of
// tokens have to be reprocessed before a recursion that find a real link appears.
func TestUnmarshalTrickyMapContainingLink(t *testing.T) {
	// Create a link; don't particularly care about its contents.
	lnk := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    0x71,
		MhType:   0x13,
		MhLength: 4,
	}}.BuildLink([]byte{1, 2, 3, 4}) // dummy value, content does not matter to this test.

	// Compose the tricky corpus.  (lnk.String "happens" to work here, although this isn't recommended or correct in general.)
	tricky := `{"/":{"/":"` + lnk.String() + `"}}`

	// Unmarshal.  Hopefully we get a map with a link in it.
	nb := basicnode.Prototype.Any.NewBuilder()
	err := Decode(nb, strings.NewReader(tricky))
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	Wish(t, n.Kind(), ShouldEqual, datamodel.Kind_Map)
	n2, err := n.LookupByString("/")
	Require(t, err, ShouldEqual, nil)
	Wish(t, n2.Kind(), ShouldEqual, datamodel.Kind_Link)
}
