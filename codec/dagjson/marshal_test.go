package dagjson

import (
	"testing"

	qt "github.com/frankban/quicktest"
	cid "github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestMarshalUndefCid(t *testing.T) {
	link, err := cid.Decode("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi")
	qt.Assert(t, err, qt.IsNil)
	node, err := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "UndefCid", qp.Link(cidlink.Link{Cid: cid.Undef}))
		qp.MapEntry(ma, "DefCid", qp.Link(cidlink.Link{Cid: link}))
	})
	qt.Assert(t, err, qt.IsNil)
	_, err = ipld.Encode(node, Encode)
	qt.Assert(t, err, qt.ErrorMatches, "encoding undefined CIDs are not supported by this codec")
}
