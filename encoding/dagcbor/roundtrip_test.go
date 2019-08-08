package dagcbor

import (
	"bytes"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
)

var fnb = fluent.WrapNodeBuilder(ipldfree.NodeBuilder())
var n = fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
	mb.Insert(knb.CreateString("plain"), vnb.CreateString("olde string"))
	mb.Insert(knb.CreateString("map"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
		mb.Insert(knb.CreateString("one"), vnb.CreateInt(1))
		mb.Insert(knb.CreateString("two"), vnb.CreateInt(2))
	}))
	mb.Insert(knb.CreateString("list"), fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
		lb.Append(vnb.CreateString("three"))
		lb.Append(vnb.CreateString("four"))
	}))
	mb.Insert(knb.CreateString("nested"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
		mb.Insert(knb.CreateString("deeper"), fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.Append(vnb.CreateString("things"))
		}))
	}))
})
var serial = "\xa4eplainkolde stringcmap\xa2cone\x01ctwo\x02dlist\x82ethreedfourfnested\xa1fdeeper\x81fthings"

func TestRoundtrip(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encoder(n, &buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, serial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := bytes.NewBufferString(serial)
		n2, err := Decoder(ipldfree.NodeBuilder(), buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n2, ShouldEqual, n)
	})
}

func TestRoundtripSimple(t *testing.T) {
	simple := fnb.CreateString("applesauce")
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encoder(simple, &buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, `japplesauce`)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := bytes.NewBufferString(`japplesauce`)
		simple2, err := Decoder(ipldfree.NodeBuilder(), buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, simple2, ShouldEqual, simple)
	})
}
