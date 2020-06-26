package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

// TestStructsContainingMaybe checks all the variations of "nullable" and "optional" on struct fields.
// It does this twice: once for the child maybes being implemented with pointers,
// and once with maybes implemented as embeds.
// The child values are scalars.
//
// Both type-level generic build and access as well as representation build and access are exercised;
// the representation used is map (the native representation for structs).
func TestStructsContainingMaybe(t *testing.T) {
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
		Wish(t, plzStr(n.LookupByString("f1")), ShouldEqual, "a")
		Wish(t, plzStr(n.LookupByString("f2")), ShouldEqual, "b")
		Wish(t, plzStr(n.LookupByString("f3")), ShouldEqual, "c")
		Wish(t, plzStr(n.LookupByString("f4")), ShouldEqual, "d")
		Wish(t, plzStr(n.LookupByString("f5")), ShouldEqual, "e")
	}
	testLookups_vvvvv_repr := func(t *testing.T, n ipld.Node) {
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, n.Length(), ShouldEqual, 5)
		Wish(t, plzStr(n.LookupByString("r1")), ShouldEqual, "a")
		Wish(t, plzStr(n.LookupByString("r2")), ShouldEqual, "b")
		Wish(t, plzStr(n.LookupByString("r3")), ShouldEqual, "c")
		Wish(t, plzStr(n.LookupByString("r4")), ShouldEqual, "d")
		Wish(t, plzStr(n.LookupByString("f5")), ShouldEqual, "e")
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
		Wish(t, plzStr(n.LookupByString("f1")), ShouldEqual, "a")
		Wish(t, plzStr(n.LookupByString("f2")), ShouldEqual, "b")
		Wish(t, erp(n.LookupByString("f3")), ShouldEqual, ipld.Null)
		Wish(t, erp(n.LookupByString("f4")), ShouldEqual, ipld.Null)
		Wish(t, plzStr(n.LookupByString("f5")), ShouldEqual, "e")
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
		Wish(t, plzStr(n.LookupByString("f1")), ShouldEqual, "a")
		Wish(t, erp(n.LookupByString("f2")), ShouldEqual, ipld.Undef)
		Wish(t, plzStr(n.LookupByString("f3")), ShouldEqual, "c")
		Wish(t, erp(n.LookupByString("f4")), ShouldEqual, ipld.Undef)
		Wish(t, plzStr(n.LookupByString("f5")), ShouldEqual, "e")
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

	// Okay, now the test actions:
	test := func(t *testing.T, ns ipld.NodeStyle, nsr ipld.NodeStyle) {
		t.Run("all fields set", func(t *testing.T) {
			t.Run("typed-create", func(t *testing.T) {
				n := build_vvvvv(t, ns)
				t.Run("typed-read", func(t *testing.T) {
					testLookups_vvvvv(t, n)
					testIteration_vvvvv(t, n)
				})
				t.Run("repr-read", func(t *testing.T) {
					testLookups_vvvvv_repr(t, n.Representation())
					testIteration_vvvvv_repr(t, n.Representation())
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				Wish(t, build_vvvvv_repr(t, nsr), ShouldEqual, build_vvvvv(t, ns))
			})
		})
		t.Run("setting nulls", func(t *testing.T) {
			t.Run("typed-create", func(t *testing.T) {
				n := build_vvzzv(t, ns)
				t.Run("typed-read", func(t *testing.T) {
					testLookups_vvzzv(t, n)
				})
				t.Run("repr-read", func(t *testing.T) {
					// nyi
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				Wish(t, build_vvzzv_repr(t, nsr), ShouldEqual, build_vvzzv(t, ns))
			})
		})
		t.Run("absent optionals", func(t *testing.T) {
			t.Run("typed-create", func(t *testing.T) {
				n := build_vuvuv(t, ns)
				t.Run("typed-read", func(t *testing.T) {
					testLookups_vuvuv(t, n)
				})
				t.Run("repr-read", func(t *testing.T) {
					testIteration_vuvuv_repr(t, n.Representation())
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				// nyi
			})
		})
		t.Run("absent trailing optionals", func(t *testing.T) {
			// Trailing optionals are especially touchy in a few details of iterators, so this gets an extra focused test.
			t.Run("typed-create", func(t *testing.T) {
				n := build_vvzuu(t, ns)
				t.Run("typed-read", func(t *testing.T) {
					// Not very interesting; still returns absent explicitly, same as 'vuvuv' scenario.
				})
				t.Run("repr-read", func(t *testing.T) {
					testIteration_vvzuu_repr(t, n.Representation())
				})
			})
			t.Run("repr-create", func(t *testing.T) {
				// nyi
			})
		})
	}

	// Do most of the type declarations.
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Stroct",
		[]schema.StructField{
			// Every field in this struct (including their order) is exercising an interesting case...
			schema.SpawnStructField("f1", ts.TypeByName("String"), false, false), // plain field.
			schema.SpawnStructField("f2", ts.TypeByName("String"), true, false),  // optional; later we have more than one optional field, nonsequentially.
			schema.SpawnStructField("f3", ts.TypeByName("String"), false, true),  // nullable; but required.
			schema.SpawnStructField("f4", ts.TypeByName("String"), true, true),   // optional and nullable; trailing optional.
			schema.SpawnStructField("f5", ts.TypeByName("String"), true, false),  // optional; and the second one in a row, trailing.
		},
		schema.SpawnStructRepresentationMap(map[string]string{
			"f1": "r1",
			"f2": "r2",
			"f3": "r3",
			"f4": "r4",
		}),
	))

	// And finally, launch tests! ...while specializing the adjunct config a bit.
	t.Run("maybe-using-embed", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = false

		prefix := "stroct"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
			test(t, getStyleByName("Stroct"), getStyleByName("Stroct.Repr"))
		})
	})
	t.Run("maybe-using-ptr", func(t *testing.T) {
		adjCfg.maybeUsesPtr["String"] = true

		prefix := "stroct2"
		pkgName := "main"
		genAndCompileAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
			test(t, getStyleByName("Stroct"), getStyleByName("Stroct.Repr"))
		})
	})
}
