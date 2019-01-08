package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"

	. "github.com/warpfork/go-wish"
)

type MutableNodeFactory func() ipld.MutableNode

func TestNodes(t *testing.T, newNode MutableNodeFactory) {
	t.Run("string node bounce", func(t *testing.T) {
		n0 := newNode()
		n0.SetString("asdf")
		Wish(t, fluent.WrapNode(n0).AsString(), ShouldEqual, "asdf")
	})
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
