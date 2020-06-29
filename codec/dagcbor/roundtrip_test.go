package dagcbor

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"testing"

	cid "github.com/ipfs/go-cid"
	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

var n = fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("plain").AssignString("olde string")
	na.AssembleEntry("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleEntry("one").AssignInt(1)
		na.AssembleEntry("two").AssignInt(2)
	})
	na.AssembleEntry("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleEntry("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleEntry("deeper").CreateList(1, func(na fluent.ListAssembler) {
			na.AssembleValue().AssignString("things")
		})
	})
})
var serial = "\xa4eplainkolde stringcmap\xa2cone\x01ctwo\x02dlist\x82ethreedfourfnested\xa1fdeeper\x81fthings"

func TestRoundtrip(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encoder(n, &buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, serial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := bytes.NewBufferString(serial)
		nb := basicnode.Prototype__Map{}.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, n)
	})
}

func TestRoundtripScalar(t *testing.T) {
	nb := basicnode.Prototype__String{}.NewBuilder()
	nb.AssignString("applesauce")
	simple := nb.Build()
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encoder(simple, &buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, `japplesauce`)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := bytes.NewBufferString(`japplesauce`)
		nb := basicnode.Prototype__String{}.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, simple)
	})
}

func TestRoundtripLinksAndBytes(t *testing.T) {
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

	var linkByteNode = fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		nva := na.AssembleEntry("Link")
		nva.AssignLink(lnk)
		nva = na.AssembleEntry("Bytes")
		bytes := make([]byte, 100)
		_, _ = rand.Read(bytes)
		nva.AssignBytes(bytes)
	})

	buf.Reset()
	err = Encoder(linkByteNode, &buf)
	Require(t, err, ShouldEqual, nil)
	nb := basicnode.Prototype__Map{}.NewBuilder()
	err = Decoder(nb, &buf)
	Require(t, err, ShouldEqual, nil)
	reconstructed := nb.Build()
	Wish(t, reconstructed, ShouldEqual, linkByteNode)
}
