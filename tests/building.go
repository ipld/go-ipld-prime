package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"

	. "github.com/warpfork/go-wish"
)

func TestBuildingScalars(t *testing.T, nb ipld.NodeBuilder) {
	t.Run("null node", func(t *testing.T) {
		n := fluent.WrapNodeBuilder(nb).CreateNull()
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Null)
		Wish(t, n.IsNull(), ShouldEqual, true)
	})
	t.Run("bool node", func(t *testing.T) {
		n := fluent.WrapNodeBuilder(nb).CreateBool(true)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Bool)
		Wish(t, n.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n).AsBool(), ShouldEqual, true)
	})
	t.Run("int node", func(t *testing.T) {
		n := fluent.WrapNodeBuilder(nb).CreateInt(17)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Int)
		Wish(t, n.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n).AsInt(), ShouldEqual, 17)
	})
	t.Run("float node", func(t *testing.T) {
		n := fluent.WrapNodeBuilder(nb).CreateFloat(0.122)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Float)
		Wish(t, n.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n).AsFloat(), ShouldEqual, 0.122)
	})
	t.Run("string node", func(t *testing.T) {
		n := fluent.WrapNodeBuilder(nb).CreateString("asdf")
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, n.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n).AsString(), ShouldEqual, "asdf")
	})
	t.Run("bytes node", func(t *testing.T) {
		n := fluent.WrapNodeBuilder(nb).CreateBytes([]byte{65, 66})
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Bytes)
		Wish(t, n.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n).AsBytes(), ShouldEqual, []byte{65, 66})
	})
}

func TestBuildingRecursives(t *testing.T, nb ipld.NodeBuilder) {
	t.Run("short list node", func(t *testing.T) {
		nb := fluent.WrapNodeBuilder(nb)
		n := nb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.Append(vnb.CreateString("asdf"))
		})
		Wish(t, fluent.WrapNode(n).LookupIndex(0).AsString(), ShouldEqual, "asdf")
	})
	t.Run("nested list node", func(t *testing.T) {
		nb := fluent.WrapNodeBuilder(nb)
		n := nb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.Append(vnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
				lb.Append(vnb.CreateString("asdf"))
			}))
			lb.Append(vnb.CreateString("quux"))
		})
		nf := fluent.WrapNode(n)
		Wish(t, nf.ReprKind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, nf.LookupIndex(0).LookupIndex(0).AsString(), ShouldEqual, "asdf")
		Wish(t, nf.LookupIndex(1).AsString(), ShouldEqual, "quux")
	})
	t.Run("long list node", func(t *testing.T) {
		nb := fluent.WrapNodeBuilder(nb)
		n := nb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.Set(9, vnb.CreateString("quux"))
			lb.Set(19, vnb.CreateString("quuux"))
		})
		nf := fluent.WrapNode(n)
		Wish(t, nf.ReprKind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, nf.Length(), ShouldEqual, 20)
		Wish(t, nf.LookupIndex(9).AsString(), ShouldEqual, "quux")
		Wish(t, nf.LookupIndex(19).AsString(), ShouldEqual, "quuux")
	})

	// todo map tests

	// todo list append tests
	//  (appends will require putting the GetNodeBuilder on Node iface!)

	// todo list append at odd sizes tests
}
