package methodsets

import (
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestViaTypedef(t *testing.T) {
	x := Thing{"alpha", "beta"}

	x.Pow()
	Wish(t, x, ShouldEqual, Thing{"base", "beta"})

	x2 := FlipTypedef(&x)
	Wish(t, x2, ShouldEqual, &Thing2ViaTypedef{"base", "beta"})

	x2.Pow()
	Wish(t, x2, ShouldEqual, &Thing2ViaTypedef{"typedef", "beta"})

	Wish(t, x, ShouldEqual, Thing{"typedef", "beta"}) // ! effects propagate back to original.

	x.Pow()
	Wish(t, x2, ShouldEqual, &Thing2ViaTypedef{"base", "beta"}) // ! and still also vice versa.

	// it's not just that we care about retaining mutability (though that's sometimes useful);
	// it's that a 'yes' to that directly implies 'yes' to "can we get this pov without any allocations".
}
