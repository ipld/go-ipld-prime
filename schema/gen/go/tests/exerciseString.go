package tests

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
)

// All these functions are here in a subpackage,
// and have to be an exported symbol,
// because we're going to call them from the generated packages.
// If they were unexported, it wouldn't work;
// and if we put them in the gen package, they bloat runtime.

func ExerciseString(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
	ns := getStyleByName("String")
	t.Run("string operations work", func(t *testing.T) {
		nb := ns.NewBuilder()
		nb.AssignString("woiu")
		n := nb.Build()
		t.Logf("%v\n", n)
	})
	t.Run("null is rejected", func(t *testing.T) {
		nb := ns.NewBuilder()
		nb.AssignNull()

	})
}
