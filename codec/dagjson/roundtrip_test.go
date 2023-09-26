package dagjson_test

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/testutil/garbage"
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
var nSorted = fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("list").CreateList(2, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignString("three")
		na.AssembleValue().AssignString("four")
	})
	na.AssembleEntry("map").CreateMap(2, func(na fluent.MapAssembler) {
		na.AssembleEntry("one").AssignInt(1)
		na.AssembleEntry("two").AssignInt(2)
	})
	na.AssembleEntry("nested").CreateMap(1, func(na fluent.MapAssembler) {
		na.AssembleEntry("deeper").CreateList(1, func(na fluent.ListAssembler) {
			na.AssembleValue().AssignString("things")
		})
	})
	na.AssembleEntry("plain").AssignString("olde string")
})
var serial = `{"list":["three","four"],"map":{"one":1,"two":2},"nested":{"deeper":["things"]},"plain":"olde string"}`

func TestRoundtrip(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagjson.Encode(n, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, serial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(serial)
		nb := basicnode.Prototype.Map.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, nSorted)
	})
}

func TestRoundtripScalar(t *testing.T) {
	nb := basicnode.Prototype__String{}.NewBuilder()
	nb.AssignString("applesauce")
	simple := nb.Build()
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagjson.Encode(simple, &buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, buf.String(), qt.Equals, `"applesauce"`)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(`"applesauce"`)
		nb := basicnode.Prototype__String{}.NewBuilder()
		err := dagjson.Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, nb.Build(), nodetests.NodeContentEquals, simple)
	})
}

func TestGarbage(t *testing.T) {
	t.Run("small garbage", func(t *testing.T) {
		seed := time.Now().Unix()
		t.Logf("randomness seed: %v\n", seed)
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < 1000; i++ {
			gbg := garbage.Generate(rnd, garbage.TargetBlockSize(1<<6))
			var buf bytes.Buffer
			err := dagjson.Encode(gbg, &buf)
			qt.Assert(t, err, qt.IsNil)
			nb := basicnode.Prototype.Any.NewBuilder()
			err = dagjson.Decode(nb, bytes.NewReader(buf.Bytes()))
			qt.Assert(t, err, qt.IsNil)
			qt.Check(t, nb.Build(), nodetests.DeepNodeContentsEquals, gbg)
		}
	})

	t.Run("medium garbage", func(t *testing.T) {
		seed := time.Now().Unix()
		t.Logf("randomness seed: %v\n", seed)
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < 100; i++ {
			gbg := garbage.Generate(rnd, garbage.TargetBlockSize(1<<16))
			var buf bytes.Buffer
			err := dagjson.Encode(gbg, &buf)
			qt.Assert(t, err, qt.IsNil)
			nb := basicnode.Prototype.Any.NewBuilder()
			err = dagjson.Decode(nb, bytes.NewReader(buf.Bytes()))
			qt.Assert(t, err, qt.IsNil)
			qt.Check(t, nb.Build(), nodetests.DeepNodeContentsEquals, gbg)
		}
	})

	t.Run("large garbage", func(t *testing.T) {
		seed := time.Now().Unix()
		t.Logf("randomness seed: %v\n", seed)
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < 10; i++ {
			gbg := garbage.Generate(rnd, garbage.TargetBlockSize(1<<20))
			var buf bytes.Buffer
			err := dagjson.Encode(gbg, &buf)
			qt.Assert(t, err, qt.IsNil)
			nb := basicnode.Prototype.Any.NewBuilder()
			err = dagjson.Decode(nb, bytes.NewReader(buf.Bytes()))
			qt.Assert(t, err, qt.IsNil)
			qt.Check(t, nb.Build(), nodetests.DeepNodeContentsEquals, gbg)
		}
	})
}
