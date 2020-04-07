package whee

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func plz(n ipld.Node, e error) ipld.Node {
	if e != nil {
		panic(e)
	}
	return n
}
func plzStr(n ipld.Node, e error) string {
	if e != nil {
		panic(e)
	}
	if s, ok := n.AsString(); ok == nil {
		return s
	} else {
		panic(ok)
	}
}
func str(n ipld.Node) string {
	if s, ok := n.AsString(); ok == nil {
		return s
	} else {
		panic(ok)
	}
}
func erp(n ipld.Node, e error) interface{} {
	if e != nil {
		return e
	}
	return n
}

// This targets both "Stroct" and "Stract",
// expecting both to be functionally equivalent
// (because they should be -- only varying in field type name, and whether the maybe of that type uses pointers).
//
// Most of what we're targetting here is if all the matrices of
// nullable and optional support are working correctly.
//
// The type-level generic builder is exercised,
// and the type-level generic accessors are exercised,
// including both lookup and length methods.
// No iterators are exercised (marshal/unmarshal are good at that).  // Okay, yes they are.
// No representations are exercised (that's a whole 'nother topic).  // Okay, reads are.  ReprBuilds aren't.
func TestGeneratedStructWithVariousFieldOptionality(t *testing.T) {
	// There's a lot of cases to cover so a shorthand labels helper funcs:
	//  - 'v' -- value in that entry
	//  - 'z' -- null in that entry
	//  - 'u' -- undefined/absent entry
	build_vvvvv := func(t *testing.T, ns ipld.NodeStyle) schema.TypedNode {
		nb := ns.NewBuilder()
		ma, err := nb.BeginMap(5)
		Require(t, err, ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f1"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("a"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f2"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("b"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f3"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("c"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f4"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("d"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f5"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("e"), ShouldEqual, nil)
		Wish(t, ma.Finish(), ShouldEqual, nil)
		return nb.Build().(schema.TypedNode)
	}
	build_vvvvv_repr := func(t *testing.T, ns ipld.NodeStyle) schema.TypedNode {
		nb := ns.NewBuilder()
		ma, err := nb.BeginMap(5)
		Require(t, err, ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r1"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("a"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r2"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("b"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r3"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("c"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r4"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("d"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f5"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("e"), ShouldEqual, nil)
		Wish(t, ma.Finish(), ShouldEqual, nil)
		return nb.Build().(schema.TypedNode)
	}
	testLookups_vvvvv := func(t *testing.T, n ipld.Node) {
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, n.Length(), ShouldEqual, 5)
		Wish(t, plzStr(n.LookupString("f1")), ShouldEqual, "a")
		Wish(t, plzStr(n.LookupString("f2")), ShouldEqual, "b")
		Wish(t, plzStr(n.LookupString("f3")), ShouldEqual, "c")
		Wish(t, plzStr(n.LookupString("f4")), ShouldEqual, "d")
		Wish(t, plzStr(n.LookupString("f5")), ShouldEqual, "e")
	}
	testLookups_vvvvv_repr := func(t *testing.T, n ipld.Node) {
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, n.Length(), ShouldEqual, 5)
		Wish(t, plzStr(n.LookupString("r1")), ShouldEqual, "a")
		Wish(t, plzStr(n.LookupString("r2")), ShouldEqual, "b")
		Wish(t, plzStr(n.LookupString("r3")), ShouldEqual, "c")
		Wish(t, plzStr(n.LookupString("r4")), ShouldEqual, "d")
		Wish(t, plzStr(n.LookupString("f5")), ShouldEqual, "e")
	}
	testIteration_vvvvv := func(t *testing.T, n ipld.Node) {
		itr := n.MapIterator()
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ := itr.Next()
		Wish(t, str(k), ShouldEqual, "f1")
		Wish(t, str(v), ShouldEqual, "a")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "f2")
		Wish(t, str(v), ShouldEqual, "b")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "f3")
		Wish(t, str(v), ShouldEqual, "c")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "f4")
		Wish(t, str(v), ShouldEqual, "d")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "f5")
		Wish(t, str(v), ShouldEqual, "e")
		Wish(t, itr.Done(), ShouldEqual, true)
		k, v, err := itr.Next()
		Wish(t, k, ShouldEqual, nil)
		Wish(t, v, ShouldEqual, nil)
		Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
	}
	testIteration_vvvvv_repr := func(t *testing.T, n ipld.Node) {
		itr := n.MapIterator()
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ := itr.Next()
		Wish(t, str(k), ShouldEqual, "r1")
		Wish(t, str(v), ShouldEqual, "a")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "r2")
		Wish(t, str(v), ShouldEqual, "b")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "r3")
		Wish(t, str(v), ShouldEqual, "c")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "r4")
		Wish(t, str(v), ShouldEqual, "d")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "f5")
		Wish(t, str(v), ShouldEqual, "e")
		Wish(t, itr.Done(), ShouldEqual, true)
		k, v, err := itr.Next()
		Wish(t, k, ShouldEqual, nil)
		Wish(t, v, ShouldEqual, nil)
		Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
	}
	build_vvzzv := func(t *testing.T, ns ipld.NodeStyle) schema.TypedNode {
		nb := ns.NewBuilder()
		ma, err := nb.BeginMap(5)
		Require(t, err, ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f1"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("a"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f2"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("b"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f3"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignNull(), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f4"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignNull(), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f5"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("e"), ShouldEqual, nil)
		Wish(t, ma.Finish(), ShouldEqual, nil)
		return nb.Build().(schema.TypedNode)
	}
	build_vvzzv_repr := func(t *testing.T, ns ipld.NodeStyle) schema.TypedNode {
		nb := ns.NewBuilder()
		ma, err := nb.BeginMap(5)
		Require(t, err, ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r1"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("a"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r2"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("b"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r3"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignNull(), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("r4"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignNull(), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f5"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("e"), ShouldEqual, nil)
		Wish(t, ma.Finish(), ShouldEqual, nil)
		return nb.Build().(schema.TypedNode)
	}
	testLookups_vvzzv := func(t *testing.T, n ipld.Node) {
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, n.Length(), ShouldEqual, 5)
		Wish(t, plzStr(n.LookupString("f1")), ShouldEqual, "a")
		Wish(t, plzStr(n.LookupString("f2")), ShouldEqual, "b")
		Wish(t, erp(n.LookupString("f3")), ShouldEqual, ipld.Null)
		Wish(t, erp(n.LookupString("f4")), ShouldEqual, ipld.Null)
		Wish(t, plzStr(n.LookupString("f5")), ShouldEqual, "e")
	}
	build_vuvuv := func(t *testing.T, ns ipld.NodeStyle) schema.TypedNode {
		nb := ns.NewBuilder()
		ma, err := nb.BeginMap(3)
		Require(t, err, ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f1"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("a"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f3"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("c"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f5"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("e"), ShouldEqual, nil)
		Wish(t, ma.Finish(), ShouldEqual, nil)
		return nb.Build().(schema.TypedNode)
	}
	testLookups_vuvuv := func(t *testing.T, n ipld.Node) {
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, n.Length(), ShouldEqual, 5)
		Wish(t, plzStr(n.LookupString("f1")), ShouldEqual, "a")
		Wish(t, erp(n.LookupString("f2")), ShouldEqual, ipld.Undef)
		Wish(t, plzStr(n.LookupString("f3")), ShouldEqual, "c")
		Wish(t, erp(n.LookupString("f4")), ShouldEqual, ipld.Undef)
		Wish(t, plzStr(n.LookupString("f5")), ShouldEqual, "e")
	}
	testIteration_vuvuv_repr := func(t *testing.T, n ipld.Node) {
		itr := n.MapIterator()
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ := itr.Next()
		Wish(t, str(k), ShouldEqual, "r1")
		Wish(t, str(v), ShouldEqual, "a")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "r3")
		Wish(t, str(v), ShouldEqual, "c")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "f5")
		Wish(t, str(v), ShouldEqual, "e")
		Wish(t, itr.Done(), ShouldEqual, true)
		k, v, err := itr.Next()
		Wish(t, k, ShouldEqual, nil)
		Wish(t, v, ShouldEqual, nil)
		Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
	}
	build_vvzuu := func(t *testing.T, ns ipld.NodeStyle) schema.TypedNode {
		nb := ns.NewBuilder()
		ma, err := nb.BeginMap(3)
		Require(t, err, ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f1"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("a"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f2"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignString("b"), ShouldEqual, nil)
		Wish(t, ma.AssembleKey().AssignString("f3"), ShouldEqual, nil)
		Wish(t, ma.AssembleValue().AssignNull(), ShouldEqual, nil)
		Wish(t, ma.Finish(), ShouldEqual, nil)
		return nb.Build().(schema.TypedNode)
	}
	testIteration_vvzuu_repr := func(t *testing.T, n ipld.Node) {
		itr := n.MapIterator()
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ := itr.Next()
		Wish(t, str(k), ShouldEqual, "r1")
		Wish(t, str(v), ShouldEqual, "a")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "r2")
		Wish(t, str(v), ShouldEqual, "b")
		Wish(t, itr.Done(), ShouldEqual, false)
		k, v, _ = itr.Next()
		Wish(t, str(k), ShouldEqual, "r3")
		Wish(t, v, ShouldEqual, ipld.Null)
		Wish(t, itr.Done(), ShouldEqual, true)
		k, v, err := itr.Next()
		Wish(t, k, ShouldEqual, nil)
		Wish(t, v, ShouldEqual, nil)
		Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
	}

	t.Run("on stroct", func(t *testing.T) {
		t.Run("type-level build and read", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				// Test building.
				n := build_vvvvv(t, _Stroct__Style{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stroct{
					f1: _String{"a"},
					f2: _String__Maybe{schema.Maybe_Value, _String{"b"}},
					f3: _String__Maybe{schema.Maybe_Value, _String{"c"}},
					f4: _String__Maybe{schema.Maybe_Value, _String{"d"}},
					f5: _String__Maybe{schema.Maybe_Value, _String{"e"}},
				})

				// Test lookup methods.
				testLookups_vvvvv(t, n)
				testIteration_vvvvv(t, n)
			})
			t.Run("setting nulls", func(t *testing.T) {
				// Test building.
				n := build_vvzzv(t, _Stroct__Style{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stroct{
					f1: _String{"a"},
					f2: _String__Maybe{schema.Maybe_Value, _String{"b"}},
					f3: _String__Maybe{schema.Maybe_Null, _String{""}},
					f4: _String__Maybe{schema.Maybe_Null, _String{""}},
					f5: _String__Maybe{schema.Maybe_Value, _String{"e"}},
				})

				// Test lookup methods.
				testLookups_vvzzv(t, n)
			})
			t.Run("not setting optionals", func(t *testing.T) {
				// Test building.
				n := build_vuvuv(t, _Stroct__Style{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stroct{
					f1: _String{"a"},
					f2: _String__Maybe{schema.Maybe_Absent, _String{""}},
					f3: _String__Maybe{schema.Maybe_Value, _String{"c"}},
					f4: _String__Maybe{schema.Maybe_Absent, _String{""}},
					f5: _String__Maybe{schema.Maybe_Value, _String{"e"}},
				})

				// Test lookup methods.
				testLookups_vuvuv(t, n)
				testIteration_vuvuv_repr(t, n.Representation())
			})
			t.Run("absent trailing optionals", func(t *testing.T) {
				// Trailing optionals are especially touchy in a few details of iterators.
				n := build_vvzuu(t, _Stract__Style{})
				testIteration_vvzuu_repr(t, n.Representation())
			})
		})
		t.Run("repr-level build and read", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				// Test building.
				n := build_vvvvv_repr(t, _Stroct__ReprStyle{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stroct{
					f1: _String{"a"},
					f2: _String__Maybe{schema.Maybe_Value, _String{"b"}},
					f3: _String__Maybe{schema.Maybe_Value, _String{"c"}},
					f4: _String__Maybe{schema.Maybe_Value, _String{"d"}},
					f5: _String__Maybe{schema.Maybe_Value, _String{"e"}},
				})

				// Test lookup methods.
				testLookups_vvvvv_repr(t, n.Representation())
				testIteration_vvvvv_repr(t, n.Representation())
			})
			t.Run("setting nulls", func(t *testing.T) {
				// Test building.
				n := build_vvzzv_repr(t, _Stroct__ReprStyle{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stroct{
					f1: _String{"a"},
					f2: _String__Maybe{schema.Maybe_Value, _String{"b"}},
					f3: _String__Maybe{schema.Maybe_Null, _String{""}},
					f4: _String__Maybe{schema.Maybe_Null, _String{""}},
					f5: _String__Maybe{schema.Maybe_Value, _String{"e"}},
				})

				// Test lookup methods.
				testLookups_vvzzv(t, n)
			})
		})
	})
	t.Run("on stract", func(t *testing.T) {
		// Tests targeting `Stract` are *identical* to tests above on `Stroct`,
		//  except for a different concrete type.
		// The point of these is to be reasonably convinced that MaybeT using pointers works,
		//  and after that feels established, we don't exhaustively test the rest.
		// (Improvable?  Should be exhaustive?  Maybe.
		//  It's mostly the concrete type asserts that make this hard.  Maybe those aren't worth having.)
		t.Run("type-level build and read", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				// Test building.
				n := build_vvvvv(t, _Stract__Style{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stract{
					f1: _Strang{"a"},
					f2: _Strang__Maybe{schema.Maybe_Value, &_Strang{"b"}},
					f3: _Strang__Maybe{schema.Maybe_Value, &_Strang{"c"}},
					f4: _Strang__Maybe{schema.Maybe_Value, &_Strang{"d"}},
					f5: _Strang__Maybe{schema.Maybe_Value, &_Strang{"e"}},
				})

				// Test lookup methods.
				testLookups_vvvvv(t, n)
				testIteration_vvvvv(t, n)
			})
			t.Run("setting nulls", func(t *testing.T) {
				// Test building.
				n := build_vvzzv(t, _Stract__Style{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stract{
					f1: _Strang{"a"},
					f2: _Strang__Maybe{schema.Maybe_Value, &_Strang{"b"}},
					f3: _Strang__Maybe{schema.Maybe_Null, nil},
					f4: _Strang__Maybe{schema.Maybe_Null, nil},
					f5: _Strang__Maybe{schema.Maybe_Value, &_Strang{"e"}},
				})

				// Test lookup methods.
				testLookups_vvzzv(t, n)
			})
			t.Run("not setting optionals", func(t *testing.T) {
				// Test building.
				n := build_vuvuv(t, _Stract__Style{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stract{
					f1: _Strang{"a"},
					f2: _Strang__Maybe{schema.Maybe_Absent, nil},
					f3: _Strang__Maybe{schema.Maybe_Value, &_Strang{"c"}},
					f4: _Strang__Maybe{schema.Maybe_Absent, nil},
					f5: _Strang__Maybe{schema.Maybe_Value, &_Strang{"e"}},
				})

				// Test lookup methods.
				testLookups_vuvuv(t, n)
				testIteration_vuvuv_repr(t, n.Representation())
			})
			t.Run("absent trailing optionals", func(t *testing.T) {
				// Trailing optionals are especially touchy in a few details of iterators.
				n := build_vvzuu(t, _Stract__Style{})
				testIteration_vvzuu_repr(t, n.Representation()) // ow
			})
		})
		// These following tests check that the renames behavior of map representations work.
		t.Run("representatiosn with renames build and read", func(t *testing.T) {
			t.Run("all fields set", func(t *testing.T) {
				// Test building.
				n := build_vvvvv_repr(t, _Stract__ReprStyle{})

				// Assert directly against expected memory state.
				Wish(t, n, ShouldEqual, &_Stract{
					f1: _Strang{"a"},
					f2: _Strang__Maybe{schema.Maybe_Value, &_Strang{"b"}},
					f3: _Strang__Maybe{schema.Maybe_Value, &_Strang{"c"}},
					f4: _Strang__Maybe{schema.Maybe_Value, &_Strang{"d"}},
					f5: _Strang__Maybe{schema.Maybe_Value, &_Strang{"e"}},
				})

				// Test lookup methods... for the type-level node.
				testLookups_vvvvv(t, n)
				testIteration_vvvvv(t, n)

				// Now test lookups and iterations for the representation.
				testLookups_vvvvv_repr(t, n.Representation())
				testIteration_vvvvv_repr(t, n.Representation())
			})
		})
	})
}
