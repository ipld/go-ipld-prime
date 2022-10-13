package json

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

var link = cid.MustParse("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi")

// mirrored in dag-json but without errors
func TestMarshalLinks(t *testing.T) {
	linkNode := basicnode.NewLink(cidlink.Link{Cid: link})
	mapNode, err := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "Lnk", qp.Node(linkNode))
	})
	qt.Assert(t, err, qt.IsNil)

	t.Run("link json", func(t *testing.T) {
		_, err := ipld.Encode(linkNode, Encode)
		qt.Assert(t, err, qt.ErrorMatches, "cannot marshal IPLD links to this codec")
	})
	t.Run("nested link json", func(t *testing.T) {
		_, err := ipld.Encode(mapNode, Encode)
		qt.Assert(t, err, qt.ErrorMatches, "cannot marshal IPLD links to this codec")
	})
}

// mirrored in dag-json but without errors
func TestMarshalBytes(t *testing.T) {
	bytsNode := basicnode.NewBytes([]byte("byte me"))
	mapNode, err := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "Byts", qp.Node(bytsNode))
	})
	qt.Assert(t, err, qt.IsNil)

	t.Run("bytes json", func(t *testing.T) {
		_, err := ipld.Encode(bytsNode, Encode)
		qt.Assert(t, err, qt.ErrorMatches, "cannot marshal IPLD bytes to this codec")
	})
	t.Run("nested bytes json", func(t *testing.T) {
		_, err := ipld.Encode(mapNode, Encode)
		qt.Assert(t, err, qt.ErrorMatches, "cannot marshal IPLD bytes to this codec")
	})
}
