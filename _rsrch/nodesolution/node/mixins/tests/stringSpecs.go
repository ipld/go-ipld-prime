package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

func SpecTestString(t *testing.T, ns ipld.NodeStyle) {
	t.Run("string node", func(t *testing.T) {
		nb := ns.NewBuilder()
		err := nb.AssignString("asdf")
		Wish(t, err, ShouldEqual, nil)
		n := nb.Build()

		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, n.IsNull(), ShouldEqual, false)
		x, err := n.AsString()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, x, ShouldEqual, "asdf")
	})
}
