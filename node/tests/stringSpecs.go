package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/datamodel"
)

func SpecTestString(t *testing.T, np datamodel.NodePrototype) {
	t.Run("string node", func(t *testing.T) {
		nb := np.NewBuilder()
		err := nb.AssignString("asdf")
		Wish(t, err, ShouldEqual, nil)
		n := nb.Build()

		Wish(t, n.Kind(), ShouldEqual, datamodel.Kind_String)
		Wish(t, n.IsNull(), ShouldEqual, false)
		x, err := n.AsString()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, x, ShouldEqual, "asdf")
	})
}
