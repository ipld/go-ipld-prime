package fluent_test

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestReflect(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		n, err := fluent.Reflect(basicnode.Prototype.Any, map[string]interface{}{
			"k1": "fine",
			"k2": "super",
			"k3": map[string]string{
				"k31": "thanks",
				"k32": "for",
				"k33": "asking",
			},
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		t.Run("CorrectContents", func(t *testing.T) {
			Wish(t, n.Length(), ShouldEqual, 3)
			Wish(t, must.String(must.Node(n.LookupByString("k1"))), ShouldEqual, "fine")
			Wish(t, must.String(must.Node(n.LookupByString("k2"))), ShouldEqual, "super")
			n := must.Node(n.LookupByString("k3"))
			Wish(t, n.Length(), ShouldEqual, 3)
			Wish(t, must.String(must.Node(n.LookupByString("k31"))), ShouldEqual, "thanks")
			Wish(t, must.String(must.Node(n.LookupByString("k32"))), ShouldEqual, "for")
			Wish(t, must.String(must.Node(n.LookupByString("k33"))), ShouldEqual, "asking")
		})
		t.Run("CorrectOrder", func(t *testing.T) {
			itr := n.MapIterator()
			k, _, _ := itr.Next()
			Wish(t, must.String(k), ShouldEqual, "k1")
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "k2")
			k, v, _ := itr.Next()
			Wish(t, must.String(k), ShouldEqual, "k3")
			itr = v.MapIterator()
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "k31")
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "k32")
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "k33")
		})
	})
	t.Run("Struct", func(t *testing.T) {
		type Woo struct {
			A string
			B string
		}
		type Whee struct {
			X string
			Z string
			M Woo
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Whee{
			X: "fine",
			Z: "super",
			M: Woo{"thanks", "really"},
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		t.Run("CorrectContents", func(t *testing.T) {
			Wish(t, n.Length(), ShouldEqual, 3)
			Wish(t, must.String(must.Node(n.LookupByString("X"))), ShouldEqual, "fine")
			Wish(t, must.String(must.Node(n.LookupByString("Z"))), ShouldEqual, "super")
			n := must.Node(n.LookupByString("M"))
			Wish(t, n.Length(), ShouldEqual, 2)
			Wish(t, must.String(must.Node(n.LookupByString("A"))), ShouldEqual, "thanks")
			Wish(t, must.String(must.Node(n.LookupByString("B"))), ShouldEqual, "really")
		})
		t.Run("CorrectOrder", func(t *testing.T) {
			itr := n.MapIterator()
			k, _, _ := itr.Next()
			Wish(t, must.String(k), ShouldEqual, "X")
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "Z")
			k, v, _ := itr.Next()
			Wish(t, must.String(k), ShouldEqual, "M")
			itr = v.MapIterator()
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "A")
			k, _, _ = itr.Next()
			Wish(t, must.String(k), ShouldEqual, "B")
		})
	})
	t.Run("NamedString", func(t *testing.T) {
		type Foo string
		type Bar struct {
			Z Foo
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Bar{"foo"})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, must.String(must.Node(n.LookupByString("Z"))), ShouldEqual, "foo")
	})
	t.Run("Interface", func(t *testing.T) {
		type Zaz struct {
			Z interface{}
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Zaz{map[string]interface{}{"wow": "wee"}})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		n, err = n.LookupByString("Z")
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, must.String(must.Node(n.LookupByString("wow"))), ShouldEqual, "wee")
	})
}
