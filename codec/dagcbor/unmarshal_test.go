package dagcbor

import (
	"runtime"
	"strings"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestFunBlocks(t *testing.T) {
	t.Run("zero length link", func(t *testing.T) {
		// This fixture has a zero length link -- not even the multibase byte (which dag-cbor insists must be zero) is there.
		buf := strings.NewReader("\x8d\x8d\x97\xd8*@")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.Equals, ErrInvalidMultibase)
	})
	t.Run("fuzz001", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9a\xff000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		if runtime.GOARCH == "386" {
			// TODO: fix refmt to properly handle 64-bit ints on 32-bit runtime
			qt.Assert(t, err.Error(), qt.Equals, "cbor: positive integer is out of length")
		} else {
			qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
		}
	})
	t.Run("fuzz002", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9a\xff000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		if runtime.GOARCH == "386" {
			// TODO: fix refmt to properly handle 64-bit ints on 32-bit
			qt.Assert(t, err.Error(), qt.Equals, "cbor: positive integer is out of length")
		} else {
			qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
		}
	})
	t.Run("fuzz003", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9f\x9f\x9f\x9f\x9f\x9f\x9f\xbb00000000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		if runtime.GOARCH == "386" {
			// TODO: fix refmt to properly handle 64-bit ints on 32-bit
			qt.Assert(t, err.Error(), qt.Equals, "cbor: positive integer is out of length")
		} else {
			qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
		}
	})
	t.Run("undef", func(t *testing.T) {
		// This fixture tests that we tolerate cbor's "undefined" token (even though it's noncanonical and you shouldn't use it),
		// and that it becomes a null in the data model level.
		buf := strings.NewReader("\xf7")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, nb.Build().Kind(), qt.Equals, datamodel.Kind_Null)
	})
	t.Run("extra bytes", func(t *testing.T) {
		buf := strings.NewReader("\xa0\x00")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.Equals, ErrTrailingBytes)
	})
}
