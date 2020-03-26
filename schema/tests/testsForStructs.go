package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
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
	var n1 ipld.Node
	t.Run("test building", func(t *testing.T) {
		t.Run("all values valid", func(t *testing.T) {
			nb := newNb(tStroct)
			ma, err := nb.BeginMap(3)
			Wish(t, err, ShouldEqual, nil)
			// Set 'f1' to a valid string.
			va, err := ma.AssembleEntry("f1")
			Wish(t, err, ShouldEqual, nil)
			Wish(t, va.AssignString("asdf"), ShouldEqual, nil)
			// Skip setting 'f2' -- it's optional.
			// Set 'f3' to null.  Nulls aren't typed.
			va, err = ma.AssembleEntry("f3")
			Wish(t, err, ShouldEqual, nil)
			Wish(t, va.AssignNull(), ShouldEqual, nil)
			// Set 'f4' to a valid string.
			va, err = ma.AssembleEntry("f4")
			Wish(t, err, ShouldEqual, nil)
			Wish(t, va.AssignString("qwer"), ShouldEqual, nil)
			Wish(t, ma.Finish(), ShouldEqual, nil)
			n1 = nb.Build()
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
			v, err := n1.LookupString("f1")
			Wish(t, err, ShouldEqual, nil)
			nb := newNb(tString)
			Require(t, nb.AssignString("asdf"), ShouldEqual, nil)
			v2 := nb.Build()
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
