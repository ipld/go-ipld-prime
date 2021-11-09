package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
)

func SpecTestString(t *testing.T, np datamodel.NodePrototype) {
	t.Run("string node", func(t *testing.T) {
		nb := np.NewBuilder()
		err := nb.AssignString("asdf")
		qt.Check(t, err, qt.IsNil)
		n := nb.Build()

		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_String)
		qt.Check(t, n.IsNull(), qt.IsFalse)
		x, err := n.AsString()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, x, qt.Equals, "asdf")
	})
}
