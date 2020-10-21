package dagcbor

import (
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestFunBlocks(t *testing.T) {
	t.Run("zero length link", func(t *testing.T) {
		// This fixture has a zero length link -- not even the multibase byte (which dag-cbor insists must be zero) is there.
		buf := strings.NewReader("\x8d\x8d\x97\xd8*@")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, ErrInvalidMultibase)
	})
	t.Run("fuzz001", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9a\xff000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, ErrAllocationBudgetExceeded)
	})
	t.Run("fuzz002", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9a\xff000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, ErrAllocationBudgetExceeded)
	})
	t.Run("fuzz003", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9f\x9f\x9f\x9f\x9f\x9f\x9f\xbb00000000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, ErrAllocationBudgetExceeded)
	})
}
