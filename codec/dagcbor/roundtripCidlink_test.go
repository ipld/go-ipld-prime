package dagcbor

import (
	"bytes"
	"io"
	"testing"

	. "github.com/warpfork/go-wish"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestRoundtripCidlink(t *testing.T) {
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    0x71,
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
