package cbor

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

var n = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
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
		err := Encode(n, &buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, serial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(serial)
		nb := basicnode.Prototype.Map.NewBuilder()
		err := Decode(nb, buf)
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
