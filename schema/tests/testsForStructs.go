package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestFoo(t *testing.T, newNb func(typ schema.Type) ipld.NodeBuilder) {
	tString := schema.SpawnString("String")
	tStroct := schema.SpawnStruct("Stroct",
		[]schema.StructField{
			schema.SpawnStructField("f1", tString, false, false),
			schema.SpawnStructField("f2", tString, true, false),
			schema.SpawnStructField("f3", tString, true, true),
			schema.SpawnStructField("f4", tString, false, false),
		},
		schema.StructRepresentation_Map{},
	)
	nbString := newNb(tString)
	nbStroct := newNb(tStroct)
	var n1 ipld.Node
	t.Run("test building", func(t *testing.T) {
		t.Run("all values valid", func(t *testing.T) {
			mb, err := nbStroct.CreateMap()
			Wish(t, err, ShouldEqual, nil)
			// Set 'f1' to a valid, typed string.
			v, _ := nbString.CreateString("asdf")
			Wish(t, mb.Insert(ipldfree.String("f1"), v), ShouldEqual, nil)
			// Skip setting 'f2' -- it's optional.
			// Set 'f3' to null.  Nulls aren't typed.
			Wish(t, mb.Insert(ipldfree.String("f3"), ipld.Null), ShouldEqual, nil)
			// Set 'f4' to a valid, typed string.
			v, _ = nbString.CreateString("qwer")
			Wish(t, mb.Insert(ipldfree.String("f4"), v), ShouldEqual, nil)
			n1, err = mb.Build()
			Wish(t, err, ShouldEqual, nil)
		})
		t.Run("wrong type rejected", func(t *testing.T) {

		})
		t.Run("invalid key rejected", func(t *testing.T) {

		})
		t.Run("missing nonoptional rejected", func(t *testing.T) {

		})
		t.Run("null nonnullable rejected", func(t *testing.T) {

		})
	})
	t.Run("reading", func(t *testing.T) {
		t.Run("regular fields", func(t *testing.T) {
			v, err := n1.TraverseField("f1")
			Wish(t, err, ShouldEqual, nil)
			v2, _ := nbString.CreateString("asdf")
			Wish(t, v, ShouldEqual, v2)
		})
		t.Run("optional absent fields", func(t *testing.T) {

		})
		t.Run("null nullable fields", func(t *testing.T) {

		})
		t.Run("invalid key", func(t *testing.T) {

		})
		t.Run("iterating", func(t *testing.T) {

		})
	})
	t.Run("test parsing", func(t *testing.T) {
		t.Run("strict order", func(t *testing.T) {
			// FUTURE
		})
	})
}
