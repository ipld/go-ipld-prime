package dagjson

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestRoundtripCidlink(t *testing.T) {
	lb := cidlink.LinkBuilder{cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x17,
		MhLength: 4,
	}}

	buf := bytes.Buffer{}
	lnk, err := lb.Build(context.Background(), ipld.LinkContext{}, n,
		func(ipld.LinkContext) (io.Writer, ipld.StoreCommitter, error) {
			return &buf, func(lnk ipld.Link) error { return nil }, nil
		},
	)
	Require(t, err, ShouldEqual, nil)

	nb := basicnode.Prototype__Any{}.NewBuilder()
	err = lnk.Load(context.Background(), ipld.LinkContext{}, nb,
		func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
			return bytes.NewReader(buf.Bytes()), nil
		},
	)
	Require(t, err, ShouldEqual, nil)
	Wish(t, nb.Build(), ShouldEqual, n)
}

// Make sure that a map that *almost* looks like a link is handled safely.
//
// This is aiming very specifically at the corner case where a minimal number of
// tokens have to be reprocessed before a recursion that find a real link appears.
func TestUnmarshalTrickyMapContainingLink(t *testing.T) {
	// Create a link; don't particularly care about its contents.
	lnk, err := cidlink.LinkBuilder{cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x17,
		MhLength: 4,
	}}.Build(context.Background(), ipld.LinkContext{}, n,
		func(ipld.LinkContext) (io.Writer, ipld.StoreCommitter, error) {
			return ioutil.Discard, func(lnk ipld.Link) error { return nil }, nil
		},
	)
	Require(t, err, ShouldEqual, nil)

	// Compose the tricky corpus.  (lnk.String "happens" to work here, although this isn't recommended or correct in general.)
	tricky := `{"/":{"/":"` + lnk.String() + `"}}`

	// Unmarshal.  Hopefully we get a map with a link in it.
	nb := basicnode.Prototype__Any{}.NewBuilder()
	err = Decoder(nb, strings.NewReader(tricky))
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
	n2, err := n.LookupByString("/")
	Require(t, err, ShouldEqual, nil)
	Wish(t, n2.Kind(), ShouldEqual, ipld.Kind_Link)
}
