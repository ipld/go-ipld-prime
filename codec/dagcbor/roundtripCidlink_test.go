package dagcbor

import (
	"bytes"
	"context"
	"io"
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
		Codec:    0x71,
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

	nb := basicnode.Style__Any{}.NewBuilder()
	err = lnk.Load(context.Background(), ipld.LinkContext{}, nb,
		func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
			return bytes.NewBuffer(buf.Bytes()), nil
		},
	)
	Require(t, err, ShouldEqual, nil)
	Wish(t, nb.Build(), ShouldEqual, n)
}
