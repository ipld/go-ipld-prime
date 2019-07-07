package dagjson

import (
	"bytes"
	"context"
	"io"
	"testing"

	. "github.com/warpfork/go-wish"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

func TestRoundtripCidlink(t *testing.T) {
	lb := cidlink.LinkBuilder{Prefix: cid.Prefix{
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
	Wish(t, err, ShouldEqual, nil)
	n2, err := lnk.Load(context.Background(), ipld.LinkContext{}, ipldfree.NodeBuilder(),
		func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
			return bytes.NewBuffer(buf.Bytes()), nil
		},
	)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, n2, ShouldEqual, n)
}
