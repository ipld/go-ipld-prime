package fluent_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
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
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		t.Run("CorrectContents", func(t *testing.T) {
			qt.Check(t, n.Length(), qt.Equals, int64(3))
			qt.Check(t, must.String(must.Node(n.LookupByString("k1"))), qt.Equals, "fine")
			qt.Check(t, must.String(must.Node(n.LookupByString("k2"))), qt.Equals, "super")
			n := must.Node(n.LookupByString("k3"))
			qt.Check(t, n.Length(), qt.Equals, int64(3))
			qt.Check(t, must.String(must.Node(n.LookupByString("k31"))), qt.Equals, "thanks")
			qt.Check(t, must.String(must.Node(n.LookupByString("k32"))), qt.Equals, "for")
			qt.Check(t, must.String(must.Node(n.LookupByString("k33"))), qt.Equals, "asking")
		})
		t.Run("CorrectOrder", func(t *testing.T) {
			itr := n.MapIterator()
			k, _, _ := itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "k1")
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "k2")
			k, v, _ := itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "k3")
			itr = v.MapIterator()
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "k31")
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "k32")
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "k33")
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
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		t.Run("CorrectContents", func(t *testing.T) {
			qt.Check(t, n.Length(), qt.Equals, int64(3))
			qt.Check(t, must.String(must.Node(n.LookupByString("X"))), qt.Equals, "fine")
			qt.Check(t, must.String(must.Node(n.LookupByString("Z"))), qt.Equals, "super")
			n := must.Node(n.LookupByString("M"))
			qt.Check(t, n.Length(), qt.Equals, int64(2))
			qt.Check(t, must.String(must.Node(n.LookupByString("A"))), qt.Equals, "thanks")
			qt.Check(t, must.String(must.Node(n.LookupByString("B"))), qt.Equals, "really")
		})
		t.Run("CorrectOrder", func(t *testing.T) {
			itr := n.MapIterator()
			k, _, _ := itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "X")
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "Z")
			k, v, _ := itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "M")
			itr = v.MapIterator()
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "A")
			k, _, _ = itr.Next()
			qt.Check(t, must.String(k), qt.Equals, "B")
		})
	})
	t.Run("NamedString", func(t *testing.T) {
		type Foo string
		type Bar struct {
			Z Foo
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Bar{"foo"})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		qt.Check(t, must.String(must.Node(n.LookupByString("Z"))), qt.Equals, "foo")
	})
	t.Run("Interface", func(t *testing.T) {
		type Zaz struct {
			Z interface{}
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Zaz{map[string]interface{}{"wow": "wee"}})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		n, err = n.LookupByString("Z")
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		qt.Check(t, must.String(must.Node(n.LookupByString("wow"))), qt.Equals, "wee")
	})
	t.Run("Bytes", func(t *testing.T) {
		n, err := fluent.Reflect(basicnode.Prototype.Any, []byte{0x1, 0x2, 0x3})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Bytes)
		b, err := n.AsBytes()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, b, qt.DeepEquals, []byte{0x1, 0x2, 0x3})
	})
	t.Run("NamedBytes", func(t *testing.T) {
		type Foo []byte
		type Bar struct {
			Z Foo
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Bar{[]byte{0x1, 0x2, 0x3}})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		n, err = n.LookupByString("Z")
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Bytes)
		b, err := n.AsBytes()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, b, qt.DeepEquals, []byte{0x1, 0x2, 0x3})
	})
	t.Run("InterfaceContainingBytes", func(t *testing.T) {
		type Zaz struct {
			Z interface{}
		}
		n, err := fluent.Reflect(basicnode.Prototype.Any, Zaz{[]byte{0x1, 0x2, 0x3}})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		n, err = n.LookupByString("Z")
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Bytes)
		b, err := n.AsBytes()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, b, qt.DeepEquals, []byte{0x1, 0x2, 0x3})
	})
}
