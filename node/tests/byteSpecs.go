package tests

import (
	"io"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
)

func SpecTestBytes(t *testing.T, np datamodel.NodePrototype) {
	t.Run("byte node", func(t *testing.T) {
		nb := np.NewBuilder()
		err := nb.AssignBytes([]byte("asdf"))
		qt.Check(t, err, qt.IsNil)
		n := nb.Build()

		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Bytes)
		qt.Check(t, n.IsNull(), qt.IsFalse)
		x, err := n.AsBytes()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, x, qt.DeepEquals, []byte("asdf"))

		lbn, ok := n.(datamodel.LargeBytesNode)
		if ok {
			str, err := lbn.AsLargeBytes()
			qt.Check(t, err, qt.IsNil)
			bytes, err := io.ReadAll(str)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, bytes, qt.DeepEquals, []byte("asdf"))
		}

	})
}
