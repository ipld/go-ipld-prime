package dagcbor

import (
	"bytes"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

var n = fluent.MustBuildMap(basicnode.Style__Map{}, 4, func(na fluent.MapAssembler) {
	na.AssembleDirectly("plain").AssignString("olde string")
	na.AssembleDirectly("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleDirectly("one").AssignInt(1)
		na.AssembleDirectly("two").AssignInt(2)
	})
	na.AssembleDirectly("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleDirectly("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleDirectly("deeper").CreateList(1, func(na fluent.ListAssembler) {
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
		nb := basicnode.Style__Map{}.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, n)
	})
}

func TestRoundtripScalar(t *testing.T) {
	nb := basicnode.Style__String{}.NewBuilder()
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
		nb := basicnode.Style__String{}.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, simple)
	})
}
