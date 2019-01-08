package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"

	. "github.com/warpfork/go-wish"
)

type MutableNodeFactory func() ipld.MutableNode

func TestScalars(t *testing.T, newNode MutableNodeFactory) {
	t.Run("null node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetNull()
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_Null)
		Wish(t, n0.IsNull(), ShouldEqual, true)
	})
	t.Run("bool node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetBool(true)
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_Bool)
		Wish(t, n0.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n0).AsBool(), ShouldEqual, true)
	})
	t.Run("int node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetInt(17)
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_Int)
		Wish(t, n0.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n0).AsInt(), ShouldEqual, 17)
	})
	t.Run("float node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetFloat(0.122)
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_Float)
		Wish(t, n0.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n0).AsFloat(), ShouldEqual, 0.122)
	})
	t.Run("string node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetString("asdf")
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, n0.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n0).AsString(), ShouldEqual, "asdf")
	})
	t.Run("bytes node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetBytes([]byte{65, 66})
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_Bytes)
		Wish(t, n0.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n0).AsBytes(), ShouldEqual, []byte{65, 66})
	})
	t.Run("link node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetBool(true)
		Wish(t, n0.Kind(), ShouldEqual, ipld.ReprKind_Bool)
		Wish(t, n0.IsNull(), ShouldEqual, false)
		Wish(t, fluent.WrapNode(n0).AsBool(), ShouldEqual, true)
	})
}

func TestRecursives(t *testing.T, newNode MutableNodeFactory) {
	t.Run("short list node bounce", func(t *testing.T) {
		n0 := newNode()
		n00 := newNode()
		n00.SetString("asdf")
		n0.SetIndex(0, n00)
		Wish(t, fluent.WrapNode(n0).TraverseIndex(0).AsString(), ShouldEqual, "asdf")
	})
	t.Run("nested list node bounce", func(t *testing.T) {
		n0 := newNode()
		n00 := newNode()
		n0.SetIndex(0, n00)
		n000 := newNode()
		n000.SetString("asdf")
		n00.SetIndex(0, n000)
		n01 := newNode()
		n01.SetString("quux")
		n0.SetIndex(1, n01)
		nf := fluent.WrapNode(n0)
		Wish(t, nf.Kind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, nf.TraverseIndex(0).TraverseIndex(0).AsString(), ShouldEqual, "asdf")
		Wish(t, nf.TraverseIndex(1).AsString(), ShouldEqual, "quux")
	})
}
