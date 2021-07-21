package dagcbor

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"

	cid "github.com/ipfs/go-cid"
	. "github.com/warpfork/go-wish"

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
var nSorted = fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleEntry("one").AssignInt(1)
		na.AssembleEntry("two").AssignInt(2)
	})
	na.AssembleEntry("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleEntry("plain").AssignString("olde string")
	na.AssembleEntry("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleEntry("deeper").CreateList(1, func(na fluent.ListAssembler) {
			na.AssembleValue().AssignString("things")
		})
	})
})
var serial = "\xa4cmap\xa2cone\x01ctwo\x02dlist\x82ethreedfoureplainkolde stringfnested\xa1fdeeper\x81fthings"

func TestRoundtrip(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encode(n, &buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, serial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(serial)
		nb := basicnode.Prototype__Map{}.NewBuilder()
		err := Decode(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, nSorted)
	})
}

func TestRoundtripScalar(t *testing.T) {
	nb := basicnode.Prototype__String{}.NewBuilder()
	nb.AssignString("applesauce")
	simple := nb.Build()
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encode(simple, &buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, `japplesauce`)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(`japplesauce`)
		nb := basicnode.Prototype__String{}.NewBuilder()
		err := Decode(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, simple)
	})
}

func TestRoundtripLinksAndBytes(t *testing.T) {
	lnk := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x71,
		MhType:   0x13,
		MhLength: 4,
	}}.BuildLink([]byte{1, 2, 3, 4}) // dummy value, content does not matter to this test.

	var linkByteNode = fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		nva := na.AssembleEntry("Link")
		nva.AssignLink(lnk)
		nva = na.AssembleEntry("Bytes")
		bytes := make([]byte, 100)
		_, _ = rand.Read(bytes)
		nva.AssignBytes(bytes)
	})

	buf := bytes.Buffer{}
	err := Encode(linkByteNode, &buf)
	Require(t, err, ShouldEqual, nil)
	nb := basicnode.Prototype__Map{}.NewBuilder()
	err = Decode(nb, &buf)
	Require(t, err, ShouldEqual, nil)
	reconstructed := nb.Build()
	Wish(t, reconstructed, ShouldEqual, linkByteNode)
}
