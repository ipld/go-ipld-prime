package dagjson

import (
	"bytes"
	"io"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestRoundtripCidlink(t *testing.T) {
	lp := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x13,
		MhLength: 4,
	}}
	lsys := cidlink.DefaultLinkSystem()

	buf := bytes.Buffer{}
	lsys.StorageWriteOpener = func(lnkCtx ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
		return &buf, func(lnk ipld.Link) error { return nil }, nil
	}
	lsys.StorageReadOpener = func(lnkCtx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
		return bytes.NewReader(buf.Bytes()), nil
	}

	lnk, err := lsys.Store(ipld.LinkContext{}, lp, n)
	Require(t, err, ShouldEqual, nil)

	n2, err := lsys.Load(ipld.LinkContext{}, lnk, basicnode.Prototype.Any)
	Require(t, err, ShouldEqual, nil)
	Wish(t, n2, ShouldEqual, n)
}

// Make sure that a map that *almost* looks like a link is handled safely.
//
// This is aiming very specifically at the corner case where a minimal number of
// tokens have to be reprocessed before a recursion that find a real link appears.
func TestUnmarshalTrickyMapContainingLink(t *testing.T) {
	// Create a link; don't particularly care about its contents.
	lnk := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x71,
		MhType:   0x13,
		MhLength: 4,
	}}.BuildLink([]byte{1, 2, 3, 4}) // dummy value, content does not matter to this test.

	// Compose the tricky corpus.  (lnk.String "happens" to work here, although this isn't recommended or correct in general.)
	tricky := `{"/":{"/":"` + lnk.String() + `"}}`

	// Unmarshal.  Hopefully we get a map with a link in it.
	nb := basicnode.Prototype__Any{}.NewBuilder()
	err := Decode(nb, strings.NewReader(tricky))
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
	n2, err := n.LookupByString("/")
	Require(t, err, ShouldEqual, nil)
	Wish(t, n2.Kind(), ShouldEqual, ipld.Kind_Link)
}
