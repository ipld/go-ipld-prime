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

var link = cid.MustParse("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi")

func TestMarshalUndefCid(t *testing.T) {
	node, err := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "UndefCid", qp.Link(cidlink.Link{Cid: cid.Undef}))
		qp.MapEntry(ma, "DefCid", qp.Link(cidlink.Link{Cid: link}))
	})
	qt.Assert(t, err, qt.IsNil)
	_, err = ipld.Encode(node, Encode)
	qt.Assert(t, err, qt.ErrorMatches, "encoding undefined CIDs are not supported by this codec")
}

// mirrored in json but with errors
func TestMarshalLinks(t *testing.T) {
	linkNode := basicnode.NewLink(cidlink.Link{Cid: link})
	mapNode, err := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "Lnk", qp.Node(linkNode))
	})
	qt.Assert(t, err, qt.IsNil)

	t.Run("link dag-json", func(t *testing.T) {
		byts, err := ipld.Encode(linkNode, Encode)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, string(byts), qt.Equals,
			`{"/":"bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}`)
	})
	t.Run("nested link dag-json", func(t *testing.T) {
		byts, err := ipld.Encode(mapNode, Encode)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, string(byts), qt.Equals,
			`{"Lnk":{"/":"bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"}}`)
	})
}

// mirrored in json but with errors
func TestMarshalBytes(t *testing.T) {
	bytsNode := basicnode.NewBytes([]byte("byte me"))
	mapNode, err := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "Byts", qp.Node(bytsNode))
	})
	qt.Assert(t, err, qt.IsNil)

	t.Run("bytes dag-json", func(t *testing.T) {
		byts, err := ipld.Encode(bytsNode, Encode)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, string(byts), qt.Equals,
			`{"/":{"bytes":"Ynl0ZSBtZQ"}}`)
	})
	t.Run("nested bytes dag-json", func(t *testing.T) {
		byts, err := ipld.Encode(mapNode, Encode)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, string(byts), qt.Equals,
			`{"Byts":{"/":{"bytes":"Ynl0ZSBtZQ"}}}`)
	})
}
