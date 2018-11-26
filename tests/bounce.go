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
		n1 := newNode()
		n1.SetString("asdf")
		Wish(t, fluent.WrapNode(n1).AsString(), ShouldEqual, "asdf")
	})
	t.Run("short array node bounce", func(t *testing.T) {
		n1 := newNode()
		n10 := newNode()
		n10.SetString("asdf")
		n1.SetIndex(0, n10)
		Wish(t, fluent.WrapNode(n1).TraverseIndex(0).AsString(), ShouldEqual, "asdf")
	})
}
