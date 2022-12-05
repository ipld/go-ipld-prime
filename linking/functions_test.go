package linking_test

import (
	"bytes"
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/multiformats/go-multicodec"
)

func TestLinkSystem_LoadHashMismatch(t *testing.T) {
	subject := cidlink.DefaultLinkSystem()
	storage := &memstore.Store{}
	subject.SetReadStorage(storage)
	subject.SetWriteStorage(storage)

	// Construct some test IPLD node.
	wantNode := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
		na.AssembleEntry("fish").AssignString("barreleye")
	})

	// Encode as raw value to be used for testing LoadRaw
	var buf bytes.Buffer
	qt.Check(t, dagcbor.Encode(wantNode, &buf), qt.IsNil)
	wantNodeRaw := buf.Bytes()

	// Store the test IPLD node and get link back.
	lctx := ipld.LinkContext{Ctx: context.TODO()}
	gotLink, err := subject.Store(lctx, cidlink.LinkPrototype{
		Prefix: cid.Prefix{
			Version:  1,
			Codec:    uint64(multicodec.DagCbor),
			MhType:   uint64(multicodec.Sha2_256),
			MhLength: -1,
		},
	}, wantNode)
	qt.Check(t, err, qt.IsNil)
	gotCidlink := gotLink.(cidlink.Link)

	// Assert all load variations return expected values for different link representations.
	for _, test := range []struct {
		name string
		link datamodel.Link
	}{
		{"datamodel.Link", gotLink},
		{"cidlink.Link", gotCidlink},
		{"&cidlink.Link", &gotCidlink},
	} {
		t.Run(test.name, func(t *testing.T) {
			gotNode, err := subject.Load(lctx, test.link, basicnode.Prototype.Any)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, ipld.DeepEqual(wantNode, gotNode), qt.IsTrue)

			gotNodeRaw, err := subject.LoadRaw(lctx, test.link)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, bytes.Equal(wantNodeRaw, gotNodeRaw), qt.IsTrue)

			gotNode, gotNodeRaw, err = subject.LoadPlusRaw(lctx, test.link, basicnode.Prototype.Any)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, ipld.DeepEqual(wantNode, gotNode), qt.IsTrue)
			qt.Check(t, bytes.Equal(wantNodeRaw, gotNodeRaw), qt.IsTrue)
		})
	}
}
